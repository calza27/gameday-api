package filebuilder

import (
	"GameDay-API/internal/models"
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
)

func BuildCsv(gameData models.GameData) ([]byte, error) {
	var b bytes.Buffer
	byteWriter := bufio.NewWriter(&b)
	w := csv.NewWriter(byteWriter)

	//player table
	var data [][]string
	data = append(data, []string{"id", "Surname", "Given", "Number"})
	for _, player := range gameData.TeamAPlayers {
		data = append(data, []string{player.Id, player.Surname, player.GivenName, fmt.Sprintf("%d", player.Number)})
	}
	data = append(data, []string{string('\u200B')})
	w.WriteAll(data)
	//clear the data
	data = [][]string{}

	//gate data
	data = append(data, []string{"id", "gameDate", "competition", "home", "homeAbbr", "venue", "opposition", "oppAbbr", "level", "round"})
	data = append(data, []string{gameData.Id, gameData.GameDate, gameData.Competition, gameData.TeamA, gameData.TeamAAbbr, gameData.Venue, gameData.TeamB, gameData.TeamBAbbr, gameData.Level, gameData.Round, "[]"})
	data = append(data, []string{string('\u200B')})
	w.WriteAll(data)
	//clear the data
	data = [][]string{}

	//scoring table
	data = append(data, []string{"id", "quarter", "team", "ScoreEvent", "goalScorer", "scoreType", "HCWorm", "LauncherNo", "typeNumber", "OpWorm"})
	for _, scoringEvent := range gameData.ScoringEvents {
		data = append(data, []string{scoringEvent.Id, fmt.Sprintf("%d", scoringEvent.Quarter), scoringEvent.Team, scoringEvent.ScoreEvent, scoringEvent.GoalScorer, scoringEvent.ScoreType, fmt.Sprintf("%d", scoringEvent.HCWorm), fmt.Sprintf("%d", scoringEvent.LauncherNumber), fmt.Sprintf("%d", scoringEvent.TypeNumber), fmt.Sprintf("%d", scoringEvent.OpWorm)})
	}
	data = append(data, []string{string('\u200B')})
	w.WriteAll(data)
	//clear the data
	data = [][]string{}

	//quarter time table
	data = append(data, []string{"id", "quarter", "time"})
	for _, quarterTime := range gameData.QuarterTimes {
		data = append(data, []string{quarterTime.Id, fmt.Sprintf("%d", quarterTime.Quarter), quarterTime.Time})
	}
	data = append(data, []string{string('\u200B')})
	w.WriteAll(data)
	//clear the data
	data = [][]string{}

	//app storage table
	data = append(data, []string{"appStorage", "data"})
	for _, appStorageEvent := range gameData.AppStorage {
		d := []string{appStorageEvent.DataType}
		for _, i := range appStorageEvent.Data {
			d = append(d, fmt.Sprintf("%d", i))
		}
		data = append(data, d)
	}
	w.WriteAll(data)

	byteWriter.Flush()
	return b.Bytes(), nil
}
