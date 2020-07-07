package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"

	"github.com/diamondburned/arikawa/gateway"
	"github.com/diamondburned/arikawa/state"
	"github.com/diamondburned/arikawa/utils/wsutil"
	"github.com/diamondburned/arikawa/voice"
	"github.com/diamondburned/arikawa/voice/voicegateway"

	"github.com/matthewpi/ayaya/dca"
	"github.com/matthewpi/ayaya/ytdl"
)

type Bot struct {
	session *state.State
	voice   *voice.Voice
}

func main() {
	var token = os.Getenv("BOT_TOKEN")
	if token == "" {
		log.Fatalln("No $BOT_TOKEN given.")
	}

	b := &Bot{}

	var err error
	b.session, err = state.New("Bot " + token)
	if err != nil {
		panic(err)
	}

	wsutil.WSDebug = func(v ...interface{}) {
		log.Println(v...)
	}

	b.session.Gateway.ErrorLog = func(err error) {
		fmt.Printf("ERROR | received error from gateway: %v\n", err)
	}

	if err := b.session.Open(); err != nil {
		panic(err)
	}

	b.session.AddHandler(b.onReady)
	b.session.AddHandler(b.onMessageCreate)

	// Add the voice repository.
	b.voice = voice.NewVoice(b.session)

	signals := make(chan os.Signal)
	signal.Notify(signals, os.Interrupt)
	<-signals
	fmt.Println("WARN  | received signal from system, shutting down..")

	b.session.Close()
}

func (b *Bot) onReady(*gateway.ReadyEvent) {
	fmt.Println("INFO  | ready!")
}

func (b *Bot) onMessageCreate(e *gateway.MessageCreateEvent) {
	fmt.Printf("DEBUG | *gateway.MessageCreateEvent: %s\n", e.ID.String())
	fmt.Printf("INFO  | %s\n", e.Message.Content)

	if e.Message.Author.Bot {
		fmt.Println("DEBUG | message was sent by a bot")
		return
	}

	if e.Message.Author.DiscordSystem {
		fmt.Println("DEBUG | message was sent by a Discord System")
		return
	}

	message := e.Message.Content
	if len(message) < 6 {
		fmt.Println("DEBUG | message has not enough content")
		return
	}

	if message[:5] != ";play" {
		fmt.Println("DEBUG | message is not the play command")
		return
	}

	videoURL := message[6:]
	if len(videoURL) < 1 {
		fmt.Println("DEBUG | message does not have a video url")
		return
	}

	fmt.Printf("INFO  | %s\n", videoURL)

	if !strings.HasPrefix(videoURL, "https://www.youtube.com/watch?v=") {
		fmt.Println("DEBUG | invalid url")
		return
	}

	videoInfo, err := ytdl.GetVideoInfo(context.Background(), videoURL)
	if err != nil {
		fmt.Printf("ERROR | failed to load video: %v\n", err)
		return
	}

	// Check the available video formats.
	formats := videoInfo.Formats.Extremes(ytdl.FormatAudioBitrateKey, true)
	if len(formats) < 1 {
		fmt.Println("ERROR | failed to get video formats")
		return
	}

	// Get the video's download url.
	downloadURL, err := ytdl.DefaultClient.GetDownloadURL(context.Background(), videoInfo, formats[0])
	if err != nil {
		fmt.Printf("ERROR | failed to get download url: %v\n", err)
		return
	}

	vs, err := b.session.VoiceState(e.GuildID, e.Author.ID)
	if err != nil {
		fmt.Printf("ERROR | failed to get voice state: %v\n", err)
		return
	}

	// Check if the user is not in a voice channel.
	if vs.ChannelID == 0 {
		fmt.Println("INFO  | user is not in a voice channel")
		return
	}

	// Get the channel information.
	channel, err := b.session.Channel(vs.ChannelID)
	if err != nil {
		fmt.Printf("ERROR | failed to get channel: %v\n", err)
		return
	}

	// TODO: Check bot's permissions to see if it can join the channel.

	// Connect to the voice channel.
	conn, err := b.voice.JoinChannel(e.GuildID, vs.ChannelID, false, false)
	if err != nil {
		fmt.Printf("ERROR | failed to join voice channel: %v\n", err)
		return
	}

	resp, err := http.Get(downloadURL.String())
	if err != nil {
		fmt.Printf("ERROR | failed to fetch video: %v\n", err)
		return
	}

	if err := conn.Speaking(voicegateway.Microphone); err != nil {
		fmt.Printf("ERROR | failed to start speaking: %v\n", err)
		return
	}

	options := dca.StdEncodeOptions
	options.Volume = 128
	options.Channels = 2
	options.FrameRate = 48000
	options.FrameDuration = 20
	options.RawOutput = true
	options.Bitrate = int(channel.VoiceBitrate)
	options.PacketLoss = 2
	options.Application = dca.AudioApplicationAudio
	options.CompressionLevel = 10
	options.BufferedFrames = 1024
	options.VBR = true

	// Start encoding the video.
	encodingSession, err := dca.EncodeMem(resp.Body, options)
	if err != nil {
		fmt.Printf("ERROR | failed to starty encode video: %v\n", err)
		return
	}

	done := make(chan error)
	_ = dca.NewStream(encodingSession, conn, done)
	defer encodingSession.Cleanup()

	if err = <-done; err != nil && err != io.EOF {
		fmt.Printf("ERROR | error occurred during stream: %v\n", err)
		return
	}

	fmt.Println("INFO  | finished playing")

	if err := conn.StopSpeaking(); err != nil {
		fmt.Printf("ERROR | an error occurred while trying to stop speaking: %v\n", err)
		return
	}

	if err := conn.Disconnect(); err != nil {
		fmt.Printf("ERROR | an error occurred while trying to disconnect: %v\n", err)
		return
	}
}
