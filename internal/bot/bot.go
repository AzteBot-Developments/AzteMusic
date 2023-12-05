package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/disgoorg/disgolink/v3/disgolink"
)

type Bot struct {
	Session  *discordgo.Session
	Lavalink disgolink.Client
	Handlers map[string]func(event *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) error
	Queues   *QueueManager
}

func GetAuthenticatedBotSession() *discordgo.Session {
	session, err := discordgo.New("Bot " + Token)
	if err != nil {
		panic(err)
	}
	return session
}

func (b *Bot) SetIntents() {
	b.Session.State.TrackVoice = true
	b.Session.Identify.Intents = discordgo.IntentGuilds | discordgo.IntentsGuildVoiceStates
}

func (b *Bot) AddHandlers() {
	b.Session.AddHandler(b.onReady)
	b.Session.AddHandler(b.onApplicationCommand)
	b.Session.AddHandler(b.onVoiceStateUpdate)
	b.Session.AddHandler(b.onVoiceServerUpdate)
}

func (b *Bot) Connect() {
	if err := b.Session.Open(); err != nil {
		panic(err)
	}
	defer b.Session.Close()
}
