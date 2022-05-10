package zbar_test

import (
	"bytes"
	"image/png"
	"io"
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/mrg0lden/go-zbar-wasm"
	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
)

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}

	return v
}

type nopCloser struct {
	io.Writer
}

func (nopCloser) Close() error {
	return nil
}

func Test_E2E(t *testing.T) {
	qr := must(qrcode.NewWith("Hello, world"))
	buf := bytes.Buffer{}
	w := standard.NewWithWriter(nopCloser{&buf}, standard.WithBuiltinImageEncoder(standard.PNG_FORMAT))
	qr.Save(w)

	img := must(png.Decode(&buf))

	s := zbar.NewScanner()
	res, err := s.ReadAll(img)
	assert.NoError(t, err)
	assert.Equal(t, []byte("Hello, world"), res[0])

}
