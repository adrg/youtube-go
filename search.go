package youtube

import (
	"encoding/json"

	"net/url"
	"strconv"
)

const gdataURL = "gdata.youtube.com/feeds/api/videos"

type SearchParams struct {
	Query      string
	Page       int
	MaxResults int
}

type response struct {
	Data struct {
		Videos []*Video `json:"items"`
	} `json:"data"`
}

func generateSearchURL(params *SearchParams) string {
	searchURL := url.URL{Host: gdataURL, Scheme: "https"}
	searchURL.RawQuery = url.Values{
		"q":            {params.Query},
		"v":            {"2.1"},
		"alt":          {"jsonc"},
		"pretty-print": {"1"},
		"orderby":      {"relevance"},
		"start-index":  {strconv.Itoa(params.Page*params.MaxResults + 1)},
		"max-results":  {strconv.Itoa(params.MaxResults)},
	}.Encode()

	return searchURL.String()
}

func Search(params *SearchParams) ([]*Video, error) {
	data, err := getURLData(generateSearchURL(params))
	if err != nil {
		return nil, err
	}

	var resp response
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	return resp.Data.Videos, nil
}
