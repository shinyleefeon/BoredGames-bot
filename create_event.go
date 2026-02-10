package main

import (
	"context"
	"time"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"BoredGames-bot/internal/database"
)


func Create_Event(s *discordgo.Session, i *discordgo.InteractionCreate, app *BotApp) {
    data := i.ApplicationCommandData()
    
    // Map options for easy access
    options := make(map[string]*discordgo.ApplicationCommandInteractionDataOption)
    for _, opt := range data.Options {
        options[opt.Name] = opt
    }

    // 1. Extract values
    name := options["name"].StringValue()
    dateStr := options["date"].StringValue()
    timeStr := options["time"].StringValue()
    
    desc := ""
    if opt, ok := options["description"]; ok {
        desc = opt.StringValue()
    }

    locName := "Voice Channel"
    if opt, ok := options["location"]; ok {
        locName = opt.StringValue()
    }

    // 2. Parse Time (Chattanooga/Eastern Time)
    layout := "2006-01-02 15:04"
    fullStr := fmt.Sprintf("%s %s", dateStr, timeStr)
    loc, _ := time.LoadLocation("America/New_York")
    scheduledTime, err := time.ParseInLocation(layout, fullStr, loc)

    if err != nil {
        s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
            Type: discordgo.InteractionResponseChannelMessageWithSource,
            Data: &discordgo.InteractionResponseData{Content: "❌ Invalid format. Use YYYY-MM-DD and HH:MM."},
        })
        return
    }

    // 3. Prepare API Parameters
    params := &discordgo.GuildScheduledEventParams{
        Name:               name,
        Description:        desc,
        ScheduledStartTime: &scheduledTime,
        PrivacyLevel:       discordgo.GuildScheduledEventPrivacyLevelGuildOnly,
    }

    // Handle Location Type (External vs Voice)
    if locName != "Voice Channel" {
        params.EntityType = discordgo.GuildScheduledEventEntityTypeExternal
        params.EntityMetadata = &discordgo.GuildScheduledEventEntityMetadata{Location: locName}
        endTime := scheduledTime.Add(4 * time.Hour) // Default 4-hour window
        params.ScheduledEndTime = &endTime
    } else {
        params.EntityType = discordgo.GuildScheduledEventEntityTypeVoice
        params.ChannelID = "1204199054578417665" // Default Voice Channel ID... Need to make this dynamic later
    }

    // 4. Call Discord API
    createdEvent, err := s.GuildScheduledEventCreate(i.GuildID, params)
    
    // 5. Final Response
    responseContent := fmt.Sprintf("✅ Event **%s** created for %s!", name, scheduledTime.Format("Mon, Jan _2 @ 3:04 PM"))
    if err != nil {
        responseContent = "❌ Error creating Discord event: " + err.Error()
    }
	if err == nil {
		ctx := context.Background()
		app.Queries.CreateEvent(ctx, database.CreateEventParams{
			DiscordEventID: createdEvent.ID,
			GuildID:        i.GuildID,
			Title:          name,
			StartTime:      scheduledTime,
		})
	}
		

    s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
        Type: discordgo.InteractionResponseChannelMessageWithSource,
        Data: &discordgo.InteractionResponseData{Content: responseContent},
    })
}