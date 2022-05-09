package uaxpl

import "testing"

func TestVersion(t *testing.T) {
	t.Run("parse mmprs", func(t *testing.T) {
		raw := "18.11.1.822.00 alpha"
		var ver Version
		err := ver.Parse(raw)
		if err != nil {
			t.Error(err)
			return
		}
		assertInt32(t, "version.major", 18, ver.Major, false)
		assertInt32(t, "version.minor", 11, ver.Minor, false)
		assertInt32(t, "version.patch", 1, ver.Patch, false)
		assertInt32(t, "version.revision", 822, ver.Revision, false)
		assertStr(t, "version.suffix", "alpha", ver.Suffix, false)
	})
	t.Run("parse mmpr", func(t *testing.T) {
		raw := "5.10.0.1"
		var ver Version
		err := ver.Parse(raw)
		if err != nil {
			t.Error(err)
			return
		}
		assertInt32(t, "version.major", 5, ver.Major, false)
		assertInt32(t, "version.minor", 10, ver.Minor, false)
		assertInt32(t, "version.patch", 0, ver.Patch, false)
		assertInt32(t, "version.revision", 1, ver.Revision, false)
	})
	t.Run("parse mmp", func(t *testing.T) {
		raw := "1.148.217"
		var ver Version
		err := ver.Parse(raw)
		if err != nil {
			t.Error(err)
			return
		}
		assertInt32(t, "version.major", 1, ver.Major, false)
		assertInt32(t, "version.minor", 148, ver.Minor, false)
		assertInt32(t, "version.patch", 217, ver.Patch, false)
	})
	t.Run("parse mm", func(t *testing.T) {
		raw := "26.07"
		var ver Version
		err := ver.Parse(raw)
		if err != nil {
			t.Error(err)
			return
		}
		assertInt32(t, "version.major", 26, ver.Major, false)
		assertInt32(t, "version.minor", 7, ver.Minor, false)
	})
	t.Run("parse m", func(t *testing.T) {
		raw := "5"
		var ver Version
		err := ver.Parse(raw)
		if err != nil {
			t.Error(err)
			return
		}
		assertInt32(t, "version.major", 5, ver.Major, false)
	})
}
