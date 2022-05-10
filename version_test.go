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
		assertInt32(t, "version.major", 18, ver.Major)
		assertInt32(t, "version.minor", 11, ver.Minor)
		assertInt32(t, "version.patch", 1, ver.Patch)
		assertInt32(t, "version.revision", 822, ver.Revision)
		assertStr(t, "version.suffix", "alpha", ver.Suffix)
	})
	t.Run("parse mmpr", func(t *testing.T) {
		raw := "5.10.0.1"
		var ver Version
		err := ver.Parse(raw)
		if err != nil {
			t.Error(err)
			return
		}
		assertInt32(t, "version.major", 5, ver.Major)
		assertInt32(t, "version.minor", 10, ver.Minor)
		assertInt32(t, "version.patch", 0, ver.Patch)
		assertInt32(t, "version.revision", 1, ver.Revision)
	})
	t.Run("parse mmp", func(t *testing.T) {
		raw := "1.148.217"
		var ver Version
		err := ver.Parse(raw)
		if err != nil {
			t.Error(err)
			return
		}
		assertInt32(t, "version.major", 1, ver.Major)
		assertInt32(t, "version.minor", 148, ver.Minor)
		assertInt32(t, "version.patch", 217, ver.Patch)
	})
	t.Run("parse mm", func(t *testing.T) {
		raw := "26.07"
		var ver Version
		err := ver.Parse(raw)
		if err != nil {
			t.Error(err)
			return
		}
		assertInt32(t, "version.major", 26, ver.Major)
		assertInt32(t, "version.minor", 7, ver.Minor)
	})
	t.Run("parse m", func(t *testing.T) {
		raw := "5"
		var ver Version
		err := ver.Parse(raw)
		if err != nil {
			t.Error(err)
			return
		}
		assertInt32(t, "version.major", 5, ver.Major)
	})
}
