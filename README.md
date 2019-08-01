# Discord Message Deleter

Discord doesn't provide an easy way to delete messages en masse both in DMs
or in a channel. This tool allows you to do that easily.

## Usage

```
$ ./discord-del -h
Usage of ./discord-del:
  -authorizationHeader string
    	Discord Authorization header to use when issuing Discord API requests.
  -channel string
    	Discord channel or DM to erase.
  -username string
    	(optional) Discord username of the logged in user
```

You'll need to grab the channel and the authorization header used to make 
requests from Discord. To do this:

1. Press Ctrl+Shift+I (Cmd+Shift+I on Mac). 
2. Go to the Network tab.
3. Press XHR.
4. Open a channel or a direct message.
5. Click on the request made by Discord.

The request should look like this:

```
Request URL: https://discordapp.com/api/v6/channels/12341234234342/messages?limit=50
```

The channel ID is the number in between `channels` and `messages`. The 
authorization header is in the Request Headers section.

## Deleting from direct messages

Please keep in mind that you'll need to provide your username WITHOUT the number
suffix if you want to delete direct messages. For instance, if your username is
Bob#1234, then you'll need to provide `-username=Bob` on the command line.

## Rate limiting

This tool is fairly fast, so you might get rate limited. If you find yourself
getting rate limited just run the tool several times. Please keep in mind that
*deleting messages on Discord is permanent*.
