package test

import (
	"fmt"
	"github.com/oneliang/frame-golang/canvas/tql"
	"testing"
	"time"
)

func TestTqlPenStroke(t *testing.T) {
	begin := time.Now().UnixMilli()
	_ = tql.GeneratePenStrokePng("stroke/pen_stroke.txt", "stroke/output.png", 60)
	end := time.Now().UnixMilli()
	fmt.Println(fmt.Sprintf("cost:%d", end-begin))
}
