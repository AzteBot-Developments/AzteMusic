package main

import (
	"context"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/disgoorg/json"
	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/disgolink/v3/lavalink"

	"github.com/AzteBot-Developments/AzteMusic/pkg/shared"
)

func (b *Bot) skip(event *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) error {

	player := b.Lavalink.ExistingPlayer(snowflake.MustParse(event.GuildID))
	if player == nil {
		return b.Session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "No player found.",
			},
		})
	}

	queue := b.Queues.Get(event.GuildID)
	if queue == nil {
		return b.Session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "No queue available.",
			},
		})
	}

	if len(queue.Tracks) == 0 {
		return b.Session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "There is no song to skip to.",
			},
		})
	}

	// Get next song on queue
	nextTrack, ok := queue.Next()
	if !ok {
		return b.Session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "An error ocurred while retrieving the next song from the queue.",
			},
		})
	}

	// Play immediately
	err := player.Update(context.TODO(), lavalink.WithTrack(nextTrack))
	if err != nil {
		return b.Session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "An error ocurred while skipping the current track.",
			},
		})
	}

	embed := shared.NewEmbed().
		SetTitle("🎵  Now Playing").
		SetDescription(
			fmt.Sprintf("`%s` (%s).", nextTrack.Info.Title, *nextTrack.Info.URI)).
		SetThumbnail(*nextTrack.Info.ArtworkURL).
		SetColor(000000)

	return b.Session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed.MessageEmbed},
		},
	})

}

func (b *Bot) shuffle(event *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) error {
	queue := b.Queues.Get(event.GuildID)
	if queue == nil {
		return b.Session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "No player found.",
			},
		})
	}

	queue.Shuffle()
	return b.Session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Queue shuffled.",
		},
	})
}

func (b *Bot) queueType(event *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) error {
	queue := b.Queues.Get(event.GuildID)
	if queue == nil {
		return b.Session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "No player found.",
			},
		})
	}

	queue.Type = QueueType(data.Options[0].Value.(string))
	return b.Session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Queue type set to `%s`.", queue.Type),
		},
	})
}

func (b *Bot) clearQueue(event *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) error {
	queue := b.Queues.Get(event.GuildID)
	if queue == nil {
		return b.Session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "No player found.",
			},
		})
	}

	queue.Clear()
	return b.Session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Queue cleared.",
		},
	})
}

func (b *Bot) queue(event *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) error {
	queue := b.Queues.Get(event.GuildID)
	if queue == nil {
		return b.Session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "No player found.",
			},
		})
	}

	if len(queue.Tracks) == 0 {
		return b.Session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "There are no songs on this queue.",
			},
		})
	}

	// Calculate the total length in time of the playlist
	var totalDurationSec int64
	for _, track := range queue.Tracks {
		totalDurationSec += track.Info.Length.Seconds()
	}

	// Get current track playing and add to embed
	currentTrack, player := b.GetCurrentTrack()

	// Build embed response for the queue response
	embed := shared.NewEmbed().
		SetTitle(fmt.Sprintf("🎵  Queue - %s", BotName)).
		SetDescription(
			fmt.Sprintf(
				"Currently playing `%s` (%s) at %s / %s.\n\nQueue Duration: %s\nThere are %d other songs in this queue.\nThe first %d tracks in the queue can be seen below.", currentTrack.Info.Title, *currentTrack.Info.URI, formatPosition(player.Position()), formatPosition(currentTrack.Info.Length), shared.FormatDuration(totalDurationSec), len(queue.Tracks), 10)).
		SetThumbnail(*currentTrack.Info.ArtworkURL).
		SetColor(000000)

	// Build a list of discordgo embed fields out of the songs on the queue
	for index, track := range queue.Tracks {
		title := fmt.Sprintf("%d. `%s` (%s)", index+1, track.Info.Title, *track.Info.URI)
		text := ""
		embed.AddField(title, text, false)
	}

	// Truncate & paginate
	embed.Truncate()

	return b.Session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed.MessageEmbed},
		},
	})
}

