package main

import (
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
	"slices"
	"strings"
	"syscall"
	"time"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var watchingServerIds []string
var watchingChannelIds []string
var publishQueue = make(chan *discordgo.MessageCreate, 100)

func main() {

	subCommand := ""
	if len(os.Args) > 1 {
		subCommand = os.Args[1]
	}

	switch subCommand {
	case "version":
		fmt.Printf("Discord Publisher Go %s, commit %s, built at %s", version, commit, date)
		os.Exit(0)
	default:
		log.SetOutput(os.Stdout)
		serverIds, channelIds, token, debug := parseFlags()
		watchingServerIds = serverIds
		watchingChannelIds = channelIds

		log.Println("Starting...")
		session := openSession(token, debug)
		defer closeSession(session)

		// Trigger reader to wait for messages and publish them
		go publishQueueReader(session)

		// Block until killed
		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		<-sc
	}
}

func publishQueueReader(s *discordgo.Session) {

	log.Println("Listening on publish queue")

	for {
		m := <-publishQueue
		log.Printf("Publishing message %v\n", m.ID)

		// If we get rate limited by Discord this will block here until we are allowed to publish again
		_, err := s.ChannelMessageCrosspost(m.ChannelID, m.ID)
		if err != nil {
			log.Println("Error publishing message", err)
		}
		time.Sleep(3 * time.Second)
	}
}

func parseFlags() (watchingServerIds []string, channelIds []string, discordToken string, debug bool) {
	serverIdsString := flag.String("servers", os.Getenv("DISCORD_SERVER_IDS"), "Discord server IDs to listen to as a CSV")
	channelIdsString := flag.String("channels", os.Getenv("DISCORD_CHANNEL_IDS"), "Discord channel IDs to listen to as a CSV")
	flag.StringVar(&discordToken, "token", os.Getenv("DISCORD_TOKEN"), "Discord bot token")
	flag.BoolVar(&debug, "debug", false, "Enable debug logging")
	flag.Parse()

	watchingServerIds = strings.Split(*serverIdsString, ",")
	if len(watchingServerIds) == 0 {
		panic("No discord server IDs, must set 'servers' flag or environment variable DISCORD_SERVER_IDS")
	}
	log.Println("Monitoring servers: ", watchingServerIds)

	watchingChannelIds = strings.Split(*channelIdsString, ",")
	if len(watchingChannelIds) == 0 {
		panic("No discord channel IDs, must set 'channels' flag or environment variable DISCORD_CHANNEL_IDS")
	}
	log.Println("Monitoring channels: ", watchingChannelIds)

	if discordToken == "" {
		panic("No discord token set, must set 'token' flag or environment variable DISCORD_TOKEN")
	}

	return watchingServerIds, watchingChannelIds, discordToken, debug
}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	log.Println("Ready...")
	_ = s.UpdateCustomStatus("Waiting for messages")
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	log.Printf("Received a messageCreate - Server: %v Channel: %v \n", m.GuildID, m.ChannelID)

	if slices.Contains(watchingServerIds, m.GuildID) && slices.Contains(watchingChannelIds, m.ChannelID) {
		log.Printf("Adding message to publish queue %v, queue size:%v\n", m.ID, len(publishQueue))
		publishQueue <- m
	}
}

func openSession(discordToken string, debug bool) *discordgo.Session {
	session, err := discordgo.New("Bot " + discordToken)
	if err != nil {
		panic(err)
	}

	if debug {
		log.Printf("Debug logging enabled")
		session.LogLevel = discordgo.LogInformational
	}

	session.AddHandler(ready)
	session.AddHandler(messageCreate)

	// Open the websocket and begin listening.
	err = session.Open()
	if err != nil {
		panic(err)
	}
	return session
}

func closeSession(session *discordgo.Session) {
	err := session.Close()
	if err != nil {
		panic(err)
	}
	log.Println("Shutting down...")
}
