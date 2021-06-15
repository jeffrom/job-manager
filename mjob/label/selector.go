package label

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

var selectorRE = regexp.MustCompile(`^(?P<name>!?\w+)\s*(?P<operator>notin|in|=|!=)?\s*\(?(?P<value>[\w ,]+)?\)?$`)
var selectorValueSplitRE = regexp.MustCompile(`, *`)

// Selectors can be used to filter queues based on their labels.
type Selectors struct {
	Names    []string            `json:"names,omitempty"`
	NotNames []string            `json:"not_names,omitempty"`
	In       map[string][]string `json:"in,omitempty"`
	NotIn    map[string][]string `json:"not_in,omitempty"`
}

func newSelectors() *Selectors {
	return &Selectors{
		In:    make(map[string][]string),
		NotIn: make(map[string][]string),
	}
}

func (s *Selectors) Len() int {
	if s == nil {
		return 0
	}
	return len(s.Names) + len(s.NotNames) + len(s.In) + len(s.NotIn)
}

// CacheKey returns a string that can be used as a cache key for the selector
// values.
func (s *Selectors) CacheKey() string {
	var b strings.Builder
	if len(s.Names) > 0 {
		b.WriteString("names:")
		b.WriteString(strings.Join(s.Names, ","))
	}
	if len(s.NotNames) > 0 {
		b.WriteString("not:")
		b.WriteString(strings.Join(s.NotNames, ","))
	}

	if len(s.In) > 0 {
		inKeys := make([]string, len(s.In))
		i := 0
		for k := range s.In {
			inKeys[i] = k
			i++
		}
		sort.Strings(inKeys)

		b.WriteString("in:")
		for _, k := range inKeys {
			b.WriteString(strings.Join(s.In[k], ","))
		}
	}

	if len(s.NotIn) > 0 {
		keys := make([]string, len(s.NotIn))
		i := 0
		for k := range s.NotIn {
			keys[i] = k
			i++
		}
		sort.Strings(keys)

		b.WriteString("notin:")
		for _, k := range keys {
			b.WriteString(strings.Join(s.NotIn[k], ","))
		}
	}

	return b.String()
}

func (s *Selectors) String() string {
	var set bool
	var b strings.Builder
	b.WriteString("Selector<")
	if names := s.Names; len(names) > 0 {
		set = true
		b.WriteString("Names: ")
		b.WriteString(strings.Join(names, ", "))
	}
	if notNames := s.NotNames; len(notNames) > 0 {
		if set {
			b.WriteString(", ")
		}
		set = true
		b.WriteString("NotNames: ")
		b.WriteString(strings.Join(notNames, ", "))
	}
	if in := s.In; len(in) > 0 {
		if set {
			b.WriteString(", ")
		}
		set = true
		b.WriteString(fmt.Sprintf("In: %v", in))
	}
	if notin := s.NotIn; len(notin) > 0 {
		if set {
			b.WriteString(", ")
		}
		set = true
		b.WriteString(fmt.Sprintf("NotIn: %v", notin))
	}
	b.WriteString(">")
	return b.String()
}

func (s *Selectors) Match(labels Labels) bool {
	if s.Len() == 0 {
		return true
	}

	// first negative matches
	for name, value := range labels {
		if valIn(name, s.NotNames) {
			return false
		}
		if notInVals, ok := s.NotIn[name]; ok {
			if valIn(value, notInVals) {
				return false
			}
		}
	}

	if len(s.Names) == 0 && len(s.In) == 0 {
		return true
	}

	// now require positive matches
	for name, value := range labels {
		if valIn(name, s.Names) {
			return true
		}
		if inVals, ok := s.In[name]; ok {
			if valIn(value, inVals) {
				return true
			}
		}
	}
	return false
}

func ParseSelectorStringArray(sels []string) (*Selectors, error) {
	sel := newSelectors()

	for i, stmt := range sels {
		if strings.TrimSpace(stmt) == "" {
			continue
		}
		match := removeEmptyParts(selectorRE.FindStringSubmatch(strings.TrimSpace(stmt)))
		// fmt.Println(match)
		if match == nil || len(match) < 2 {
			return nil, fmt.Errorf("label: invalid selector (part %d): %q", i, stmt)
		}
		if len(match) > 4 {
			return nil, fmt.Errorf("label: invalid selector (part %d): %q", i, stmt)
		}

		name := match[1]
		if len(match) == 2 {
			not := len(name) > 0 && name[0] == '!'
			if not {
				sel.NotNames = append(sel.NotNames, name[1:])
			} else {
				sel.Names = append(sel.Names, name)
			}
			continue
		}

		operator, valStr := match[2], match[3]
		vals := trimPartSpaces(selectorValueSplitRE.Split(valStr, -1))
		if len(vals) == 0 {
			return nil, fmt.Errorf("label: invalid selector values (part %d): %q (%q)", i, stmt, valStr)
		}

		switch operator {
		case "in", "=":
			sel.In[name] = append(sel.In[name], vals...)
		case "notin", "!=":
			sel.NotIn[name] = append(sel.NotIn[name], vals...)
		default:
			return nil, fmt.Errorf("label: invalid operator %q (part %d): %q", operator, i, stmt)
		}
	}
	sort.Strings(sel.Names)
	sort.Strings(sel.NotNames)
	for _, v := range sel.In {
		sort.Strings(v)
	}
	for _, v := range sel.NotIn {
		sort.Strings(v)
	}
	return sel, nil
}

func ParseSelectors(s string) (*Selectors, error) {
	sels := SplitSelectors(s)
	// fmt.Printf("%q\n", sels)
	return ParseSelectorStringArray(sels)
}

func SplitSelectors(s string) []string {
	var parts []string
	prev := 0
	depth := 0
	for i, ch := range s {
		if ch == '(' {
			depth++
		}
		if depth > 0 {
			if ch == ')' {
				depth--
			}
			continue
		}

		if ch == ',' && (i-prev) != 0 { // nosemgrep: dgryski.semgrep-go.oddcompare.odd-comparison
			parts = append(parts, s[prev:i])
			prev = i + 1
		}
	}
	if prev < len(s) {
		parts = append(parts, s[prev:])
	}
	return parts
}

func removeEmptyParts(parts []string) []string {
	if len(parts) == 0 {
		return parts
	}

	var next []string
	for _, part := range parts {
		if strings.TrimSpace(part) == "" {
			continue
		}
		next = append(next, part)
	}

	return next
}

func trimPartSpaces(parts []string) []string {
	for i, part := range parts {
		parts[i] = strings.TrimSpace(part)
	}

	return parts
}

func valIn(val string, vals []string) bool {
	for _, v := range vals {
		if val == v {
			return true
		}
	}
	return false
}
