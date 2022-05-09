package zbar

import (
	"github.com/tetratelabs/wazero/api"
)

type symbol struct {
	mem   api.Memory
	ptr32 uint32
}

func newSymbol(ptr uint32, mem api.Memory) symbol {
	return symbol{
		ptr32: ptr >> 2,
		mem:   mem,
	}
}

func (s *symbol) read(b []byte, offset uint32) (uint32, bool) {
	lenD, _ := s.mem.ReadUint32Le(ctx, s.ptr32+4)
	ptr, _ := s.mem.ReadUint32Le(ctx, s.ptr32+5)
	if offset+uint32(len(b)) < lenD {
		b2, _ := s.mem.Read(ctx, ptr+offset, uint32(len(b)))
		copy(b, b2)
		return uint32(len(b)), true
	}

	b2, _ := s.mem.Read(ctx, ptr+offset, lenD-offset)
	copy(b, b2)
	return lenD - offset, false
}

func (s symbol) next() (symbol, bool) {
	ptr, ok := s.mem.ReadUint32Le(ctx, s.ptr32+11)
	if !ok || ptr == 0 {
		return symbol{}, false
	}

	return newSymbol(ptr, s.mem), true
}

type img struct {
	ptr uint32
}
