package main

import (
	"context"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Called once the Discord servers confirm a succesful connection.
func (b *Bot) onReady(s *discordgo.Session, event *discordgo.Ready) {

	// Initial Lavalink service connection setup
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	b.AddLavalinkNode(ctx)

	// Any initial setup for the music service !
	// i.e join designated server, play designated playlist, etc.

	// Set the playing status
	if StatusText != "" {
		s.UpdateGameStatus(0, StatusText)
	}

	// Join designated channel
	if DesignatedChannelId != "" {
		if err := s.ChannelVoiceJoinManual(GuildId, DesignatedChannelId, false, false); err != nil {
			log.Fatalf("Could not join designated voice channel (onReady): %v", err)
		}

		// Play designated playlist on loop, FOREVER :')
		if DesignatedPlaylistUrl != "" {
			if err := b.PlayOnStartupFromUrl(event, DesignatedPlaylistUrl); err != nil {
				log.Fatalf("Could not play designated playlist (onReady): %v", err)
			}
		}
	}

}
