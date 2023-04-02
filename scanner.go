package zbar

import (
	"errors"
	"image"
	"io"
	"runtime"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"golang.org/x/exp/slices"
)

type Scanner struct {
	mod     api.Module
	rt      wazero.Runtime
	scanner uint32
}

type ScannerConfig struct {
	Config Config
	Value  int32
}

var (
	ErrNoSymbolsFound    = errors.New("zbar: no symbols were found")
	ErrFailedToSetConfig = errors.New("zbar: failed to set config")
)

func NewScanner() *Scanner {
	return must(NewScannerWithConfig(map[SymbolType][]ScannerConfig{
		SymbolType_QRCODE: {
			ScannerConfig{Config_ENABLE, 1},
			ScannerConfig{Config_ADD_CHECK, 1},
			ScannerConfig{Config_EMIT_CHECK, 0},
			ScannerConfig{Config_BINARY, 1},
		},
	}))
}

func NewScannerWithConfig(cfg map[SymbolType][]ScannerConfig) (*Scanner, error) {

	rt, zbar := newZbarInstance()

	res := must(zbar.ExportedFunction("ImageScanner_create").
		Call(ctx))

	s := Scanner{
		mod:     zbar,
		rt:      rt,
		scanner: uint32(res[0]),
	}

	defer runtime.SetFinalizer(&s, (*Scanner).destroy)

	var err error

	for t, cfgs := range cfg {
		for _, cfg := range cfgs {
			err = s.SetConfig(t, cfg)
			if err != nil {
				return nil, err
			}
		}
	}

	return &s, nil
}

func (s *Scanner) SetConfig(t SymbolType, cfg ScannerConfig) error {
	res, err := s.mod.ExportedFunction("ImageScanner_set_config").
		Call(ctx, uint64(s.scanner), uint64(t), uint64(cfg.Config), uint64(cfg.Value))

	if err != nil {
		return err
	}

	if res[0] != 0 {
		return ErrFailedToSetConfig
	}

	return nil
}

func (s *Scanner) ReadAll(img image.Image) ([][]byte, error) {
	zbarImg := s.createImage(img)

	err := s.scan(zbarImg)
	if err != nil {
		return nil, err
	}

	var data [][]byte
	symbol, next := s.getSymbols(zbarImg)

	for next {
		data = append(data, slices.Clone(symbol.readAll()))
		symbol, next = symbol.next()
	}

	s.mod.ExportedFunction("Image_destroy").
		Call(ctx, uint64(zbarImg.ptr))

	return data, nil

}

func (s *Scanner) Reader(img image.Image) (*Reader, error) {
	zbarImg := s.createImage(img)

	err := s.scan(zbarImg)
	if err != nil {
		return nil, err
	}

	symbol, ok := s.getSymbols(zbarImg)

	if !ok {
		return nil, ErrNoSymbolsFound
	}

	return &Reader{
		s:   symbol,
		mod: s.mod,
		img: zbarImg,
	}, nil

}

func (s *Scanner) scan(i img) error {
	res, err := s.mod.ExportedFunction("ImageScanner_scan").
		Call(ctx, uint64(s.scanner), uint64(i.ptr))
	if err != nil {
		return err
	}

	if int32(res[0]) < 0 {
		return errors.New("zbar: an unexpected error has happened")
	}

	return nil
}

func (s *Scanner) createImage(i image.Image) img {

	bounds := i.Bounds()

	buf := s.malloc(uint32(bounds.Dx() * bounds.Dy()))

	switch i := i.(type) {
	case *image.Gray:
		s.mod.Memory().Write(buf, i.Pix)
	default:
		s.mod.Memory().Write(buf, toGray(i).Pix)
	}

	var zbarImg img

	res := must(s.mod.ExportedFunction("Image_create").
		Call(ctx,
			uint64(bounds.Dx()), // width
			uint64(bounds.Dy()), // height
			0x30303859,          // format: "Y800"
			uint64(buf),
			uint64(bounds.Dx()*bounds.Dy()),
			0,
		))

	zbarImg.ptr = uint32(res[0])

	return zbarImg
}

func (s *Scanner) malloc(size uint32) uint32 {
	res := must(s.mod.ExportedFunction("malloc").
		Call(ctx, uint64(size)))
	return uint32(res[0])
}

func (s *Scanner) getSymbols(i img) (symbol, bool) {
	res := must(s.mod.ExportedFunction("Image_get_symbols").
		Call(ctx, uint64(i.ptr)))

	if res[0] == 0 {
		return symbol{}, false
	}

	res = must(s.mod.ExportedFunction("SymbolSet_get_first").
		Call(ctx, res[0]))
	return newSymbol(uint32(res[0]), s.mod), true
}

func (s *Scanner) destroy() {
	s.mod.ExportedFunction("ImageScanner_destory").
		Call(ctx, uint64(s.scanner))
	s.rt.Close(ctx)
}

// panics if img is nil
func toGray(img image.Image) *image.Gray {
	var (
		b = img.Bounds()
		g = image.NewGray(b)
	)

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			g.Set(x, y, img.At(x, y))
		}
	}
	return g
}

var _ io.ReadCloser = &Reader{}

type Reader struct {
	mod           api.Module
	s             symbol
	offset        uint32
	img           img
	firstReadDone bool
}

// Read reads a symbol into b
//
// Make sure not to Call (*Reader).Next() before
// receiving io.EOF. Otherwise, the remaining data
// will be skipped.
// The reader may return
//
// len(b) MUST NOT be larger than math.MaxUint32,
// otherwise unexpected behavior may occur.
func (r *Reader) Read(b []byte) (int, error) {
	if r.s.ptr == 0 {
		return 0, io.EOF
	}

	n, hasData := r.s.read(b, r.offset)
	r.offset += n
	if hasData {
		return int(n), nil
	}

	return int(n), io.EOF
}

func (r *Reader) Next() (ok bool) {
	if !r.firstReadDone {
		r.firstReadDone = true
		return true
	}
	r.s, ok = r.s.next()
	return ok
}

func (r *Reader) Close() error {
	if r.img.ptr == 0 {
		return nil
	}
	// zbar will call free on the buffer
	// module.c:55:44
	r.mod.ExportedFunction("Image_destroy").
		Call(ctx, uint64(r.img.ptr))
	r.img.ptr = 0
	return nil
}
