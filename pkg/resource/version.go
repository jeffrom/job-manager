package resource

import "fmt"

type Version struct {
	v int32
}

func NewVersion(v int32) *Version {
	return &Version{
		// v: strconv.FormatInt(int64(v), 10),
		v: v,
	}
}

func (v *Version) String() string { return fmt.Sprintf("v%d", v.v) }
func (v *Version) Strict() string { return fmt.Sprint(v.v) }
func (v *Version) Raw() int32     { return v.v }
func (v *Version) Inc()           { v.v++ }
