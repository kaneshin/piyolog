package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/kaneshin/piyolog"
)

func main() {
	body, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	data, err := piyolog.Parse(string(body))
	if err != nil {
		log.Fatal(err)
	}

	daily := data.Entries[len(data.Entries)-1]
	fmt.Printf("Daily Report: %s\nBaby: %s (Birthday: %s)\n",
		daily.Date.Format(time.DateOnly),
		daily.Baby.Name,
		daily.Baby.DateOfBirth.Format(time.DateOnly),
	)

	milks := []piyolog.FormulaLog{}
	pees := []piyolog.PeeLog{}
	count := 0
	sum := 0
	unit := ""
	for _, plog := range daily.Logs {
		switch v := plog.(type) {
		case piyolog.FormulaLog:
			milks = append(milks, v)
			// to calculate the formula average
			sum += v.Amount
			count++
			unit = v.Unit
		case piyolog.PeeLog:
			pees = append(pees, v)
		}
	}
	fmt.Println("\n-- Milk Stats --")
	for _, milk := range milks {
		fmt.Printf("- %s\n", milk)
	}
	fmt.Printf("-> Avg: %.2f%s\n", float64(sum)/float64(count), unit)
	fmt.Println("\n-- Pee Stats --")
	for _, pee := range pees {
		fmt.Printf("- %s\n", pee)
	}
	fmt.Println("\n-- Comment --\n", daily.Journal)
}
