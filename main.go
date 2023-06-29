// Copyright 2023 Shannon Wynter
//
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.

package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime/debug"
	"time"

	"github.com/freman/ghostport"
)

const (
	logsAPI       = "https://www.snapengage.com/api/v2/%s/logs?widgetId=%s&start=%s&end=%s"
	dateFormat    = "2006-01-02"
	dirDateFormat = "2006" + string(os.PathSeparator) + "01" + string(os.PathSeparator) + "02"
)

func main() {
	flgToken := strEnvFlag("token", "SNAPENGAGE_TOKEN", "Snapengage API Token")
	flgOrg := strEnvFlag("org", "SNAPENGAGE_ORG", "Snapengage Organisation ID")
	flgWidget := strEnvFlag("widget", "SNAPENGAGE_WIDGET", "Snapengage Widget ID")
	flgStart := strEnvFlag("start", "SNAPENGAGE_START", "Start exporting from what date YYYY-MM-DD")
	flgEnd := strEnvFlag("end", "SNAPENGAGE_END", "Stop exporting at date YYYY-MM-DD - Defaults to TODAY if unspecified")
	flgOutput := strEnvFlag("output", "SNAPENGAGE_OUTPUT", "Directory to write cases to")
	flgByDate := flag.Bool("bydate", false, "Store by date, makes for directories with smaller numbers of files")

	flgHelp := flag.Bool("help", false, "Displays help message")

	flag.Parse()

	if *flgHelp {
		flag.PrintDefaults()
		return
	}

	chkRequiredFlags(
		*flgToken, "token",
		*flgOrg, "org",
		*flgWidget, "widget",
		*flgStart, "start",
		*flgOutput, "output",
	)

	startDate := parseDateFlag("start", *flgStart)
	endDate := time.Now()

	if *flgEnd != "" {
		endDate = parseDateFlag("end", *flgEnd)
	}

	toStr := *flgOutput
	if *flgByDate {
		toStr = filepath.Join(toStr, "YYYY", "MM", "DD")
	}

	fmt.Printf("Exporting from %s to %s writing to %s\n",
		startDate.Format(dateFormat),
		endDate.Format(dirDateFormat),
		toStr,
	)

	if err := os.Mkdir(*flgOutput, 0755); err != nil && !errors.Is(err, os.ErrExist) {
		fmt.Println("Unable to create output path", err.Error())
		os.Exit(1)
	}

	gt := ghostport.New()

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("An unexpected panic!")

			debug.PrintStack()
			fmt.Println(gt)
			fmt.Println(r)

			os.Exit(2)
		}
	}()

	c := &client{
		Client:  &http.Client{Timeout: time.Minute, Transport: gt},
		token:   *flgToken,
		retries: 5,
	}

	var last time.Time

	queryURL := fmt.Sprintf(logsAPI, *flgOrg, *flgWidget, startDate.Format(dateFormat), endDate.Format(dateFormat))
	jsonResp, err := c.fetch(queryURL)

	for {
		if err != nil {
			fmt.Println(err)
			if !time.Now().IsZero() {
				fmt.Println("Last seen date:", last.Format(dateFormat))
			}
			os.Exit(1)
		}

		for _, v := range jsonResp.Cases {
			last = saveFile(*flgOutput, v, *flgByDate)
		}

		if !last.IsZero() {
			fmt.Print(" ", last.Format(dateFormat))
		}
		fmt.Println(" :D")

		if jsonResp.LinkToNextSetOfResults == "" {
			break
		}

		time.Sleep(time.Millisecond * 105)

		jsonResp, err = c.fetch(jsonResp.LinkToNextSetOfResults)
	}
}

func saveFile(dir string, c Case, byDate bool) time.Time {
	if byDate {
		dir = filepath.Join(dir, c.Created.Format(dirDateFormat))
		if err := os.MkdirAll(dir, 0755); err != nil && !errors.Is(err, os.ErrExist) {
			fmt.Println("Unable to create output path", err.Error())
			os.Exit(1)
		}
	}

	fn := filepath.Join(dir, c.ID+".json")
	if s, err := os.Stat(fn); err == nil && s.Size() > 0 {
		fmt.Print(",")
		return c.Created
	}

	f, err := os.Create(fn)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if err := json.NewEncoder(f).Encode(c); err != nil {
		panic(err)
	}

	fmt.Print(".")

	return c.Created
}
