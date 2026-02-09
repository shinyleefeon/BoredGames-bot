package main

import (
	"fmt"
	"os"
	"os/signal"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"strings"
	"database/sql"
	//"time"
	"BoredGames-bot/internal/database"
)

type BotApp struct {
	Session *discordgo.Session
	DB *sql.DB
}

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
	db, err := database.InitDB(dbPath)
	if err != nil {
		fmt.Printf("Database initialization failed: %v\n", err)
		return
	}
	defer db.Close()


	//make new bot
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		return
	}

	dg.AddHandler(messageCreate)
	dg.AddHandler(handleInteractions)
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

	//go app.StartReminderLoop()

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

//func (app *BotApp) StartReminderLoop() {}