package site

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"strings"
)

const (
	ogWidth  = 1200
	ogHeight = 630
)

// OG cards are deliberately rendered with a tiny embedded bitmap alphabet.
// That keeps the build dependency-free and makes the bytes independent of the
// fonts installed on the machine running the build.
func init() {
	RegisterContent("a5-og", func(data PageData, routes *RouteSet) {
		routes.Register(Route{
			Name: "a5-og-images",
			Output: RouteOutput{Expand: func(buildData PageData) []Output {
				outputs := make([]Output, 0, len(buildData.Posts))
				for _, post := range buildData.Posts {
					post := post
					outputs = append(outputs, ByteOutput(
						"journal/"+post.Slug+"/og.png",
						renderOGImage(post),
					))
				}
				return outputs
			}},
		})
	})
}

func renderOGImage(post Post) []byte {
	canvas := image.NewRGBA(image.Rect(0, 0, ogWidth, ogHeight))
	background := color.RGBA{R: 246, G: 243, B: 236, A: 255}
	ink := color.RGBA{R: 29, G: 29, B: 27, A: 255}
	muted := color.RGBA{R: 110, G: 106, B: 96, A: 255}
	green := color.RGBA{R: 36, G: 92, B: 74, A: 255}
	draw.Draw(canvas, canvas.Bounds(), &image.Uniform{C: background}, image.Point{}, draw.Src)

	// The left rule and corner mark echo the site's restrained green links and
	// provide a stable visual anchor when a title is short.
	draw.Draw(canvas, image.Rect(72, 72, 80, ogHeight-72), &image.Uniform{C: green}, image.Point{}, draw.Src)
	draw.Draw(canvas, image.Rect(80, 72, 112, 80), &image.Uniform{C: green}, image.Point{}, draw.Src)

	drawOGText(canvas, 128, 88, "SIDURI / HUMAN IN THE LOOP", 4, muted)

	titleLines := wrapOGTitle(post.Title, 27)
	y := 195
	for _, line := range titleLines {
		drawOGText(canvas, 128, y, line, 8, ink)
		y += 72
	}

	if post.Series != "" {
		drawOGText(canvas, 128, 484, "SERIES / "+post.Series, 4, green)
	}
	drawOGText(canvas, 128, 540, post.Date, 4, muted)

	var encoded bytes.Buffer
	if err := png.Encode(&encoded, canvas); err != nil {
		// image/png's encoder only fails for an invalid writer. A bytes.Buffer is
		// valid, so returning nil would hide a programmer error less clearly.
		panic(err)
	}
	return encoded.Bytes()
}

func wrapOGTitle(title string, width int) []string {
	words := strings.Fields(strings.ToUpper(title))
	if len(words) == 0 {
		return []string{"UNTITLED"}
	}
	var lines []string
	line := words[0]
	for _, word := range words[1:] {
		if len(line)+1+len(word) <= width {
			line += " " + word
			continue
		}
		lines = append(lines, line)
		line = word
	}
	lines = append(lines, line)
	if len(lines) > 4 {
		lines = lines[:4]
		last := lines[3]
		if len(last) > width-3 {
			last = last[:width-3]
		}
		lines[3] = last + "..."
	}
	return lines
}

func drawOGText(dst *image.RGBA, x, y int, text string, scale int, ink color.Color) {
	for _, r := range strings.ToUpper(text) {
		if r == '—' || r == '–' {
			r = '-'
		}
		glyph, ok := ogGlyphs[r]
		if !ok {
			glyph = ogGlyphs['?']
		}
		for row, pattern := range glyph {
			for column, pixel := range pattern {
				if pixel != '#' {
					continue
				}
				x0 := x + column*scale
				y0 := y + row*scale
				draw.Draw(dst, image.Rect(x0, y0, x0+scale, y0+scale), &image.Uniform{C: ink}, image.Point{}, draw.Src)
			}
		}
		x += (len(glyph[0]) + 1) * scale
	}
}

