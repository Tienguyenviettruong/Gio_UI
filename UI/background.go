package UI

import (
	"bytes"
	"image"
	"image/draw"
	"image/png"
	"io/ioutil"
	"log"

	"gioui.org/layout"
	"gioui.org/op/paint"
)

func drawImageBackground(gtx layout.Context) {
	data, err := ioutil.ReadFile("asset/background.png")
	if err != nil {
		log.Fatal(err)
	}

	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		log.Fatal(err)
	}

	dst := image.NewRGBA(image.Rect(0, 0, gtx.Constraints.Max.X, gtx.Constraints.Max.Y))
	draw.Draw(dst, dst.Bounds(), img, image.Point{}, draw.Over)

	paint.NewImageOp(dst).Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
}
