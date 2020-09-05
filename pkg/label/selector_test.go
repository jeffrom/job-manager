package label

import "testing"

func TestSelector(t *testing.T) {
	labels := Labels{
		"cool": "nice",
	}

	tcs := []struct {
		in    string
		match bool
	}{
		{
			in:    "cool",
			match: true,
		},
		{
			in:    "!naw",
			match: true,
		},
		{
			in:    "cool,!naw",
			match: true,
		},
		{
			in:    "!naw,cool",
			match: true,
		},
		{
			in:    "!cool",
			match: false,
		},
		{
			in:    "!naw,!cool",
			match: false,
		},
		{
			in:    "cool in (nice)",
			match: true,
		},
		{
			in:    "!naw,cool in (nice)",
			match: true,
		},
		{
			in:    "cool in (nice, hella)",
			match: true,
		},
		{
			in:    "cool = nice",
			match: true,
		},
		{
			in:    "cool= nice",
			match: true,
		},
		{
			in:    "cool =nice",
			match: true,
		},
		{
			in:    "cool=nice",
			match: true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.in, func(t *testing.T) {
			sel, err := ParseSelectors(tc.in)
			if err != nil {
				t.Fatal(err)
			}

			match := sel.Match(labels)
			if match != tc.match {
				t.Fatalf("expected match to be %v", tc.match)
			}
		})
	}
}
