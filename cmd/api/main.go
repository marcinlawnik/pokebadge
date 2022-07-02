package main

import (
	"github.com/shopspring/decimal"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"
	"pokebadge/internal/credly"
	image_overlayer "pokebadge/internal/image-overlayer"
)

func main() {
	//Get Most recent badge
	getBadgeUC := credly.NewGetLatestCertBadgeUseCase(
		credly.GetLatestCertBadgeUseCaseInput{
			Client:   credly.NewCredlyClient(http.Client{}),
			Username: "",
		},
	)

	badge, _ := getBadgeUC.Do()

	//Download badge image
	err := DownloadFile("badge.png", badge.ImageURL)
	if err != nil {
		panic(err)
	}
	badgeFile, err := os.Open("badge.png")
	if err != nil {
		log.Fatalf("failed to open: %s", err)
	}

	badgeImage, err := png.Decode(badgeFile)
	if err != nil {
		log.Fatalf("failed to decode: %s", err)
	}
	defer badgeFile.Close()

	//downoad bg image
	err = DownloadFile("bg.png", "https://raw.githubusercontent.com/marcinlawnik/pokebadge/main/assets/badge-template.png")
	if err != nil {
		panic(err)
	}

	bgFile, err := os.Open("bg.png")
	if err != nil {
		log.Fatalf("failed to open: %s", err)
	}

	bgImage, err := png.Decode(bgFile)
	if err != nil {
		log.Fatalf("failed to decode: %s", err)
	}
	defer bgFile.Close()

	overlayImageUC := image_overlayer.OverlayImagesUseCase{
		Input: image_overlayer.OverlayImagesUseCaseInput{
			BackgroundImage:      bgImage,
			ForegroundImage:      badgeImage,
			ForegroundImageScale: decimal.NewFromFloat(0.25),
			PositionX:            600,
			PositionY:            50,
		},
	}

	resultImage, err := overlayImageUC.Do()
	if err != nil {
		return
	}

	//Save output image
	result, err := os.Create("result.png")
	if err != nil {
		log.Fatalf("failed to create: %s", err)
	}
	err = png.Encode(result, resultImage)
	if err != nil {
		panic(err)
	}

	defer result.Close()

}

func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
