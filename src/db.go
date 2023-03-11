package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/table"
	gds "github.com/gleicon/go-descriptive-statistics"

	"github.com/boltdb/bolt"
)

func openOrCreateDB(path string) (*bolt.DB, error) {
	dotdir := filepath.Join(GetEnvHomeDir(), path)

	os.MkdirAll(dotdir, os.ModePerm)
	dbfile := filepath.Join(dotdir, "timeseries.db")

	db, err := bolt.Open(dbfile, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}
	return db, nil

}

func listValuesForLabel(metricName string, days int,
	printDP bool, db *bolt.DB) {

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(metricName))
		if b == nil {
			fmt.Println("This label or metric does not exists: ", metricName)
			os.Exit(-1)

		}
		c := b.Cursor()
		var values gds.Enum
		datapoints := 0
		firstDataPoint := ""
		lastDataPoint := ""
		rows := []table.Row{}

		if days == -1 {
			for k, v := c.First(); k != nil; k, v = c.Next() {
				rows = append(rows, table.Row{string(k), string(v)})
				if firstDataPoint == "" {
					firstDataPoint = string(k)
				}
				lastDataPoint = string(k)

				vf, _ := strconv.ParseFloat(string(v), 64)
				values = append(values, vf)
				if printDP {
					fmt.Println(string(k), string(v))
				}
				datapoints++
			}
		} else {
			min := []byte(time.Now().AddDate(0, 0, days*-1).Format(time.RFC3339))
			max := []byte(time.Now().AddDate(0, 0, 0).Format(time.RFC3339))
			for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {
				rows = append(rows, table.Row{string(k), string(v)})
				if firstDataPoint == "" {
					firstDataPoint = string(k)
				}
				lastDataPoint = string(k)

				vf, _ := strconv.ParseFloat(string(v), 64)
				values = append(values, vf)
				if printDP {
					fmt.Println(string(k), string(v))
				}
				datapoints++
			}

		}
		headers := fmt.Sprintf("Timeseries name: %s\n", metricName)
		headers += fmt.Sprintf("Max: %.2f \n", values.Percentile(100))
		headers += fmt.Sprintf("99 percentile: %.2f \n", values.Percentile(99))
		headers += fmt.Sprintf("First datapoint: %s\n", firstDataPoint)
		headers += fmt.Sprintf("Last datapoint: %s\n", lastDataPoint)
		headers += fmt.Sprintf("Datapoints: %d\n", datapoints)

		columns := []table.Column{
			{Title: "Timestamp", Width: 30},
			{Title: "Value", Width: 10},
		}

		RenderTable(&columns, &rows, headers)

		return nil
	})

}

func listAllValuesForLabel(metricName string, days int, db *bolt.DB) {
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(metricName))
		if b == nil {
			fmt.Println("This label or metric does not exists: ", metricName)
			os.Exit(-1)

		}
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			fmt.Printf("key=%s, value=%s\n", k, v)
		}

		return nil
	})

}

func ingestMetric(db *bolt.DB) {
	inMetric := bufio.NewReader(os.Stdin)
	line, _, err := inMetric.ReadLine()
	lm := string(line)

	if strings.Contains(lm, "=") != true {
		log.Println("Metric shoud be in name=value format")
		log.Println(lm)
		os.Exit(-1)
	}
	metric := strings.Split(lm, "=")
	label := metric[0]
	//value, err := strconv.ParseFloat(metric[1], 64)
	value := metric[1]
	if err != nil {
		log.Println("Wrong value type: ", err)
		os.Exit(-1)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(label))
		if err != nil {
			return fmt.Errorf("Error creating bucket %s: %s", label, err)
		}
		now := time.Now().Format(time.RFC3339)
		err = b.Put([]byte(now), []byte(value))
		return nil
	})

}
