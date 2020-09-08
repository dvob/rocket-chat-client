# rocket-chat-client

## Usage
To use the Rocket.Chat client you have to configure the server URL, a user ID and a token.
The user ID and a token can be created under https://<YOUR_ROCKET_CHAT_SERVER>/account/tokens

```
export RC_URL=https://your-rocket-chat-server.com
export RC_USER_ID=...
export RC_TOKEN=...
```

Then you can:
```shell
# list users
rocket-chat-client list user

# list channels
rocket-chat-client list channel

# send message
rocket-chat-client send "@username" "Hello World!"
rocket-chat-client send "#channel-name" "Hello World!"
```