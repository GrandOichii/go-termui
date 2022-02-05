package termui

// cct: curses color text

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	nc "github.com/rthornton128/goncurses"
)

// cct format example:
// "${red-black}Hello, world!"

const (
	rawcctregex = `\$\{([\w|-]+)\}([^\$]*)`
)

var (
	cctregex = regexp.MustCompile(rawcctregex)

	colors = map[string]int16{
		"red":     nc.C_RED,
		"blue":    nc.C_BLUE,
		"green":   nc.C_GREEN,
		"black":   nc.C_BLACK,
		"yellow":  nc.C_YELLOW,
		"cyan":    nc.C_CYAN,
		"magenta": nc.C_MAGENTA,
		"white":   nc.C_WHITE,
		"gray":    245,
		"pink":    219,
		"orange":  202,
	}
	colorMap = map[string]nc.Char{}
)

type CCTMessage struct {
	strings []string
	colors  []nc.Char
}

func (m CCTMessage) pair(i int) (string, nc.Char) {
	return m.strings[i], m.colors[i]
}

func (m CCTMessage) pairCount() int {
	return len(m.strings)
}

func (m CCTMessage) Length() int {
	result := 0
	for _, s := range m.strings {
		result += len(s)
	}
	return result
}

func (m CCTMessage) Draw(win *nc.Window, y, x int, attrs ...nc.Char) {
	for i := 0; i < m.pairCount(); i++ {
		s, color := m.pair(i)
		put(win, y, x, s, append(attrs, color)...)
		x += len(s)
	}
}

func parseColors(colorPair string) (nc.Char, error) {
	originalColorPair := colorPair
	if !strings.ContainsRune(colorPair, '-') {
		colorPair += "-normal"
	}
	result, has := colorMap[colorPair]
	if has {
		// color pair already initialized
		return result, nil
	}
	// initializing color pair
	split := strings.Split(colorPair, "-")
	if len(split) != 2 {
		return 0, fmt.Errorf("termui - %v is not valid cct color pair", originalColorPair)
	}
	fg := split[0]
	fgres, has := colors[fg]
	if !has {
		if fg == "normal" {
			fgres = -1
		} else {
			return 0, fmt.Errorf("termui - can't recognize color %v in color pair %v", fg, colorPair)
		}
	}

	bg := split[1]
	bgres, has := colors[bg]
	if !has {
		if bg == "normal" {
			bgres = -1
		} else {
			return 0, fmt.Errorf("termui - can't recognize color %v in color pair %v", bg, colorPair)
		}
	}
	pairI := int16(len(colorMap) + 1)
	err := nc.InitPair(pairI, fgres, bgres)
	if err != nil {
		return 0, err
	}
	result = nc.ColorPair(pairI)
	colorMap[colorPair] = result
	return result, nil
}

func ToCCTMessage(line string) (*CCTMessage, error) {
	if !strings.HasPrefix(line, "$") {
		line = "${normal}" + line
	}
	result := CCTMessage{}
	matches := cctregex.FindAllStringSubmatch(line, -1)
	result.strings = make([]string, 0, len(matches))
	result.colors = make([]nc.Char, 0, len(matches))
	for _, match := range matches {
		colorPair := match[1]
		s := match[2]
		colorKey, err := parseColors(colorPair)
		if err != nil {
			return nil, err
		}
		result.strings = append(result.strings, s)
		result.colors = append(result.colors, colorKey)
	}
	os.WriteFile("file.txt", []byte(fmt.Sprintf("*%v*", result)), 0755)
	return &result, nil
}

func GetCCTs(lines []string) ([]*CCTMessage, error) {
	result := make([]*CCTMessage, 0, len(lines))
	for _, line := range lines {
		cct, err := ToCCTMessage(line)
		if err != nil {
			return nil, err
		}
		result = append(result, cct)
	}
	return result, nil
}

func Foo(message string) {
	matches := cctregex.FindAllStringSubmatch(message, -1)
	for _, match := range matches {
		fmt.Println("---")
		fmt.Println(match[1])
		fmt.Println(match[2])
		// for _, smatch := range match {
		// 	fmt.Println(smatch)
		// }
	}
	fmt.Println("---")
}

func init() {
	// add all remaining colors
	for i := 10; i < 250; i++ {
		colors[strconv.Itoa(i)] = int16(i)
	}
}

func initColors() {
	err := nc.StartColor()
	// nc.Flash()
	if err != nil {
		panic(err)
	}
	err = nc.UseDefaultColors()
	if err != nil {
		panic(err)
	}
}
