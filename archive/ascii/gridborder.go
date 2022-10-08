package ascii

// BorderMap is the configuration of the characters used for the border of a grid.
//
// https://en.wikipedia.org/wiki/Box_Drawing
type BorderMap struct {
	Horizontal, Vertical, Cross                rune
	TopLeft, TopRight, BottomLeft, BottomRight rune
	HorizontalAndTop, HorizontalAndBottom      rune
	VerticalAndLeft, VerticalAndRight          rune
}

var DefaultBorderMap = BorderMap{
	Horizontal:          '─',
	Vertical:            '│',
	Cross:               '┼',
	TopLeft:             '┌',
	TopRight:            '┐',
	BottomLeft:          '└',
	BottomRight:         '┘',
	HorizontalAndTop:    '┬',
	HorizontalAndBottom: '┴',
	VerticalAndLeft:     '├',
	VerticalAndRight:    '┤',
}
