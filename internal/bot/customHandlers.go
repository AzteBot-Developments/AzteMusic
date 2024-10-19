package main

import (
	"context"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Called once the Discord servers confirm a succesful connection.
func (b *Bot) onReady(s *discordgo.Session, event *discordgo.Ready) {

	// Initial lavalink setup unless it was setup already
	if !b.HasLavaLinkClient {
		b.SetupLavalink()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		b.AddLavalinkNode(ctx)
	}

	// Any initial setup for the music service !!
	// i.e join designated server, play designated playlist, etc.

	// Set the playing status
	if StatusText != "" {
		s.UpdateGameStatus(0, StatusText)
	}

}
