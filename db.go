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

	gds "github.com/gleicon/go-descriptive-statistics"

	"github.com/boltdb/bolt"
	"github.com/joliv/spark"
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

func listValuesForLabel(metricName string, chart bool, db *bolt.DB) {
	listValuesForLabelPerDays(metricName, 1, chart, db)
}

func listValuesForLabelPerDays(metricName string, days int, printChart bool, db *bolt.DB) {
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(metricName))
		if b == nil {
			fmt.Println("This label or metric does not exists: ", metricName)
			os.Exit(-1)

		}
		c := b.Cursor()
		min := []byte(time.Now().AddDate(0, 0, days*-1).Format(time.RFC3339))
		max := []byte(time.Now().AddDate(0, 0, 0).Format(time.RFC3339))
		if printChart {
			// last date and a sparkline with max(value) printed
			var values gds.Enum
			var ndp string
			var lastVal string
			for k, v := c.First(); k != nil; k, v = c.Next() {
				// only plots diff values as the space is limited
				if string(v) == lastVal {
					continue
				}
				lastVal = string(v)
				vf, _ := strconv.ParseFloat(string(v), 64)
				t, err := time.Parse(time.RFC3339, string(k))

				if err != nil {
					fmt.Println(err)
				}
				values = append(values, vf)

				ndp = t.String()
			}
			fmt.Printf("Max: %.2f MB\n", values.Percentile(100)/1024/1024)
			fmt.Println("Newest datapoint: ", ndp)
			sparkline := spark.Line(values)
			fmt.Println(sparkline)

		} else {
			for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {
				fmt.Println(string(k), string(v))
			}
		}
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
