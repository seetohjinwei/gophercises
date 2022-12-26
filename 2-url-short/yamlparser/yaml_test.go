package yaml

import (
	"reflect"
	"testing"

	"github.com/seetohjinwei/gophercises/2-url-short/url"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name string
		data string
		want []url.Url
	}{
		{
			name: "single link",
			data: `
- shorts: [a]
  full: full_a
`,
			want: []url.Url{
				{Short: "a", Full: "full_a"},
			},
		},
		{
			name: "multiple links",
			data: `
- shorts:
  - a
  full: full_a
- shorts:
  - b
  full: full_b
`,
			want: []url.Url{
				{Short: "a", Full: "full_a"},
				{Short: "b", Full: "full_b"},
			},
		},
		{
			name: "multiple shorts",
			data: `
- shorts:
  - a1
  - a2
  - a3
  full: full_a
`,
			want: []url.Url{
				{Short: "a1", Full: "full_a"},
				{Short: "a2", Full: "full_a"},
				{Short: "a3", Full: "full_a"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := Parse([]byte(test.data))
			assertEquals(t, got, test.want)
		})
	}
}

func assertEquals(t testing.TB, got, want []url.Url) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v want %v", got, want)
	}
}
