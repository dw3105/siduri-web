package site

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/a-h/templ"
	"github.com/dw3105/siduri-web/internal/assets"
)

// The font is deliberately named even while the binary is absent. A later
// operator can place a real, subsetted WOFF2 at this path without changing
// the document head or route code.
const fallbackFont = "siduri-subset.woff2"

//go:embed assets_a4.css
var criticalCSS []byte

func init() {
	RegisterHead(HeadFragment{
		Name: "a4-assets",
		Render: func(PageData) templ.Component {
			return templ.Raw(a4HeadHTML())
		},
	})
	RegisterContent("a4-assets", registerA4Assets)
}

func registerA4Assets(data PageData, routes *RouteSet) {
	root := a4ProjectRoot()
	images, err := assets.Discover(root)
	if err != nil {
		panic(fmt.Errorf("a4 assets: %w", err))
	}

	bySource := make(map[string]assets.Image, len(images))
	for _, image := range images {
		bySource[image.Source] = image
	}
	for _, post := range data.Posts {
		if _, err := assets.RewriteMarkdownImages(post.Body, bySource); err != nil {
			panic(fmt.Errorf("a4 assets: post %q: %w", post.Slug, err))
		}
	}

	fontFiles, err := a4FontFiles(root)
	if err != nil {
		panic(fmt.Errorf("a4 fonts: %w", err))
	}
	outputs := make([]Output, 0, len(images)+len(fontFiles))
	for _, image := range images {
		image := image
		outputs = append(outputs, ByteOutput(image.OutputPath, image.Bytes))
	}
	for _, font := range fontFiles {
		font := font
		outputs = append(outputs, ByteOutput(filepath.ToSlash(filepath.Join("fonts", font.name)), font.data))
	}
	routes.Register(Route{
		Name: "a4-processed-assets",
		Output: RouteOutput{Expand: func(PageData) []Output {
			return outputs
		}},
	})
}

type a4FontFile struct {
	name string
	data []byte
}

func a4FontFiles(root string) ([]a4FontFile, error) {
	directory := filepath.Join(root, "static", "fonts")
	entries, err := os.ReadDir(directory)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var files []a4FontFile
	for _, entry := range entries {
		if entry.IsDir() || strings.ToLower(filepath.Ext(entry.Name())) != ".woff2" {
			continue
		}
		data, err := os.ReadFile(filepath.Join(directory, entry.Name()))
		if err != nil {
			return nil, err
		}
		files = append(files, a4FontFile{name: entry.Name(), data: data})
	}
	sort.Slice(files, func(i, j int) bool { return files[i].name < files[j].name })
	if len(files) > 2 {
		return nil, fmt.Errorf("at most two WOFF2 files are allowed, found %d", len(files))
	}
	return files, nil
}

func a4FontPreloads() []string {
	files, err := a4FontFiles(a4ProjectRoot())
	if err != nil {
		panic(fmt.Errorf("a4 fonts: %w", err))
	}
	if len(files) == 0 {
		return []string{fallbackFont}
	}
	preloads := make([]string, 0, len(files))
	for _, file := range files {
		preloads = append(preloads, file.name)
	}
	return preloads
}

func a4HeadHTML() string {
	var head strings.Builder
	head.WriteString(`<style data-critical="a4">`)
	head.Write(criticalCSS)
	head.WriteString(`</style>`)
	for _, font := range a4FontPreloads() {
		head.WriteString(`<link rel="preload" href="/fonts/`)
		head.WriteString(font)
		head.WriteString(`" as="font" type="font/woff2" crossorigin>`)
	}
	return head.String()
}

// a4ProjectRoot finds the checkout when the build command is run from the
// repository root or from a Go package below it. Build's public registration
// seam intentionally carries no root parameter, so this is kept local and
// deterministic rather than changing PageData or shared build code.
func a4ProjectRoot() string {
	working, err := os.Getwd()
	if err != nil {
		return "."
	}
	for directory := working; ; directory = filepath.Dir(directory) {
		if isA4ProjectRoot(directory) {
			return directory
		}
		parent := filepath.Dir(directory)
		if parent == directory {
			break
		}
	}
	return working
}

func isA4ProjectRoot(directory string) bool {
	_, contentErr := os.Stat(filepath.Join(directory, "content", "posts"))
	_, tagsErr := os.Stat(filepath.Join(directory, "content", "tags.yml"))
	return contentErr == nil && tagsErr == nil
}
