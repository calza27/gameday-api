package main

import (
	"GameDay-API/internal/filebuilder"
	"GameDay-API/internal/models"
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

const jsonObject = "{\"Id\":\"9f255d6b-e410-48c7-9aff-53c543623e76\",\"GameDate\":\"2024-06-20\",\"Competition\":\"APS\",\"TeamA\":\"Haileybury\",\"TeamAAbbr\":\"HY\",\"Venue\":\"Haileybury\",\"TeamB\":\"ScotchCollege\",\"TeamBAbbr\":\"SC\",\"Level\":\"1\",\"Round\":\"1\",\"TeamAPlayers\":[{\"Id\":\"9f255d6b-e410-48c7-9aff-53c543623e76\",\"Surname\":\"Doe\",\"GivenName\":\"John\",\"Number\":1},{\"Id\":\"9f255d6b-e410-48c7-9aff-53c543623e76\",\"Surname\":\"Dans\",\"GivenName\":\"Tim\",\"Number\":2},{\"Id\":\"9f255d6b-e410-48c7-9aff-53c543623e76\",\"Surname\":\"Crosby\",\"GivenName\":\"Sidney\",\"Number\":87}],\"ScoringEvents\":[{\"Id\":\"9f255d6b-e410-48c7-9aff-53c543623e76\",\"Quarter\":1,\"Team\":\"Haileybury\",\"ScoreEvent\":\"1\",\"GoalScorer\":\"2\",\"ScoreType\":\"4\",\"HCWorm\":8,\"LauncherNumber\":16,\"TypeNumber\":32,\"OpWorm\":64},{\"Id\":\"9f255d6b-e410-48c7-9aff-53c543623e76\",\"Quarter\":1,\"Team\":\"Scotch\",\"ScoreEvent\":\"3\",\"GoalScorer\":\"6\",\"ScoreType\":\"12\",\"HCWorm\":24,\"LauncherNumber\":48,\"TypeNumber\":96,\"OpWorm\":192}],\"QuarterTimes\":[{\"Id\":\"9f255d6b-e410-48c7-9aff-53c543623e76\",\"Quarter\":1,\"Time\":\"30:05\"},{\"Id\":\"9f255d6b-e410-48c7-9aff-53c543623e76\",\"Quarter\":2,\"Time\":\"29:55\"},{\"Id\":\"9f255d6b-e410-48c7-9aff-53c543623e76\",\"Quarter\":3,\"Time\":\"28:59\"},{\"Id\":\"9f255d6b-e410-48c7-9aff-53c543623e76\",\"Quarter\":4,\"Time\":\"31:01\"}],\"AppStorage\":[{\"DataType\":\"Kick\",\"Data\":[1,2,3,4,5]},{\"DataType\":\"Scrap\",\"Data\":[5,3,1]}]}"

func Test_export(t *testing.T) {
	var gameData models.GameData
	err := json.Unmarshal([]byte(jsonObject), &gameData)
	if err != nil {
		panic(err)
	}

	bytes, err := filebuilder.BuildPdf(gameData)
	if err != nil {
		panic(err)
	}

	if len(bytes) == 0 {
		t.Errorf("Expected bytes to be greater than 0")
	}

	file, err := os.Create("test.pdf")
	if err != nil {
		panic(err)
	}
	file.Write(bytes)
	file.Close()
	fmt.Printf("PDF file created: %s\n", file.Name())
}
