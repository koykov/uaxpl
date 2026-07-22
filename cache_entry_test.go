package uaxpl

import (
	"bytes"
	"io"
	"testing"
	"time"

	"github.com/koykov/bitset"
	"github.com/koykov/entry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCacheEntry_MarshalUnmarshal(t *testing.T) {
	tests := []struct {
		name  string
		entry CacheEntry
	}{
		{
			name: "full data",
			entry: CacheEntry{
				Bitset:          bitset.Bitset(0b10101010),
				ClientType:      ClientType(1),
				ClientName64:    entry.Entry64(123456789),
				ClientVersion64: entry.Entry64(987654321),
				EngineName64:    entry.Entry64(111111111),
				EngineVersion64: entry.Entry64(222222222),
				DeviceType:      DeviceType(42),
				BrandName64:     entry.Entry64(333333333),
				ModelName64:     entry.Entry64(444444444),
				OSName64:        entry.Entry64(555555555),
				OSVersion64:     entry.Entry64(666666666),
				Data:            []byte("test data for cache entry"),
				Hkey:            uint64(777777777),
				Timestamp:       time.Now().Unix(),
			},
		},
		{
			name: "empty data",
			entry: CacheEntry{
				Bitset:          bitset.Bitset(0),
				ClientType:      ClientType(0),
				ClientName64:    entry.Entry64(0),
				ClientVersion64: entry.Entry64(0),
				EngineName64:    entry.Entry64(0),
				EngineVersion64: entry.Entry64(0),
				DeviceType:      DeviceType(0),
				BrandName64:     entry.Entry64(0),
				ModelName64:     entry.Entry64(0),
				OSName64:        entry.Entry64(0),
				OSVersion64:     entry.Entry64(0),
				Hkey:            uint64(0),
				Timestamp:       0,
			},
		},
		{
			name: "large data",
			entry: CacheEntry{
				Bitset:          bitset.Bitset(0xFFFFFFFFFFFFFFFF),
				ClientType:      ClientType(255),
				ClientName64:    entry.Entry64(0xFFFFFFFFFFFFFFFF),
				ClientVersion64: entry.Entry64(0xFFFFFFFFFFFFFFFF),
				EngineName64:    entry.Entry64(0xFFFFFFFFFFFFFFFF),
				EngineVersion64: entry.Entry64(0xFFFFFFFFFFFFFFFF),
				DeviceType:      DeviceType(0xFFFF),
				BrandName64:     entry.Entry64(0xFFFFFFFFFFFFFFFF),
				ModelName64:     entry.Entry64(0xFFFFFFFFFFFFFFFF),
				OSName64:        entry.Entry64(0xFFFFFFFFFFFFFFFF),
				OSVersion64:     entry.Entry64(0xFFFFFFFFFFFFFFFF),
				Data:            bytes.Repeat([]byte("x"), 1024),
				Hkey:            uint64(0xFFFFFFFFFFFFFFFF),
				Timestamp:       0x7FFFFFFFFFFFFFFF,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			size := tt.entry.Size()
			expsz := tt.entry.minsz() + len(tt.entry.Data)
			assert.Equal(t, expsz, size, "Size() returns wrong value")

			buf := make([]byte, size)
			n, err := tt.entry.MarshalTo(buf)
			require.NoError(t, err, "MarshalTo failed")
			assert.Equal(t, size, n, "MarshalTo wrote wrong number of bytes")

			var r CacheEntry
			err = r.Unmarshal(buf)
			require.NoError(t, err, "Unmarshal failed")

			assert.Equal(t, tt.entry.Bitset, r.Bitset, "Bitset mismatch")
			assert.Equal(t, tt.entry.ClientType, r.ClientType, "ClientType mismatch")
			assert.Equal(t, tt.entry.ClientName64, r.ClientName64, "ClientName64 mismatch")
			assert.Equal(t, tt.entry.ClientVersion64, r.ClientVersion64, "ClientVersion64 mismatch")
			assert.Equal(t, tt.entry.EngineName64, r.EngineName64, "EngineName64 mismatch")
			assert.Equal(t, tt.entry.EngineVersion64, r.EngineVersion64, "EngineVersion64 mismatch")
			assert.Equal(t, tt.entry.DeviceType, r.DeviceType, "DeviceType mismatch")
			assert.Equal(t, tt.entry.BrandName64, r.BrandName64, "BrandName64 mismatch")
			assert.Equal(t, tt.entry.ModelName64, r.ModelName64, "ModelName64 mismatch")
			assert.Equal(t, tt.entry.OSName64, r.OSName64, "OSName64 mismatch")
			assert.Equal(t, tt.entry.OSVersion64, r.OSVersion64, "OSVersion64 mismatch")
			assert.Equal(t, tt.entry.Hkey, r.Hkey, "Hkey mismatch")
			assert.Equal(t, tt.entry.Timestamp, r.Timestamp, "Timestamp mismatch")
			assert.Equal(t, tt.entry.Data, r.Data, "Data mismatch")
		})
	}
}

