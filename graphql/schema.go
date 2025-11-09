package graphql

import (
	"encoding/json"
	"net/http"

	"news-graphql/news"

	"github.com/graphql-go/graphql"
)

type SchemaServer struct {
	schema graphql.Schema
	client *news.NewsClient
}

func NewSchema(client *news.NewsClient) *SchemaServer {
	articleType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Article",
		Fields: graphql.Fields{
			"source":      &graphql.Field{Type: graphql.String},
			"author":      &graphql.Field{Type: graphql.String},
			"title":       &graphql.Field{Type: graphql.String},
			"description": &graphql.Field{Type: graphql.String},
			"url":         &graphql.Field{Type: graphql.String},
			"urlToImage":  &graphql.Field{Type: graphql.String},
			"publishedAt": &graphql.Field{Type: graphql.String},
			"content":     &graphql.Field{Type: graphql.String},
		},
	})

	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"topHeadlines": &graphql.Field{
				Type: graphql.NewList(articleType),
				Args: graphql.FieldConfigArgument{
					"country":  &graphql.ArgumentConfig{Type: graphql.String},
					"q":        &graphql.ArgumentConfig{Type: graphql.String},
					"pageSize": &graphql.ArgumentConfig{Type: graphql.Int},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					country, _ := p.Args["country"].(string)
					q, _ := p.Args["q"].(string)
					pageSize, _ := p.Args["pageSize"].(int)
					return client.TopHeadlines(country, q, pageSize)
				},
			},
		},
	})

	schema, _ := graphql.NewSchema(graphql.SchemaConfig{
		Query: rootQuery,
	})

	return &SchemaServer{schema: schema, client: client}
}

func (s *SchemaServer) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Query     string                 `json:"query"`
			Variables map[string]interface{} `json:"variables"`
		}
		if r.Method == "GET" {
			// allow quick browser test via ?query=...
			req.Query = r.URL.Query().Get("query")
		} else {
			dec := json.NewDecoder(r.Body)
			if err := dec.Decode(&req); err != nil {
				http.Error(w, "invalid json", http.StatusBadRequest)
				return
			}
		}

		result := graphql.Do(graphql.Params{
			Schema:         s.schema,
			RequestString:  req.Query,
			VariableValues: req.Variables,
			Context:        r.Context(),
		})
		if len(result.Errors) > 0 {
			w.WriteHeader(http.StatusBadRequest)
		}
		w.Header().Set("Content-Type", "application/json")
		enc := json.NewEncoder(w)
		_ = enc.Encode(result)
	})
}
