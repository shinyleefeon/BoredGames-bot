package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"BoredGames-bot/internal/database"
)

var (
	SheetsService *sheets.Service
	SpreadsheetID string
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

func PullCurrentEvents(s *discordgo.Session, guildID string) []discordgo.GuildScheduledEvent {
	events, err := s.GuildScheduledEvents(guildID, false)
	if err != nil {
		fmt.Printf("Error fetching scheduled events: %v\n", err)
		return []discordgo.GuildScheduledEvent{}
	}
	
	// Convert []*GuildScheduledEvent to []GuildScheduledEvent
	result := make([]discordgo.GuildScheduledEvent, len(events))
	for i, event := range events {
		if event != nil {
			result[i] = *event
		}
	}
	return result
}

func InitGoogleSheets() {
	ctx := context.Background()
	
	credsFile := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if credsFile == "" {
		log.Fatal("GOOGLE_APPLICATION_CREDENTIALS environment variable is not set")
	}
	
	var err error
	SheetsService, err = sheets.NewService(ctx, option.WithCredentialsFile(credsFile))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	// Load the Spreadsheet ID from environment variables
	SpreadsheetID = os.Getenv("SPREADSHEET_ID")
	if SpreadsheetID == "" {
		log.Fatal("SPREADSHEET_ID environment variable is not set")
	}

	fmt.Println("✅ Connected to Google Sheets API")
}

func GetBoardGames() ([]BoardGame, error) {
	// Define the range to read. "Sheet1!A2:J" means "From A2 to the end of J"
	readRange := "Sheet1!A2:J"

	resp, err := SheetsService.Spreadsheets.Values.Get(SpreadsheetID, readRange).Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve data from sheet: %v", err)
	}

	var games []BoardGame

	for _, row := range resp.Values {
		// Ensure the row has enough columns to avoid "index out of range"
		if len(row) < 10 {
			continue // Skip incomplete rows
		}

		// Helper to safely convert interface{} to string and int
		name := fmt.Sprintf("%v", row[0])
		category := fmt.Sprintf("%v", row[1])
		minP := toInt(row[2])
		maxP := toInt(row[3])
		playTime := toInt(row[4])
		desc := fmt.Sprintf("%v", row[5])
		winner := fmt.Sprintf("%v", row[6])
		
		// Handle Booleans (TRUE/FALSE strings)
		played := row[7] == "TRUE"
		liked := row[8] == "TRUE"
		
		link := fmt.Sprintf("%v", row[9])

		game := BoardGame{
			Name:           name,
			Category:       category,
			MinPlayers:     minP,
			MaxPlayers:     maxP,
			Playtime:       playTime,
			Description:    desc,
			PreviousWinner: winner,
			PlayedYet:      played,
			LikedIt:        liked,
			RulesLink:      link,
		}
		games = append(games, game)
	}

	return games, nil
}

// Helper function to handle string-to-int conversion safely

func toInt(val interface{}) int {
	str := fmt.Sprintf("%v", val)
	i, _ := strconv.Atoi(str)
	return i
}


func SyncGames(db *sql.DB, sheetGames []BoardGame) error {
    // Create a queries object wrapping the DB connection
    q := database.New(db)
    ctx := context.Background()

    changes := 0
    for _, game := range sheetGames {
        // sqlc generated this strict struct for us!
        params := database.AddBoardGameParams{
            Name:           game.Name,
            Category:       sql.NullString{String: game.Category, Valid: game.Category != ""},
            MinPlayers:     sql.NullInt64{Int64: int64(game.MinPlayers), Valid: game.MinPlayers != 0},
            MaxPlayers:     sql.NullInt64{Int64: int64(game.MaxPlayers), Valid: game.MaxPlayers != 0},
            PlayTime:       sql.NullInt64{Int64: int64(game.Playtime), Valid: game.Playtime != 0},
            Description:    sql.NullString{String: game.Description, Valid: game.Description != ""},
            PreviousWinner: sql.NullString{String: game.PreviousWinner, Valid: game.PreviousWinner != ""},
            PlayedYet:      sql.NullBool{Bool: game.PlayedYet, Valid: true},
            LikedIt:        sql.NullBool{Bool: game.LikedIt, Valid: true},
            RulesLink:      sql.NullString{String: game.RulesLink, Valid: game.RulesLink != ""},
        }

        // Execute the upsert
        _, err := q.AddBoardGame(ctx, params)
        if err != nil {
            log.Printf("Failed to sync game '%s': %v", game.Name, err)
            continue
        }
        changes++
    }

    log.Printf("✅ Database Sync Complete: Processed %d games.", changes)
    return nil
}