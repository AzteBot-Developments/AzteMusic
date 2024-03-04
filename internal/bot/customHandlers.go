package main

import (
	"context"
	"log"
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

	// Any initial setup for the music service !
	// i.e join designated server, play designated playlist, etc.

	// Set the playing status
	if StatusText != "" {
		s.UpdateGameStatus(0, StatusText)
	}

	// Join designated channel
	repeatPlaylistCount := 3

	if DesignatedChannelId != "" {
		if err := s.ChannelVoiceJoinManual(GuildId, DesignatedChannelId, false, false); err != nil {
			log.Fatalf("Could not join designated voice channel (onReady): %v", err)
		}

		// Play designated playlist on loop, FOREVER :')
		if DesignatedPlaylistUrl != "" {
			if err := b.PlayOnStartupFromSource(event, DesignatedPlaylistUrl, repeatPlaylistCount); err != nil {
				log.Fatalf("Could not play designated playlist (onReady): %v", err)
			}
		}

		// Also run a cron to check whether there is anything playing - if there isn't, shuffle and play the designated playlist
		var numSec int = 15
		ticker := time.NewTicker(time.Duration(numSec) * time.Second)
		quit := make(chan struct{})
		go func() {
			for {
				select {
				case <-ticker.C:
					serverQueue := b.Queues.Get(GuildId)
					if len(serverQueue.Tracks) == 0 || !ServiceIsPlayingTrack(b, GuildId) {
						if err := b.PlayOnStartupFromSource(event, DesignatedPlaylistUrl, repeatPlaylistCount); err != nil {
							log.Printf("Could not play designated playlist (onReady CRON): %v", err)
						}
					}
				case <-quit:
					ticker.Stop()
					return
				}
			}
		}()
	}

}