func (b *Bot) pause(event *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) error {
	player := b.Lavalink.ExistingPlayer(snowflake.MustParse(event.GuildID))
	if player == nil {
		return b.Session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "No player found.",
			},
		})
	}

	if err := player.Update(context.TODO(), lavalink.WithPaused(!player.Paused())); err != nil {
		return b.Session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Error while pausing: `%s`", err),
			},
		})
	}

	status := "playing"
	if player.Paused() {
		status = "paused"
	}

	return b.Session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Player is now %s.", status),
		},
	})
}

func (b *Bot) stop(event *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) error {
	player := b.Lavalink.ExistingPlayer(snowflake.MustParse(event.GuildID))
	if player == nil {
		return b.Session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "No player found.",
			},
		})
	}

	if err := b.Session.ChannelVoiceJoinManual(event.GuildID, "", false, false); err != nil {
		return b.Session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Error while disconnecting: `%s`.", err),
			},
		})
	}

	return b.Session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Player stopped.",
		},
	})
}

func (b *Bot) nowPlaying(event *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) error {
	player := b.Lavalink.ExistingPlayer(snowflake.MustParse(event.GuildID))
	if player == nil {
		return b.Session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "No player found.",
			},
		})
	}

	track := player.Track()
	if track == nil {
		return b.Session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "No track found.",
			},
		})
	}

	embed := shared.NewEmbed().
		SetTitle("🎵  Now Playing").
		SetDescription(
			fmt.Sprintf("`%s` (%s).\n%s / %s", track.Info.Title, *track.Info.URI, formatPosition(player.Position()), formatPosition(track.Info.Length))).
		SetThumbnail(*track.Info.ArtworkURL).
		SetColor(000000)

	return b.Session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed.MessageEmbed},
		},
	})
}

func formatPosition(position lavalink.Duration) string {
	if position == 0 {
		return "0:00"
	}
	return fmt.Sprintf("%d:%02d", position.Minutes(), position.SecondsPart())
}

