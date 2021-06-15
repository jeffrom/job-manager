package resource

import (
	"database/sql/driver"
	"encoding/json"
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

func (v *Version) Equals(other *Version) bool {
	if v == nil || other == nil {
		return (v == nil) == (other == nil)
	}

	return v.v == other.v
}

func (v Version) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"v%d"`, v.v)), nil
}

func (v *Version) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	n, err := strconv.ParseInt(s[1:], 10, 32)
	if err != nil {
		return err
	}
	v.v = int32(n)
	return nil
}

func (v *Version) Scan(value interface{}) error {
	if value == nil {
		*v = Version{}
		return nil
	}

	*v = Version{v: int32(value.(int64))}
	return nil
}

func (v *Version) Value() (driver.Value, error) {
	if v == nil {
		return 0, nil
	}
	return v.v, nil
}
