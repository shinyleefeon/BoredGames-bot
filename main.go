package main

import (
	"fmt"
	"os"
	"os/signal"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)


func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
	}

	token := os.Getenv("DiscordBotToken")
	if token == "" {
		fmt.Println("No Discord bot token provided in environment variable 'DiscordBotToken'")
		return
	}


	//make new bot
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		return
	}

	dg.AddHandler(messageCreate)
	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentMessageContent
	//open a websocket connection to Discord
	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening connection: ", err)
		return
	}
	defer dg.Close()

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
	if m.Content == "!clanker" {
		s.ChannelMessageSend(m.ChannelID, "Beep boop I'm a bot")
	} else if m.Content == "!Clanker seize him and take him to the penis explosion chamber" {
		file, err := os.Open("Resources/MODS.gif")
		if err != nil {
			fmt.Println("Error opening MODS.gif:", err)
			return
		}
		defer file.Close()
		_, _ = s.ChannelFileSend(m.ChannelID, "MODS.gif", file)
	}
}