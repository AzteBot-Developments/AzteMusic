package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/snowflake/v2"
)

type Bot struct {
	Session           *discordgo.Session
	Lavalink          disgolink.Client
	Handlers          map[string]func(event *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) error
	Queues            *QueueManager
	HasLavaLinkClient bool
}

func NewBot() *Bot {
	newBot := &Bot{
		Queues: &QueueManager{
			queues: make(map[string]*Queue),
		},
	}
	newBot.HasLavaLinkClient = false
	return newBot
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

func (b *Bot) SetupLavalink() {
	b.Lavalink = disgolink.New(snowflake.MustParse(b.Session.State.User.ID),
		disgolink.WithListenerFunc(b.onPlayerPause),
		disgolink.WithListenerFunc(b.onPlayerResume),
		disgolink.WithListenerFunc(b.onTrackStart),
		disgolink.WithListenerFunc(b.onTrackEnd),
		disgolink.WithListenerFunc(b.onTrackException),
		disgolink.WithListenerFunc(b.onTrackStuck),
		disgolink.WithListenerFunc(b.onWebSocketClosed),
	)
	b.HasLavaLinkClient = true
}

func (b *Bot) AddVoiceHandlers() {
	b.Session.AddHandler(b.onApplicationCommand)
	b.Session.AddHandler(b.onVoiceStateUpdate)
	b.Session.AddHandler(b.onVoiceServerUpdate)
}
