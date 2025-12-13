package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

type Item struct {
	Link   string   `json:"link"`
	Skills []string `json:"skills"`
}

var BadData = []string{
	"ability to",
	"experience",
	"bachelor",
	"master",
	"degree",
	"responsible",
	"ensure",
	"understanding of",
	"united states",
}

func main() {
	verbose := flag.Bool("verbose", false, "enable verbose logging")
	full := flag.Bool("full", false, "enable full file processing")

	flag.Parse()

	if *verbose {
		fmt.Println("Verbose mode enabled")
	}

	filepath := "dataset/job_skills_sample.csv"

	if *full {
		fmt.Println("Verbose mode enabled")
		filepath = "dataset/job_skills.csv"
	}

	fmt.Println("Program running...")

	f, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("Error closing file: %v", err)
		}
	}()

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
			Link: record[0],
		}

		for skill := range strings.SplitSeq(record[1], ",") {
			words := strings.Fields(skill)
			if len(words) > 10 {
				continue
			}

			skill = strings.TrimSuffix(skill, ",")
			skill = strings.TrimSpace(skill)

			lower := strings.ToLower(skill)
			skip := false
			for _, bad := range BadData {
				if strings.Contains(lower, bad) {
					skip = true
					break
				}
			}
			if skip {
				continue
			}

			item.Skills = append(item.Skills, skill)
		}

		if *verbose {
			fmt.Printf("Link: %s - Number of Skills: %v \n", item.Link, len(item.Skills))
		}

		items = append(items, item)

	}
	fmt.Println(len(items))
	elapsed := time.Since(start)
	fmt.Printf("Duration: %vms\n", elapsed.Milliseconds())
	fmt.Printf("Duration: %v\n", elapsed)

	skills := map[string]bool{}
	skillCount := 0

	for _, v := range items {
		for _, skill := range v.Skills {
			skillCount++
			skills[skill] = true
		}
	}

	fmt.Println("Unique skills count:", len(skills))
	fmt.Println("Total skills count:", skillCount)
	elapsed = time.Since(start)
	fmt.Printf("Duration: %vms\n", elapsed.Milliseconds())
	fmt.Printf("Duration: %v\n", elapsed)

	w, err := os.Create("output/unique_skills.csv")
	if err != nil {
		log.Fatalf("Failed to create output file: %v", err)
	}
	defer func() {
		if err := w.Close(); err != nil {
			log.Printf("Error closing file: %v", err)
		}
	}()

	header := []string{"id", "skill"}

	writer := csv.NewWriter(w)
	if err := writer.Write(header); err != nil {
		log.Fatalf("Failed to write header to CSV: %v", err)
	}

	id := 1
	for skill := range skills {
		cleanedSkill := strings.TrimSpace(skill)
		row := []string{fmt.Sprintf("%d", id), cleanedSkill}
		if err := writer.Write(row); err != nil {
			log.Fatalf("Failed to write row to CSV: %v", err)
		}
		id++
	}

	if err := writer.Error(); err != nil {
		log.Fatalf("Error writing to CSV: %v", err)
	}

	log.Println("Unique skills written to output/unique_skills.csv successfully")
}
