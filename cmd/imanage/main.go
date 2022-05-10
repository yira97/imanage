package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/yira97/imanage/internal/pkg/transcode"
)

var (
	flagWorkspace                string
	flagTranscode                bool
	flagTranscodeQuality         int
	flagTranscodeTargetWebp      bool
	flagTranscodeTargetAvif      bool
	flagTranscodeIncludeMetaData bool
	flagTranscodeUseLibwebp      bool
	flagTranscodeUseLibavif      bool
)

const (
	defaultTranscodeQuality int = 90
)

func init() {
	// 指定工作区目录
	flag.StringVar(&flagWorkspace, "workspace", "./imanage_data", "workspace directory, eg: --workspace /PATH/TO/YOUR_WORKSPACE")
	// 转码命令
	flag.BoolVar(&flagTranscode, "transcode", false, "perform transcode, eg: --transcode to_webp to_avif")
	// 指定转码到webp
	flag.BoolVar(&flagTranscodeTargetWebp, "to_webp", true, "transcode images to webP format")
	// 指定转码到avif
	flag.BoolVar(&flagTranscodeTargetAvif, "to_avif", false, "transcode images to AVIF format")
	// 指定转码质量
	flag.IntVar(&flagTranscodeQuality, "quality", defaultTranscodeQuality, "assign the quality of the output image")
	// 指定转码包含元数据
	flag.BoolVar(&flagTranscodeIncludeMetaData, "include_metadata", false, "whether or not include the metadata")
	// 指定使用libwebp
	flag.BoolVar(&flagTranscodeUseLibwebp, "use_libwebp", false, "use libwebp which is in your system")
	// 指定使用libavif
	flag.BoolVar(&flagTranscodeUseLibwebp, "use_libavif", false, "use libavif which is in your system")
}

func GetInputDir() string {
	return path.Join(flagWorkspace, "input")
}

func GetOutputDir() string {
	return path.Join(flagWorkspace, "output")
}

// WorkspaceInit initialize the workspace which will be used later
func WorkspaceInit() {
	if err := os.MkdirAll(GetInputDir(), os.ModePerm); err != nil {
		log.Fatal(err)
	}

	if err := os.MkdirAll(GetOutputDir(), os.ModePerm); err != nil {
		log.Fatal(err)
	}
}

func main() {
	flag.Parse()
	WorkspaceInit()

	if flagTranscode {
		t := transcode.NewTranscoder()

		// setup transcoder
		opts := &transcode.Options{Quality: &flagTranscodeQuality}
		if flagTranscodeIncludeMetaData {
			opts.Metadata = true
		}
		if flagTranscodeUseLibwebp {
			opts.UseLibwebp = true
		}
		if flagTranscodeUseLibavif {
			opts.UseLibavif = true
		}
		t.Setup(opts)

		t.SetOutput(GetOutputDir())
		if flagTranscodeTargetAvif {
			t.AddTarget(transcode.AVIF)
		}
		if flagTranscodeTargetWebp {
			t.AddTarget(transcode.WEBP)
		}

		entries, err := os.ReadDir(GetInputDir())
		if err != nil {
			log.Fatal("failed to read workspace input")
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			f := path.Join(GetInputDir(), entry.Name())
			t.SetInput(f)
			err := t.Execute()
			if err != nil {
				log.Printf("transcode failed: [%s]: %v", f, err)
			}
		}

		fmt.Println("Transcode complete!")
	}
}
