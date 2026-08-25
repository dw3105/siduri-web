package site

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

var requiredPostFields = []string{"title", "slug", "date", "summary", "plain_summary", "tags", "draft"}

// Post is the deliberately small content contract shared by the loader and templates.
type Post struct {
	Title        string
	Slug         string
	Date         string
	Updated      string
	Summary      string
	PlainSummary string
	Tags         []string
	Draft        bool
	Series       string
	Part         string
	Body         string
}

type Content struct {
	Posts []Post
	Tags  []string
}

func LoadContent(root string) (Content, error) {
	tags, err := loadTags(filepath.Join(root, "content", "tags.yml"))
	if err != nil {
		return Content{}, err
	}
	posts, err := loadPosts(filepath.Join(root, "content", "posts"), tags)
	if err != nil {
		return Content{}, err
	}
	return Content{Posts: posts, Tags: tags}, nil
}

func loadTags(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read tag vocabulary %s: %w", path, err)
	}
	var tags []string
	inList := false
	seen := make(map[string]bool)
	for _, raw := range strings.Split(string(data), "\n") {
		line := strings.TrimSpace(raw)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if line == "tags:" {
			inList = true
			continue
		}
		if !inList || !strings.HasPrefix(line, "-") {
			return nil, fmt.Errorf("%s: expected a tag list item", path)
		}
		tag := strings.TrimSpace(strings.TrimPrefix(line, "-"))
		tag = unquote(tag)
		if tag == "" || seen[tag] {
			return nil, fmt.Errorf("%s: empty or duplicate tag %q", path, tag)
		}
		seen[tag] = true
		tags = append(tags, tag)
	}
	if len(tags) == 0 {
		return nil, fmt.Errorf("%s: tag vocabulary is empty", path)
	}
	return tags, nil
}

func loadPosts(dir string, vocabulary []string) ([]Post, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read posts directory %s: %w", dir, err)
	}
	allowed := make(map[string]bool, len(vocabulary))
	for _, tag := range vocabulary {
		allowed[tag] = true
	}
	posts := make([]Post, 0, len(entries))
	slugs := make(map[string]string)
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".md" {
			continue
		}
		path := filepath.Join(dir, entry.Name())
		post, err := parsePost(path)
		if err != nil {
			return nil, err
		}
		if previous, exists := slugs[post.Slug]; exists {
			return nil, fmt.Errorf("%s: duplicate slug %q also used by %s", path, post.Slug, previous)
		}
		slugs[post.Slug] = path
		for _, tag := range post.Tags {
			if !allowed[tag] {
				return nil, fmt.Errorf("%s: unknown tag %q", path, tag)
			}
		}
		posts = append(posts, post)
	}
	sort.Slice(posts, func(i, j int) bool {
		if posts[i].Date != posts[j].Date {
			return posts[i].Date > posts[j].Date
		}
		return posts[i].Slug < posts[j].Slug
	})
	return posts, nil
}

func parsePost(path string) (Post, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Post{}, fmt.Errorf("read post %s: %w", path, err)
	}
	lines := strings.Split(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n")
	if len(lines) == 0 || strings.TrimSpace(lines[0]) != "---" {
		return Post{}, fmt.Errorf("%s: frontmatter must start with ---", path)
	}
	end := -1
	for i := 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == "---" {
			end = i
			break
		}
	}
	if end < 0 {
		return Post{}, fmt.Errorf("%s: frontmatter has no closing ---", path)
	}
	fields, err := parseFrontmatter(lines[1:end], path)
	if err != nil {
		return Post{}, err
	}
	for _, field := range requiredPostFields {
		if _, ok := fields[field]; !ok {
			return Post{}, fmt.Errorf("%s: missing required frontmatter field %q", path, field)
		}
	}
	post := Post{
		Title:        fields["title"],
		Slug:         fields["slug"],
		Date:         fields["date"],
		Updated:      fields["updated"],
		Summary:      fields["summary"],
		PlainSummary: fields["plain_summary"],
		Tags:         splitList(fields["tags"]),
		Series:       fields["series"],
		Part:         fields["part"],
		Body:         strings.TrimSpace(strings.Join(lines[end+1:], "\n")),
	}
	if post.Title == "" || post.Slug == "" || post.Summary == "" || post.PlainSummary == "" {
		return Post{}, fmt.Errorf("%s: title, slug, summary, and plain_summary must not be empty", path)
	}
	if !validSlug(post.Slug) {
		return Post{}, fmt.Errorf("%s: invalid slug %q", path, post.Slug)
	}
	if _, err := time.Parse("2006-01-02", post.Date); err != nil {
		return Post{}, fmt.Errorf("%s: date must be YYYY-MM-DD: %w", path, err)
	}
	if post.Updated != "" {
		if _, err := time.Parse("2006-01-02", post.Updated); err != nil {
			return Post{}, fmt.Errorf("%s: updated must be YYYY-MM-DD: %w", path, err)
		}
	}
	draft, err := strconv.ParseBool(fields["draft"])
	if err != nil {
		return Post{}, fmt.Errorf("%s: draft must be true or false", path)
	}
	post.Draft = draft
	if len(post.Tags) == 0 {
		return Post{}, fmt.Errorf("%s: tags must contain at least one tag", path)
	}
	return post, nil
}

