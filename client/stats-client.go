package client

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type StatsClient struct {
	BaseURL    string
	httpClient *http.Client
}

type Players struct {
	Data []Player
}

type Player struct {
	PlayerId       int    `json:"playerId"`
	SkaterFullName string `json:"skaterFullName"`
	TeamAbbrevs    string `json:"teamAbbrevs"`
	GamesPlayed    int    `json:"gamesPlayed"`
	Goals          int    `json:"goals"`
	Assists        int    `json:"assists"`
	Points         int    `json:"points"`
}

type Goalie struct {
	PlayerId            int     `json:"playerId"`
	GoalieFullName      string  `json:"goalieFullName"`
	TeamAbbrevs         string  `json:"teamAbbrevs"`
	GamesPlayed         int     `json:"gamesPlayed"`
	Goals               int     `json:"goals"`
	Assists             int     `json:"assists"`
	Points              int     `json:"points"`
	GoalsAgainstAverage float32 `json:"goalsAgainstAverage"`
	SavePct             float32 `json:"savePct"`
}

type Goalies struct {
	Data []Goalie
}

func (c *StatsClient) GetPlayerStats() (Players, error) {
	path := "/stats/rest/en/skater/summary?isAggregate=false&isGame=false&sort=%5B%7B%22property%22:%22points%22,%22direction%22:%22DESC%22%7D,%7B%22property%22:%22goals%22,%22direction%22:%22DESC%22%7D,%7B%22property%22:%22assists%22,%22direction%22:%22DESC%22%7D%5D&start=0&limit=50&factCayenneExp=gamesPlayed%3E=1&cayenneExp=active%3D1%20and%20gameTypeId=2%20and%20nationalityCode=%22FIN%22%20and%20seasonId%3C=20192020%20and%20seasonId%3E=20192020"
	var players Players
	body, err := c.get(path)
	if err != nil {
		return players, err
	}
	err = json.Unmarshal(body, &players)
	return players, err
}

func (c *StatsClient) GetGoalieStats() (Goalies, error) {
	path := "/stats/rest/en/goalie/summary?isAggregate=false&isGame=false&sort=%5B%7B%22property%22:%22wins%22,%22direction%22:%22DESC%22%7D,%7B%22property%22:%22savePct%22,%22direction%22:%22DESC%22%7D%5D&start=0&limit=50&factCayenneExp=gamesPlayed%3E=1&cayenneExp=gameTypeId=2%20and%20nationalityCode=%22FIN%22%20and%20seasonId%3C=20192020%20and%20seasonId%3E=20192020"
	var goalies Goalies
	body, err := c.get(path)
	if err != nil {
		return goalies, err
	}
	err = json.Unmarshal(body, &goalies)
	return goalies, err
}

func (c *StatsClient) get(path string) ([]byte, error) {
	u := c.BaseURL + path
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	res, getErr := c.httpClient.Do(req)
	if getErr != nil {
		return nil, getErr
	}
	return ioutil.ReadAll(res.Body)
}

func NewStatsClient() StatsClient {
	httpClient := http.DefaultClient
	baseUrl := "https://api.nhle.com"
	return StatsClient{httpClient: httpClient, BaseURL: baseUrl}
}
