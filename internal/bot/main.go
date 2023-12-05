package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"regexp"
	"syscall"
	"time"
)

var (
	urlPattern    = regexp.MustCompile("^https?://[-a-zA-Z0-9+&@#/%?=~_|!:,.;]*[-a-zA-Z0-9+&@#/%=~_|]?")
	searchPattern = regexp.MustCompile(`^(.{2})search:(.+)`)

	Token   = os.Getenv("BOT_TOKEN")
	GuildId = os.Getenv("GUILD_ID")

	DesignatedChannelId   = os.Getenv("DESIGNATED_VOICE_CHANNEL_ID")
	DesignatedPlaylistUrl = os.Getenv("DESIGNATED_PLAYLIST_URL")
	StatusText            = os.Getenv("STATUS_TEXT")

	b = &Bot{
		Queues: &QueueManager{
			queues: make(map[string]*Queue),
		},
	}
)

func main() {

	// Retrieve an authenticated Discord bot session through the token provided as an env variable
	b.Session = GetAuthenticatedBotSession()

	// Set the required intents for the bot's operation and what states it tracks
	b.SetIntents()

	// Register the handlers for the Discord session (onReady, onVoiceUpdate, etc.)
	b.AddHandlers()

	// Connect the authenticated bot session to the Discord servers
	b.Connect()

	// Register the bot's slash commands (play, shuffle, skip, etc.)
	b.RegisterCommands()

	// Connect to the associated LavaLink server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	b.AddLavalinkNode(ctx)

	log.Printf("Discord bot is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}
