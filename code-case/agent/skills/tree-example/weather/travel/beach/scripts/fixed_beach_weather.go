package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

type BeachWeatherResult struct {
	Place       string `json:"place"`
	Weather     string `json:"weather"`
	Temperature string `json:"temperature"`
	Wind        string `json:"wind"`
	UV          string `json:"uv"`
	Risk        string `json:"risk"`
	Advice      string `json:"advice"`
}

func main() {
	place := flag.String("place", "", "beach or city name")
	flag.Parse()

	if *place == "" {
		fmt.Fprintln(os.Stderr, "place is required")
		os.Exit(1)
	}

	result := BeachWeatherResult{
		Place:       *place,
		Weather:     "晴到多云",
		Temperature: "26°C",
		Wind:        "海风 3 级",
		UV:          "较强",
		Risk:        "长时间暴晒风险较高",
		Advice:      "适合海边游玩，建议带防晒用品和薄外套。",
	}

	if err := json.NewEncoder(os.Stdout).Encode(result); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
