package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

type DrivingWeatherResult struct {
	City       string `json:"city"`
	Weather    string `json:"weather"`
	Visibility string `json:"visibility"`
	RoadRisk   string `json:"road_risk"`
	Advice     string `json:"advice"`
}

func main() {
	city := flag.String("city", "", "city name")
	flag.Parse()

	if *city == "" {
		fmt.Fprintln(os.Stderr, "city is required")
		os.Exit(1)
	}

	result := DrivingWeatherResult{
		City:       *city,
		Weather:    "多云",
		Visibility: "正常",
		RoadRisk:   "路面风险较低",
		Advice:     "可以正常开车通勤，建议预留 10 分钟缓冲时间。",
	}

	if err := json.NewEncoder(os.Stdout).Encode(result); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
