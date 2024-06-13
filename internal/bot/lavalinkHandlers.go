package main

import (
	"context"
	"log"

	"github.com/AzteBot-Developments/AzteMusic/pkg/shared"
	"github.com/bwmarrin/discordgo"
	"github.com/disgoorg/snowflake/v2"
)

func (b *Bot) onVoiceStateUpdate(session *discordgo.Session, event *discordgo.VoiceStateUpdate) {
	if event.UserID != session.State.User.ID {
		return
	}

	var channelID *snowflake.ID
	if event.ChannelID != "" {
		id := snowflake.MustParse(event.ChannelID)
		channelID = &id
	}
	b.Lavalink.OnVoiceStateUpdate(context.TODO(), snowflake.MustParse(event.GuildID), channelID, event.SessionID)
	if event.ChannelID == "" {
		b.Queues.Delete(event.GuildID)
	}
}

func (b *Bot) onVoiceServerUpdate(session *discordgo.Session, event *discordgo.VoiceServerUpdate) {
	b.Lavalink.OnVoiceServerUpdate(context.TODO(), snowflake.MustParse(event.GuildID), event.Token, event.Endpoint)
}

func (b *Bot) onApplicationCommand(session *discordgo.Session, event *discordgo.InteractionCreate) {

	data := event.ApplicationCommandData()

	// If allowed roles are configured, only allow a user with one of these roles to execute an app command
	// The app commands which require role permissions are defined here
	if shared.StringInSlice(data.Name, RestrictedCommands) && len(AllowedRoles) > 0 {
		if event.Type == discordgo.InteractionApplicationCommand {
			// Check if the user has the allowed role
			hasAllowedRole := false
			for _, role := range event.Member.Roles {
				roleObj, err := session.State.Role(event.GuildID, role)
				if err != nil {
					log.Println("Error getting role:", err)
					return
				}
				if shared.StringInSlice(roleObj.Name, AllowedRoles) {
					hasAllowedRole = true
				}
				if hasAllowedRole {
					break
				}
			}

			if !hasAllowedRole {
				// If the user doesn't have the allowed role, send a response
				session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "You do not have the required role to use this command.",
					},
				})
				return
			}
		}
	}

	handler, ok := b.Handlers[data.Name]
	if !ok {
		log.Println("unknown command: ", data.Name)
		return
	}
	if err := handler(event, data); err != nil {
		log.Printf("error ocurred for %s: %v\n", data.Name, err)
		return
	}
}
