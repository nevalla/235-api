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

var ca = cache.New(1*time.Minute, 10*time.Minute)

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {

		if x, found := fromCache("response"); found {
			var cachedResponse *Response
			cachedResponse = x
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

func fromCache(key string) (*Response, bool) {
	var r *Response
	if x, found := ca.Get("response"); found {
		fmt.Println("Found response from cache")
		r = x.(*Response)
		return r, found
	}
	return r, false
}

func toCache(k string, x *Response, d time.Duration) {
	ca.Set("response", x, d)
}
