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
	names []string
	nots  []string
	in    map[string][]string
	notin map[string][]string
}

func newSelectors() *Selectors {
	return &Selectors{
		in:    make(map[string][]string),
		notin: make(map[string][]string),
	}
}

func (s Selectors) Match(labels Labels) bool {
	// first negative matches
	for name, value := range labels {
		if valIn(name, s.nots) {
			return false
		}
		if notInVals, ok := s.notin[name]; ok {
			if valIn(value, notInVals) {
				return false
			}
		}
	}

	if len(s.names) == 0 && len(s.in) == 0 {
		return true
	}

	// now require positive matches
	for name, value := range labels {
		if valIn(name, s.names) {
			return true
		}
		if inVals, ok := s.in[name]; ok {
			if valIn(value, inVals) {
				return true
			}
		}
	}
	return false
}

func ParseSelectors(s string) (*Selectors, error) {
	sel := newSelectors()
	statements := strings.Split(s, ",")
	for i, stmt := range statements {
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
				sel.nots = append(sel.nots, name[1:])
			} else {
				sel.names = append(sel.names, name)
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
			sel.in[name] = append(sel.in[name], vals...)
		case "notin", "!=":
			sel.notin[name] = append(sel.notin[name], vals...)
		default:
			return nil, fmt.Errorf("label: invalid operator %q (part %d): %q", operator, i, stmt)
		}
	}
	return sel, nil
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
