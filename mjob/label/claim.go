package label

import (
	"fmt"
	"sort"
	"strings"
)

// Claims provide a facility for prioritizing jobs execution onto consumers
// with matching claims. Queues can be configured with a claim window, a
// duration during which only consumers with matching claims can dequeue the
// job. After the duration elapses, dequeue behavior returns to normal.
type Claims map[string][]string

func (c Claims) Format() []string {
	if c == nil {
		return nil
	}
	var claims []string
	for k, c := range c {
		for _, v := range c {
			claims = append(claims, fmt.Sprintf("%s=%s", k, v))
		}
	}
	sort.Strings(claims)
	return claims
}

func (c Claims) Match(other Claims) bool {
	// fmt.Println("match start:", c, other)
	if c == nil || other == nil {
		return (c == nil) == (other == nil)
	}
	if len(c) != len(other) {
		return false
	}
	for k, v := range c {
		otherV := other[k]
		if len(v) != len(otherV) {
			return false
		}
		for i, claim := range v {
			otherClaim := otherV[i]
			if claim != otherClaim {
				return false
			}
		}
	}
	return true
}

func ParseClaims(claims []string) (Claims, error) {
	if len(claims) == 0 {
		return nil, nil
	}
	c := make(Claims)
	for _, cl := range claims {
		parts := strings.SplitN(cl, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("label: invalid claim: %q", cl)
		}
		c[parts[0]] = append(c[parts[0]], parts[1])
	}
	return c, nil
}
