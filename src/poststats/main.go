/*
	Name: poststats
	Description: A postfix log statistics tool
	Author: Enrico Bianchi
*/

package main

import (
	"bufio"
	"compress/gzip"
	"encoding/csv"
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type Mails struct {
	MailID string
	Date   time.Time
	Size   int
}

type Stats struct {
	Size    int
	Counted int
}

const VERSION = "0.0.1"

func aggregate(mails []Mails) map[time.Time]Stats {
	result := make(map[time.Time]Stats)

	for _, mail := range mails {
		if val, ok := result[mail.Date]; ok {
			val.Counted += 1
			val.Size += mail.Size

			result[mail.Date] = val
		} else {
			result[mail.Date] = Stats{
				Size:    mail.Size,
				Counted: 1,
			}
		}
	}

	return result
}

func parse(year int, line, queue string) Mails {
	var err error

	output := Mails{}

	layout := "Jan  2 15:04:05"

	if strings.Contains(line, strings.Join([]string{queue, "/qmgr"}, "")) &&
		!(strings.Contains(line, "removed") ||
			strings.Contains(line, "warning: ") ||
			strings.Contains(line, "status=expired") ||
			strings.Contains(line, "skipped")) {
		date := line[:15]
		output.Date, err = time.Parse(layout, date)
		if err != nil {
			log.Fatal(err, "\n\t", line)
		}
		output.Date = output.Date.AddDate(year-output.Date.Year(), 0, 0)

		parseline := line[strings.Index(line, ": ")+2:]
		split1 := strings.Split(parseline, ": ")
		if len(split1) < 2 {
			log.Fatal("Malformed line: "+line, "\n\t", line)
		}
		output.MailID = split1[0]

		split2 := strings.Split(split1[1], ",")
		if len(split2) < 3 {
			log.Fatal("Not enough values: "+strings.Join(split2, ","), "\n\t", line)
		}

		index := strings.Index(split2[1], "=") + 1
		if index == -1 {
			log.Fatal("Malformed line: "+split2[1], "\n\t", line)
		}
		size := split2[1][index:]

		output.Size, err = strconv.Atoi(size)
		if err != nil {
			log.Fatal(err, "\n\t", line)
		}
	}

	return output
}

func getyear(filename string) int {
	split := strings.Split(filename, "-")
	if len(split) < 2 {
		result := time.Now().Year()
		return result
	} else {
		layout := "20060102"
		date := strings.Split(split[1], ".")[0]
		result, err := time.Parse(layout, date)
		if err != nil {
			log.Fatal(err)
		}
		return result.Year()
	}
}

func process(filename, queue *string) []Mails {
	var result []Mails
	var scanner *bufio.Scanner

	file, err := os.Open(*filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	if strings.HasSuffix(*filename, ".gz") {
		gzfile, err := gzip.NewReader(file)
		if err != nil {
			log.Fatal(err)
		}
		defer gzfile.Close()

		scanner = bufio.NewScanner(gzfile)
	} else {
		scanner = bufio.NewScanner(file)
	}

	year := getyear(*filename)

	for scanner.Scan() {
		line := scanner.Text()

		processed := parse(year, line, *queue)
		if processed.MailID != "" {
			result = append(result, processed)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return result
}

func save(stats map[time.Time]Stats, output string, appendfile bool) {
	var file *os.File
	var err error

	if appendfile {
		file, err = os.OpenFile(output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	} else {
		file, err = os.Create(output)
	}

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	for key, value := range stats {
		err := writer.Write([]string{
			key.String(),
			strconv.Itoa(value.Counted),
			strconv.Itoa(value.Size),
		})

		if err != nil {
			log.Fatal(err)
		}
		writer.Flush()
	}
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	queue := kingpin.Arg("queue", "Queue to process").Required().String()
	file := kingpin.Arg("filename", "Logfile to process").Required().String()
	output := kingpin.Flag("output", "Set the csv file to save data").Short('O').Default("output.csv").String()
	appendfile := kingpin.Flag("appendfile", "Append to existent csv").Short('A').Bool()

	kingpin.Version(VERSION)
	kingpin.Parse()

	processed := process(file, queue)
	aggregated := aggregate(processed)
	save(aggregated, *output, *appendfile)
}
