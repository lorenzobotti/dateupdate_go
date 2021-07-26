package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"
)

const tmzArea = "Europe"
const tmzLocation = "Rome"

func main() {
	help, format, sudo, area, location := parseFlags()
	if help {
		fmt.Fprintln(os.Stderr, helpScreen)
		os.Exit(0)
	}

	result, err := getDateFromApi(area, location)
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

	if res.StatusCode != 200 {
		return ApiResult{}, fmt.Errorf("getDateFromApi: server responded with %d", res.StatusCode)
	}

	result := ApiResult{}
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return ApiResult{}, err
	}

	return result, nil
}

func parseFlags() (help, formatCommand, sudo bool, area, location string) {
	flag.BoolVar(&help, "help", false, "")
	flag.BoolVar(&help, "h", false, "")
	flag.BoolVar(&formatCommand, "command", false, "")
	flag.BoolVar(&formatCommand, "c", false, "")
	flag.BoolVar(&sudo, "sudo", false, "")
	flag.BoolVar(&sudo, "s", false, "")

	flag.StringVar(&area, "area", tmzArea, "")
	flag.StringVar(&area, "a", tmzArea, "")
	flag.StringVar(&location, "location", tmzLocation, "")
	flag.StringVar(&location, "l", tmzLocation, "")

	flag.Usage = func() {}
	flag.Parse()

	if sudo {
		formatCommand = true
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
	return fmt.Sprintf("%02d%02d%02d %02d:%02d",
		r.DateTime.Year(),
		r.DateTime.Month(),
		r.DateTime.Day(),
		r.DateTime.Hour(),
		r.DateTime.Minute(),
	)
}

const helpScreen = `dateupdate by Lorenzo Botti
This program calls the WorldTimeAPI (worldtimeapi.org) and formats its
responde the way the date command accepts it. It's meant to help when
fucking retard Manjaro can't figure out what time it is

Flags:
--help, -h:
    Show this help screen
    
--command, -c
    Format the output as a date command. Example: date --set "20021227 21:23"
    
--sudo, -s
    Format the output as a date command with sudo. Example: sudo date --set "20021227 21:23"

--location, -l
    Set the location of the timezone in the call to the API. Example: Europe

--area, -a
    Set the area of the timezone in the call to the API. Example: Rome`
