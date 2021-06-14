// Package label contains function managing labels and selectors.
package label

import (
	"sort"
	"strings"
)

// Labels can be used to filter queues. Their format is:
// "KEY=VALUE[,KEY=VALUE...]"
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

func (l Labels) String() string {
	parts := make([]string, len(l))
	i := 0
	for k, v := range l {
		parts[i] = k + "=" + v
		i++
	}
	sort.Strings(parts)
	return strings.Join(parts, ",")
}

// ParseStringArray reads an array of formatted labels, such as commandline
// flags.
func ParseStringArray(args []string) (Labels, error) {
	l := make(Labels)
	for _, lbl := range args {
		parts := strings.SplitN(lbl, "=", 2)
		key, val := parts[0], parts[1]
		l[key] = val
	}
	return l, nil
}
