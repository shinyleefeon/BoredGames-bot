package main

import (
	"fmt"
	"os"
	"os/signal"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"strings"
	"database/sql"
	"time"
	"context"
	"log"
	"BoredGames-bot/internal/database"
)

type BotApp struct {
	Session *discordgo.Session
	DB *sql.DB
	Queries database.Querier
}

type BoardGame struct {
	Name           string
	Category       string
	MinPlayers     int
	MaxPlayers     int
	Playtime       int
	Description    string
	PreviousWinner string
	PlayedYet      bool
	LikedIt        bool
	RulesLink	  string
}

var gameLibrary []BoardGame

var commands = []*discordgo.ApplicationCommand{
	{
        Name:        "schedule",
        Description: "Schedule an event",
        Options: []*discordgo.ApplicationCommandOption{
            {
                Type:        discordgo.ApplicationCommandOptionString,
                Name:        "name",
                Description: "Name of the board game (e.g. Terraforming Mars)",
                Required:    true,
            },
            {
                Type:        discordgo.ApplicationCommandOptionString,
                Name:        "date",
                Description: "Date (YYYY-MM-DD)",
                Required:    true,
            },
            {
                Type:        discordgo.ApplicationCommandOptionString,
                Name:        "time",
                Description: "Time (HH:MM in 24h format)",
                Required:    true,
            },
            {
                Type:        discordgo.ApplicationCommandOptionString,
                Name:        "description",
                Description: "Optional: Details about the session",
                Required:    false,
            },
            {
                Type:        discordgo.ApplicationCommandOptionString,
                Name:        "location",
                Description: "Optional: Where we are playing (defaults to Voice)",
                Required:    false,
            },
        },
	},
	{
		Name: "streak",
		Description: "Check your current attendance streak",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name: "user",
				Type: discordgo.ApplicationCommandOptionUser,
				Description: "The user to check the streak for",
				Required: true,
			},
		},
	},
}



func main() {
	// Load environment
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
	}

	token := os.Getenv("DiscordBotToken")
	if token == "" {
		fmt.Println("No Discord bot token provided in environment variable 'DiscordBotToken'")
		return
	}
	defGuild := os.Getenv("DefaultGuildID")
	if defGuild == "" {
		fmt.Println("No default guild ID provided in environment variable 'DefaultGuildID'")
		return
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./boredgames.db" // Default path if not provided
	}
	fmt.Printf("Using database path: %s\n", dbPath)
	
	// run db migrations
	dbConn, err := database.InitDB(dbPath)
	if err != nil {
		fmt.Printf("Database initialization failed: %v\n", err)
		return
	}
	defer dbConn.Close()

	InitGoogleSheets()

	//make new bot
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		return
	}

	// Create app before registering handlers
	app := &BotApp{
        Session: dg,
        DB:      dbConn,
        Queries: database.New(dbConn),
    }

	dg.AddHandler(messageCreate)
	dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		handleInteractions(s, i, app)
	})
	//Declare intents for discord
	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentMessageContent | discordgo.IntentsGuildScheduledEvents

	//open a websocket connection to Discord
	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening connection: ", err)
		return
	}
	defer dg.Close()

	// This wipes all global commands if you ever need a fresh start
	oldCommands, _ := dg.ApplicationCommands(dg.State.User.ID, defGuild)
	for _, v := range oldCommands {
		dg.ApplicationCommandDelete(dg.State.User.ID, defGuild, v.ID)
	}

	//register commands
	fmt.Println("Registering commands...")
	for _, v := range commands {
		_, err := dg.ApplicationCommandCreate(dg.State.User.ID, defGuild, v)
		if err != nil {
			fmt.Printf("Cannot create '%v' command: %v\n", v.Name, err)
		}
	}

	currentEvents := PullCurrentEvents(dg, defGuild)
	if len(currentEvents) > 0 {
		for _, event := range currentEvents {
			_, err := app.Queries.GetEventByEventID(context.Background(), event.ID)
			if err != nil {
				// Event doesn't exist in database, create it
				app.Queries.CreateEvent(context.Background(), database.CreateEventParams{
					DiscordEventID: event.ID,
					GuildID:        event.GuildID,
					Title:          event.Name,
					StartTime:      event.ScheduledStartTime,
				})
			}
		}
	}
		
	
	BoardGames, err := GetBoardGames()
	if err != nil {
		fmt.Printf("Error getting board games: %v\n", err)
		return
	}

	err = SyncGames(dbConn, BoardGames)
	if err != nil {
		fmt.Printf("Error syncing board games: %v\n", err)
		return
	}

	go app.StartReminderLoop()

	fmt.Println("Bot is now running. Press CTRL-C to exit.")

	//wait for a signal to quit
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)
	<-sc
}


