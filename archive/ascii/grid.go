package ascii

import (
	"fmt"
	"strings"
)

type Grid struct {
	// template defines the format and spacing of the grid
	// where each character represents a cell in the grid
	// and each line represents a row in the grid.
	template []string

	// data is the data to be rendered on the grid
	data map[rune]string
}

func (g *Grid) templateSize() (int64, int64) {
	var width int
	for _, line := range g.template {
		if len(line) > width {
			width = len(line)
		}
	}
	return int64(width), int64(len(g.template))
}

// NewGrid creates a new Grid
func NewGrid(template []string, data map[rune]string) *Grid {
	return &Grid{
		template: template,
		data:     data,
	}
}

// Island is a contiguous area of the grid.
//
// For example, the following grid template has 3 islands:
//
//	tmpl := []string{
//		"abbb",
//	    "accc",
//	}
//
// ...Where a is an island, b is an island, and c is an island.
//
// For another example, the following grid template has 4 islands
// using the same aliases as the previous example:
//
//	tmpl := []string{
//		"abbb",
//	    "accc",
//	    "abbb",
//	}
//
// ...Where a is an island, c is an island, and b is 2 islands.
type Island struct {
	Name    rune
	Content string
	Borders *BorderMap

	Padding, Border, Margin                              *int64
	PaddingTop, PaddingRight, PaddingBottom, PaddingLeft *int64
	BorderTop, BorderRight, BorderBottom, BorderLeft     *int64
	MarginTop, MarginRight, MarginBottom, MarginLeft     *int64

	// TODO: add more options
	// - [ ] alignment
	// - [ ] text transform to apply to substr
}

type RenderEnvironment struct {
	Width, Height int64
	Islands       []*Island
}

func (g *Grid) ValidateEnv(env RenderEnvironment) error {
	tmplWidth, tmplHeight := g.templateSize()

	// Validate the Width of the environment.
	switch {
	case env.Width < 0:
	case env.Width%tmplWidth != 0:
	case env.Width < int64(tmplWidth):
		return fmt.Errorf("invalid width: %d", env.Width)
	}

	// Validate the Height of the environment.
	switch {
	case env.Height < 0:
	case env.Height%tmplHeight != 0:
	case env.Height < int64(tmplHeight):
		return fmt.Errorf("invalid height: %d", env.Height)
	}

	// Validate the Islands of the environment.
	var uniqueIslands int
	for _, island := range env.Islands {
		if _, ok := g.data[island.Name]; !ok {
			return fmt.Errorf("invalid island name: %q", island.Name)
		}
		uniqueIslands++
	}
	if uniqueIslands != len(g.data) {
		return fmt.Errorf("invalid number of islands: %d", uniqueIslands)
	}

	return nil
}

type renderData []struct {
	col, row int64
	content  string
}

func (g *Grid) calculate(env *RenderEnvironment) renderData {
	var data renderData
	// tmplWidth, tmplHeight := g.templateSize()
	return data
}

func (g *Grid) Render(result *[]string, opts *RenderEnvironment) {
	var col, row int64

	for _, insertion := range g.calculate(opts) {
		for insertion.row > row {
			(*result)[row] += strings.Repeat(" ", int(opts.Width-col))
			col = 0
			row++
		}

		for insertion.col > col {
			(*result)[row] += " "
			col++
		}

		(*result)[row] += insertion.content
		col += int64(len(insertion.content))
	}
}
