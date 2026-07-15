package tmdb

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type TMDbClient struct {
	apiKey     string
	httpClient http.Client
}

func NewTMDbClient(apiKey string, httpClient http.Client) *TMDbClient {
	return &TMDbClient{
		apiKey:     apiKey,
		httpClient: httpClient,
	}
}

type Movie struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	PosterPath  string  `json:"poster_path"`
	ReleaseDate string  `json:"release_date"`
	Overview    string  `json:"overview"`
	Popularity  float64 `json:"popularity"`
	VoteAverage float64 `json:"vote_average"`
}
type SearchResponse struct {
	Results []Movie `json:"results"`
}

type CastMember struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	ProfilePath string `json:"profile_path"`
}

type Cast struct {
	Cast []CastMember `json:"cast"`
}

type CrewMember struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	ProfilePath string `json:"profile_path"`
	Job         string `json:"job"`
}

type CreditsResponse struct {
	Cast []CastMember `json:"cast"`
	Crew []CrewMember `json:"crew"`
}

func (tm *TMDbClient) SearchMovie(query string) ([]Movie, error) {
	url := fmt.Sprintf("https://api.themoviedb.org/3/search/movie?api_key=%s&query=%s", tm.apiKey, query)
	resp, err := tm.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var searchResponse SearchResponse
	err = json.Unmarshal(bodyBytes, &searchResponse)
	if err != nil {
		return nil, err
	}
	return searchResponse.Results, err

}

func (tm *TMDbClient) GetMovieCredits(movieID int) (*CreditsResponse, error) {
	url := fmt.Sprintf("https://api.themoviedb.org/3/movie/%d/credits?api_key=%s", movieID, tm.apiKey)

	resp, err := tm.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var creditResponse CreditsResponse
	err = json.Unmarshal(bodyBytes, &creditResponse)
	if err != nil {
		return nil, err
	}
	return &creditResponse, err
}
