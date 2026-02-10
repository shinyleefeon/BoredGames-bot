package main

import (
	
	

	"github.com/bwmarrin/discordgo"
)

func handleInteractions(s *discordgo.Session, i *discordgo.InteractionCreate, app *BotApp) {
    if i.Type != discordgo.InteractionApplicationCommand {
        return
    }

    switch i.ApplicationCommandData().Name {
    case "schedule":
        Create_Event(s, i, app)
	case "streak":
		Check_Streak(s, i, app)
    // Future cases go here:
    // case "record-win":
    //     Record_Win(s, i)
	}
}