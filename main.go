package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/kaneshin/go-piyolog/piyolog"
)

func main() {
	body, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	daily, err := piyolog.ParseDaily(string(body))
	if err != nil {
		log.Fatal(err)
	}
	for _, plog := range daily.Logs {
		switch v := plog.(type) {
		case piyolog.FormulaLog:
			fmt.Printf("%s %s %s\n", v.CreatedAt().Format("15:04"), v.Type(), v.Amount)
		case piyolog.SolidLog:
			fmt.Printf("%s %s\n", v.CreatedAt().Format("15:04"), v.Type())
		}
	}
}
