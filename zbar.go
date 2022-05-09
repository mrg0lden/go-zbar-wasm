package zbar

import (
	"context"

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
)

func init() {
	must(wasi.InstantiateSnapshotPreview1(ctx, r))
	must(r.NewModuleBuilder("env").
		ExportFunction("emscripten_notify_memory_growth", func(int32) {}).
		Instantiate(ctx))
	zbarCompiled = must(r.CompileModule(ctx, zbarWasm))
}

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