func TestCacheEntry_MarshalUnmarshal_RoundTrip(t *testing.T) {
	original := &CacheEntry{
		Bitset:          bitset.Bitset(0b11110000),
		ClientType:      ClientType(5),
		ClientName64:    entry.Entry64(12345),
		ClientVersion64: entry.Entry64(67890),
		EngineName64:    entry.Entry64(11111),
		EngineVersion64: entry.Entry64(22222),
		DeviceType:      DeviceType(100),
		BrandName64:     entry.Entry64(33333),
		ModelName64:     entry.Entry64(44444),
		OSName64:        entry.Entry64(55555),
		OSVersion64:     entry.Entry64(66666),
		Data:            []byte("round trip test data"),
		Hkey:            uint64(99999),
		Timestamp:       time.Now().Unix(),
	}

	buf := make([]byte, original.Size())
	n, err := original.MarshalTo(buf)
	require.NoError(t, err)
	assert.Equal(t, original.Size(), n)

	var r CacheEntry
	err = r.Unmarshal(buf)
	require.NoError(t, err)

	assert.Equal(t, original.Bitset, r.Bitset)
	assert.Equal(t, original.ClientType, r.ClientType)
	assert.Equal(t, original.ClientName64, r.ClientName64)
	assert.Equal(t, original.ClientVersion64, r.ClientVersion64)
	assert.Equal(t, original.EngineName64, r.EngineName64)
	assert.Equal(t, original.EngineVersion64, r.EngineVersion64)
	assert.Equal(t, original.DeviceType, r.DeviceType)
	assert.Equal(t, original.BrandName64, r.BrandName64)
	assert.Equal(t, original.ModelName64, r.ModelName64)
	assert.Equal(t, original.OSName64, r.OSName64)
	assert.Equal(t, original.OSVersion64, r.OSVersion64)
	assert.Equal(t, original.Hkey, r.Hkey)
	assert.Equal(t, original.Timestamp, r.Timestamp)
	assert.Equal(t, original.Data, r.Data)
}

func TestCacheEntry_ErrorCases(t *testing.T) {
	t.Run("marshal short buffer", func(t *testing.T) {
		e := CacheEntry{
			Data: []byte("test"),
		}
		buf := make([]byte, e.Size()-1)
		n, err := e.MarshalTo(buf)
		assert.Equal(t, 0, n)
		assert.ErrorIs(t, err, io.ErrShortBuffer)
	})

	t.Run("unmarshal short data", func(t *testing.T) {
		e := CacheEntry{}
		data := make([]byte, e.minsz()-1)
		err := e.Unmarshal(data)
		assert.ErrorIs(t, err, io.ErrUnexpectedEOF)
	})

	t.Run("unmarshal incomplete data", func(t *testing.T) {
		e := CacheEntry{
			Data: []byte("test data that is longer than announced"),
		}
		buf := make([]byte, e.Size())
		n, _ := e.MarshalTo(buf)

		truncated := buf[:n-5]
		var r CacheEntry
		err := r.Unmarshal(truncated)
		assert.ErrorIs(t, err, io.ErrUnexpectedEOF)
	})
}

