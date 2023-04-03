package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/danvolchek/ridepanda-notifier/internal"
	"github.com/danvolchek/ridepanda-notifier/internal/matching"
	"github.com/danvolchek/ridepanda-notifier/internal/notifications"
	"github.com/danvolchek/ridepanda-notifier/internal/ridePandaAPI"
	"github.com/danvolchek/ridepanda-notifier/internal/worker"
	"github.com/tidwall/jsonc"
	"os"
)

// TODO: config field validation

type config struct {
	Worker        worker.Config        `json:"worker"`
	RidePandaAPI  ridePandaAPI.Config  `json:"ridePandaAPI"`
	Matcher       matching.Config      `json:"matcher"`
	Notifications notifications.Config `json:"notifications"`
}

func loadConfig() (config, error) {
	data, err := os.ReadFile("config.json")
	if err != nil {
		return config{}, fmt.Errorf("could not read config file: %v", err)
	}

	var cfg config
	err = json.Unmarshal(jsonc.ToJSON(data), &cfg)
	if err != nil {
		return config{}, fmt.Errorf("could not parse config file: %v", err)
	}

	return cfg, nil
}

func main() {
	once := flag.Bool("once", false, "check once and exit")
	flag.Parse()

	log := internal.NewLogger("main")

	log.Println("Starting!")
	cfg, err := loadConfig()
	if err != nil {
		log.Println("could not load config:", err)
		os.Exit(1)
	}

	client := ridePandaAPI.NewClient(cfg.RidePandaAPI)
	matcher := matching.NewMatcher(cfg.Matcher)
	notifier := notifications.NewNotifier(cfg.Notifications)
	wrkr := worker.New(cfg.Worker, client, matcher, notifier)

	if *once {
		wrkr.CheckOnce()
	} else {
		wrkr.Start()
	}

	log.Println("Finished!")
}
