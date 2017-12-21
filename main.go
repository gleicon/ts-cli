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
		listValuesForLabel(*labelToPrint, *labelToChart, db)
	} else {
		// ingest metric
		ingestMetric(db)
	}
}
