package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"
)

const tmzArea = "Europe"
const tmzLocation = "Rome"

func main() {
	help, format, sudo := parseFlags()
	if help {
		fmt.Fprintln(os.Stderr, helpScreen)
		os.Exit(0)
	}

	result, err := getDateFromApi(tmzArea, tmzLocation)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	output := result.UnixFormat()

	if format {
		output = `date --set "` + output + `"`
	}

	if sudo {
		output = `sudo ` + output
	}

	fmt.Println(output)
}

func getDateFromApi(area, location string) (ApiResult, error) {
	const dateApiUrl = `http://worldtimeapi.org/api/timezone/%s/%s`

	url := fmt.Sprintf(dateApiUrl, area, location)
	res, err := http.Get(url)
	if err != nil {
		return ApiResult{}, err
	}

	result := ApiResult{}
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return ApiResult{}, err
	}

	return result, nil
}

func parseFlags() (help, formatCommand, sudo bool) {
	if len(os.Args) <= 1 {
		return
	}

	for _, flag := range os.Args[1:] {
		if flag == "--help" || flag == "-h" {
			help = true
		}

		if flag == "--command" || flag == "-c" {
			formatCommand = true
		}

		if flag == "--sudo" || flag == "-s" {
			sudo = true
		}
	}

	return
}

// abbreviation     CEST
// day_of_year      207
// dst_from         2021-03-28T01:00:00+00:00
// raw_offset       3600
// dst_offset       3600
// utc_datetime     2021-07-26T12:59:59.739761+00:00
// week_number      30
// client_ip        37.162.183.219
// dst      true
// dst_until        2021-10-31T01:00:00+00:00
// utc_offset       +02:00
// datetime         2021-07-26T14:59:59.739761+02:00
// day_of_week      1
// timezone         Europe/Rome
// unixtime         1.627304399e+09

type ApiResult struct {
	Abbreviation string `json:""`
	Timezone     string `json:"timezone"`

	ClientIp net.IP `json:"client_ip"`

	Dst        bool      `json:"dst"`
	DstStarted time.Time `json:"dst_from"`
	DstUntil   time.Time `json:"dst_until"`
	// DstOffset time.Duration `json:""`
	// UtcOffset time.Duration

	DayOfYear  int `json:"day_of_year"`
	DayOfWeek  int `json:"day_of_week"`
	WeekNumber int `json:"week_number"`

	DateTime time.Time `json:"datetime"`
	// UnixTime int
}

func (r ApiResult) UnixFormat() string {
	return fmt.Sprintf("%02d%02d%02d %d:%d",
		r.DateTime.Year(),
		r.DateTime.Month(),
		r.DateTime.Day(),
		r.DateTime.Hour(),
		r.DateTime.Second(),
	)
}

const helpScreen = `dateupdate by Lorenzo Botti
	This program calls the WorldTimeAPI (worldtimeapi.org) and formats its
	responde the way the date command accepts it. It's meant to help when
	fucking retard Manjaro can't figure out what time it is`
