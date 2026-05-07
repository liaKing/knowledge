package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

type WeatherResult struct {
	City        string `json:"city"`
	Weather     string `json:"weather"`
	Temperature string `json:"temperature"`
	Wind        string `json:"wind"`
	Advice      string `json:"advice"`
}

func main() {
	city := flag.String("city", "", "city name")
	flag.Parse()

	if *city == "" {
		fmt.Fprintln(os.Stderr, "city is required")
		os.Exit(1)
	}

	result := WeatherResult{
		City:        *city,
		Weather:     "晴",
		Temperature: "26°C",
		Wind:        "东南风 2 级",
		Advice:      "适合外出，建议注意防晒并适量补水。",
	}

	if err := json.NewEncoder(os.Stdout).Encode(result); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
