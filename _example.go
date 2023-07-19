package main

import (
	"fmt"
	"image/color"
	"os"
)

func main() {
	// 使用默认参数创建Watermarker
	wm1, err := textWatermark.NewWatermarker()
	if err != nil {
		fmt.Println("创建Watermarker出错:", err)
		return
	}

	// 使用自定义参数创建Watermarker
	customColor := color.RGBA{R: 255, G: 0, B: 0, A: 80} // 红色
	wm2, err := textWatermark.NewWatermarker(
		textWatermark.WithWatermarkString("自定义水印"),
		textWatermark.WithWatermarkColor(customColor),
		textWatermark.WithSkew(30),
		textWatermark.WithStepX(250),
		textWatermark.WithStepY(130),
		textWatermark.WithFontPath("customFont.ttf"),
	)
	if err != nil {
		fmt.Println("创建自定义Watermarker出错:", err)
		return
	}

	// 打开要添加水印的图像文件
	file, err := os.Open("path_to_image.jpg")
	if err != nil {
		fmt.Println("打开图像文件出错:", err)
		return
	}
	defer file.Close()

	// 为图像添加水印
	watermarkedData, err := wm2.AddWatermark(file, "image/jpeg")
	if err != nil {
		fmt.Println("添加水印出错:", err)
		return
	}

	// 保存处理后的图像到当前目录
	outputFile, err := os.Create("watermarked_image.jpg")
	if err != nil {
		fmt.Println("创建输出文件出错:", err)
		return
	}
	defer outputFile.Close()

	_, err = outputFile.Write(watermarkedData)
	if err != nil {
		fmt.Println("写入输出文件出错:", err)
		return
	}

	fmt.Println("水印图像已成功保存为 watermarked_image.jpg")
}
