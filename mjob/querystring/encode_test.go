package querystring

import (
	"fmt"
	"net/url"
	"reflect"
	"testing"
	"time"
)

type Nested struct {
	A   SubNested  `json:"a"`
	B   *SubNested `json:"b"`
	Ptr *SubNested `json:"ptr,omitempty"`
}

type SubNested struct {
	Value string `json:"value"`
}

func TestValues_types(t *testing.T) {
	str := "string"
	strPtr := &str
	timeVal := time.Date(2000, 1, 1, 12, 34, 56, 0, time.UTC)

	tests := []struct {
		in   interface{}
		want url.Values
	}{
		{
			// basic primitives
			struct {
				A string
				B int
				C uint
				D float32
				E bool
			}{},
			url.Values{
				"A": {""},
				"B": {"0"},
				"C": {"0"},
				"D": {"0"},
				"E": {"false"},
			},
		},
		{
			// pointers
			struct {
				A *string
				B *int
				C **string
				D *time.Time
			}{
				A: strPtr,
				C: &strPtr,
				D: &timeVal,
			},
			url.Values{
				"A": {str},
				"B": {""},
				"C": {str},
				"D": {"2000-01-01T12:34:56Z"},
			},
		},
		{
			// slices and arrays
			struct {
				A []string
				B []string `json:",comma"`
				C []string `json:",space"`
				D [2]string
				E [2]string `json:",comma"`
				F [2]string `json:",space"`
				G []*string `json:",space"`
				H []bool    `json:",int,space"`
				I []string  `json:",brackets"`
				J []string  `json:",semicolon"`
				K []string  `json:",numbered"`
			}{
				A: []string{"a", "b"},
				B: []string{"a", "b"},
				C: []string{"a", "b"},
				D: [2]string{"a", "b"},
				E: [2]string{"a", "b"},
				F: [2]string{"a", "b"},
				G: []*string{&str, &str},
				H: []bool{true, false},
				I: []string{"a", "b"},
				J: []string{"a", "b"},
				K: []string{"a", "b"},
			},
			url.Values{
				"A":   {"a", "b"},
				"B":   {"a,b"},
				"C":   {"a b"},
				"D":   {"a", "b"},
				"E":   {"a,b"},
				"F":   {"a b"},
				"G":   {"string string"},
				"H":   {"1 0"},
				"I[]": {"a", "b"},
				"J":   {"a;b"},
				"K0":  {"a"},
				"K1":  {"b"},
			},
		},
		{
			// other types
			struct {
				A time.Time
				B time.Time `json:",unix"`
				C bool      `json:",int"`
				D bool      `json:",int"`
			}{
				A: time.Date(2000, 1, 1, 12, 34, 56, 0, time.UTC),
				B: time.Date(2000, 1, 1, 12, 34, 56, 0, time.UTC),
				C: true,
				D: false,
			},
			url.Values{
				"A": {"2000-01-01T12:34:56Z"},
				"B": {"946730096"},
				"C": {"1"},
				"D": {"0"},
			},
		},
		{
			struct {
				Nest Nested `json:"nest"`
			}{
				Nested{
					A: SubNested{
						Value: "that",
					},
				},
			},
			url.Values{
				"nest[a][value]": {"that"},
				"nest[b]":        {""},
			},
		},
		{
			struct {
				Nest Nested `json:"nest"`
			}{
				Nested{
					Ptr: &SubNested{
						Value: "that",
					},
				},
			},
			url.Values{
				"nest[a][value]":   {""},
				"nest[b]":          {""},
				"nest[ptr][value]": {"that"},
			},
		},
		{
			nil,
			url.Values{},
		},
	}

	for i, tt := range tests {
		v, err := Values(tt.in)
		if err != nil {
			t.Errorf("%d. Values(%q) returned error: %v", i, tt.in, err)
		}

		if !reflect.DeepEqual(tt.want, v) {
			t.Errorf("%d. Values(%q) returned %v, want %v", i, tt.in, v, tt.want)
		}
	}
}

func TestValues_omitEmpty(t *testing.T) {
	str := ""
	s := struct {
		a string
		A string
		B string  `json:",omitempty"`
		C string  `json:"-"`
		D string  `json:"omitempty"` // actually named omitempty, not an option
		E *string `json:",omitempty"`
	}{E: &str}

	v, err := Values(s)
	if err != nil {
		t.Errorf("Values(%v) returned error: %v", s, err)
	}

	want := url.Values{
		"A":         {""},
		"omitempty": {""},
		"E":         {""}, // E is included because the pointer is not empty, even though the string being pointed to is
	}
	if !reflect.DeepEqual(want, v) {
		t.Errorf("Values(%v) returned %v, want %v", s, v, want)
	}
}

type A struct {
	B
}

type B struct {
	C string
}

type D struct {
	B
	C string
}

type e struct {
	B
	C string
}

type F struct {
	e
}

func TestValues_embeddedStructs(t *testing.T) {
	tests := []struct {
		in   interface{}
		want url.Values
	}{
		{
			A{B{C: "foo"}},
			url.Values{"C": {"foo"}},
		},
		{
			D{B: B{C: "bar"}, C: "foo"},
			url.Values{"C": {"foo", "bar"}},
		},
		{
			F{e{B: B{C: "bar"}, C: "foo"}}, // With unexported embed
			url.Values{"C": {"foo", "bar"}},
		},
	}

	for i, tt := range tests {
		v, err := Values(tt.in)
		if err != nil {
			t.Errorf("%d. Values(%q) returned error: %v", i, tt.in, err)
		}

		if !reflect.DeepEqual(tt.want, v) {
			t.Errorf("%d. Values(%q) returned %v, want %v", i, tt.in, v, tt.want)
		}
	}
}

func TestValues_invalidInput(t *testing.T) {
	_, err := Values("")
	if err == nil {
		t.Errorf("expected Values() to return an error on invalid input")
	}
}

type EncodedArgs []string

func (m EncodedArgs) EncodeValues(key string, v *url.Values) error {
	for i, arg := range m {
		v.Set(fmt.Sprintf("%s.%d", key, i), arg)
	}
	return nil
}

func TestValues_Marshaler(t *testing.T) {
	s := struct {
		Args EncodedArgs `json:"arg"`
	}{[]string{"a", "b", "c"}}
	v, err := Values(s)
	if err != nil {
		t.Errorf("Values(%q) returned error: %v", s, err)
	}

	want := url.Values{
		"arg.0": {"a"},
		"arg.1": {"b"},
		"arg.2": {"c"},
	}
	if !reflect.DeepEqual(want, v) {
		t.Errorf("Values(%q) returned %v, want %v", s, v, want)
	}
}

func TestValues_MarshalerWithNilPointer(t *testing.T) {
	s := struct {
		Args *EncodedArgs `json:"arg"`
	}{}
	v, err := Values(s)
	if err != nil {
		t.Errorf("Values(%v) returned error: %v", s, err)
	}

	want := url.Values{}
	if !reflect.DeepEqual(want, v) {
		t.Errorf("Values(%v) returned %v, want %v", s, v, want)
	}
}

func TestTagParsing(t *testing.T) {
	name, opts := parseTag("field,foobar,foo")
	if name != "field" {
		t.Fatalf("name = %q, want field", name)
	}
	for _, tt := range []struct {
		opt  string
		want bool
	}{
		{"foobar", true},
		{"foo", true},
		{"bar", false},
		{"field", false},
	} {
		if opts.Contains(tt.opt) != tt.want {
			t.Errorf("Contains(%q) = %v", tt.opt, !tt.want)
		}
	}
}
