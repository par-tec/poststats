/*
	Name: mailgraph
	Description: A mail grapher from csv file
	Author: Enrico Bianchi
*/

package main

import (
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
	"time"
)

const VERSION = "0.0.1"

type Data struct {
	Date  *time.Time
	Count int
	Size  int
}

func read(input string) []Data {
	var result []Data

	// TODO read file and return result
	return result
}

func write(output string, data []Data) {
	// TODO: write result data to graph
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	input := kingpin.Arg("input", "CSV file to graph").Required().String()
	output := kingpin.Arg("output", "PNG file to save result").Required().String()

	kingpin.Version(VERSION)
	kingpin.Parse()

	data := read(*input)
	write(*output, data)
}
