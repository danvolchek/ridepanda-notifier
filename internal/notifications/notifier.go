package notifications

import (
	"errors"
	"fmt"
	"github.com/danvolchek/ridepanda-notifier/internal"
	"github.com/danvolchek/ridepanda-notifier/internal/matching"
	"github.com/danvolchek/ridepanda-notifier/internal/ridePandaAPI"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type Notifier struct {
	Config Config

	log *log.Logger
}

func NewNotifier(config Config) *Notifier {
	return &Notifier{
		Config: config,
		log:    internal.NewLogger("notifier"),
	}
}

func (n *Notifier) Notify(matches []matching.Match) error {
	var errs []string

	for _, match := range matches {
		message := "The following bikes matched your target:\n\nhttps://get.ridepanda.com/amazon/wizard/vehicles\n\n" + strings.Join(internal.Map(match.Vehicles, ridePandaAPI.Vehicle.NameWithVariants), "\n")
		title := fmt.Sprintf("Bikes: %s is available!", match.Name)

		n.log.Printf("Notifying: '%s' -> '%s'\n", title, strings.ReplaceAll(message, "\n", " "))
		err := sendMessage(n.Config.ServerURL, n.Config.AppToken, title, message)
		if err != nil {
			errs = append(errs, err.Error())
		}
	}
	if len(errs) != 0 {
		return errors.New(strings.Join(errs, ", "))
	}

	return nil
}

func (n *Notifier) NotifyNothing() error {
	message := "I'm still working!"
	title := "Bikes: No matches"

	n.log.Println("Notifying: no matches but still working properly")

	return sendMessage(n.Config.ServerURL, n.Config.AppToken, title, message)
}

func (n *Notifier) NotifyError(err error) error {
	title := "Bikes: Something went wrong"

	n.log.Printf("Notifying: about an error: '%s' -> '%s'\n", title, err.Error())

	return sendMessage(n.Config.ServerURL, n.Config.AppToken, title, err.Error())
}

func sendMessage(loc, token, title, message string) error {
	_, err := http.PostForm(fmt.Sprintf("%s/message?token=%s", loc, token),
		url.Values{"message": {message}, "title": {title}, "priority": {"7"}})

	return err
}