func (b *Bot) play(event *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) error {
	identifier := data.Options[0].StringValue()
	if !urlPattern.MatchString(identifier) && !searchPattern.MatchString(identifier) {
		identifier = lavalink.SearchTypeYouTube.Apply(identifier)
	}

	voiceState, err := b.Session.State.VoiceState(event.GuildID, event.Member.User.ID)
	if err != nil {
		return b.Session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "I need to be connected to a voice channel before I can play any songs.",
			},
		})
	}

	if err := b.Session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	}); err != nil {
		return err
	}

	player := b.Lavalink.Player(snowflake.MustParse(event.GuildID))
	queue := b.Queues.Get(event.GuildID)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var toPlay *lavalink.Track
	b.Lavalink.BestNode().LoadTracksHandler(ctx, identifier, disgolink.NewResultHandler(
		func(track lavalink.Track) {
			// Embed build
			embed := shared.NewEmbed().
				SetTitle("🎵  Loading Track").
				SetDescription(
					fmt.Sprintf("`%s` (%s).\nTrack Duration: %s", track.Info.Title, *track.Info.URI, formatPosition(track.Info.Length))).
				SetThumbnail(*track.Info.ArtworkURL).
				SetColor(000000)

			// Interaction response
			_, _ = b.Session.InteractionResponseEdit(event.Interaction, &discordgo.WebhookEdit{
				Embeds: &[]*discordgo.MessageEmbed{embed.MessageEmbed},
			})

			// Queue handling
			if player.Track() == nil {
				toPlay = &track
			} else {
				queue.Add(track)
			}
		},
		func(playlist lavalink.Playlist) {
			// Calculate total length of loaded playlist
			var totalDurationSec int64
			for _, track := range playlist.Tracks {
				totalDurationSec += track.Info.Length.Seconds()
			}

			// Embed build
			embed := shared.NewEmbed().
				SetTitle(fmt.Sprintf("🎵  Loading Playlist `%s` with `%d` tracks", playlist.Info.Name, len(playlist.Tracks))).
				SetDescription(
					fmt.Sprintf("Playlist Duration: %s.\nFirst track in playlist: `%s` (%s)", shared.FormatDuration(totalDurationSec), playlist.Tracks[0].Info.Title, *playlist.Tracks[0].Info.URI)).
				SetThumbnail(*playlist.Tracks[0].Info.ArtworkURL).
				SetColor(000000)

			// Interaction response
			_, _ = b.Session.InteractionResponseEdit(event.Interaction, &discordgo.WebhookEdit{
				Embeds: &[]*discordgo.MessageEmbed{embed.MessageEmbed},
			})

			// Queue handling
			if player.Track() == nil {
				toPlay = &playlist.Tracks[0]
				queue.Add(playlist.Tracks[1:]...)
			} else {
				queue.Add(playlist.Tracks...)
			}
		},
		func(tracks []lavalink.Track) {
			// Embed build
			embed := shared.NewEmbed().
				SetTitle("🎵  Loading Track").
				SetDescription(
					fmt.Sprintf("`%s` (%s).\nTrack Duration: %s", tracks[0].Info.Title, *tracks[0].Info.URI, formatPosition(tracks[0].Info.Length))).
				SetThumbnail(*tracks[0].Info.ArtworkURL).
				SetColor(000000)

			// Interaction response
			_, _ = b.Session.InteractionResponseEdit(event.Interaction, &discordgo.WebhookEdit{
				Embeds: &[]*discordgo.MessageEmbed{embed.MessageEmbed},
			})

			// Queue handling
			if player.Track() == nil {
				toPlay = &tracks[0]
			} else {
				queue.Add(tracks[0])
			}
		},
		func() {
			_, _ = b.Session.InteractionResponseEdit(event.Interaction, &discordgo.WebhookEdit{
				Content: json.Ptr(fmt.Sprintf("Nothing found for: `%s`", identifier)),
			})
		},
		func(err error) {
			_, _ = b.Session.InteractionResponseEdit(event.Interaction, &discordgo.WebhookEdit{
				Content: json.Ptr(fmt.Sprintf("Error while looking up query: `%s`", err)),
			})
		},
	))
	if toPlay == nil {
		return nil
	}

	if err := b.Session.ChannelVoiceJoinManual(event.GuildID, voiceState.ChannelID, false, false); err != nil {
		return err
	}

	b.Session.UpdateGameStatus(0, toPlay.Info.Title)

	return player.Update(context.TODO(), lavalink.WithTrack(*toPlay))
}

func (b *Bot) help(event *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) error {

	embed := shared.NewEmbed().
		SetTitle(fmt.Sprintf("🎵  `%s` Slash Commands Guide", BotName)).
		SetDescription(fmt.Sprintf("See below the available slash commands for `%s`.", BotName)).
		SetThumbnail("https://i.postimg.cc/262tK7VW/148c9120-e0f0-4ed5-8965-eaa7c59cc9f2-2.jpg").
		SetColor(000000)

	// Build a list of discordgo embed fields out of the available slash commands
	for _, command := range Commands {

		text := command.Description
		title := fmt.Sprintf("`/%s`", command.Name)

		if len(command.Options) > 0 {
			for _, param := range command.Options {
				var required string
				if param.Required {
					required = "required"
				} else {
					required = "optional"
				}
				title += fmt.Sprintf(" `[%s (%s) - %s]`", param.Name, required, param.Description)
			}
		}

		embed.AddField(title, text, false)
	}

	return b.Session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed.MessageEmbed},
		},
	})
}

