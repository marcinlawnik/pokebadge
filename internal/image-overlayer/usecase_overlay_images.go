package image_overlayer

import (
	"github.com/shopspring/decimal"
	xdraw "golang.org/x/image/draw"
	"image"
	"image/draw"
)

type OverlayImagesUseCase struct {
	Input OverlayImagesUseCaseInput
}

type OverlayImagesUseCaseInput struct {
	BackgroundImage      image.Image
	ForegroundImage      image.Image
	ForegroundImageScale decimal.Decimal
	PositionX            int
	PositionY            int
}

func (c OverlayImagesUseCase) Do() (image.Image, error) {
	b := c.Input.BackgroundImage.Bounds()
	result := image.NewRGBA(b)
	draw.Draw(result, b, c.Input.BackgroundImage, image.ZP, draw.Src)

	offset := image.Pt(c.Input.PositionX, c.Input.PositionY)

	// Set the expected size that you want:
	resizedBadge := image.NewRGBA(image.Rect(
		0,
		0,
		int(decimal.NewFromInt(int64(c.Input.ForegroundImage.Bounds().Max.X)).Mul(c.Input.ForegroundImageScale).IntPart()),
		int(decimal.NewFromInt(int64(c.Input.ForegroundImage.Bounds().Max.Y)).Mul(c.Input.ForegroundImageScale).IntPart()),
	))

	// Resize:
	xdraw.BiLinear.Scale(resizedBadge, resizedBadge.Rect, c.Input.ForegroundImage, c.Input.ForegroundImage.Bounds(), draw.Over, nil)

	draw.Draw(result, resizedBadge.Bounds().Add(offset), resizedBadge, image.ZP, draw.Over)

	return result, nil
}