func TestCacheEntry_ZeroValue(t *testing.T) {
	var e CacheEntry
	assert.Equal(t, uint64(0), e.Hkey)
	assert.Equal(t, int64(0), e.Timestamp)
	assert.Nil(t, e.Data)
	assert.Equal(t, e.minsz(), e.Size())
}

func TestCacheEntry_FromCtx(t *testing.T) {
	ctx := &Ctx{
		Bitset:          bitset.Bitset(0x12345678),
		clientType:      ClientType(10),
		clientName64:    entry.Entry64(111),
		clientVersion64: entry.Entry64(222),
		engineName64:    entry.Entry64(333),
		engineVersion64: entry.Entry64(444),
		deviceType:      DeviceType(55),
		brandName64:     entry.Entry64(666),
		modelName64:     entry.Entry64(777),
		osName64:        entry.Entry64(888),
		osVersion64:     entry.Entry64(999),
		buf:             []byte("context data"),
	}

	var e CacheEntry
	e.FromCtx(ctx)

	assert.Equal(t, ctx.Bitset, e.Bitset)
	assert.Equal(t, ctx.clientType, e.ClientType)
	assert.Equal(t, ctx.clientName64, e.ClientName64)
	assert.Equal(t, ctx.clientVersion64, e.ClientVersion64)
	assert.Equal(t, ctx.engineName64, e.EngineName64)
	assert.Equal(t, ctx.engineVersion64, e.EngineVersion64)
	assert.Equal(t, ctx.deviceType, e.DeviceType)
	assert.Equal(t, ctx.brandName64, e.BrandName64)
	assert.Equal(t, ctx.modelName64, e.ModelName64)
	assert.Equal(t, ctx.osName64, e.OSName64)
	assert.Equal(t, ctx.osVersion64, e.OSVersion64)
	assert.Equal(t, ctx.buf, e.Data)
}

func BenchmarkCacheEntry_MarshalTo(b *testing.B) {
	e := CacheEntry{
		Bitset:          bitset.Bitset(0b10101010),
		ClientType:      ClientType(1),
		ClientName64:    entry.Entry64(123456789),
		ClientVersion64: entry.Entry64(987654321),
		EngineName64:    entry.Entry64(111111111),
		EngineVersion64: entry.Entry64(222222222),
		DeviceType:      DeviceType(42),
		BrandName64:     entry.Entry64(333333333),
		ModelName64:     entry.Entry64(444444444),
		OSName64:        entry.Entry64(555555555),
		OSVersion64:     entry.Entry64(666666666),
		Data:            bytes.Repeat([]byte("benchmark data"), 100),
		Hkey:            uint64(777777777),
		Timestamp:       time.Now().Unix(),
	}

	buf := make([]byte, e.Size())

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = e.MarshalTo(buf)
	}
}

func BenchmarkCacheEntry_Unmarshal(b *testing.B) {
	e := CacheEntry{
		Bitset:          bitset.Bitset(0b10101010),
		ClientType:      ClientType(1),
		ClientName64:    entry.Entry64(123456789),
		ClientVersion64: entry.Entry64(987654321),
		EngineName64:    entry.Entry64(111111111),
		EngineVersion64: entry.Entry64(222222222),
		DeviceType:      DeviceType(42),
		BrandName64:     entry.Entry64(333333333),
		ModelName64:     entry.Entry64(444444444),
		OSName64:        entry.Entry64(555555555),
		OSVersion64:     entry.Entry64(666666666),
		Data:            bytes.Repeat([]byte("benchmark data"), 100),
		Hkey:            uint64(777777777),
		Timestamp:       time.Now().Unix(),
	}

	buf := make([]byte, e.Size())
	_, _ = e.MarshalTo(buf)

	var r CacheEntry

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = r.Unmarshal(buf)
	}
}

func BenchmarkCacheEntry_Size(b *testing.B) {
	e := CacheEntry{
		Data: bytes.Repeat([]byte("benchmark data"), 100),
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e.Size()
	}
}
