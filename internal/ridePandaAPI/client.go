package ridePandaAPI

import (
	"context"
	"github.com/danvolchek/ridepanda-notifier/internal"
	"github.com/shurcooL/graphql"
	"golang.org/x/exp/slices"
	"log"
	"net/http"
	"strings"
)

type roundTripper struct {
	underlying http.RoundTripper

	config Config
}

func (r *roundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	// required otherwise request fails
	req.Header["rp-ff-id"] = []string{r.config.RpFfId}
	req.Header["rp-org"] = []string{r.config.RpOrg}

	// just in case
	req.Header.Set("User-Agent", r.config.UserAgent)

	//fmt.Println(req)

	return r.underlying.RoundTrip(req)
}

type Client struct {
	config *Config
	client *graphql.Client

	log *log.Logger
}

func NewClient(config Config) *Client {
	httpClient := &http.Client{
		Transport: &roundTripper{
			config:     config,
			underlying: http.DefaultTransport,
		},
	}

	client := graphql.NewClient(config.ServerUrl, httpClient)

	return &Client{
		config: &config,
		client: client,
		log:    internal.NewLogger("api"),
	}
}

func (r *Client) GetVehicles() ([]Vehicle, error) {
	var query struct {
		Vehicles []Vehicle `graphql:"vehicles(hubId: $hubId)"`
	}

	err := r.client.Query(context.Background(), &query, map[string]interface{}{
		"hubId": r.config.HubId,
	})

	if err != nil {
		r.log.Printf("graphql query failed: %s", err)
		return nil, err
	}

	vehicles := query.Vehicles
	slices.SortFunc(vehicles, func(a, b Vehicle) bool {
		return strings.Compare(a.Name(), b.Name()) == -1
	})

	r.log.Printf("The RidePanda API returned %d vehicles (%d in stock)\n", len(vehicles), len(internal.Filter(vehicles, func(v Vehicle) bool { return v.InStock })))

	return query.Vehicles, nil
}
