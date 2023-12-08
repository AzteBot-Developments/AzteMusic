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

func (b *Bot) GetCurrentTrack() (*lavalink.Track, disgolink.Player) {
	player := b.Lavalink.ExistingPlayer(snowflake.MustParse(GuildId))
	if player == nil {
		return nil, nil
	}

	track := player.Track()
	if track == nil {
		return nil, nil
	}

	return track, player
}

func (b *Bot) AddToQueueFromSource(url string, repeatCount int) {
	playlistUrl := url

	if !urlPattern.MatchString(playlistUrl) && !searchPattern.MatchString(playlistUrl) {
		playlistUrl = lavalink.SearchTypeYouTube.Apply(playlistUrl)
	}

	queue := b.Queues.Get(GuildId)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	b.Lavalink.BestNode().LoadTracksHandler(ctx, playlistUrl, disgolink.NewResultHandler(
		func(track lavalink.Track) {
			queue.Add(track)
		},
		func(playlist lavalink.Playlist) {
			// Repeat the queue `repeatCount` times
			for i := 0; i < repeatCount; i++ {
				queue.Add(playlist.Tracks[0:]...)
			}
		},
		func(tracks []lavalink.Track) {
			queue.Add(tracks[0])
		},
		nil,
		nil,
	))
}

// Plays a YT track or playlist from the given source URL.
func (b *Bot) PlayOnStartupFromSource(event *discordgo.Ready, url string, repeatCount int) error {

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
				for i := 0; i < repeatCount; i++ {
					queue.Add(playlist.Tracks[0:]...)
				}
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
