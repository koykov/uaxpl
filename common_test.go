package uaxpl

import (
	"testing"
)

func testGetRWH(stage *clientDS) string {
	if rwh, ok := stage.Headers["http-x-requested-with"]; ok {
		return rwh
	}
	if rwh, ok := stage.Headers["X-Requested-With"]; ok {
		return rwh
	}
	if rwh, ok := stage.Headers["x-requested-with"]; ok {
		return rwh
	}
	return ""
}

func assertStr(tb testing.TB, stage, a, b string) bool {
	if a != b {
		tb.Errorf("%s mismatch: need '%s' got '%s'", stage, b, a)
		return false
	}
	return true
}

func assertInt32(tb testing.TB, stage string, a, b int32) bool {
	if a != b {
		tb.Errorf("%s mismatch: need %d got %d", stage, b, a)
		return false
	}
	return true
}
