package ridePandaAPI

type Config struct {
	// The hub id to query vehicles for. Can be retrieved through another graphQL query.
	HubId string `json:"hubId"`

	// No idea; necessary otherwise requests fail.
	RpFfId string `json:"rpFfId"`

	// An organization id; necessary otherwise requests fail.
	RpOrg string `json:"rpOrg"`

	// The user agent to make requests as.
	UserAgent string `json:"userAgent"`

	// The URL of the GraphQL server.
	ServerUrl string `json:"serverUrl"`
}
