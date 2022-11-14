package markdown

import (
	"fmt"
	"strings"
	"testing"
)

func TestConvert(t *testing.T) {
	type test struct {
		in  string
		out string
	}

	tests := []test{
		{
			in: trimLF(`
This is a suboptimal approach with room for insane improvement. 

### Test markdown

- [x] Hello world
- [ ] [Hello world](https://acmcsuf.com)
- [ ] Hello world: <https://acmcsuf.com>
- [ ] Hello!!
`),
			out: trimLF(`
This is a suboptimal approach with room for insane improvement.

### Test markdown

- ☑ Hello world
- ☐ [Hello world](https://acmcsuf.com)
- ☐ Hello world: https://acmcsuf.com
- ☐ Hello!!
`),
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprint(i+1), func(t *testing.T) {
			got := Convert(test.in, "")
			if got != test.out {
				t.Errorf("unexpected output (got/want):\n" +
					got + "\n" +
					"----------------\n" +
					test.out)
			}
		})
	}
}

func trimLF(s string) string {
	return strings.Trim(s, "\n")
}
