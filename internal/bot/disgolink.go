package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"strconv"

	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/snowflake/v2"
	"github.com/joho/godotenv"
)

var (
	_ = godotenv.Load("./cmd/music-service/.env")

	NodeName      = os.Getenv("LAVALINK_NODE_NAME")
	NodeAddress   = os.Getenv("LAVALINK_NODE_ADDRESS")
	NodePassword  = os.Getenv("LAVALINK_NODE_PASSWORD")
	NodeSecure, _ = strconv.ParseBool(os.Getenv("LAVALINK_NODE_SECURE"))
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
		slog.Error("failed to get node version", slog.Any("err", err))
		os.Exit(1)
	}
	log.Printf("lavalink node version: %s", version)
}
