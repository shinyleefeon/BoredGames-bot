package main

import (
	
	

	"github.com/bwmarrin/discordgo"
)

func handleInteractions(s *discordgo.Session, i *discordgo.InteractionCreate) {
    if i.Type != discordgo.InteractionApplicationCommand {
        return
    }

    switch i.ApplicationCommandData().Name {
    case "schedule":
        Create_Event(s, i) // Calling the Create_Event function
    // Future cases go here:
    // case "record-win":
    //     Record_Win(s, i)
    }
}