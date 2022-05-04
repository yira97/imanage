package main

import (
	"flag"
	"log"

	"github.com/yira97/imanage/internal/pkg/transcode"
)

var (
	flagWorkspace           string
	flagTranscode           bool
	flagTranscodeTargetWebp bool
	flagTranscodeTargetAvif bool
)

func init() {
	flag.StringVar(&flagWorkspace, "workspace", "./imanage_data", "workspace directory, eg: --workspace /PATH/TO/YOUR_WORKSPACE")
	flag.BoolVar(&flagTranscode, "transcode", false, "perform transcode, eg: --transcode")
	flag.BoolVar(&flagTranscodeTargetWebp, "to_avif", false, "")
}

func main() {
	flag.Parse()

	t := transcode.NewTranscoder()
	t.Setup(nil)
	t.SetInput("/Users/yiran/Downloads/ticket_すみだ水族館.pdf")
	t.SetOutput("/Users/yiran/Desktop")
	// t.AddTarget(transcode.AVIF)
	t.AddTarget(transcode.WEBP)
	err := t.Execute()
	if err != nil {
		log.Fatal("error")
	}

	// check path exist
	// if path not exist
	// create it
	// if created failed
	//exit
	// make inbox folder
	// make data folder
	// make data/webp folder
	// make data/avif folder

	// ### cmd:
}
