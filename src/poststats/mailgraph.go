/*
	Name: mailgraph
	Description: A mail grapher from csv file
	Author: Enrico Bianchi
*/

package main

import (
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
)

const VERSION = "0.0.1"

func read(input string) {
	// TODO read file and return result
}

func write(output string) {
	// TODO: write result data to graph
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	input := kingpin.Arg("input", "CSV file to graph").Required().String()
	output := kingpin.Arg("output", "PNG file to save result").Required().String()

	kingpin.Version(VERSION)
	kingpin.Parse()

	read(*input)
	write(*output)
}
