package utils

import (
	"encoding/base64"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"log"
	"os"
	"strings"

	"github.com/404name/termui-demo/resource"
)

func CovertImg(base64Img string) image.Image {
	image, _, err := image.Decode(base64.NewDecoder(base64.StdEncoding, strings.NewReader(base64Img)))
	if err != nil {
		log.Fatalf("failed to decode gopher image: %v", err)
		return nil
	}
	return image
}

// 读取不到默认返回LOGO
func LoadImg(path string) image.Image {
	f, err := os.Open(path)
	if err != nil {
		fmt.Println(err.Error())
		return CovertImg(resource.LOGO)
	}
	// decode图片
	m, err := png.Decode(f)
	if err != nil {

		fmt.Println(err.Error())
		return CovertImg(resource.LOGO)

	}
	return m
}
