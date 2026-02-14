package zbar

import (
	"github.com/tetratelabs/wazero/api"
)

type symbol struct {
	mod       api.Module
	len, data uint32
	ptr       uint32
}

func newSymbol(ptr uint32, zbar api.Module) symbol {
	s := symbol{
		ptr: ptr,
		mod: zbar,
	}

	res := must(s.mod.ExportedFunction("Symbol_get_data_length").
		Call(ctx, uint64(s.ptr)))
	s.len = uint32(res[0])
	res = must(s.mod.ExportedFunction("Symbol_get_data").
		Call(ctx, uint64(s.ptr)))
	s.data = uint32(res[0])

	return s
}

func (s *symbol) read(b []byte, offset uint32) (uint32, bool) {
	if offset+uint32(len(b)) < s.len {
		b2, _ := s.mod.Memory().Read(s.data+offset, uint32(len(b)))
		copy(b, b2)
		return uint32(len(b)), true
	}

	b2, _ := s.mod.Memory().Read(s.data+offset, s.len-offset)
	copy(b, b2)
	return s.len - offset, false
}

func (s *symbol) readAll() []byte {
	b, ok := s.mod.Memory().Read(s.data, s.len)
	if !ok {
		return nil
	}

	return b
}

func (s symbol) next() (symbol, bool) {
	res := must(s.mod.ExportedFunction("Symbol_next").
		Call(ctx, uint64(s.ptr)))

	if res[0] == 0 {
		return symbol{}, false
	}

	return newSymbol(uint32(res[0]), s.mod), true
}

type img struct {
	ptr uint32
}
