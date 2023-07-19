package textWatermark

import (
	"bytes"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/golang/freetype"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"

	"github.com/golang/freetype/truetype"
)

// Watermarker 结构体，封装了添加水印所需的属性
type Watermarker struct {
	watermarkString string         // 水印文字
	watermarkColor  *image.Uniform // 水印颜色
	skew            float64        // 倾斜度
	stepX, stepY    int            // X轴和Y轴的步长
	font            *truetype.Font // 字体
}

type WatermarkerOptions struct {
	WatermarkString *string
	WatermarkColor  *color.RGBA
	Skew            *float64
	StepX, StepY    *int
	FontPath        *string
}

type WatermarkerOptionSetter func(*WatermarkerOptions)

func WithWatermarkString(s string) WatermarkerOptionSetter {
	return func(opt *WatermarkerOptions) {
		opt.WatermarkString = &s
	}
}

func WithWatermarkColor(c color.RGBA) WatermarkerOptionSetter {
	return func(opt *WatermarkerOptions) {
		opt.WatermarkColor = &c
	}
}

func WithSkew(s float64) WatermarkerOptionSetter {
	return func(opt *WatermarkerOptions) {
		opt.Skew = &s
	}
}

func WithStepX(x int) WatermarkerOptionSetter {
	return func(opt *WatermarkerOptions) {
		opt.StepX = &x
	}
}

func WithStepY(y int) WatermarkerOptionSetter {
	return func(opt *WatermarkerOptions) {
		opt.StepY = &y
	}
}

func WithFontPath(path string) WatermarkerOptionSetter {
	return func(opt *WatermarkerOptions) {
		opt.FontPath = &path
	}
}

func NewWatermarker(setters ...WatermarkerOptionSetter) (*Watermarker, error) {
	opt := &WatermarkerOptions{}

	for _, setter := range setters {
		setter(opt)
	}

	watermarkString := "Default Watermark" // 默认值
	if opt.WatermarkString != nil {
		watermarkString = *opt.WatermarkString
	}

	watermarkColor := color.RGBA{R: 128, G: 128, B: 128, A: 80} // 默认灰色，半透明
	if opt.WatermarkColor != nil {
		watermarkColor = *opt.WatermarkColor
	}

	skew := 22.5 // 默认倾斜度
	if opt.Skew != nil {
		skew = *opt.Skew
	}

	stepX := 240 // 默认X轴步长
	if opt.StepX != nil {
		stepX = *opt.StepX
	}

	stepY := 120 // 默认Y轴步长
	if opt.StepY != nil {
		stepY = *opt.StepY
	}

	fontPath := "SourceHanSansCN-Bold.ttf" // 默认字体路径
	if opt.FontPath != nil {
		fontPath = *opt.FontPath
	}

	fontBytes, err := os.ReadFile(fontPath)
	if err != nil {
		return nil, fmt.Errorf("读取字体文件出错: %w", err)
	}

	fnt, err := truetype.Parse(fontBytes)
	if err != nil {
		return nil, fmt.Errorf("解析字体出错: %w", err)
	}

	return &Watermarker{
		watermarkString: watermarkString,
		watermarkColor:  image.NewUniform(watermarkColor),
		skew:            skew,
		stepX:           stepX,
		stepY:           stepY,
		font:            fnt,
	}, nil
}

// GifAddWaterMark 为GIF图像添加水印
func (w *Watermarker) GifAddWaterMark(imgFile io.Reader) ([]byte, error) {
	img, err := gif.DecodeAll(imgFile)
	if err != nil {
		return nil, fmt.Errorf("解码GIF出错: %w", err)
	}

	newGIF := &gif.GIF{}
	for i, frame := range img.Image {
		newImg := image.NewPaletted(frame.Bounds(), frame.Palette)
		draw.Draw(newImg, newImg.Bounds(), frame, frame.Bounds().Min, draw.Src)

		w.writeWatermark(newImg)

		newGIF.Image = append(newGIF.Image, newImg)
		newGIF.Delay = append(newGIF.Delay, img.Delay[i])
		newGIF.Disposal = append(newGIF.Disposal, img.Disposal[i])
	}

	var buf bytes.Buffer
	err = gif.EncodeAll(&buf, newGIF)
	if err != nil {
		return nil, fmt.Errorf("编码GIF出错: %w", err)
	}
	return buf.Bytes(), nil
}

// ImageAddWaterMark 为普通图像添加水印
func (w *Watermarker) ImageAddWaterMark(imgFile io.Reader, format string) ([]byte, error) {
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return nil, fmt.Errorf("解码图像出错: %w", err)
	}

	newImg := image.NewRGBA(img.Bounds())
	draw.Draw(newImg, newImg.Bounds(), img, img.Bounds().Min, draw.Src)

	w.writeWatermark(newImg)

	var buf bytes.Buffer
	switch format {
	case "jpeg", "jpg":
		err = jpeg.Encode(&buf, newImg, nil)
	case "png":
		err = png.Encode(&buf, newImg)
	default:
		return nil, fmt.Errorf("不支持的图像格式: %s", format)
	}
	if err != nil {
		return nil, fmt.Errorf("编码图像出错: %w", err)
	}
	return buf.Bytes(), nil
}

// writeWatermark 是一个私有方法，用于在图像上写入水印
func (w *Watermarker) writeWatermark(newImg draw.Image) {
	for y := -w.stepY; y <= newImg.Bounds().Max.Y+w.stepY; y += w.stepY {
		for x := -w.stepX; x <= newImg.Bounds().Max.X+w.stepX; x += w.stepX {
			offsetX := 0
			if (y/w.stepY)%2 == 1 {
				offsetX = w.stepX / 2
			}

			c := freetype.NewContext()
			c.SetDPI(72)
			c.SetFont(w.font)
			c.SetFontSize(12)
			c.SetSrc(w.watermarkColor)

			textImg := image.NewRGBA(image.Rect(0, 0, len(w.watermarkString)*int(c.PointToFixed(12)>>6), int(c.PointToFixed(12*1.5)>>6)))
			c.SetClip(textImg.Bounds())
			c.SetDst(textImg)

			pt := freetype.Pt(0, int(c.PointToFixed(12)>>6))
			_, _ = c.DrawString(w.watermarkString, pt)

			rotated := imaging.Rotate(textImg, w.skew, image.Transparent)
			draw.Draw(newImg, rotated.Bounds().Add(image.Pt(x+offsetX, y)), rotated, image.Pt(0, 0), draw.Over)
		}
	}
}

// AddWatermark 为文件添加水印
func (w *Watermarker) AddWatermark(file io.Reader, contentType string) ([]byte, error) {
	var fileType string
	switch contentType {
	case "image/jpeg":
		fileType = "jpeg"
	case "image/jpg":
		fileType = "jpg"
	case "image/png":
		fileType = "png"
	case "image/gif":
		fileType = "gif"
	default:
		return nil, fmt.Errorf("不支持的文件类型: %s", contentType)
	}

	switch fileType {
	case "gif":
		return w.GifAddWaterMark(file)
	default:
		return w.ImageAddWaterMark(file, fileType)
	}
}
