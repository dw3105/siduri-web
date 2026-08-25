package site

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Tool is the content contract for one entry under content/tools/<slug>/.
// Tool content deliberately lives outside the post model: tools are not
// journal entries and never enter feeds, tag pages, or the journal index.
type Tool struct {
	Title      string
	Slug       string
	Date       string
	Summary    string
	Why        string
	Language   string
	Status     string
	Repository string
	Install    string
	Example    string
	Screenshot string
}

var toolStatuses = []string{"active", "maintained", "abandoned"}

var (
	toolURLReferencePattern    = regexp.MustCompile(`(?i)/tools/([a-z0-9]+(?:-[a-z0-9]+)*)/?`)
	toolMarkerReferencePattern = regexp.MustCompile(`(?i)(?:\[\[|tool:|tools:)\s*([a-z0-9]+(?:-[a-z0-9]+)*)\s*(?:\]\])?`)
)

func init() {
	RegisterHead(HeadFragment{
		Name:   "tools-a6-css",
		Render: func(PageData) interfaceComponent { return toolsStylesheetHead() },
	})
	RegisterContent("tools-a6", registerToolsA6)
}

func registerToolsA6(data PageData, routes *RouteSet) {
	tools, err := loadTools(toolContentDir())
	if err != nil {
		registerToolBuildError(routes, err)
		return
	}

	if _, err := postsByTool(tools, data.Posts); err != nil {
		registerToolBuildError(routes, err)
		return
	}

	routes.Register(Route{
		Name: "tools-index-a6",
		Output: RouteOutput{Expand: func(buildData PageData) []Output {
			return []Output{PageOutput(
				"tools/index.html",
				func() interfaceComponent {
					return toolsIndexPage(tools, buildToolPosts(tools, buildData.Posts), buildData, "", "", renderHeadFragments(buildData))
				},
			)}
		}},
	})

	routes.Register(Route{
		Name: "tools-pages-a6",
		Output: RouteOutput{Expand: func(buildData PageData) []Output {
			buildPostsByTool, buildErr := postsByTool(tools, buildData.Posts)
			if buildErr != nil {
				panic(buildErr)
			}
			outputs := make([]Output, 0, len(tools))
			for _, tool := range tools {
				tool := tool
				mentioned := append([]Post(nil), buildPostsByTool[tool.Slug]...)
				outputs = append(outputs, PageOutput(
					"tools/"+tool.Slug+"/index.html",
					func() interfaceComponent {
						return toolPage(tool, mentioned, buildData, renderHeadFragments(buildData))
					},
				))
			}
			return outputs
		}},
	})

	routes.Register(Route{
		Name: "tools-filter-pages-a6",
		Output: RouteOutput{Expand: func(buildData PageData) []Output {
			return toolFilterOutputs(tools, buildData, buildToolPosts(tools, buildData.Posts))
		}},
	})

	routes.Register(Route{
		Name: "tools-stylesheet-a6",
		Output: RouteOutput{Expand: func(PageData) []Output {
			return []Output{ByteOutput("tools_a6.css", readToolStylesheet())}
		}},
	})
}

type toolBuildError struct {
	err error
}

func (failure toolBuildError) Render(context.Context, io.Writer) error {
	return failure.err
}

func registerToolBuildError(routes *RouteSet, err error) {
	routes.Register(Route{
		Name: "tools-a6-validation-error",
		Output: RouteOutput{Expand: func(PageData) []Output {
			return []Output{PageOutput("tools/index.html", func() interfaceComponent {
				return toolBuildError{err: err}
			})}
		}},
	})
}

func toolContentDir() string {
	if workingDir, err := os.Getwd(); err == nil {
		if dir := findToolsDir(workingDir); dir != "" {
			return dir
		}
	}
	// `go test` starts in internal/site, while `go run ./cmd/siduri` starts in
	// the project root. The source location is a stable checkout fallback for
	// the former without making the route depend on a package-global root.
	if _, source, _, ok := runtime.Caller(0); ok {
		if dir := findToolsDir(filepath.Dir(source)); dir != "" {
			return dir
		}
	}
	return filepath.Join("content", "tools")
}

