package ascii

import "fmt"

func ExampleRender() {
	grid := NewGrid([]string{
		"ab",
		"ab",
	}, map[rune]string{
		'A': "Hello",
		'B': "World",
	})

	var result *[]string
	grid.Render(result, &RenderEnvironment{
		Width:  4,
		Height: 4,
	})

	for _, line := range *result {
		fmt.Println(line)
	}

	// Output:
	// HeWo
	// llrl
}
