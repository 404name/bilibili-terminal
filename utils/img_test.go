package utils

import (
	"testing"

	"github.com/404name/termui-demo/resource"
)

func TestLoadImg(t *testing.T) {
	println(LoadImg("." + resource.BaseImg))
}
