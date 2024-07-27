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

Create an app in Discord, see [here](https://discord.com/developers/docs/quick-start/getting-started) for how.

When selecting permissions in OAuth2 menu you should select `bot` and then `Manage Messages` and `Send Messages` permissions.

Then please ensure that the channels you wish to have messages published from allow these permissions for your bot.

Channels you specify using the `Channel IDs` parameter must be announcement channels.

## Releasing

Create a new git tag and run [goreleaser](https://goreleaser.com/)

```shell
goreleaser release
```

## Future Improvements TODO

- Queue per channel so rate limit doesn't block all channels