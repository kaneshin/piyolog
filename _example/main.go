package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

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

	count := 0
	sum := 0
	unit := ""
	daily := data.Entries[0]
	for _, plog := range daily.Logs {
		switch v := plog.(type) {
		case piyolog.FormulaLog:
			// print only formula log
			fmt.Printf("%s\n", v)
			// to calculate the formula average
			sum += v.Amount
			count++
			unit = v.Unit
		}
	}
	fmt.Printf("Avg: %.2f%s\n", float64(sum)/float64(count), unit)
}
