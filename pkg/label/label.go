// Package label contains function managing labels and selectors.
package label

import "strings"

type Labels map[string]string

func (l Labels) Equals(other Labels) bool {
	if len(l) != len(other) {
		return false
	}
	for k, v := range l {
		if other[k] != v {
			return false
		}
	}
	return true
}

func ParseStringArray(labels []string) (Labels, error) {
	l := make(Labels)
	for _, lbl := range labels {
		parts := strings.SplitN(lbl, "=", 2)
		key, val := parts[0], parts[1]
		l[key] = val
	}
	return l, nil
}