func findToolsDir(start string) string {
	current, err := filepath.Abs(start)
	if err != nil {
		return ""
	}
	for {
		candidate := filepath.Join(current, "content", "tools")
		if info, statErr := os.Stat(candidate); statErr == nil && info.IsDir() {
			return candidate
		}
		parent := filepath.Dir(current)
		if parent == current {
			return ""
		}
		current = parent
	}
}

func readToolStylesheet() []byte {
	styles := []string{}
	if workingDir, err := os.Getwd(); err == nil {
		styles = append(styles, filepath.Join(filepath.Dir(findToolsDir(workingDir)), "..", "static", "tools_a6.css"))
	}
	if _, source, _, ok := runtime.Caller(0); ok {
		styles = append(styles, filepath.Join(filepath.Dir(source), "..", "..", "static", "tools_a6.css"))
	}
	for _, path := range styles {
		if data, err := os.ReadFile(filepath.Clean(path)); err == nil {
			return data
		}
	}
	panic("site: read static/tools_a6.css: file not found")
}

func loadTools(dir string) ([]Tool, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read tools directory %s: %w", dir, err)
	}
	tools := make([]Tool, 0, len(entries))
	seen := make(map[string]string)
	for _, entry := range entries {
		if !entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		tool, parseErr := parseToolDirectory(filepath.Join(dir, entry.Name()))
		if parseErr != nil {
			return nil, parseErr
		}
		if previous, exists := seen[tool.Slug]; exists {
			return nil, fmt.Errorf("%s: duplicate tool slug %q also used by %s", filepath.Join(dir, entry.Name()), tool.Slug, previous)
		}
		seen[tool.Slug] = entry.Name()
		tools = append(tools, tool)
	}
	if len(tools) == 0 {
		return nil, fmt.Errorf("%s: no tool directories found", dir)
	}
	sort.Slice(tools, func(i, j int) bool {
		if tools[i].Date != tools[j].Date {
			return tools[i].Date > tools[j].Date
		}
		return tools[i].Slug < tools[j].Slug
	})
	return tools, nil
}

func parseToolDirectory(dir string) (Tool, error) {
	var path string
	for _, name := range []string{"index.md", "tool.md"} {
		candidate := filepath.Join(dir, name)
		if _, err := os.Stat(candidate); err == nil {
			path = candidate
			break
		}
	}
	if path == "" {
		entries, err := os.ReadDir(dir)
		if err != nil {
			return Tool{}, fmt.Errorf("read tool directory %s: %w", dir, err)
		}
		var markdownFiles []string
		for _, entry := range entries {
			if !entry.IsDir() && filepath.Ext(entry.Name()) == ".md" {
				markdownFiles = append(markdownFiles, entry.Name())
			}
		}
		if len(markdownFiles) == 1 {
			path = filepath.Join(dir, markdownFiles[0])
		}
	}
	if path == "" {
		return Tool{}, fmt.Errorf("%s: expected index.md or tool.md", dir)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return Tool{}, fmt.Errorf("read tool %s: %w", path, err)
	}
	lines := strings.Split(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n")
	if len(lines) == 0 || strings.TrimSpace(lines[0]) != "---" {
		return Tool{}, fmt.Errorf("%s: frontmatter must start with ---", path)
	}
	end := -1
	for i := 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == "---" {
			end = i
			break
		}
	}
	if end < 0 {
		return Tool{}, fmt.Errorf("%s: frontmatter has no closing ---", path)
	}
	fields, err := parseToolFrontmatter(lines[1:end], path)
	if err != nil {
		return Tool{}, err
	}
	tool := Tool{
		Title:      fields["title"],
		Slug:       fields["slug"],
		Date:       fields["date"],
		Summary:    fields["summary"],
		Why:        fields["why"],
		Language:   fields["language"],
		Status:     fields["status"],
		Repository: fields["repository"],
		Install:    fields["install"],
		Example:    fields["example"],
		Screenshot: fields["screenshot"],
	}
	if tool.Why == "" {
		tool.Why = strings.TrimSpace(strings.Join(lines[end+1:], "\n"))
	}
	if err := validateTool(tool, path); err != nil {
		return Tool{}, err
	}
	if filepath.Base(dir) != tool.Slug {
		return Tool{}, fmt.Errorf("%s: directory name must match tool slug %q", path, tool.Slug)
	}
	return tool, nil
}

