package uaxpl

import "testing"

func assertCVS(tb testing.TB, a, b ClientType) {
	if a != b {
		tb.Errorf("client version mismatch: need '%s' got '%s'", b, a)
	}
}

func assertStr(tb testing.TB, stage, a, b string) {
	if a != b {
		tb.Errorf("%s mismatch: need '%s' got '%s'", stage, b, a)
	}
}

func assertInt32(tb testing.TB, stage string, a, b int32) {
	if a != b {
		tb.Errorf("%s mismatch: need %d got %d", stage, b, a)
	}
}
