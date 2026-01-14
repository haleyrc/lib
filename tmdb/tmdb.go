package tmdb

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	apiURL   = "https://api.themoviedb.org/3"
	htmlURL  = "https://www.themoviedb.org"
	imageURL = "https://image.tmdb.org/t/p/w1280"

	authenticationPath = "/authentication"
	movieCreditsPath   = "/movie/%d/credits"
	movieDetailPath    = "/movie/%d"
	searchMoviePath    = "/search/movie"
)

// MovieURL returns the URL of the movie's page on themoviedb.org.
func MovieURL(id int) string {
	return fmt.Sprintf("%s/movie/%d", htmlURL, id)
}

// SearchURL returns a URL for searching movies by query on themoviedb.org.
func SearchURL(query string) string {
	params := url.Values{}
	params.Add("query", query)
	params.Add("language", "en-US")
	return fmt.Sprintf("%s/search/movie?%s", htmlURL, params.Encode())
}

// Client makes calls to the TMDB API.
type Client struct {
	token string
}

// NewClient returns a Client that authenticates with the given API access
// token. It returns an error if token is empty.
func NewClient(token string) (*Client, error) {
	if token == "" {
		return nil, fmt.Errorf("tmdb: new client: token is required")
	}
	c := &Client{
		token: token,
	}
	return c, nil
}

// AuthenticateResponse is the response from the authentication endpoint.
type AuthenticateResponse struct {
	Success bool `json:"success"`
}

// Authenticate validates the access token. If it is valid, `Success` will be
// `true`, otherwise it will be `false`.
func (c *Client) Authenticate(ctx context.Context) (*AuthenticateResponse, error) {
	req, err := c.newRequest(ctx, authenticationPath)
	if err != nil {
		return nil, fmt.Errorf("authenticating: %w", err)
	}

	var resp AuthenticateResponse
	if err := c.do(req, &resp); err != nil {
		return nil, fmt.Errorf("authenticating: %w", err)
	}

	return &resp, nil
}

// GetMovieCreditsResponse is the response from the movie credits endpoint.
type GetMovieCreditsResponse struct {
	ID   int `json:"id"`
	Cast []struct {
		Adult              bool    `json:"adult"`
		Gender             int     `json:"gender"`
		ID                 int     `json:"id"`
		KnownForDepartment string  `json:"known_for_department"`
		Name               string  `json:"name"`
		OriginalName       string  `json:"original_name"`
		Popularity         float64 `json:"popularity"`
		ProfilePath        string  `json:"profile_path"`
		CastID             int     `json:"cast_id"`
		Character          string  `json:"character"`
		CreditID           string  `json:"credit_id"`
		Order              int     `json:"order"`
	} `json:"cast"`
}

// GetMovieCredits returns the cast for the movie with the given TMDB ID.
func (c *Client) GetMovieCredits(ctx context.Context, id int) (*GetMovieCreditsResponse, error) {
	req, err := c.newRequest(ctx, fmt.Sprintf(movieCreditsPath, id))
	if err != nil {
		return nil, fmt.Errorf("getting movie credits: %w", err)
	}

	var resp GetMovieCreditsResponse
	if err := c.do(req, &resp); err != nil {
		return nil, fmt.Errorf("getting movie credits: %w", err)
	}

	return &resp, nil
}

