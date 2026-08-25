package site

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strings"
	"time"
)

const rssContentNamespace = "http://purl.org/rss/1.0/modules/content/"

func init() {
	RegisterContent("a3-machine-readable", func(data PageData, routes *RouteSet) {
		posts := publishedPosts(data.Posts)
		routes.Register(Route{
			Name: "a3-machine-readable",
			Output: RouteOutput{Expand: func(PageData) []Output {
				return []Output{
					ByteOutput("feed.xml", renderRSS(posts)),
					ByteOutput("feed.json", renderJSONFeed(posts)),
					ByteOutput("sitemap.xml", renderSitemap(posts)),
					ByteOutput("robots.txt", renderRobots()),
					ByteOutput("llms.txt", renderLLMs(posts)),
				}
			}},
		})
	})
}

func renderRSS(posts []Post) []byte {
	posts = publishedPosts(posts)
	var output bytes.Buffer
	encoder := xml.NewEncoder(&output)
	writeXMLToken(encoder.EncodeToken(xml.ProcInst{
		Target: "xml",
		Inst:   []byte(`version="1.0" encoding="UTF-8"`),
	}))

	writeXMLToken(encoder.EncodeToken(xml.StartElement{
		Name: xml.Name{Local: "rss"},
		Attr: []xml.Attr{
			{Name: xml.Name{Local: "version"}, Value: "2.0"},
			{Name: xml.Name{Local: "xmlns:content"}, Value: rssContentNamespace},
		},
	}))
	writeXMLToken(encoder.EncodeToken(xml.StartElement{Name: xml.Name{Local: "channel"}}))
	writeXMLText(encoder, "title", "Siduri")
	writeXMLText(encoder, "link", siteURL("/"))
	writeXMLText(encoder, "description", "A journal about shipping software with agents while keeping a person responsible for the result.")
	for _, post := range posts {
		body := renderMarkdown(post.Body)
		writeXMLToken(encoder.EncodeToken(xml.StartElement{Name: xml.Name{Local: "item"}}))
		writeXMLText(encoder, "title", post.Title)
		writeXMLText(encoder, "link", canonicalURL(post))
		writeXMLToken(encoder.EncodeToken(xml.StartElement{
			Name: xml.Name{Local: "guid"},
			Attr: []xml.Attr{{Name: xml.Name{Local: "isPermaLink"}, Value: "true"}},
		}))
		writeXMLToken(encoder.EncodeToken(xml.CharData([]byte(canonicalURL(post)))))
		writeXMLToken(encoder.EncodeToken(xml.EndElement{Name: xml.Name{Local: "guid"}}))
		writeXMLText(encoder, "pubDate", rssDate(post.Date))
		writeXMLText(encoder, "description", body)
		writeXMLText(encoder, "content:encoded", body)
		writeXMLToken(encoder.EncodeToken(xml.EndElement{Name: xml.Name{Local: "item"}}))
	}
	writeXMLToken(encoder.EncodeToken(xml.EndElement{Name: xml.Name{Local: "channel"}}))
	writeXMLToken(encoder.EncodeToken(xml.EndElement{Name: xml.Name{Local: "rss"}}))
	writeXMLToken(encoder.Flush())
	return append(output.Bytes(), '\n')
}

func renderJSONFeed(posts []Post) []byte {
	posts = publishedPosts(posts)
	items := make([]jsonFeedItem, 0, len(posts))
	for _, post := range posts {
		body := renderMarkdown(post.Body)
		modified := post.Updated
		items = append(items, jsonFeedItem{
			ID:            canonicalURL(post),
			URL:           canonicalURL(post),
			Title:         post.Title,
			Summary:       post.Summary,
			ContentHTML:   body,
			DatePublished: feedDate(post.Date),
			DateModified:  optionalFeedDate(modified),
		})
	}
	feed := jsonFeed{
		Version:     "https://jsonfeed.org/version/1.1",
		Title:       "Siduri",
		HomePageURL: siteURL("/"),
		FeedURL:     siteURL("/feed.json"),
		Items:       items,
	}
	data, err := json.Marshal(feed)
	if err != nil {
		panic(fmt.Sprintf("site: encode JSON Feed: %v", err))
	}
	return append(data, '\n')
}

type jsonFeed struct {
	Version     string         `json:"version"`
	Title       string         `json:"title"`
	HomePageURL string         `json:"home_page_url"`
	FeedURL     string         `json:"feed_url"`
	Items       []jsonFeedItem `json:"items"`
}

