package uaxpl

import (
	"strconv"
	"strings"

	"github.com/koykov/byteconv"
)

type Version struct {
	Major    int32
	Minor    int32
	Patch    int32
	Revision int32
	Suffix   string

	p bool
}

func (v *Version) Parse(s string) error {
	if v.p {
		return nil
	}
	if len(s) > 0 {
		for i := 0; i < 4; i++ {
			idx := strings.Index(s, ".")
			if idx != -1 || len(s) > 0 {
				if idx == -1 {
					idx = len(s)
				}
				c := s[:idx]
				d, err := strconv.Atoi(c)
				if err != nil {
					return err
				}
				switch i {
				case 0:
					v.Major = int32(d)
				case 1:
					v.Minor = int32(d)
				case 2:
					v.Patch = int32(d)
				case 3:
					v.Revision = int32(d)
				}
				if idx < len(s) {
					s = s[idx+1:]
				} else {
					s = s[:0]
				}
			}
		}
		if len(s) > 0 {
			idx := strings.IndexAny(s, " -")
			if idx != -1 {
				v.Suffix = s[idx+1:]
			}
		}
	}
	v.p = true
	return nil
}

func (v Version) Write(dst []byte) []byte {
	var ok bool
	if v.Major > 0 {
		dst = strconv.AppendInt(dst, int64(v.Major), 10)
		ok = true
	}
	if v.Minor > 0 || ok {
		dst = append(dst, '.')
		dst = strconv.AppendInt(dst, int64(v.Minor), 10)
		ok = true
	}
	if v.Patch > 0 || ok {
		dst = append(dst, '.')
		dst = strconv.AppendInt(dst, int64(v.Patch), 10)
		ok = true
	}
	if v.Revision > 0 || ok {
		dst = append(dst, '.')
		dst = strconv.AppendInt(dst, int64(v.Revision), 10)
		ok = true
	}
	if len(v.Suffix) > 0 {
		dst = append(dst, ' ')
		dst = append(dst, v.Suffix...)
	}
	if len(dst) == 0 {
		dst = append(dst, '0')
	}
	return dst
}

func (v Version) String() string {
	return byteconv.B2S(v.Write(nil))
}

func (v *Version) Reset() {
	v.Major, v.Minor, v.Patch, v.Revision, v.Suffix = 0, 0, 0, 0, ""
	v.p = false
}

func getMajor(ver string) string {
	p := strings.Index(ver, ".")
	if p == -1 {
		return ""
	}
	return ver[:p]
}
