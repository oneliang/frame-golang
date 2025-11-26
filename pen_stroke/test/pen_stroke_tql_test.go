package test

import (
	"fmt"
	"github.com/oneliang/frame-golang/pen_stroke/tql"
	"testing"
	"time"
)

func TestTqlPenStroke(t *testing.T) {
	begin := time.Now().UnixMilli()
	_ = tql.GeneratePenStrokePng("stroke/_0_1.txt", "stroke/output_1.png", 60)
	end := time.Now().UnixMilli()
	fmt.Println(fmt.Sprintf("cost:%d", end-begin))
}
