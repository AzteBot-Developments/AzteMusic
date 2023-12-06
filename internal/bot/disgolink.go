package main

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/joho/godotenv"
)

var (
	_ = godotenv.Load(".env")

	NodeName      = os.Getenv("LAVALINK_NODE_NAME")
	NodeAddress   = os.Getenv("LAVALINK_NODE_ADDRESS")
	NodePassword  = os.Getenv("LAVALINK_NODE_PASSWORD")
	NodeSecure, _ = strconv.ParseBool(os.Getenv("LAVALINK_NODE_SECURE"))
)

func (b *Bot) AddLavalinkNode(ctx context.Context) {
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
