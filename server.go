package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/nevalla/235-api/client"
	"github.com/patrickmn/go-cache"
)

type Response struct {
	Schedule client.Schedule `json:"schedule"`
	Players  []client.Player `json:"players"`
	Goalies  []client.Goalie `json:"goalies"`
}

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {

		if x, found := fromCache("response"); found {
			var cachedResponse *Response
			fmt.Println("Found from cache")
			cachedResponse = x.(*Response)
			return c.JSONPretty(http.StatusOK, &cachedResponse, "  ")
		} else {
			var response Response
			statsClient := client.NewStatsClient()
			scheduleClient := client.NewScheduleClient()

			players, playerError := statsClient.GetPlayerStats()
			if playerError == nil {
				response.Players = players.Data
			}
			goalies, goalieError := statsClient.GetGoalieStats()
			if goalieError == nil {
				response.Goalies = goalies.Data
			}

			schedule, scheduleError := scheduleClient.GetSchedule()

			if scheduleError == nil {
				response.Schedule = schedule
			}
			toCache("response", &response, cache.DefaultExpiration)
			return c.JSONPretty(http.StatusOK, response, "  ")
		}
	})
	e.Logger.Fatal(e.Start(":1323"))
}

func fromCache(key string) (interface{}, bool) {
	c := cache.New(1*time.Minute, 10*time.Minute)
	var cachedResponse *Response
	if x, found := c.Get("response"); found {
		fmt.Println("Found from cache")
		cachedResponse = x.(*Response)
		return cachedResponse, found
	}
	return cachedResponse, false
}

func toCache(k string, x interface{}, d time.Duration) {
	c := cache.New(1*time.Minute, 10*time.Minute)
	c.Set("response", &x, d)
}
