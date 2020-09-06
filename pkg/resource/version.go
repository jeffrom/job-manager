package resource

import "strconv"

type Version struct {
	v string
}

func NewVersion(v int32) *Version {
	return &Version{
		v: strconv.FormatInt(int64(v), 10),
	}
}

func (v *Version) String() string {
	return "v" + string(v.v)
}
