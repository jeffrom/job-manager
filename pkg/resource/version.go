package resource

import (
	"fmt"
	"strconv"
)

type Version struct {
	v int32
}

func NewVersion(v int32) *Version {
	return &Version{
		// v: strconv.FormatInt(int64(v), 10),
		v: v,
	}
}

func NewVersionFromString(s string) (*Version, error) {
	if len(s) > 0 && s[0] == 'v' {
		s = s[1:]
	}
	if s == "" {
		return NewVersion(0), nil
	}

	v, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return nil, err
	}
	return &Version{v: int32(v)}, nil
}

func (v *Version) String() string { return fmt.Sprintf("v%d", v.v) }
func (v *Version) Strict() string { return fmt.Sprint(v.v) }
func (v *Version) Raw() int32     { return v.v }
func (v *Version) Inc()           { v.v++ }
