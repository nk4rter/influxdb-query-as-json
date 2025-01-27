package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	flags "github.com/jessevdk/go-flags"
)

func main() {
	var opts struct {
		Url   string `long:"url"   short:"u" description:"URL that InfluxDB is bound to" required:"true"`
		Org   string `long:"org"   short:"o" description:"Organization name"             required:"true"`
		Token string `long:"token" short:"t" description:"Access token"                  required:"true"`
		File  string `long:"file"  short:"f" description:"Query file, '-' means stdin"   required:"true"`
	}
	_, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(1)
	}

	var file string
	if opts.File == "-" {
		file = "/dev/stdin"
	} else {
		file = opts.File
	}
	query, err := os.ReadFile(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Failed to read query: %s\n", err)
		os.Exit(1)
	}

	client := influxdb2.NewClient(opts.Url, opts.Token)
	defer client.Close()

	queryAPI := client.QueryAPI(opts.Org)

	result, err := queryAPI.Query(context.Background(), string(query))
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Database query failed: %s\n", err)
		os.Exit(1)
	}

	for result.Next() {
		values := result.Record().Values()
		delete(values, "table")
		delete(values, "result")

		jsonData, err := json.Marshal(values)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: JSON formatting failed: %s\n", err)
		}
		fmt.Printf("%s\n", jsonData)
	}
	if result.Err() != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Query parsing error: %s\n", result.Err().Error())
		os.Exit(1)
	}
}
