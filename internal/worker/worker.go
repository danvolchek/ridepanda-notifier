package worker

import (
	"fmt"
	"github.com/danvolchek/ridepanda-notifier/internal"
	"github.com/danvolchek/ridepanda-notifier/internal/matching"
	"github.com/danvolchek/ridepanda-notifier/internal/notifications"
	"github.com/danvolchek/ridepanda-notifier/internal/ridePandaAPI"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"time"
)

type Worker struct {
	config   Config
	client   *ridePandaAPI.Client
	matcher  *matching.Matcher
	notifier *notifications.Notifier

	attemptNum int

	log *log.Logger
}

func init() {
	rand.Seed(time.Now().UnixMilli())
}

func New(config Config, client *ridePandaAPI.Client, matcher *matching.Matcher, notifier *notifications.Notifier) *Worker {
	if config.Jitter.underlying >= config.CheckFrequency.underlying {
		panic("worker: jitter must be less than check frequency")
	}

	return &Worker{
		config:   config,
		client:   client,
		matcher:  matcher,
		notifier: notifier,
		log:      internal.NewLogger("worker"),
	}
}

func (w *Worker) CheckOnce() {
	w.checkAndNotify(true)
}

func (w *Worker) Start() {
	w.checkAndNotify(true)

	var next time.Time

	check := func() {
		calculatedJitter := time.Duration(rand.Float64() * float64(w.config.Jitter.underlying))
		w.log.Printf("Waiting for %s before checking...", calculatedJitter)
		time.Sleep(calculatedJitter)
		w.attemptNum += 1
		w.checkAndNotify(false)
		w.log.Printf("Next check is at %v (send nothing: %v)\n", next.Format(time.RFC3339), w.willNotifyNothing())
		w.log.Println()

		if w.attemptNum == w.config.NotifyNothing {
			w.attemptNum = 0
		}
	}

	now := time.Now()
	next = time.Date(now.Year(), now.Month(), now.Day(), w.config.StartHour, w.config.StartMinute, 0, 0, time.Local)
	for next.Before(now) {
		next = next.Add(w.config.CheckFrequency.underlying)
		w.attemptNum = (w.attemptNum + 1) % w.config.NotifyNothing
	}

	w.log.Printf("Next check is at %v (send nothing: %v)\n", next.Format(time.RFC3339), w.willNotifyNothing())
	w.log.Println()

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt)

	timer := time.NewTimer(next.Sub(time.Now()))

	for {
		select {
		case <-stopChan:
			w.log.Println("Received interrupt, exiting")
			return
		case <-timer.C:
			next = next.Add(w.config.CheckFrequency.underlying)
			go check()
			timer.Reset(next.Sub(time.Now()))
		}
	}
}

func (w *Worker) willNotifyNothing() bool {
	return (w.attemptNum+1)%w.config.NotifyNothing == 0
}

func (w *Worker) checkAndNotify(forceNothing bool) {
	err := w.checkForVehicles(forceNothing)
	if err != nil {
		w.log.Printf("[error] something went wrong checking for vehicles: %s", err)
		err = w.notifier.NotifyError(err)
		if err != nil {
			w.log.Println("Failed to notify about an error:", err)
		}
	}
}

func (w *Worker) checkForVehicles(forceNothing bool) error {
	w.log.Println("Checking for vehicles")
	vehicles, err := w.client.GetVehicles()
	if err != nil {
		return fmt.Errorf("couldn't get vehicles: %v", err)
	}

	matches, err := w.matcher.Match(vehicles)
	if err != nil {
		return fmt.Errorf("couldn't match vehicles: %v", err)
	}

	if len(matches) > 0 {
		err = w.notifier.Notify(matches)
		if err != nil {
			return fmt.Errorf("couldn't notify about vehicles: %v", err)
		}
	}

	if len(matches) == 0 && (forceNothing || w.attemptNum%w.config.NotifyNothing == 0) {
		err = w.notifier.NotifyNothing()
		if err != nil {
			return fmt.Errorf("couldn't notify about nothing: %v", err)
		}
	}

	return nil
}