func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	//ignore messages from the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}
	if strings.HasPrefix(strings.ToLower(m.Content), "!clanker") {
		args := strings.Fields(m.Content)
		if len(args) < 2 {
			s.ChannelMessageSend(m.ChannelID, "Hey that's our word!")
			return
		}
		switch strings.ToLower(args[1]) {





		case "seize":
			if strings.ToLower(m.Content) == "!clanker seize him and take him to the penis explosion chamber" {
				file, err := os.Open("Resources/MODS.gif")
				if err != nil {
					fmt.Println("Error opening MODS.gif:", err)
					return
				}
				defer file.Close()
				_, _ = s.ChannelFileSend(m.ChannelID, "MODS.gif", file)
			}
		
		/*case "create_event":
			err := Create_Event(args[2:])
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error %s \n Create_Event Usage: ", err))

		}*/
	}

	
	
	
	
	
}
}

func (app *BotApp) StartReminderLoop() {
	ctx := context.Background()
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		events, err := app.Queries.GetEventsForReminder(ctx)
		if err != nil {
			log.Printf("Error getting events for reminder: %v", err)
		}

		for _, event := range events {
			log.Printf("Reminder for event: %s starting at %v", event.Title, event.StartTime)
			
			interestedUsers, err := app.Session.GuildScheduledEventUsers(event.GuildID, event.DiscordEventID, 100, false, "", "")
			if err != nil {
				log.Printf("Error fetching interested users for event %d: %v", event.ID, err)
				continue
			}

			for _, user := range interestedUsers {
				dbUser, err := app.Queries.GetUserByDiscordID(ctx, user.User.ID)
				if err != nil {
					// User doesn't exist, create them
					dbUser, err = app.Queries.CreateUser(ctx, database.CreateUserParams{
						DiscordID: user.User.ID,
						Username:  user.User.Username,
					})
					if err != nil {
						log.Printf("Error creating user for Discord ID %s: %v", user.User.ID, err)
						continue
					}
				} else {
					app.Queries.IncrementStreak(ctx, dbUser.ID)
				}
				
				
				err = app.Queries.AddParticipant(ctx, database.AddParticipantParams{
					EventID: event.ID,
					UserID:  dbUser.ID,
				})
				if err != nil {
					log.Printf("Error adding participant for event %d and user %s: %v", event.ID, user.User.ID, err)
					continue
				}
			}

			participantDiscordIDs, err := app.Queries.GetParticipants(ctx, event.ID)
			if err != nil {
				log.Printf("Error getting participants for event %d: %v", event.ID, err)
				continue
			}
			fmt.Printf("Sending reminders to participants: %v\n", participantDiscordIDs)
			for _, discordID := range participantDiscordIDs {
				ch, err := app.Session.UserChannelCreate(discordID)
				if err != nil {
					log.Printf("Error creating DM channel for user %s: %v", discordID, err)
					continue
				}
				msg := fmt.Sprintf("🔔 **Reminder:** '%s' starts in 30 minutes!", event.Title)
                _, err = app.Session.ChannelMessageSend(ch.ID, msg)
                if err != nil {
                     log.Printf("Failed to send DM: %v", err)
                }
            }
			
			// Reset streaks for users who aren't participating
			allUsers, err := app.Queries.GetAllUsers(ctx)
			if err != nil {
				log.Printf("Error getting all users: %v", err)
			} else {
				// Create a map of participating discord IDs for quick lookup
				participantMap := make(map[string]bool)
				for _, discordID := range participantDiscordIDs {
					participantMap[discordID] = true
				}
				
				// Reset streak for users not participating
				for _, user := range allUsers {
					if !participantMap[user.DiscordID] {
						err := app.Queries.ResetStreak(ctx, user.ID)
						if err != nil {
							log.Printf("Error resetting streak for user %s: %v", user.DiscordID, err)
						} else {
							log.Printf("Reset streak for non-participating user: %s", user.Username)
						}
					}
				}
			}
			
			if err := app.Queries.MarkReminderSent(ctx, event.ID); err != nil {
				log.Printf("Error marking reminder sent for event %d: %v", event.ID, err)
			}
		
		
		
		
		}
	}
}