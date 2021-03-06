package bookmap

import (
	"image"
	"image/color"
	"math"

	font "github.com/lian/gonky/font/terminus"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/llgcode/draw2d/draw2dkit"
)

func colourGradientor(p float64, begin, end color.RGBA) color.RGBA {
	if p > 1.0 {
		p = 1.0
	}
	w := p*2 - 1
	w1 := (w + 1) / 2.0
	w2 := 1 - w1

	r := uint8(float64(begin.R)*w1 + float64(end.R)*w2)
	g := uint8(float64(begin.G)*w1 + float64(end.G)*w2)
	b := uint8(float64(begin.B)*w1 + float64(end.B)*w2)

	return color.RGBA{R: r, G: g, B: b, A: 0xff}
}

func (g *Graph) DrawTradeDots(gc *draw2dimg.GraphicContext, x, rowHeight, pricePosition, priceSteps, maxSizeHisto float64) {
	// trade ask dots
	cx := x
	y := 0.0
	maxIdx := len(g.Timeslots)
	for idx := maxIdx - 1; idx >= 0; idx-- {
		slot := g.Timeslots[idx]
		cx -= float64(g.SlotWidth)
		if cx < 0 {
			break
		}
		if slot.isEmpty() || slot.AskTradeSize == 0 {
			continue
		}

		y = ((pricePosition - slot.AskPrice) / priceSteps) * rowHeight

		startAngle := 0 * (math.Pi / 180.0)
		angle := 360 * (math.Pi / 180.0)

		xx := (cx + (float64(g.SlotWidth) / 2))
		t := (slot.AskTradeSize / (maxSizeHisto * 0.8))
		if t > 1.0 {
			t = 1.0
		}
		size := 4 + float64(t*15)
		gc.ArcTo(xx, y, size, size, startAngle, angle)
		gc.SetFillColor(g.Green)
		gc.Fill()
	}

	// trade bid dots
	cx = x
	y = 0.0
	maxIdx = len(g.Timeslots)
	for idx := maxIdx - 1; idx >= 0; idx-- {
		slot := g.Timeslots[idx]
		cx -= float64(g.SlotWidth)
		if cx < 0 {
			break
		}
		if slot.isEmpty() || slot.BidTradeSize == 0 {
			continue
		}

		y = ((pricePosition - slot.BidPrice) / priceSteps) * rowHeight

		startAngle := 0 * (math.Pi / 180.0)
		angle := 360 * (math.Pi / 180.0)

		xx := (cx + (float64(g.SlotWidth) / 2))
		t := (slot.BidTradeSize / (maxSizeHisto * 0.8))
		if t > 1.0 {
			t = 1.0
		}
		size := 4 + float64(t*15)
		gc.ArcTo(xx, y, size, size, startAngle, angle)
		gc.SetFillColor(g.Red)
		gc.Fill()
	}
}

func (g *Graph) DrawBidAskLines(gc *draw2dimg.GraphicContext, x, rowHeight, pricePosition, priceSteps float64) {
	gc.SetLineWidth(2.0)

	// ask line
	gc.SetStrokeColor(g.Red)
	cx := x
	y := 0.0
	start := true
	maxIdx := len(g.Timeslots)
	for idx := maxIdx - 1; idx >= 0; idx-- {
		slot := g.Timeslots[idx]

		cx -= float64(g.SlotWidth)
		if cx < 0 {
			break
		}
		if slot.isEmpty() || slot.AskPrice == 0.0 {
			gc.Stroke()
			start = true
			continue
		}

		y = ((pricePosition - slot.AskPrice) / priceSteps) * rowHeight

		if start {
			start = false
			gc.MoveTo(cx+float64(g.SlotWidth), y)
		} else {
			gc.LineTo(cx+float64(g.SlotWidth), y)
		}
		gc.LineTo(cx, y)
	}
	gc.Stroke()

	// bid line
	gc.SetStrokeColor(g.Green)
	cx = x
	y = 0.0
	start = true
	maxIdx = len(g.Timeslots)
	for idx := maxIdx - 1; idx >= 0; idx-- {
		slot := g.Timeslots[idx]
		cx -= float64(g.SlotWidth)
		if cx < 0 {
			break
		}
		if slot.isEmpty() || slot.BidPrice == 0.0 {
			gc.Stroke()
			start = true
			continue
		}

		y = ((pricePosition - slot.BidPrice) / priceSteps) * rowHeight

		if start {
			start = false
			gc.MoveTo(cx+float64(g.SlotWidth), y)
		} else {
			gc.LineTo(cx+float64(g.SlotWidth), y)
		}
		gc.LineTo(cx, y)
	}
	gc.Stroke()
	gc.SetLineWidth(1.0)
}

func (g *Graph) DrawTimeslots(gc *draw2dimg.GraphicContext, x, rowsCount, rowHeight, pricePosition, priceSteps, maxSizeHisto float64) {
	cx := x

	maxIdx := len(g.Timeslots)
	for idx := maxIdx - 1; idx >= 0; idx-- {
		slot := g.Timeslots[idx]
		//fmt.Println("slot", slot.From, slot.To, slot.Stats == nil)

		cx -= float64(g.SlotWidth)
		if cx < 0 {
			break
		}

		if len(slot.Rows) == 0 {
			slot.GenerateRows(rowsCount, pricePosition, priceSteps)
			slot.Refill()
		} else {
			if idx >= (maxIdx - 3) { // only need to refill last/current two
				slot.Refill()
			}
		}

		x1 := cx
		x2 := cx + float64(g.SlotWidth)

		for i, row := range slot.Rows {
			strength := (row.Size / maxSizeHisto)
			if strength > 0 {
				y := float64(i) * rowHeight
				draw2dkit.Rectangle(gc, x1, y, x2, y+rowHeight)
				gc.SetFillColor(colourGradientor(strength, g.Fg1, g.Bg1))
				gc.Fill()
			}
		}
	}
}

func (g *Graph) DrawTimeline(gc *draw2dimg.GraphicContext, image *image.RGBA, x, y float64) {
	cx := x

	maxIdx := len(g.Timeslots)
	for idx := maxIdx - 1; idx >= 0; idx-- {
		slot := g.Timeslots[idx]

		cx -= float64(g.SlotWidth)
		if cx < 0 {
			break
		}

		if math.Mod(float64(idx), 30) == 0 {
			/*
				gc.SetLineWidth(1.0)
				gc.SetFillColor(g.Bg1)
				gc.MoveTo(cx, 0)
				gc.LineTo(cx, y)
				gc.Fill()
			*/
			font.DrawString(image, int(cx), int(y), slot.From.Format("15:04:05"), g.Fg1)
		}
	}
}
