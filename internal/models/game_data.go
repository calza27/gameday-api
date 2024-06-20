package models

type GameData struct {
	Id            string            `json:"Id"`
	GameDate      string            `json:"GameDate"`
	Competition   string            `json:"Competition"`
	TeamA         string            `json:"TeamA"`
	TeamAAbbr     string            `json:"TeamAAbbr"`
	Venue         string            `json:"Venue"`
	TeamB         string            `json:"TeamB"`
	TeamBAbbr     string            `json:"TeamBAbbr"`
	Level         string            `json:"Level"`
	Round         string            `json:"Round"`
	TeamAPlayers  []Player          `json:"TeamAPlayers"`
	ScoringEvents []ScoringEvent    `json:"ScoringEvents"`
	QuarterTimes  []QuarterTime     `json:"QuarterTimes"`
	AppStorage    []AppStorageEvent `json:"AppStorage"`
}

type Player struct {
	Id        string `json:"Id"`
	Surname   string `json:"Surname"`
	GivenName string `json:"GivenName"`
	Number    int    `json:"Number"`
}

type ScoringEvent struct {
	Id             string `json:"Id"`
	Quarter        int    `json:"Quarter"`
	Team           string `json:"Team"`
	ScoreEvent     string `json:"ScoreEvent"`
	GoalScorer     string `json:"GoalScorer"`
	ScoreType      string `json:"ScoreType"`
	HCWorm         int    `json:"HCWorm"`
	LauncherNumber int    `json:"LauncherNumber"`
	TypeNumber     int    `json:"TypeNumber"`
	OpWorm         int    `json:"OpWorm"`
}

type QuarterTime struct {
	Id      string `json:"Id"`
	Quarter int    `json:"Quarter"`
	Time    string `json:"Time"`
}

type AppStorageEvent struct {
	DataType string `json:"DataType"`
	Data     []int  `json:"Data"`
}
