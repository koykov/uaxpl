package uaxpl

import (
	"encoding/binary"
	"io"

	"github.com/koykov/bitset"
	"github.com/koykov/bytealg"
	"github.com/koykov/entry"
)

type CacheEntry struct {
	bitset.Bitset

	ClientType      ClientType
	ClientName64    entry.Entry64
	ClientVersion64 entry.Entry64

	EngineName64    entry.Entry64
	EngineVersion64 entry.Entry64

	DeviceType  DeviceType
	BrandName64 entry.Entry64
	ModelName64 entry.Entry64

	OSName64    entry.Entry64
	OSVersion64 entry.Entry64

	Data []byte

	Hkey      uint64
	Timestamp int64
}

func (r *CacheEntry) FromCtx(ctx *Ctx) {
	r.Bitset = ctx.Bitset
	r.ClientType = ctx.clientType
	r.ClientName64 = ctx.clientName64
	r.ClientVersion64 = ctx.clientVersion64
	r.EngineName64 = ctx.engineName64
	r.EngineVersion64 = ctx.engineVersion64
	r.DeviceType = ctx.deviceType
	r.BrandName64 = ctx.brandName64
	r.ModelName64 = ctx.modelName64
	r.OSName64 = ctx.osName64
	r.OSVersion64 = ctx.osVersion64
	r.Data = bytealg.Copy(ctx.src)
}

func (r *CacheEntry) Size() int {
	return r.minsz() + len(r.Data)
}

func (r *CacheEntry) minsz() int {
	return 8 + // Bitset
		1 + // ClientType
		8 + // ClientName64
		8 + // ClientVersion64
		8 + // EngineName64
		8 + // EngineVersion64
		2 + // DeviceType
		8 + // BrandName64
		8 + // ModelName64
		8 + // OSName64
		8 + // OSVersion64
		8 + // Hkey
		8 + // Timestamp
		4 // len(Data)
}

func (r *CacheEntry) MarshalTo(buf []byte) (int, error) {
	if len(buf) < r.Size() {
		return 0, io.ErrShortBuffer
	}
	var off int

	binary.LittleEndian.PutUint64(buf[off:off+8], uint64(r.Bitset))
	off += 8

	buf[off] = byte(r.ClientType)
	off++
	binary.LittleEndian.PutUint64(buf[off:off+8], uint64(r.ClientName64))
	off += 8
	binary.LittleEndian.PutUint64(buf[off:off+8], uint64(r.ClientVersion64))
	off += 8
	binary.LittleEndian.PutUint64(buf[off:off+8], uint64(r.EngineName64))
	off += 8
	binary.LittleEndian.PutUint64(buf[off:off+8], uint64(r.EngineVersion64))
	off += 8

	binary.LittleEndian.PutUint16(buf[off:off+2], uint16(r.DeviceType))
	off += 2
	binary.LittleEndian.PutUint64(buf[off:off+8], uint64(r.BrandName64))
	off += 8
	binary.LittleEndian.PutUint64(buf[off:off+8], uint64(r.ModelName64))
	off += 8
	binary.LittleEndian.PutUint64(buf[off:off+8], uint64(r.OSName64))
	off += 8
	binary.LittleEndian.PutUint64(buf[off:off+8], uint64(r.OSVersion64))
	off += 8

	binary.LittleEndian.PutUint64(buf[off:off+8], r.Hkey)
	off += 8
	binary.LittleEndian.PutUint64(buf[off:off+8], uint64(r.Timestamp))
	off += 8

	binary.LittleEndian.PutUint32(buf[off:off+4], uint32(len(r.Data)))
	off += 4
	copy(buf[off:], r.Data)
	off += len(r.Data)

	return off, nil
}

func (r *CacheEntry) Unmarshal(data []byte) error {
	if len(data) < r.minsz() {
		return io.ErrUnexpectedEOF
	}

	var off int
	r.Bitset = bitset.Bitset(binary.LittleEndian.Uint64(data[off : off+8]))
	off += 8

	r.ClientType = ClientType(data[off])
	off++
	r.ClientName64 = entry.Entry64(binary.LittleEndian.Uint64(data[off : off+8]))
	off += 8
	r.ClientVersion64 = entry.Entry64(binary.LittleEndian.Uint64(data[off : off+8]))
	off += 8
	r.EngineName64 = entry.Entry64(binary.LittleEndian.Uint64(data[off : off+8]))
	off += 8
	r.EngineVersion64 = entry.Entry64(binary.LittleEndian.Uint64(data[off : off+8]))
	off += 8

	r.DeviceType = DeviceType(binary.LittleEndian.Uint16(data[off : off+2]))
	off += 2
	r.BrandName64 = entry.Entry64(binary.LittleEndian.Uint64(data[off : off+8]))
	off += 8
	r.ModelName64 = entry.Entry64(binary.LittleEndian.Uint64(data[off : off+8]))
	off += 8
	r.OSName64 = entry.Entry64(binary.LittleEndian.Uint64(data[off : off+8]))
	off += 8
	r.OSVersion64 = entry.Entry64(binary.LittleEndian.Uint64(data[off : off+8]))
	off += 8

	r.Hkey = binary.LittleEndian.Uint64(data[off : off+8])
	off += 8
	r.Timestamp = int64(binary.LittleEndian.Uint64(data[off : off+8]))
	off += 8

	ln := binary.LittleEndian.Uint32(data[off : off+4])
	off += 4
	if len(data) < off+int(ln) {
		return io.ErrUnexpectedEOF
	}
	r.Data = append(r.Data[:0], data[off:off+int(ln)]...)

	return nil
}
