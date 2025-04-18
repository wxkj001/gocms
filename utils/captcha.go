package utils

import (
	"bytes"
	"math/rand"
	"time"

	svg "github.com/ajstarks/svgo"
)

func GenerateSVG(width, height int) ([]byte, string) {
	rand.Seed(time.Now().UnixNano())

	var svgContent bytes.Buffer
	canvas := svg.New(&svgContent)
	canvas.Start(width, height)
	canvas.Rect(0, 0, width, height, "fill:white")

	// 添加干扰线
	for i := 0; i < 5; i++ {
		canvas.Line(rand.Intn(width), rand.Intn(height), rand.Intn(width), rand.Intn(height), "stroke:rgb(150,150,150);stroke-width:1")
	}

	// 添加干扰点 
	for i := 0; i < 30; i++ {
		canvas.Circle(rand.Intn(width), rand.Intn(height), 1, "fill:rgb(200,200,200)")
	}

	// 生成4位字母数字验证码
	chars := []rune("ABCDEFGHJKLMNPQRSTUVWXYZ23456789")
	code := make([]rune, 4)
	for i := range code {
		code[i] = chars[rand.Intn(len(chars))]
	}
	codeStr := string(code)

	// 绘制验证码文本
	canvas.Text(width/2, height/2+5, codeStr, "text-anchor:middle; font-size:24px; font-weight:bold; fill:rgb(80,80,80); font-family:Arial")

	canvas.End()

	return svgContent.Bytes(), codeStr
}
