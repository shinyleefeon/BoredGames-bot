package main

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"fmt"
)

func Check_Streak(s *discordgo.Session, i *discordgo.InteractionCreate, app *BotApp) {
	ctx := context.Background()
	data := i.ApplicationCommandData()
	
	// Get the user from the command option
	var discordID string
	if len(data.Options) > 0 && data.Options[0].Name == "user" {
		discordID = data.Options[0].UserValue(s).ID
	} else {
		discordID = i.Member.User.ID
	}

	// First get the user from database
	dbUser, err := app.Queries.GetUserByDiscordID(ctx, discordID)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "❌ User not found in database."},
		})
		return
	}

	// Get the streak
	streakResult, err := app.Queries.GetUserStreak(ctx, dbUser.ID)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "❌ Error fetching attendance streak: " + err.Error()},
		})
		return
	}

	streak := int64(0)
	if streakResult.Valid {
		streak = streakResult.Int64
	}
	
	responseContent := fmt.Sprintf("Current attendance streak: **%d**", streak)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: responseContent},
	})
}