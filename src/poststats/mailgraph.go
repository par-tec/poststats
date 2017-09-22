/*
	Name: mailgraph
	Description: A mail grapher from csv file
	Author: Enrico Bianchi
*/

package main

import (
	"encoding/csv"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"gopkg.in/alecthomas/kingpin.v2"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

const VERSION = "0.0.1"
const SIZESTAT = "Size"
const COUNTSTAT = "Count"

type Data struct {
	Date  time.Time
	Count int
	Size  int
}

func read(input string) []Data {
	var result []Data
	var date time.Time
	var count, size int
	var err error

	layout := "2006-01-02 15:04:05 -0700 MST"

	file, err := os.Open(input)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	csvfile := csv.NewReader(file)

	for {
		record, err := csvfile.Read()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				log.Fatal(err)
			}
		}

		if date, err = time.Parse(layout, record[0]); err != nil {
			log.Fatal("Cannot parse date", err)
		}

		if count, err = strconv.Atoi(record[1]); err != nil {
			log.Fatal("Cannot parse counted emails", err)
		}

		if size, err = strconv.Atoi(record[1]); err != nil {
			log.Fatal("Cannot parse size emails", err)
		}

		result = append(result, Data{
			Date:  date,
			Count: count,
			Size:  size,
		})
	}

	return result
}

func write(generate, output string, data []Data) {
	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	p.Title.Text = "Mail statistics"
	p.X.Label.Text = "Date"
	p.Y.Label.Text = generate

	if err = plotutil.AddLinePoints(p, getpoints(data, generate)); err != nil {
		log.Fatal(err)
	}

	if err := p.Save(4*vg.Centimeter, 4*vg.Centimeter, output); err != nil {
		log.Fatal(err)
	}
}

func getpoints(data []Data, generate string) plotter.XYs {
	pts := make(plotter.XYs, len(data))

	for i := range pts {
		pts[i].X = float64(data[i].Date.Unix())
		if generate == COUNTSTAT {
			pts[i].Y = float64(data[i].Count)
		} else {
			pts[i].Y = float64(data[i].Size)
		}
	}

	return pts
}

func main() {
	var generate string

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	input := kingpin.Arg("input", "CSV file to graph").Required().String()
	output := kingpin.Arg("output", "PNG file to save result").Required().String()
	size := kingpin.Flag("size", "Generate graphic for size statistics").Short('S').Bool()
	count := kingpin.Flag("count", "Generate graphic for count statistics").Short('C').Bool()

	kingpin.Version(VERSION)
	kingpin.Parse()

	if *size && *count {
		kingpin.FatalUsage("Cannot generate graphic for size and count statistics in the same time")
	}

	if *size {
		generate = SIZESTAT
	} else if *count {
		generate = COUNTSTAT
	} else {
		kingpin.FatalUsage("Specify count or size grap statistics")
	}

	data := read(*input)
	write(generate, *output, data)
}
