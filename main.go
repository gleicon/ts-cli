package main

import (
	"flag"
	"log"
	"os"
)

var DOTDIR = ".ts-cli"

func main() {

	labelToPrint := flag.String("p", "", "Label to print data")
	labelToChart := flag.Bool("c", false, "chart instead of list")
	printDatapoints := flag.Bool("d", false, "print all datapoints")
	daysInterval := flag.Int("i", -1, "days interval")
	flag.Usage = help
	flag.Parse()

	db, err := openOrCreateDB(DOTDIR)
	defer db.Close()
	if err != nil {
		log.Println(err)
		os.Exit(-1)
	}

	if *labelToPrint != "" {
		// print metric data
		listValuesForLabel(*labelToPrint, *labelToChart, *daysInterval,
			*printDatapoints, db)
	} else {
		// ingest metric
		ingestMetric(db)
	}
}
