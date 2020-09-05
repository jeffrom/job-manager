package label

import (
	"fmt"
	"regexp"
	"strings"
)

// var selectorRE = regexp.MustCompile(`^(?P<name>!?\w+) *(?P<operator>!=|=in|notin)? *(?P<value>\(?[\w ,]+\)?)?$`)
var selectorRE = regexp.MustCompile(`^(?P<name>!?\w+)\s*(?P<operator>notin|in|=|!=)?\s*\(?(?P<value>[\w ,]+)?\)?$`)
var selectorValueSplitRE = regexp.MustCompile(`, *`)

// type selectorValue struct {
// 	not   bool
// 	in    []string
// 	notin []string
// }

// type Selectors map[string]selectorValue

type Selectors struct {
	Names    []string
	NotNames []string
	In       map[string][]string
	NotIn    map[string][]string
}

func newSelectors() *Selectors {
	return &Selectors{
		In:    make(map[string][]string),
		NotIn: make(map[string][]string),
	}
}

func (s Selectors) String() string {
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

func (s Selectors) Match(labels Labels) bool {
	if len(s.Names) == 0 && len(s.NotNames) == 0 && len(s.In) == 0 && len(s.NotIn) == 0 {
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

		if ch == ',' && (i-prev) != 0 {
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
