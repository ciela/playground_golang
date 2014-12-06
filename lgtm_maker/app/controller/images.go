package controller

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/ciela/playground_golang/lgtm_maker/aws"

	"code.google.com/p/go-uuid/uuid"
	"github.com/mitchellh/goamz/s3"
	"github.com/naoina/kocha"
)

const (
	// 最大画像容量(4MB)
	MaxBytes = 4194304
	// MIMEs
	JpegCT = "image/jpeg"
	PngCT  = "image/png"
	GifCT  = "image/gif"
)

var lgtmImg image.Image

func init() {
	log.Println("Reading default LGTM image...")

	lgtmReader, err := os.Open("assets/lgtm.png")
	if err != nil {
		log.Fatalln(err)
	}
	defer lgtmReader.Close()

	lgtmImg, _, err = image.Decode(lgtmReader)
	if err != nil {
		log.Fatalln(err)
	}
}

type (
	RGBADrawer     func(i *image.Image) (err error)
	PalettedDrawer func(i *image.Paletted) (dst *image.Paletted, err error)

	Images struct {
		*kocha.DefaultController
	}
)

var drawLGTMWithRGBA = func(i *image.Image) (err error) {
	rect := (*i).Bounds()
	rgbaImg := image.NewRGBA(rect)
	draw.Draw(rgbaImg, rect, *i, rect.Min, draw.Src)
	//TODO adaptive resize
	draw.Draw(rgbaImg, lgtmImg.Bounds(), lgtmImg, rect.Min, draw.Over)
	*i = rgbaImg
	return
}

var drawLGTMWithPaletted = func(i *image.Paletted) (dst *image.Paletted, err error) {
	rect := (*i).Bounds()

	hist := make(map[color.Color]int)
	for y := 0; y < rect.Dy(); y++ {
		for x := 0; x < rect.Dx(); x++ {
			hist[i.At(x, y)]++
		}
	}

	var r1, r2 color.Color // the top 2 rare colors
	lastV := rect.Size().X * rect.Size().Y
	for k, v := range hist {
		if v < lastV {
			r2 = r1
			r1 = k
		}
		lastV = v
	}

	if len(hist) > 254 {
		delete(hist, r1)
		delete(hist, r2)
	} else if len(hist) > 255 {
		delete(hist, r1)
	}

	var newPalette []color.Color
	for k, _ := range hist {
		newPalette = append(newPalette, k)
	}
	log.Printf("Image size: %v, Length of original palette: %v, Length of new palette: %v, Calculated colors: %v\n", rect.Size(), len(i.Palette), len(newPalette), len(hist))

	dst = image.NewPaletted(rect, append(newPalette, color.Black, color.White))
	draw.Draw(dst, rect, i, rect.Min, draw.Src)
	draw.Draw(dst, lgtmImg.Bounds(), lgtmImg, rect.Min, draw.Over)
	return
}

func (im *Images) GET(c *kocha.Context) kocha.Result {
	// FIXME: auto-generated by kocha
	return kocha.Render(c)
}

func (im *Images) POST(c *kocha.Context) kocha.Result {
	//リクエスト容量が大きかったら最初から弾く
	if c.Request.ContentLength > MaxBytes {
		return kocha.RenderError(c, http.StatusBadRequest, "Size of your reqest is too large")
	}

	//FormDataの取得
	f, h, err := c.Request.FormFile("image")
	if err != nil {
		return kocha.RenderError(c, http.StatusBadRequest, "Request has not been accepted")
	}
	defer f.Close()

	//ヘッダ情報からの形式取得
	var imgf string
	ct := h.Header["Content-Type"][0]
	switch ct {
	case JpegCT:
		imgf = "jpeg"
	case GifCT:
		imgf = "gif"
	case PngCT:
		imgf = "png"
	default:
		return kocha.RenderError(c, http.StatusBadRequest, "Image format must be either jpeg, gif or png")
	}

	//画像のデコードとフォーマットバリデーション
	img, fm, err := image.Decode(f)
	if err != nil {
		return kocha.RenderError(c, http.StatusBadRequest, "Error has occured when decoding the reqested image")
	} else if fm != imgf {
		return kocha.RenderError(c, http.StatusBadRequest, "The format of actual image is not the same as the header")
	}

	//指定のDrawerを使って描画しつつエンコード
	b := new(bytes.Buffer)
	switch ct {
	case JpegCT:
		err = encodeJPEG(&img, b, drawLGTMWithRGBA)
	case GifCT:
		err = encodeGIF(&f, b, drawLGTMWithPaletted)
	case PngCT:
		err = encodePNG(&img, b, drawLGTMWithRGBA)
	}
	if err != nil {
		return kocha.RenderError(c, http.StatusBadRequest, "Error has occured when encoding the requested image")
	} else if b.Len() > MaxBytes {
		return kocha.RenderError(c, http.StatusBadRequest, "Actual size of requested image is too large")
	}

	// 配置用のパス決めてS3に配置
	p := uuid.New() //ver4
	if err = aws.LgtmBucket.Put(p, b.Bytes(), ct, s3.PublicRead); err != nil {
		return kocha.RenderError(c, http.StatusInternalServerError, "An error has occured when uploading image: "+p)
	}

	// TODO DBにIDを保存

	return kocha.Render(c, kocha.Data{"imagePath": p})
}

func (im *Images) PUT(c *kocha.Context) kocha.Result {
	// FIXME: auto-generated by kocha
	//iid := c.Params.Get("imageId")
	return kocha.Render(c)
}

func (im *Images) DELETE(c *kocha.Context) kocha.Result {
	// FIXME: auto-generated by kocha
	//iid := c.Params.Get("imageId")
	return kocha.Render(c)
}

func encodeJPEG(img *image.Image, b *bytes.Buffer, draw RGBADrawer) (err error) {
	if err = draw(img); err != nil {
		log.Println(err.Error())
		return
	}
	err = jpeg.Encode(b, *img, &jpeg.Options{Quality: 100})
	return
}

func encodePNG(img *image.Image, b *bytes.Buffer, draw RGBADrawer) (err error) {
	if err = draw(img); err != nil {
		log.Println(err.Error())
		return
	}
	err = png.Encode(b, *img)
	return
}

func encodeGIF(f *multipart.File, b *bytes.Buffer, draw PalettedDrawer) (err error) {
	(*f).Seek(0, 0) //Seekerのリセット
	gImg, err := gif.DecodeAll(*f)
	if err != nil {
		log.Println("Decoding error: " + err.Error())
		return
	}
	var frames []*image.Paletted
	for _, p := range gImg.Image { // []*image.Palleted
		newP, err := draw(p)
		if err != nil {
			log.Println("Drawing error: " + err.Error())
			return err
		}
		frames = append(frames, newP)
	}
	g := &gif.GIF{
		Delay:     gImg.Delay,
		Image:     frames,
		LoopCount: gImg.LoopCount,
	}
	err = gif.EncodeAll(b, g)
	return
}
