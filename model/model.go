package model

var (
	Video *VideoDetail
)

func Init() {
	Video = &VideoDetail{}
	Video.Init()
}