func readingTime(body string) int {
	words := len(strings.Fields(body))
	minutes := (words + 199) / 200
	if minutes < 1 {
		return 1
	}
	return minutes
}

func parseFrontmatter(lines []string, path string) (map[string]string, error) {
	fields := make(map[string]string)
	known := map[string]bool{
		"title": true, "slug": true, "date": true, "summary": true,
		"plain_summary": true, "tags": true, "draft": true, "series": true, "part": true, "updated": true,
	}
	for i := 0; i < len(lines); i++ {
		trimmed := strings.TrimSpace(lines[i])
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		colon := strings.IndexByte(trimmed, ':')
		if colon <= 0 {
			return nil, fmt.Errorf("%s: malformed frontmatter line %q", path, trimmed)
		}
		key := strings.TrimSpace(trimmed[:colon])
		if !known[key] {
			return nil, fmt.Errorf("%s: unknown frontmatter field %q", path, key)
		}
		value := strings.TrimSpace(trimmed[colon+1:])
		if key == "tags" && value == "" {
			var items []string
			for i++; i < len(lines); i++ {
				item := strings.TrimSpace(lines[i])
				if item == "" {
					continue
				}
				if !strings.HasPrefix(item, "-") {
					i--
					break
				}
				items = append(items, unquote(strings.TrimSpace(strings.TrimPrefix(item, "-"))))
			}
			fields[key] = strings.Join(items, "\x1f")
			continue
		}
		fields[key] = unquote(value)
	}
	return fields, nil
}

func splitList(value string) []string {
	if strings.Contains(value, "\x1f") {
		return nonEmpty(strings.Split(value, "\x1f"))
	}
	value = strings.TrimSpace(value)
	if strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]") {
		value = strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(value, "["), "]"))
		if value == "" {
			return nil
		}
		parts := strings.Split(value, ",")
		for i := range parts {
			parts[i] = unquote(strings.TrimSpace(parts[i]))
		}
		return nonEmpty(parts)
	}
	if value == "" {
		return nil
	}
	return []string{unquote(value)}
}

func nonEmpty(items []string) []string {
	result := make([]string, 0, len(items))
	for _, item := range items {
		if item != "" {
			result = append(result, item)
		}
	}
	return result
}

func unquote(value string) string {
	if len(value) >= 2 && ((value[0] == '"' && value[len(value)-1] == '"') || (value[0] == '\'' && value[len(value)-1] == '\'')) {
		if value[0] == '"' {
			if unquoted, err := strconv.Unquote(value); err == nil {
				return unquoted
			}
		}
		return value[1 : len(value)-1]
	}
	return value
}

func validSlug(slug string) bool {
	if slug == "" || strings.HasPrefix(slug, "-") || strings.HasSuffix(slug, "-") {
		return false
	}
	for _, r := range slug {
		if (r < 'a' || r > 'z') && (r < '0' || r > '9') && r != '-' {
			return false
		}
	}
	return true
}
