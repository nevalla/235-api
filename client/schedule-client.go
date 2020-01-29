package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Schedule struct {
	Dates []Date `json:"dates"`
}

type Date struct {
	Date  string `json:"date"`
	Games []Game `json:"games"`
}

type Game struct {
	GameDate time.Time `json:"gameDate"`
	Teams    Teams     `json:"teams"`
}

type Status struct {
	Code          int    `json:"statusCode"`
	DetailedState string `json:"detailedState"`
}

type Teams struct {
	Away GameTeam `json:"away"`
	Home GameTeam `json:"home"`
}

type GameTeam struct {
	Team  Team `json:"team"`
	Score int  `json:"score"`
}

type Team struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	Abbreviation string `json:"abbreviation"`
	ShortName    string `json:"shortName"`
}

type ScheduleClient struct {
	BaseURL    string
	httpClient *http.Client
}

func NewScheduleClient() ScheduleClient {
	httpClient := http.DefaultClient
	baseUrl := "https://statsapi.web.nhl.com/api/v1/schedule"
	return ScheduleClient{httpClient: httpClient, BaseURL: baseUrl}
}

func (c *ScheduleClient) GetSchedule() (Schedule, error) {
	fmt.Println("Fetching schedule")
	path := "?date=2020-01-28&expand=schedule.teams,schedule.scoringplays"
	var schedule Schedule
	body, err := c.get(path)
	if err != nil {
		return schedule, err
	}
	err = json.Unmarshal(body, &schedule)
	return schedule, err
}

func (c *ScheduleClient) get(path string) ([]byte, error) {
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
