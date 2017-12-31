## Time Series for Command LIne

Simple cli to gather, list and plot timeseries. I've built it to track LTE usage from my router as in 

```
$ curl my-wifi-router/cgi?a | grep Usage | ts-cli
```

Also used to track cassandra compaction over time. ANything that you use with wc to periodically show a value can be stored as a timeseries and reported.

That will pass usage=bytes to ts-cli, aggregate under the usage metric and I could follow it up on terminal or send an email w/ it. Used to quickly track metrics from logs while developing locally.

	ts-cli -> Time Series for Command LIne

	Stores and chart a timeseries for any label and value on terminal

	Data stored using BoltDB at ~/.ts-cli/timeseries.db

	echo 'mymetric=10' | ts-cli to store a new point using the current time. data format is "label=value"

	ts-cli -p <label> to print the points, e.g. ts-cli -p mymetric for the label mymetric

	Extra parameters: -c for chart, -d for all datapoints and -i <days> to select how many days back from now() you want. 

### Build

$ make clean

$ make all

### License

MIT
