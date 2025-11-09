package news

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type NewsClient struct {
	apiKey string
	client *http.Client
}

type Article struct {
	Source      string `json:"source"`
	Author      string `json:"author"`
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	URLToImage  string `json:"urlToImage"`
	PublishedAt string `json:"publishedAt"`
	Content     string `json:"content"`
}

type newsAPIResponse struct {
	Status       string `json:"status"`
	TotalResults int    `json:"totalResults"`
	Articles     []struct {
		Source struct {
			ID   interface{} `json:"id"`
			Name string      `json:"name"`
		} `json:"source"`
		Author      string `json:"author"`
		Title       string `json:"title"`
		Description string `json:"description"`
		URL         string `json:"url"`
		URLToImage  string `json:"urlToImage"`
		PublishedAt string `json:"publishedAt"`
		Content     string `json:"content"`
	} `json:"articles"`
}

func NewNewsClient(apiKey string) *NewsClient {
	return &NewsClient{
		apiKey: apiKey,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (nc *NewsClient) TopHeadlines(country, q string, pageSize int) ([]Article, error) {
	if pageSize <= 0 {
		pageSize = 20
	}
	// Always use 'everything' endpoint instead of 'top-headlines'
	base := "https://newsapi.org/v2/everything"
	if q == "" {
		q = "India" // default search
	}
	params := url.Values{}
	params.Set("q", q)
	params.Set("pageSize", fmt.Sprintf("%d", pageSize))

	req, err := http.NewRequest("GET", base+"?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Api-Key", nc.apiKey)

	resp, err := nc.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var nr newsAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&nr); err != nil {
		return nil, err
	}

	out := make([]Article, 0, len(nr.Articles))
	for _, a := range nr.Articles {
		out = append(out, Article{
			Source:      a.Source.Name,
			Author:      a.Author,
			Title:       a.Title,
			Description: a.Description,
			URL:         a.URL,
			URLToImage:  a.URLToImage,
			PublishedAt: a.PublishedAt,
			Content:     a.Content,
		})
	}
	return out, nil
}
