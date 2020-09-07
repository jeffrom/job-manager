package label

import (
	"fmt"
	"strings"
)

// Claims provide consumers the ability to signal cache availability to
// job-manager.
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
	// sort.Strings(claims)
	return claims
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
