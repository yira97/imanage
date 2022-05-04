package transcode

import (
	"bytes"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
	"path"
	"strings"

	"github.com/Kagami/go-avif"
	"github.com/chai2010/webp"
)

type Options struct {
	Quality      *int
	OutputPrefix *string
}

type ImageFormat string

const (
	JPEG ImageFormat = "jpeg"
	PNG  ImageFormat = "png"
	GIF  ImageFormat = "gif"
	WEBP ImageFormat = "webp"
	AVIF ImageFormat = "avif"
)

// isSupportedInputFormat 判断一个文件名（不包含路径，纯文件名）是否是支持的图像格式。
func isSupportedInputFormat(filename string) bool {
	frags := strings.Split(filename, ".")
	if len(frags) == 1 {
		return false
	}
	switch ImageFormat(strings.ToLower(frags[len(frags)-1])) {
	// AVIF AND WEBP is not supported yet.
	case JPEG, PNG, GIF:
		return true
	default:
		return false
	}
}

type Transcoder interface {
	Setup(opts *Options)
	AddTarget(ImageFormat)
	Targets() []ImageFormat
	SetInput(file string)
	SetOutput(dir string)
	Execute() error
}

func NewTranscoder() Transcoder {
	return &transCoderTool{
		targets: make(map[ImageFormat]transcoderAdapter),
	}
}

// transCoderTool 实现了 Transcoder
type transCoderTool struct {
	opts      Options
	targets   map[ImageFormat]transcoderAdapter
	inputFile string
	outputDir string
}

func (t *transCoderTool) Setup(opts *Options) {
	if opts == nil {
		return
	}
	t.opts = *opts
}

func (t *transCoderTool) AddTarget(i ImageFormat) {
	switch i {
	case AVIF:
		t.targets[i] = &avifAdapter{}
	case WEBP:
		t.targets[i] = &webpAdapter{}
	}

}

func (t *transCoderTool) Targets() []ImageFormat {
	slice := make([]ImageFormat, 0)

	for k := range t.targets {
		slice = append(slice, k)
	}

	return slice
}

func (t *transCoderTool) SetInput(file string) {
	t.inputFile = file
}

func (t *transCoderTool) SetOutput(dir string) {
	t.outputDir = dir
}

func (t *transCoderTool) Execute() error {
	srcPath := t.inputFile
	src, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("Can't open sorce file: %w", err)
	}
	base := path.Base(t.inputFile)
	if !isSupportedInputFormat(base) {
		return fmt.Errorf("unsupported input format")
	}
	dstPathBuilder := strings.Builder{}
	if t.opts.OutputPrefix != nil {
		dstPathBuilder.WriteString(*t.opts.OutputPrefix)
	}
	baseWithoutSuffix := strings.TrimRight(base, path.Ext(t.inputFile))
	dstPathBuilder.WriteString(baseWithoutSuffix)
	dstPathWithoutSuffix := dstPathBuilder.String()

	for format, adapter := range t.targets {
		dst, err := os.Create(dstPathWithoutSuffix + string(format))
		if err != nil {
			return fmt.Errorf("Can't create destination file: %w", err)
		}
		img, _, err := image.Decode(src)
		if err != nil {
			return fmt.Errorf("image decode error: %w", err)
		}
		if err = adapter.Write(img, dst, &t.opts); err != nil {
			return err
		}
	}

	return nil
}

type transcoderAdapter interface {
	Write(img image.Image, out io.Writer, opts *Options) error
}

type webpAdapter struct{}

type avifAdapter struct{}

// https://pkg.go.dev/github.com/chai2010/webp
func (a *webpAdapter) Write(img image.Image, out io.Writer, opts *Options) error {
	var buf bytes.Buffer

	if err := webp.Encode(&buf, img, &webp.Options{Lossless: true}); err != nil {
		return err
	}
	if _, err := out.Write(buf.Bytes()); err != nil {
		return err
	}

	return nil
}

// https://pkg.go.dev/github.com/Kagami/go-avif
func (a *avifAdapter) Write(img image.Image, out io.Writer, opts *Options) error {
	if err := avif.Encode(out, img, nil); err != nil {
		return fmt.Errorf("Can't encode source image: %w", err)
	}

	return nil
}
