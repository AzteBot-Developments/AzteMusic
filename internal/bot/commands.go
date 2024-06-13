package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var Commands = []*discordgo.ApplicationCommand{
	{
		Name:        "play",
		Description: "Plays a song",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "identifier",
				Description: "The song link or search query",
				Required:    true,
			},
		},
	},
	{
		Name:        "play_default",
		Description: "Plays the default Azteca Essentials playlist",
	},
	{
		Name:        "loop",
		Description: "Plays the currently playing song in a loop, until stopped",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "loop-modifier",
				Description: "Select the action to use for the loop function.",
				Required:    true,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{
						Name:  "Loop Current Song",
						Value: "start",
					},
					{
						Name:  "Stop Loop",
						Value: "stop",
					},
				},
			},
		},
	},
	{
		Name:        "pause",
		Description: "Pauses the current song",
	},
	{
		Name:        "skip",
		Description: "Skips the current song",
	},
	{
		Name:        "help",
		Description: fmt.Sprintf("Returns a slash commands guide for the %s", BotName),
	},
	{
		Name:        "now-playing",
		Description: "Shows the current playing song",
	},
	{
		Name:        "stop",
		Description: "Stops the current song and stops the player",
	},
	{
		Name:        "shuffle",
		Description: "Shuffles the current queue",
	},
	{
		Name:        "queue",
		Description: "Shows the current queue",
	},
	{
		Name:        "clear-queue",
		Description: "Clears the current queue",
	},
	{
		Name:        "queue-type",
		Description: "Sets the queue type",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "type",
				Description: "The queue type",
				Required:    true,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{
						Name:  "default",
						Value: "default",
					},
					{
						Name:  "repeat-track",
						Value: "repeat-track",
					},
					{
						Name:  "repeat-queue",
						Value: "repeat-queue",
					},
				},
			},
		},
	},
}

func (b *Bot) RegisterCommands() {
	b.Handlers = map[string]func(event *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) error{
		"play":         b.play,
		"play_default": b.play_default,
		"pause":        b.pause,
		"skip":         b.skip,
		"now-playing":  b.nowPlaying,
		"stop":         b.stop,
		"queue":        b.queue,
		"clear-queue":  b.clearQueue,
		"queue-type":   b.queueType,
		"shuffle":      b.shuffle,
		"loop":         b.loop,
		"help":         b.help,
	}

	if _, err := b.Session.ApplicationCommandBulkOverwrite(b.Session.State.User.ID, GuildId, Commands); err != nil {
		panic(err)
	}
}