// GetMovieDetailResponse is the response from the movie detail endpoint.
type GetMovieDetailResponse struct {
	Adult               bool   `json:"adult"`
	BackdropPath        string `json:"backdrop_path"`
	BelongsToCollection *struct {
		ID           int    `json:"id"`
		Name         string `json:"name"`
		PosterPath   string `json:"poster_path"`
		BackdropPath string `json:"backdrop_path"`
	} `json:"belongs_to_collection"`
	Budget int `json:"budget"`
	Genres []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"genres"`
	Homepage            string   `json:"homepage"`
	ID                  int      `json:"id"`
	IMDBID              string   `json:"imdb_id"`
	OriginCountry       []string `json:"origin_country"`
	OriginalLanguage    string   `json:"original_language"`
	OriginalTitle       string   `json:"original_title"`
	Overview            string   `json:"overview"`
	Popularity          float64  `json:"popularity"`
	PosterPath          string   `json:"poster_path"`
	ProductionCompanies []struct {
		ID            int    `json:"id"`
		LogoPath      string `json:"logo_path"`
		Name          string `json:"name"`
		OriginCountry string `json:"origin_country"`
	} `json:"production_companies"`
	ProductionCountries []struct {
		ISO31661 string `json:"iso_3166_1"`
		Name     string `json:"name"`
	} `json:"production_countries"`
	ReleaseDate     string `json:"release_date"`
	Revenue         int    `json:"revenue"`
	Runtime         int    `json:"runtime"`
	SpokenLanguages []struct {
		EnglishName string `json:"english_name"`
		ISO6391     string `json:"iso_639_1"`
		Name        string `json:"name"`
	} `json:"spoken_languages"`
	Status      string  `json:"status"`
	Tagline     string  `json:"tagline"`
	Title       string  `json:"title"`
	Video       bool    `json:"video"`
	VoteAverage float64 `json:"vote_average"`
	VoteCount   int     `json:"vote_count"`
}

// PosterURL returns the full URL for the movie's poster image.
func (gmdr *GetMovieDetailResponse) PosterURL() string {
	return imageURL + gmdr.PosterPath
}

// GetMovieDetail returns detailed information for the movie with the given TMDB
// ID.
func (c *Client) GetMovieDetail(ctx context.Context, id int) (*GetMovieDetailResponse, error) {
	req, err := c.newRequest(ctx, fmt.Sprintf(movieDetailPath, id))
	if err != nil {
		return nil, fmt.Errorf("getting movie detail: %w", err)
	}

	var resp GetMovieDetailResponse
	if err := c.do(req, &resp); err != nil {
		return nil, fmt.Errorf("getting movie detail: %w", err)
	}

	return &resp, nil
}

// SearchMovieResponse is the response from the movie search endpoint.
type SearchMovieResponse struct {
	Page    int `json:"page"`
	Results []struct {
		Adult            bool    `json:"adult"`
		BackdropPath     string  `json:"backdrop_path"`
		GenreIDs         []int   `json:"genre_ids"`
		ID               int     `json:"id"`
		OriginalLanguage string  `json:"original_language"`
		OriginalTitle    string  `json:"original_title"`
		Overview         string  `json:"overview"`
		Popularity       float64 `json:"popularity"`
		PosterPath       string  `json:"poster_path"`
		ReleaseDate      string  `json:"release_date"`
		Title            string  `json:"title"`
		Video            bool    `json:"video"`
		VoteAverage      float64 `json:"vote_average"`
		VoteCount        int     `json:"vote_count"`
	} `json:"results"`
	TotalPages   int `json:"total_pages"`
	TotalResults int `json:"total_results"`
}

// SearchMovie returns a list of movie results that match the provided title.
func (c *Client) SearchMovie(ctx context.Context, title string) (*SearchMovieResponse, error) {
	req, err := c.newRequest(ctx, searchMoviePath)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	params := url.Values{}
	params.Add("query", title)
	req.URL.RawQuery = params.Encode()

	var resp SearchMovieResponse
	if err := c.do(req, &resp); err != nil {
		return nil, fmt.Errorf("searching movie: %w", err)
	}

	return &resp, nil
}

func (c *Client) newRequest(ctx context.Context, path string) (*http.Request, error) {
	url := apiURL + path
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+c.token)

	return req, nil
}

func (c *Client) do(req *http.Request, resp any) error {
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("making request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("handling response: %s", http.StatusText(res.StatusCode))
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %w", err)
	}

	if err := json.Unmarshal(body, &resp); err != nil {
		return fmt.Errorf("parsing response body: %w", err)
	}

	return nil
}
