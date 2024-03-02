package main

import (
	"context"
	"fmt"
	"log"

	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/disgolink/v3/lavalink"
)

func (b *Bot) onPlayerPause(player disgolink.Player, event lavalink.PlayerPauseEvent) {
	fmt.Printf("onPlayerPause: %v\n", event)
	b.Session.UpdateGameStatus(0, StatusText)
}

func (b *Bot) onPlayerResume(player disgolink.Player, event lavalink.PlayerResumeEvent) {
	fmt.Printf("onPlayerResume: %v\n", event)
	b.Session.UpdateGameStatus(0, player.Track().Info.Title)
}

func (b *Bot) onTrackStart(player disgolink.Player, event lavalink.TrackStartEvent) {
	fmt.Printf("onTrackStart: %v\n", event)
	b.Session.UpdateGameStatus(0, event.Track.Info.Title)
}

func (b *Bot) onTrackEnd(player disgolink.Player, event lavalink.TrackEndEvent) {

	fmt.Printf("onTrackEnd: %v\n", event)

	b.Session.UpdateGameStatus(0, StatusText)

	if !event.Reason.MayStartNext() {
		return
	}

	queue := b.Queues.Get(event.GuildID().String())

	// in the case of the radio service, we can check here whether the queue is empty
	// if it is, play form url again
	if len(queue.Tracks) < 2 && DesignatedPlaylistUrl != "" && DesignatedChannelId != "" {
		b.AddToQueueFromSource(DesignatedPlaylistUrl, 3)
	}

	var (
		nextTrack lavalink.Track
		ok        bool
	)
	switch queue.Type {
	case QueueTypeNormal:
		nextTrack, ok = queue.Next()

	case QueueTypeRepeatTrack:
		nextTrack = event.Track

	case QueueTypeRepeatQueue:
		queue.Add(event.Track)
		nextTrack, ok = queue.Next()
	}

	if !ok {
		// retry to play designated playlist
		if DesignatedPlaylistUrl != "" && DesignatedChannelId != "" {
			b.AddToQueueFromSource(DesignatedPlaylistUrl, 3)
		} else {
			// No tracks on the queue, or could not play next, so can safely disconnect from the VC to save resources.
			if err := b.Session.ChannelVoiceJoinManual(GuildId, "", false, false); err != nil {
				fmt.Printf("[onTrackEnd] Error ocurred when disconnecting from VC: %v", err)
			}
			return
		}
	}

	if err := player.Update(context.TODO(), lavalink.WithTrack(nextTrack)); err != nil {
		log.Fatal("Failed to play next track: ", err)
	}
}

func (b *Bot) onTrackException(player disgolink.Player, event lavalink.TrackExceptionEvent) {
	fmt.Printf("onTrackException: %v\n", event)
	b.Session.UpdateGameStatus(0, StatusText)
}

func (b *Bot) onTrackStuck(player disgolink.Player, event lavalink.TrackStuckEvent) {
	fmt.Printf("onTrackStuck: %v\n", event)
	b.Session.UpdateGameStatus(0, StatusText)
}

func (b *Bot) onWebSocketClosed(player disgolink.Player, event lavalink.WebSocketClosedEvent) {
	fmt.Printf("onWebSocketClosed: %v\n", event)
	b.Session.UpdateGameStatus(0, StatusText)
}

func (b *Bot) onUnknownEvent(player disgolink.Player, event lavalink.UnknownEvent) {
	fmt.Printf("onWebSocketClosed: %v\n", event)
	b.Session.UpdateGameStatus(0, StatusText)
}