func (b *Bot) play_default(event *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) error {

	identifier := DesignatedPlaylistUrl
	if !urlPattern.MatchString(identifier) && !searchPattern.MatchString(identifier) {
		identifier = lavalink.SearchTypeYouTube.Apply(identifier)
	}

	if err := b.Session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	}); err != nil {
		return err
	}

	player := b.Lavalink.Player(snowflake.MustParse(event.GuildID))
	queue := b.Queues.Get(event.GuildID)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var toPlay *lavalink.Track
	b.Lavalink.BestNode().LoadTracksHandler(ctx, identifier, disgolink.NewResultHandler(
		func(track lavalink.Track) {
			// Embed build
			embed := shared.NewEmbed().
				SetTitle("🎵  Loading Default Track").
				SetDescription(
					fmt.Sprintf("`%s` (%s).\nTrack Duration: %s", track.Info.Title, *track.Info.URI, formatPosition(track.Info.Length))).
				SetThumbnail(*track.Info.ArtworkURL).
				SetColor(000000)

			// Interaction response
			_, _ = b.Session.InteractionResponseEdit(event.Interaction, &discordgo.WebhookEdit{
				Embeds: &[]*discordgo.MessageEmbed{embed.MessageEmbed},
			})

			// Queue handling
			if player.Track() == nil {
				toPlay = &track
			} else {
				queue.Add(track)
			}
		},
		func(playlist lavalink.Playlist) {
			// Calculate total length of loaded playlist
			var totalDurationSec int64
			for _, track := range playlist.Tracks {
				totalDurationSec += track.Info.Length.Seconds()
			}

			// Embed build
			embed := shared.NewEmbed().
				SetTitle(fmt.Sprintf("🎵  Loading Default Playlist `%s` with `%d` tracks", playlist.Info.Name, len(playlist.Tracks))).
				SetDescription(
					fmt.Sprintf("Playlist Duration: %s.\nFirst track in playlist: `%s` (%s)", shared.FormatDuration(totalDurationSec), playlist.Tracks[0].Info.Title, *playlist.Tracks[0].Info.URI)).
				SetThumbnail(*playlist.Tracks[0].Info.ArtworkURL).
				SetColor(000000)

			// Interaction response
			_, _ = b.Session.InteractionResponseEdit(event.Interaction, &discordgo.WebhookEdit{
				Embeds: &[]*discordgo.MessageEmbed{embed.MessageEmbed},
			})

			// Queue handling
			if player.Track() == nil {
				toPlay = &playlist.Tracks[0]
				queue.Add(playlist.Tracks[1:]...)
			} else {
				queue.Add(playlist.Tracks...)
			}
		},
		func(tracks []lavalink.Track) {
			// Embed build
			embed := shared.NewEmbed().
				SetTitle("🎵  Loading Default Track").
				SetDescription(
					fmt.Sprintf("`%s` (%s).\nTrack Duration: %s", tracks[0].Info.Title, *tracks[0].Info.URI, formatPosition(tracks[0].Info.Length))).
				SetThumbnail(*tracks[0].Info.ArtworkURL).
				SetColor(000000)

			// Interaction response
			_, _ = b.Session.InteractionResponseEdit(event.Interaction, &discordgo.WebhookEdit{
				Embeds: &[]*discordgo.MessageEmbed{embed.MessageEmbed},
			})

			// Queue handling
			if player.Track() == nil {
				toPlay = &tracks[0]
			} else {
				queue.Add(tracks[0])
			}
		},
		func() {
			_, _ = b.Session.InteractionResponseEdit(event.Interaction, &discordgo.WebhookEdit{
				Content: json.Ptr(fmt.Sprintf("Nothing found for: `%s`", identifier)),
			})
		},
		func(err error) {
			_, _ = b.Session.InteractionResponseEdit(event.Interaction, &discordgo.WebhookEdit{
				Content: json.Ptr(fmt.Sprintf("Error while looking up query: `%s`", err)),
			})
		},
	))
	if toPlay == nil {
		return nil
	}

	if err := b.Session.ChannelVoiceJoinManual(GuildId, DesignatedChannelId, false, false); err != nil {
		return err
	}

	b.Session.UpdateGameStatus(0, toPlay.Info.Title)

	return player.Update(context.TODO(), lavalink.WithTrack(*toPlay))
}

func (b *Bot) loop(event *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) error {

	loopModifier := data.Options[0].StringValue()

	player := b.Lavalink.Player(snowflake.MustParse(event.GuildID))
	queue := b.Queues.Get(event.GuildID)

	switch loopModifier {
	case "start":
		var currentlyPlaying *lavalink.Track = player.Track()
		if currentlyPlaying == nil {
			return fmt.Errorf("no song is currently playing in order to loop it")
		}
		const count = 512
		for range count {
			queue.Add(*currentlyPlaying)
		}
	case "stop":
		queue.Clear()
	}

	return nil
}
