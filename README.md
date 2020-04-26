# Ayaya
A demo music bot utilizing the [arikawa](https://github.com/diamondburned/arikawa) library.

## Usage
~~Currently, you need to have my pending voice support pull request cloned to `../arikawa` depending
on where you clone this repo.~~

[arikawa](https://github.com/diamondburned/arikawa) has [merged voice support](https://github.com/diamondburned/arikawa/commit/ccf4c69801f63c561319e5459ea2c09a0f8a41ad) and [fixed some problems](https://github.com/diamondburned/arikawa/commit/51e88a47b2e21f7683cc157b36c752ccd9221f96) with my early implementation, you are now fine to use the latest release of arikawa instead of cloning my fork.

```shell script
# Clone the repositories
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
