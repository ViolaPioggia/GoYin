package tools

import (
	"bytes"
	"github.com/disintegration/imaging"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"image"
)

func GetVideoCover(videoPath string, coverPath string) error {
	buf := bytes.NewBuffer(nil)
	err := ffmpeg.Input(videoPath).
		Filter("select", ffmpeg.Args{"gte(n,1)"}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf).Run()
	if err != nil {
		return err
	}
	var img image.Image
	if img, err = imaging.Decode(buf); err != nil {
		return err
	}
	if err = imaging.Save(img, coverPath); err != nil {
		return err
	}
	return nil
}
