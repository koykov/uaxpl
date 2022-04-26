package uaxpl

import "testing"

func assertCVS(tb testing.TB, a, b ClientType) bool {
	if a != b {
		tb.Errorf("client version mismatch: need '%s' got '%s'", b, a)
		return false
	}
	return true
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