func parseToolFrontmatter(lines []string, path string) (map[string]string, error) {
	fields := make(map[string]string)
	known := map[string]bool{
		"title": true, "slug": true, "date": true, "summary": true, "why": true,
		"language": true, "status": true, "repository": true, "repo": true,
		"install": true, "example": true, "command": true, "screenshot": true,
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
		if value == "|" || value == ">" {
			var block []string
			for i++; i < len(lines); i++ {
				if strings.TrimSpace(lines[i]) == "" {
					block = append(block, "")
					continue
				}
				if lines[i][0] != ' ' && lines[i][0] != '\t' {
					i--
					break
				}
				block = append(block, strings.TrimSpace(lines[i]))
			}
			value = strings.TrimSpace(strings.Join(block, "\n"))
		}
		if key == "repo" {
			key = "repository"
		}
		if key == "command" {
			key = "example"
		}
		fields[key] = unquote(value)
	}
	return fields, nil
}

func validateTool(tool Tool, path string) error {
	for _, field := range []struct {
		name  string
		value string
	}{
		{"title", tool.Title}, {"slug", tool.Slug}, {"date", tool.Date},
		{"summary", tool.Summary}, {"language", tool.Language}, {"status", tool.Status},
		{"repository", tool.Repository}, {"install", tool.Install},
	} {
		if strings.TrimSpace(field.value) == "" {
			return fmt.Errorf("%s: missing required tool field %q", path, field.name)
		}
	}
	if tool.Example == "" && tool.Screenshot == "" {
		return fmt.Errorf("%s: tool needs an example or screenshot", path)
	}
	if !validSlug(tool.Slug) {
		return fmt.Errorf("%s: invalid tool slug %q", path, tool.Slug)
	}
	if _, err := time.Parse("2006-01-02", tool.Date); err != nil {
		return fmt.Errorf("%s: date must be YYYY-MM-DD: %w", path, err)
	}
	validStatus := false
	for _, status := range toolStatuses {
		if tool.Status == status {
			validStatus = true
			break
		}
	}
	if !validStatus {
		return fmt.Errorf("%s: unknown tool status %q", path, tool.Status)
	}
	return nil
}

func postsByTool(tools []Tool, posts []Post) (map[string][]Post, error) {
	known := make(map[string]bool, len(tools))
	for _, tool := range tools {
		known[tool.Slug] = true
	}
	result := make(map[string][]Post, len(tools))
	for _, post := range posts {
		for _, slug := range toolReferences(post.Body) {
			if !known[slug] {
				return nil, fmt.Errorf("post %q references unknown tool %q", post.Slug, slug)
			}
			if !containsString(postToolSlugs(result[slug]), post.Slug) {
				result[slug] = append(result[slug], post)
			}
		}
	}
	return result, nil
}

func postToolSlugs(posts []Post) []string {
	result := make([]string, 0, len(posts))
	for _, post := range posts {
		result = append(result, post.Slug)
	}
	return result
}

func buildToolPosts(tools []Tool, posts []Post) map[string][]Post {
	result, err := postsByTool(tools, posts)
	if err != nil {
		panic(err)
	}
	return result
}

func toolReferences(body string) []string {
	seen := make(map[string]bool)
	var result []string
	add := func(slug string) {
		if !seen[slug] {
			seen[slug] = true
			result = append(result, slug)
		}
	}
	for _, match := range toolURLReferencePattern.FindAllStringSubmatch(body, -1) {
		add(strings.ToLower(match[1]))
	}
	for _, match := range toolMarkerReferencePattern.FindAllStringSubmatch(body, -1) {
		add(strings.ToLower(match[1]))
	}
	sort.Strings(result)
	return result
}

func containsString(items []string, value string) bool {
	for _, item := range items {
		if item == value {
			return true
		}
	}
	return false
}

