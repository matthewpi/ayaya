# ayaya
A demo music bot utilizing the [arikawa](https://github.com/diamondburned/arikawa) library.

## Using
```shell script
# Clone the repository
git clone https://github.com/matthewpi/ayaya.git

# Run the bot
cd ayaya
make && BOT_TOKEN="<BOT TOKEN>" ./ayaya
```

The command is `;play <URL>`, the bot will not send any messages, it will only log to the terminal.

**NOTE: The current URL checker requires that the URL starts with
`https://www.youtube.com/watch?v=`, make sure you don't give `youtu.be` or regular `http://`
links**

## Notes
* The included `dca` folder is a fork of [dca](https://github.com/jonas747/dca) by [jonas747](https://github.com/jonas747)
and is included because it supports the [arikawa](https://github.com/diamondburned/arikawa) voice connection.
* The included `ytdl` folder is a fork of [ytdl](https://github.com/rylio/ytdl) which has [PR #92](https://github.com/rylio/ytdl/pull/92)
applied to it because it wasn't merged into the official library.
