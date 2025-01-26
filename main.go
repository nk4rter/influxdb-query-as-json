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
		Bucket      string `long:"bucket"      short:"b" description:"TODO" required:"true"`
		Url         string `long:"url"         short:"u" description:"TODO" required:"true"`
		Org         string `long:"org"         short:"o" description:"TODO" required:"true"`
		Token       string `long:"token"       short:"t" description:"TODO" required:"true"`
		Measurement string `long:"measurement" short:"m" description:"TODO" required:"true"`
	}
	_, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(1)
	}

	query := fmt.Sprintf(`
from(bucket: "%s")
  |> range(start: 0)
  |> last()
  |> filter(fn: (r) => r["_measurement"] == "%s")
`, opts.Bucket, opts.Measurement)

	client := influxdb2.NewClient(opts.Url, opts.Token)
	defer client.Close()

	queryAPI := client.QueryAPI(opts.Org)

	data := map[string]interface{}{}

	result, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Database query failed: %s\n", err)
		os.Exit(1)
	}

	for result.Next() {
		data[result.Record().Field()] = result.Record().Value()
		data["_time"] = result.Record().Time().UnixNano()
	}
	if result.Err() != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Query parsing error: %s\n", result.Err().Error())
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: JSON formatting failed: %s\n", err)
	}
	fmt.Printf("%s\n", jsonData)
}