func toolFilterOutputs(tools []Tool, data PageData, postsByTool map[string][]Post) []Output {
	statuses := append([]string(nil), toolStatuses...)
	languages := toolLanguages(tools)
	outputs := []Output{
		PageOutput("tools/filters/all/index.html", func() interfaceComponent {
			return toolFilterFragment(tools, postsByTool, data, "", "")
		}),
	}
	for _, status := range statuses {
		status := status
		outputs = append(outputs,
			PageOutput("tools/status/"+status+"/index.html", func() interfaceComponent {
				return toolsIndexPage(tools, postsByTool, data, status, "", renderHeadFragments(data))
			}),
			PageOutput("tools/filters/status-"+status+"/index.html", func() interfaceComponent {
				return toolFilterFragment(tools, postsByTool, data, status, "")
			}),
		)
	}
	for _, language := range languages {
		language := language
		slug := filterPathPart(language)
		outputs = append(outputs,
			PageOutput("tools/language/"+slug+"/index.html", func() interfaceComponent {
				return toolsIndexPage(tools, postsByTool, data, "", language, renderHeadFragments(data))
			}),
			PageOutput("tools/filters/language-"+slug+"/index.html", func() interfaceComponent {
				return toolFilterFragment(tools, postsByTool, data, "", language)
			}),
		)
		for _, status := range statuses {
			status, language := status, language
			outputs = append(outputs,
				PageOutput("tools/status/"+status+"/language/"+slug+"/index.html", func() interfaceComponent {
					return toolsIndexPage(tools, postsByTool, data, status, language, renderHeadFragments(data))
				}),
				PageOutput("tools/filters/status-"+status+"-language-"+slug+"/index.html", func() interfaceComponent {
					return toolFilterFragment(tools, postsByTool, data, status, language)
				}),
			)
		}
	}
	return outputs
}

func toolLanguages(tools []Tool) []string {
	seen := make(map[string]bool)
	var languages []string
	for _, tool := range tools {
		language := strings.TrimSpace(tool.Language)
		key := strings.ToLower(language)
		if !seen[key] {
			seen[key] = true
			languages = append(languages, language)
		}
	}
	sort.Slice(languages, func(i, j int) bool {
		return strings.ToLower(languages[i]) < strings.ToLower(languages[j])
	})
	return languages
}

func filterPathPart(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	var out strings.Builder
	for _, r := range value {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			out.WriteRune(r)
		} else if out.Len() > 0 {
			out.WriteByte('-')
		}
	}
	return strings.Trim(out.String(), "-")
}

func filteredTools(tools []Tool, status, language string) []Tool {
	filtered := make([]Tool, 0, len(tools))
	for _, tool := range tools {
		if status != "" && tool.Status != status {
			continue
		}
		if language != "" && !strings.EqualFold(tool.Language, language) {
			continue
		}
		filtered = append(filtered, tool)
	}
	return filtered
}

// toolFilterLinkPaths is kept in Go so the no-JavaScript href and the htmx
// fragment endpoint cannot drift apart.
func toolFilterLinkPaths(status, language string) (string, string) {
	if status == "" && language == "" {
		return "/tools/", "/tools/filters/all/index.html"
	}
	full := "/tools/"
	fragment := "/tools/filters/"
	if status != "" {
		full += "status/" + status + "/"
		fragment += "status-" + status
	}
	if language != "" {
		if status != "" {
			full += "language/"
			fragment += "-language-"
		} else {
			full += "language/"
			fragment += "language-"
		}
		full += filterPathPart(language) + "/"
		fragment += filterPathPart(language)
	}
	full += "index.html"
	fragment += "/index.html"
	return full, fragment
}

func toolStatusLabel(status string) string {
	if status == "" {
		return "All"
	}
	return strings.ToUpper(status[:1]) + status[1:]
}

func toolCountString(count int) string {
	return strconv.Itoa(count)
}

func toolFilterHref(status, language string) string {
	href, _ := toolFilterLinkPaths(status, language)
	return href
}

func toolFilterEndpoint(status, language string) string {
	_, endpoint := toolFilterLinkPaths(status, language)
	return endpoint
}