var ogGlyphs = map[rune][]string{
	' ': {"...", "...", "...", "...", "...", "...", "..."},
	'-': {".....", ".....", ".....", "#####", ".....", ".....", "....."},
	'.': {".....", ".....", ".....", ".....", ".....", ".##..", ".##.."},
	'/': {"....#", "...#.", "...#.", "..#..", ".#...", ".#...", "#...."},
	':': {".....", ".##..", ".##..", ".....", ".##..", ".##..", "....."},
	'?': {"###..", "...#.", "..#..", ".#...", ".#...", ".....", ".#..."},
	'A': {".###.", "#...#", "#...#", "#####", "#...#", "#...#", "#...#"},
	'B': {"####.", "#...#", "#...#", "####.", "#...#", "#...#", "####."},
	'C': {".####", "#....", "#....", "#....", "#....", "#....", ".####"},
	'D': {"####.", "#...#", "#...#", "#...#", "#...#", "#...#", "####."},
	'E': {"#####", "#....", "#....", "####.", "#....", "#....", "#####"},
	'F': {"#####", "#....", "#....", "####.", "#....", "#....", "#...."},
	'G': {".####", "#....", "#....", "#.###", "#...#", "#...#", ".####"},
	'H': {"#...#", "#...#", "#...#", "#####", "#...#", "#...#", "#...#"},
	'I': {"#####", "..#..", "..#..", "..#..", "..#..", "..#..", "#####"},
	'J': {"..###", "...#.", "...#.", "...#.", "#..#.", "#..#.", ".##.."},
	'K': {"#...#", "#..#.", "#.#..", "##...", "#.#..", "#..#.", "#...#"},
	'L': {"#....", "#....", "#....", "#....", "#....", "#....", "#####"},
	'M': {"#...#", "##.##", "#.#.#", "#.#.#", "#...#", "#...#", "#...#"},
	'N': {"#...#", "##..#", "##..#", "#.#.#", "#..##", "#..##", "#...#"},
	'O': {".###.", "#...#", "#...#", "#...#", "#...#", "#...#", ".###."},
	'P': {"####.", "#...#", "#...#", "####.", "#....", "#....", "#...."},
	'Q': {".###.", "#...#", "#...#", "#...#", "#.#.#", "#..#.", ".##.#"},
	'R': {"####.", "#...#", "#...#", "####.", "#.#..", "#..#.", "#...#"},
	'S': {".####", "#....", "#....", ".###.", "....#", "....#", "####."},
	'T': {"#####", "..#..", "..#..", "..#..", "..#..", "..#..", "..#.."},
	'U': {"#...#", "#...#", "#...#", "#...#", "#...#", "#...#", ".###."},
	'V': {"#...#", "#...#", "#...#", "#...#", "#...#", ".#.#.", "..#.."},
	'W': {"#...#", "#...#", "#...#", "#.#.#", "#.#.#", "##.##", "#...#"},
	'X': {"#...#", "#...#", ".#.#.", "..#..", ".#.#.", "#...#", "#...#"},
	'Y': {"#...#", "#...#", ".#.#.", "..#..", "..#..", "..#..", "..#.."},
	'Z': {"#####", "....#", "...#.", "..#..", ".#...", "#....", "#####"},
	'0': {".###.", "#...#", "#..##", "#.#.#", "##..#", "#...#", ".###."},
	'1': {"..#..", ".##..", "..#..", "..#..", "..#..", "..#..", ".###."},
	'2': {".###.", "#...#", "....#", "...#.", "..#..", ".#...", "#####"},
	'3': {"####.", "....#", "....#", ".###.", "....#", "....#", "####."},
	'4': {"...#.", "..##.", ".#.#.", "#..#.", "#####", "...#.", "...#."},
	'5': {"#####", "#....", "#....", "####.", "....#", "....#", "####."},
	'6': {".###.", "#....", "#....", "####.", "#...#", "#...#", ".###."},
	'7': {"#####", "....#", "...#.", "..#..", ".#...", ".#...", ".#..."},
	'8': {".###.", "#...#", "#...#", ".###.", "#...#", "#...#", ".###."},
	'9': {".###.", "#...#", "#...#", ".####", "....#", "....#", ".###."},
}
