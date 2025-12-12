package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type Item struct {
	Link   string   `json:"link"`
	Skills []string `json:"skills"`
}

func main() {
	fmt.Println("Hello, World!")

	f, err := os.Open("dataset/job_skills_sample.csv")
	if err != nil {
		panic(err)
	}

	defer f.Close()

	var items []Item
	r := csv.NewReader(f)
	start := time.Now()

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		item := Item{
			Link:   record[0],
			Skills: strings.Split(record[1], ","),
		}

		fmt.Printf("Link: %s - Number of Skills: %v \n", item.Link, len(item.Skills))

		items = append(items, item)

	}
	fmt.Println(len(items))
	elapsed := time.Since(start)
	fmt.Printf("Duration: %vms\n", elapsed.Milliseconds())
	fmt.Printf("Duration: %v\n", elapsed)
}
