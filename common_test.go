package uaxpl

import (
	"encoding/json"
	"os"
	"testing"
)

func testGetRWH(stage *browserDS) string {
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

func testLoadBrowserDS(filename string) ([]browserDS, error) {
	var ds []browserDS
	contents, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(contents, &ds); err != nil {
		return nil, err
	}
	return ds, nil
}

func assertCVS(tb testing.TB, a, b ClientType, suppress bool) bool {
	if a != b {
		if !suppress {
			tb.Errorf("client version mismatch: need '%s' got '%s'", b, a)
		}
		return false
	}
	return true
}

func assertStr(tb testing.TB, stage, a, b string, suppress bool) bool {
	if a != b {
		if !suppress {
			tb.Errorf("%s mismatch: need '%s' got '%s'", stage, b, a)
		}
		return false
	}
	return true
}

func assertVerStr(tb testing.TB, stage, a, b string, suppress bool) bool {
	if a != b && len(b) > 0 {
		if !suppress {
			tb.Errorf("%s mismatch: need '%s' got '%s'", stage, b, a)
		}
		return false
	}
	return true
}

func assertInt32(tb testing.TB, stage string, a, b int32, suppress bool) bool {
	if a != b {
		if !suppress {
			tb.Errorf("%s mismatch: need %d got %d", stage, b, a)
		}
		return false
	}
	return true
}
