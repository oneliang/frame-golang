package tql

import (
	"encoding/json"
	"fmt"
	"github.com/oneliang/util-golang/constants"
	"github.com/oneliang/util-golang/file"
	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/renderers"
	"math"
	"strings"
)

type PenDot struct {
	F         float64 `json:"f"`
	LineWidth float64 `json:"lineWidth"`
	T         int     `json:"t"`
	Type      string  `json:"type"`
	X         float64 `json:"x"`
	Y         float64 `json:"y"`
}
type PenStroke struct {
	Action    int       `json:"action"`
	BookId    int       `json:"bookId"`
	Color     int       `json:"color"`
	ColorHex  string    `json:"colorHex"`
	Id        string    `json:"id"`
	LineWidth float64   `json:"lineWidth"`
	List      []*PenDot `json:"list"`
	Page      int       `json:"page"`
	PenType   int       `json:"penType"`
	Time      int64     `json:"time"`
}

type CoordinateInfo struct {
	MinX float64
	MinY float64
	MaxX float64
	MaxY float64
}

func parsePenStroke(filename string) ([]*PenStroke, *CoordinateInfo) {
	var penStrokeList []*PenStroke
	var coordinateInfo = &CoordinateInfo{
		MinX: -1,
		MinY: -1,
		MaxX: -1,
		MaxY: -1,
	}
	_ = file.ReadFileContentEachLine(filename, func(content string) bool {
		fixContent := strings.Trim(content, constants.STRING_BLANK)
		var penStroke PenStroke
		err := json.Unmarshal([]byte(fixContent), &penStroke)
		if err != nil {
			return false
		}
		penStrokeList = append(penStrokeList, &penStroke)
		for _, penDot := range penStroke.List {
			//initialize all value
			if coordinateInfo.MinX < 0 {
				coordinateInfo.MinX = penDot.X
			}
			if coordinateInfo.MinY < 0 {
				coordinateInfo.MinY = penDot.Y
			}
			if coordinateInfo.MaxX < 0 {
				coordinateInfo.MaxX = penDot.X
			}
			if coordinateInfo.MaxY < 0 {
				coordinateInfo.MaxY = penDot.Y
			}
			//compare
			if coordinateInfo.MinX >= penDot.X {
				coordinateInfo.MinX = penDot.X
			}
			if coordinateInfo.MinY >= penDot.Y {
				coordinateInfo.MinY = penDot.Y
			}
			if coordinateInfo.MaxX < penDot.X {
				coordinateInfo.MaxX = penDot.X
			}
			if coordinateInfo.MaxY < penDot.Y {
				coordinateInfo.MaxY = penDot.Y
			}
		}
		return true
	})
	return penStrokeList, coordinateInfo
}

const (
	TYPE_PEN_DOWN = "PEN_DOWN"
	TYPE_PEN_MOVE = "PEN_MOVE"
	TYPE_PEN_UP   = "PEN_UP"
)

func GeneratePenStrokePng(penStrokeFile string, outputFile string, margin float64) error {
	penStrokeList, coordinateInfo := parsePenStroke(penStrokeFile)

	var width = math.Ceil(coordinateInfo.MaxX - coordinateInfo.MinX + 2*margin)
	var height = math.Ceil(coordinateInfo.MaxY - coordinateInfo.MinY + 2*margin)

	fmt.Println(fmt.Sprintf("width:%f, height:%f", width, height))

	c := canvas.New(width, height)
	ctx := canvas.NewContext(c)
	ctx.SetView(canvas.Identity.Scale(1.0, -1.0))
	//ctx.SetS
	ctx.SetFillColor(canvas.White)
	ctx.SetStrokeColor(canvas.Black)
	ctx.SetStrokeWidth(1.0)
	offsetX := margin
	offsetY := margin
	//draw margin
	ctx.MoveTo(0, 0)
	ctx.LineTo(width, 0)
	ctx.LineTo(width, height)
	//ctx.LineTo(0, height)
	//ctx.LineTo(0, 0)

	for _, penStroke := range penStrokeList {
		var lastX, lastY float64
		//fmt.Println(len(penStroke.List))
		for dotIndex, penDot := range penStroke.List {
			if dotIndex == 0 && penDot.Type != TYPE_PEN_DOWN {
				break
			}
			x := penDot.X - coordinateInfo.MinX + offsetX
			y := penDot.Y - coordinateInfo.MinY + offsetY
			switch penDot.Type {
			case TYPE_PEN_DOWN:
				ctx.MoveTo(x, y)
				lastX = x
				lastY = y
			case TYPE_PEN_MOVE:
				ctx.QuadTo(lastX, lastY, (lastX+x)/2, (lastY+y)/2)
				lastX = x
				lastY = y
			case TYPE_PEN_UP:
				ctx.QuadTo(lastX, lastY, (lastX+x)/2, (lastY+y)/2)
				//lastX = x
				//lastY = y
			}
		}
		//break
	}
	ctx.FillStroke()
	c.Fit(0)
	return renderers.Write(outputFile, c, canvas.DPI(48))
}
