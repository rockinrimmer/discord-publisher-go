# Discord Publisher Go

A simple go cli application to monitor Discord announcement channels and publish 
messages automatically.

>Note: Discord has a 10 publish limit per channel per hour. When over this limit messages
> will queue until allowed to publish again.
 
## Running

To find the current version run:
```shell
discord-publisher-go version
```

### Parameters

| Parameter     | Flag     | Environment Variable | Description                               |
|---------------|----------|----------------------|-------------------------------------------|
| Server IDs    | servers  | DISCORD_SERVER_IDS   | Discord server IDs to listen to as a CSV  |
| Channel IDs   | channels | DISCORD_CHANNEL_IDS  | Discord channel IDs to listen to as a CSV |
| Discord Token | token    | DISCORD_TOKEN        | Discord bot token                         |
| Debug Logging | debug    |                      | Enable debug logging                      |

### Setup

TBD

## Releasing

Create a new git tag and run [goreleaser](https://goreleaser.com/)

```shell
goreleaser release
```

## Future Improvements TODO

- Queue per channel so rate limit doesn't block all channels