type jsonFeedItem struct {
	ID            string `json:"id"`
	URL           string `json:"url"`
	Title         string `json:"title"`
	Summary       string `json:"summary,omitempty"`
	ContentHTML   string `json:"content_html"`
	DatePublished string `json:"date_published"`
	DateModified  string `json:"date_modified,omitempty"`
}

func renderSitemap(posts []Post) []byte {
	posts = publishedPosts(posts)
	var output bytes.Buffer
	encoder := xml.NewEncoder(&output)
	writeXMLToken(encoder.EncodeToken(xml.ProcInst{
		Target: "xml",
		Inst:   []byte(`version="1.0" encoding="UTF-8"`),
	}))
	writeXMLToken(encoder.EncodeToken(xml.StartElement{
		Name: xml.Name{Local: "urlset"},
		Attr: []xml.Attr{{
			Name:  xml.Name{Local: "xmlns"},
			Value: "http://www.sitemaps.org/schemas/sitemap/0.9",
		}},
	}))
	for _, path := range []string{"/", "/journal/", "/about/", "/contact/"} {
		writeXMLToken(encoder.EncodeToken(xml.StartElement{Name: xml.Name{Local: "url"}}))
		writeXMLText(encoder, "loc", siteURL(path))
		writeXMLToken(encoder.EncodeToken(xml.EndElement{Name: xml.Name{Local: "url"}}))
	}
	for _, post := range posts {
		writeXMLToken(encoder.EncodeToken(xml.StartElement{Name: xml.Name{Local: "url"}}))
		writeXMLText(encoder, "loc", canonicalURL(post))
		writeXMLText(encoder, "lastmod", lastModifiedDate(post))
		writeXMLToken(encoder.EncodeToken(xml.EndElement{Name: xml.Name{Local: "url"}}))
	}
	writeXMLToken(encoder.EncodeToken(xml.EndElement{Name: xml.Name{Local: "urlset"}}))
	writeXMLToken(encoder.Flush())
	return append(output.Bytes(), '\n')
}

func renderRobots() []byte {
	return []byte("User-agent: *\nAllow: /\nSitemap: " + siteURL("/sitemap.xml") + "\n")
}

func renderLLMs(posts []Post) []byte {
	posts = publishedPosts(posts)
	var output strings.Builder
	output.WriteString("# Siduri\n\n")
	output.WriteString("> A journal about shipping software with agents while keeping a person responsible for the result.\n\n")
	output.WriteString("## Pages\n\n")
	for _, path := range []string{"/", "/journal/", "/about/", "/contact/"} {
		output.WriteString("- ")
		output.WriteString(siteURL(path))
		output.WriteByte('\n')
	}
	output.WriteString("\n## Feeds\n\n")
	output.WriteString("- ")
	output.WriteString(siteURL("/feed.xml"))
	output.WriteByte('\n')
	output.WriteString("- ")
	output.WriteString(siteURL("/feed.json"))
	output.WriteString("\n\n## Posts\n\n")
	for _, post := range posts {
		fmt.Fprintf(&output, "- [%s](%s): %s\n", post.Title, canonicalURL(post), post.PlainSummary)
	}
	return []byte(output.String())
}

func writeXMLText(encoder *xml.Encoder, name, value string) {
	writeXMLToken(encoder.EncodeToken(xml.StartElement{Name: xml.Name{Local: name}}))
	writeXMLToken(encoder.EncodeToken(xml.CharData([]byte(value))))
	writeXMLToken(encoder.EncodeToken(xml.EndElement{Name: xml.Name{Local: name}}))
}

func writeXMLToken(err error) {
	if err != nil {
		panic(fmt.Sprintf("site: encode XML: %v", err))
	}
}

func feedDate(date string) string {
	parsed, err := time.Parse("2006-01-02", date)
	if err != nil {
		return date
	}
	return parsed.UTC().Format(time.RFC3339)
}

func optionalFeedDate(date string) string {
	if date == "" {
		return ""
	}
	return feedDate(date)
}

func lastModifiedDate(post Post) string {
	if post.Updated != "" {
		return post.Updated
	}
	return post.Date
}

func rssDate(date string) string {
	parsed, err := time.Parse("2006-01-02", date)
	if err != nil {
		return date
	}
	return parsed.UTC().Format(time.RFC1123Z)
}
