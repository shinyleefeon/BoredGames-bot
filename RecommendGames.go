package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"BoredGames-bot/internal/database"
)

type GameRecommendation struct {
	Name        string
	Description string
	RulesLink   string
}

func Recommend_game(s *discordgo.Session, i *discordgo.InteractionCreate, app *BotApp, ){
    data := i.ApplicationCommandData()

	options := make(map[string]*discordgo.ApplicationCommandInteractionDataOption)
	for _, opt := range data.Options {
		options[opt.Name] = opt
	}

	players := options["players"].IntValue()
	
	playedYet := ""
	notSpecified := true
	if opt, ok := options["played_yet"]; ok {
		if opt.BoolValue() {
			playedYet = "played"
		} else {
			playedYet = "not played"
		}
		notSpecified = false
	}

	ctx := context.Background()
	var recommendedGames []GameRecommendation
	var err error

	params := database.RecommendRandomBoardGameParams{
		MinPlayers: sql.NullInt64{Int64: players, Valid: true},
		MaxPlayers: sql.NullInt64{Int64: players, Valid: true},
	}

	if notSpecified {
		rows, e := app.Queries.RecommendRandomBoardGame(ctx, params)
		err = e
		for _, row := range rows {
			desc := ""
			if row.Description.Valid {
				desc = row.Description.String
			}
			rulesLink := ""
			if row.RulesLink.Valid {
				rulesLink = row.RulesLink.String
			}
			recommendedGames = append(recommendedGames, GameRecommendation{
				Name:        row.Name,
				Description: desc,
				RulesLink:   rulesLink,
			})
		}
	} else if playedYet == "played" {
		playedParams := database.RecomendPlayedBoardGameParams{
			MinPlayers: sql.NullInt64{Int64: players, Valid: true},
			MaxPlayers: sql.NullInt64{Int64: players, Valid: true},
		}
		rows, e := app.Queries.RecomendPlayedBoardGame(ctx, playedParams)
		err = e
		for _, row := range rows {
			desc := ""
			if row.Description.Valid {
				desc = row.Description.String
			}
			rulesLink := ""
			if row.RulesLink.Valid {
				rulesLink = row.RulesLink.String
			}
			recommendedGames = append(recommendedGames, GameRecommendation{
				Name:        row.Name,
				Description: desc,
				RulesLink:   rulesLink,
			})
		}
	} else if playedYet == "not played" {
		unplayedParams := database.RecomendUnplayedBoardGameParams{
			MinPlayers: sql.NullInt64{Int64: players, Valid: true},
			MaxPlayers: sql.NullInt64{Int64: players, Valid: true},
		}
		rows, e := app.Queries.RecomendUnplayedBoardGame(ctx, unplayedParams)
		err = e
		for _, row := range rows {
			desc := ""
			if row.Description.Valid {
				desc = row.Description.String
			}
			rulesLink := ""
			if row.RulesLink.Valid {
				rulesLink = row.RulesLink.String
			}
			recommendedGames = append(recommendedGames, GameRecommendation{
				Name:        row.Name,
				Description: desc,
				RulesLink:   rulesLink,
			})
		}
	}

	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "❌ Error fetching game recommendations."},
		})
		return
	}

	if len(recommendedGames) == 0 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "❌ Uh Oh. No games found matching the criteria."},
		})
		return
	}

	responseContent := "Recommended Games:\n"
	for _, game := range recommendedGames {
		responseContent += fmt.Sprintf("- **%s**: %s\n Link: %s\n", game.Name, game.Description, game.RulesLink)
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: responseContent},
	})

	
}
