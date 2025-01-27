# Usage example

Let's say we have a bucket "example" with the following data.

```txt
myMeasurement,tag=tag1 stringValue="hello",intValue=123i 1000000000
myMeasurement,tag=tag1 stringValue="world",intValue=546i 2000000000
```

Here's an example of how to get the data in JSON format:

example_query.flux:

```flux
from(bucket: "example")
    |> range(start: 0)
    |> filter(fn: (r) => r._measurement == "myMeasurement")
    |> pivot(rowKey: ["_time"], columnKey: ["_field"], valueColumn: "_value")
    |> drop(columns: ["_start", "_stop", "_measurement"])
```

Performing the query:

```sh
influxdb-query-as-json \
    -u=http://localhost:8086 \
    -o=MyOrg \
    -t=MyToken_XXXXXXXXXXXXXXXX \
    -f=example_query.flux | jq
```

Resulting output:

```json
{
  "_time": "1970-01-01T00:00:01Z",
  "intValue": 123,
  "stringValue": "hello",
  "tag": "tag1"
}
{
  "_time": "1970-01-01T00:00:02Z",
  "intValue": 546,
  "stringValue": "world",
  "tag": "tag1"
}
```
