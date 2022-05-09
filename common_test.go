package uaxpl

import "testing"

func testGetRWH(stage *dsBrowser) string {
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

func assertVerStr(tb testing.TB, stage, a, b string) bool {
	if a != b && len(b) > 0 {
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
