package main

import (
	"context"
	"log"

	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/snowflake/v2"
)

func (b *Bot) AddLavalinkNode(ctx context.Context) {
	b.Lavalink = disgolink.New(snowflake.MustParse(b.Session.State.User.ID),
		disgolink.WithListenerFunc(b.onPlayerPause),
		disgolink.WithListenerFunc(b.onPlayerResume),
		disgolink.WithListenerFunc(b.onTrackStart),
		disgolink.WithListenerFunc(b.onTrackEnd),
		disgolink.WithListenerFunc(b.onTrackException),
		disgolink.WithListenerFunc(b.onTrackStuck),
		disgolink.WithListenerFunc(b.onWebSocketClosed),
	)

	node, err := b.Lavalink.AddNode(ctx, disgolink.NodeConfig{
		Name:     NodeName,
		Address:  NodeAddress,
		Password: NodePassword,
		Secure:   NodeSecure,
	})
	if err != nil {
		panic(err)
	}
	version, err := node.Version(ctx)
	if err != nil {
		panic(err)
	}
	log.Printf("lavalink node version: %s", version)
}
