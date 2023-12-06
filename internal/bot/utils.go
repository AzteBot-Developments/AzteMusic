package main

import (
	"context"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/disgolink/v3/lavalink"
	"github.com/disgoorg/snowflake/v2"
)

// Plays a YT track or playlist from the given source URL.
func (b *Bot) PlayOnStartupFromUrl(event *discordgo.Ready, url string, repeatCount int) error {

	playlistUrl := url

	if !urlPattern.MatchString(playlistUrl) && !searchPattern.MatchString(playlistUrl) {
		playlistUrl = lavalink.SearchTypeYouTube.Apply(playlistUrl)
	}

	player := b.Lavalink.Player(snowflake.MustParse(GuildId))
	queue := b.Queues.Get(GuildId)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var toPlay *lavalink.Track
	b.Lavalink.BestNode().LoadTracksHandler(ctx, playlistUrl, disgolink.NewResultHandler(
		func(track lavalink.Track) {
			if player.Track() == nil {
				toPlay = &track
			} else {
				queue.Add(track)
			}
		},
		func(playlist lavalink.Playlist) {
			if player.Track() == nil {
				toPlay = &playlist.Tracks[0]
				queue.Add(playlist.Tracks[1:]...)
				// Repeat the queue `repeatCount` times
				// TODO
			} else {
				queue.Add(playlist.Tracks...)
			}
		},
		func(tracks []lavalink.Track) {
			if player.Track() == nil {
				toPlay = &tracks[0]
			} else {
				queue.Add(tracks[0])
			}
		},
		nil,
		nil,
	))
	if toPlay == nil {
		return nil
	}

	if err := b.Session.ChannelVoiceJoinManual(GuildId, DesignatedChannelId, false, false); err != nil {
		log.Fatalf("Could not join channel (2) at startup: %v", err)
		return err
	}

	return player.Update(context.TODO(), lavalink.WithTrack(*toPlay))
}
