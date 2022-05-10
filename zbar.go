package zbar

import (
	"context"
	"image"
	"sync"

	_ "embed"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/wasi"
)

var (
	//go:embed zbar.wasm
	zbarWasm     []byte
	zbarCompiled wazero.CompiledCode
	ctx          = context.Background()
	r            = wazero.NewRuntimeWithConfig(wazero.NewRuntimeConfigJIT())
	pool         = sync.Pool{New: func() any { return NewScanner() }}
)

func init() {
	must(wasi.InstantiateSnapshotPreview1(ctx, r))
	must(r.NewModuleBuilder("env").
		ExportFunction("emscripten_notify_memory_growth", func(int32) {}).
		Instantiate(ctx))
	zbarCompiled = must(r.CompileModule(ctx, zbarWasm))
}

func ReadAll(img image.Image) ([][]byte, error) {
	s := pool.Get().(*Scanner)
	defer pool.Put(s)

	return s.ReadAll(img)
}

func Read(img image.Image) (*Reader, error) {
	s := pool.Get().(*Scanner)
	defer pool.Put(s)

	return s.Reader(img)
}

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
