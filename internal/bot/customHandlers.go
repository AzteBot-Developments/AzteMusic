package main

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/disgoorg/snowflake/v2"
)

// Called once the Discord servers confirm a succesful connection.
func (b *Bot) onReady(s *discordgo.Session, event *discordgo.Ready) {

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

		// Play designated playlist
		if DesignatedPlaylistUrl != "" {
			_ = b.Lavalink.Player(snowflake.MustParse(GuildId))
			if err := b.playOnStartupFromUrl(event, DesignatedPlaylistUrl); err != nil {
				log.Fatalf("Could not play designated playlist (onReady): %v", err)
			}
		}
	}

}
