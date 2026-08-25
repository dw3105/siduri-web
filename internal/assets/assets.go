// Package assets contains the build-time asset transformations used by the
// site. It deliberately has no dependency on the site package so it can be
// tested independently of route registration.
package assets

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"html"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

const contentHashLength = 16

// Image is the processed representation of one source image. Bytes are the
// bytes that must be emitted; Source is retained for diagnostics and tests.
// The output is always a deterministic PNG because the standard library does
// not provide WebP or AVIF encoders.
type Image struct {
	Source       string
	OutputPath   string
	URL          string
	Bytes        []byte
	Width        int
	Height       int
	SourceFormat string
	Format       string
	Hash         string
}

// ProcessImage decodes source and re-encodes it with the standard library's
// deterministic, best-compression PNG encoder. The content hash is calculated
// from the processed bytes, never from source.
func ProcessImage(sourceName string, source []byte) (Image, error) {
	sourceName = normalizeSourceName(sourceName)
	if sourceName == "" {
		return Image{}, fmt.Errorf("image source name is empty")
	}
	if len(source) == 0 {
		return Image{}, fmt.Errorf("image %q is empty", sourceName)
	}

	config, format, err := image.DecodeConfig(bytes.NewReader(source))
	if err != nil {
		return Image{}, fmt.Errorf("decode image %q: %w", sourceName, err)
	}
	if config.Width <= 0 || config.Height <= 0 {
		return Image{}, fmt.Errorf("image %q has invalid dimensions %dx%d", sourceName, config.Width, config.Height)
	}
	decoded, _, err := image.Decode(bytes.NewReader(source))
	if err != nil {
		return Image{}, fmt.Errorf("decode image %q: %w", sourceName, err)
	}

	var processed bytes.Buffer
	encoder := png.Encoder{CompressionLevel: png.BestCompression}
	if err := encoder.Encode(&processed, decoded); err != nil {
		return Image{}, fmt.Errorf("encode image %q as PNG: %w", sourceName, err)
	}
	data := append([]byte(nil), processed.Bytes()...)
	hash := sha256.Sum256(data)
	fullHash := hex.EncodeToString(hash[:])
	base := safeBaseName(sourceName)
	filename := fmt.Sprintf("%s-%s.png", base, fullHash[:contentHashLength])
	outputPath := path.Join("assets", filename)

	return Image{
		Source:       sourceName,
		OutputPath:   outputPath,
		URL:          "/" + outputPath,
		Bytes:        data,
		Width:        config.Width,
		Height:       config.Height,
		SourceFormat: format,
		Format:       "png",
		Hash:         fullHash,
	}, nil
}

// ProcessFile reads and processes one image file.
func ProcessFile(root, filename string) (Image, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return Image{}, fmt.Errorf("read image %s: %w", filename, err)
	}
	relative, err := filepath.Rel(root, filename)
	if err != nil {
		return Image{}, fmt.Errorf("relate image %s to %s: %w", filename, root, err)
	}
	return ProcessImage(filepath.ToSlash(relative), data)
}

// Discover processes every supported image below root/assets. Non-image
// source files are ignored, but image formats that this dependency-free
// pipeline cannot decode fail loudly so they can never be copied unprocessed.
func Discover(root string) ([]Image, error) {
	assetsRoot := filepath.Join(root, "assets")
	if _, err := os.Stat(assetsRoot); err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("stat asset directory %s: %w", assetsRoot, err)
	}

	var images []Image
	err := filepath.WalkDir(assetsRoot, func(filename string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			return nil
		}
		extension := strings.ToLower(filepath.Ext(entry.Name()))
		if !knownImageExtension(extension) {
			return nil
		}
		image, err := ProcessFile(root, filename)
		if err != nil {
			return err
		}
		images = append(images, image)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("discover assets: %w", err)
	}
	sort.Slice(images, func(i, j int) bool { return images[i].Source < images[j].Source })
	return images, nil
}

func knownImageExtension(extension string) bool {
	switch extension {
	case ".png", ".jpg", ".jpeg", ".gif":
		return true
	case ".avif", ".bmp", ".heic", ".tif", ".tiff", ".webp", ".svg":
		return true
	default:
		return false
	}
}

func normalizeSourceName(source string) string {
	source = strings.TrimSpace(strings.ReplaceAll(source, "\\", "/"))
	source = strings.TrimPrefix(source, "/")
	source = strings.TrimPrefix(source, "./")
	if source == "" {
		return ""
	}
	clean := path.Clean(source)
	if clean == "." || clean == ".." || strings.HasPrefix(clean, "../") {
		return ""
	}
	return clean
}

func safeBaseName(source string) string {
	base := strings.TrimSuffix(path.Base(source), path.Ext(source))
	var result strings.Builder
	for _, r := range base {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9':
			result.WriteRune(r)
		case r == '-', r == '_':
			result.WriteRune(r)
		default:
			result.WriteByte('-')
		}
	}
	if result.Len() == 0 {
		return "image"
	}
	return strings.Trim(result.String(), "-")
}

var markdownImagePattern = regexp.MustCompile(`!\[([^\]]*)\]\(\s*<?([^\s)>]+)>?(?:\s+["'][^)]*["'])?\s*\)`)

// RewriteMarkdownImages replaces local Markdown images with safe HTML image
// tags. The first image is treated as potentially above the fold; subsequent
// images receive loading="lazy". Every generated tag has dimensions from the
// decoded source and a content-hashed URL.
func RewriteMarkdownImages(markdown string, images map[string]Image) (string, error) {
	var rewriteErr error
	imageNumber := 0
	result := markdownImagePattern.ReplaceAllStringFunc(markdown, func(match string) string {
		if rewriteErr != nil {
			return match
		}
		parts := markdownImagePattern.FindStringSubmatch(match)
		alt := strings.TrimSpace(parts[1])
		source := normalizeSourceName(parts[2])
		if alt == "" {
			rewriteErr = fmt.Errorf("image %q has no alt text", parts[2])
			return match
		}
		if !strings.HasPrefix(source, "assets/") {
			rewriteErr = fmt.Errorf("image %q is not a local asset", parts[2])
			return match
		}
		processed, ok := images[source]
		if !ok {
			rewriteErr = fmt.Errorf("unprocessed image %q", parts[2])
			return match
		}
		imageNumber++
		return ImageTag(processed, alt, imageNumber > 1)
	})
	if rewriteErr != nil {
		return "", rewriteErr
	}
	return result, nil
}

// ImageTag renders an image with explicit dimensions and a processed,
// content-hashed source URL. Set lazy for images below the fold.
func ImageTag(image Image, alt string, lazy bool) string {
	var tag strings.Builder
	tag.WriteString(`<img src="`)
	tag.WriteString(html.EscapeString(image.URL))
	tag.WriteString(`" alt="`)
	tag.WriteString(html.EscapeString(alt))
	tag.WriteString(`" width="`)
	tag.WriteString(fmt.Sprintf("%d", image.Width))
	tag.WriteString(`" height="`)
	tag.WriteString(fmt.Sprintf("%d", image.Height))
	if lazy {
		tag.WriteString(`" loading="lazy`)
	}
	tag.WriteString(`" decoding="async" />`)
	return tag.String()
}
