package site

import (
	"encoding/json"
	"sort"
)

// siteBaseURL is unconfirmed until the operator selects the production domain.
const siteBaseURL = "https://siduri.ai"

func init() {
	RegisterHead(HeadFragment{
		Name: "a3-meta",
		Render: func(data PageData) interfaceComponent {
			return a3Head(data)
		},
	})
}

func siteURL(path string) string {
	return siteBaseURL + path
}

func canonicalURL(post Post) string {
	return siteURL("/journal/" + post.Slug + "/")
}

func publishedPosts(posts []Post) []Post {
	result := make([]Post, 0, len(posts))
	for _, post := range posts {
		if !post.Draft {
			result = append(result, post)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Date != result[j].Date {
			return result[i].Date > result[j].Date
		}
		return result[i].Slug < result[j].Slug
	})
	return result
}

func metaDescription(posts []Post) string {
	posts = publishedPosts(posts)
	if len(posts) > 0 {
		return posts[0].Summary
	}
	return "A journal about shipping software with agents while keeping a person responsible for the result."
}

type schemaGraph struct {
	Context string         `json:"@context"`
	Graph   []schemaEntity `json:"@graph"`
}

type schemaEntity struct {
	Type             string           `json:"@type"`
	ID               string           `json:"@id,omitempty"`
	Name             string           `json:"name,omitempty"`
	URL              string           `json:"url,omitempty"`
	Headline         string           `json:"headline,omitempty"`
	Description      string           `json:"description,omitempty"`
	DatePublished    string           `json:"datePublished,omitempty"`
	DateModified     string           `json:"dateModified,omitempty"`
	Author           *schemaReference `json:"author,omitempty"`
	MainEntityOfPage *schemaReference `json:"mainEntityOfPage,omitempty"`
}

type schemaReference struct {
	ID string `json:"@id"`
}

func structuredData(posts []Post) string {
	posts = publishedPosts(posts)
	entities := make([]schemaEntity, 0, len(posts)+1)
	entities = append(entities, schemaEntity{
		Type: "Person",
		ID:   siteURL("/#person"),
		Name: "Siduri",
		URL:  siteURL("/about/"),
	})
	for _, post := range posts {
		modified := post.Updated
		if modified == "" {
			modified = post.Date
		}
		entities = append(entities, schemaEntity{
			Type:             "Article",
			ID:               canonicalURL(post),
			Headline:         post.Title,
			Description:      post.Summary,
			DatePublished:    post.Date,
			DateModified:     modified,
			Author:           &schemaReference{ID: siteURL("/#person")},
			MainEntityOfPage: &schemaReference{ID: canonicalURL(post)},
		})
	}
	data, err := json.Marshal(schemaGraph{
		Context: "https://schema.org",
		Graph:   entities,
	})
	if err != nil {
		panic("site: encode JSON-LD")
	}
	return string(data)
}
