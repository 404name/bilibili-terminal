package config

type Config struct {
	Output   Output   `mapstructure:"output" json:"output" yaml:"output"`
	BasePath BasePath `mapstructure:"base-path" json:"base-path" yaml:"base-path"`
}

func NewDefaultConfig() Config {
	return Config{
		Output{
			OutputImgPath:   "./output/img/img%04d.png",
			OutputAudioPath: "./output/audio/audio.mp3",
			OutputVideoPath: "./output/video/video.mp4",
			OutputLogPath:   "./output/app.log",
		},
		BasePath{
			VideoCoverImg: "./resource/images/video_cover_img.png,",
			VideoStopImg:  "./resource/images/video_stop_img.png",
		},
	}
}
