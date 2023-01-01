package config

type Output struct {
	OutputImgPath   string `mapstructure:"output-img-path" json:"output-img-path" yaml:"output-img-path"`
	OutputAudioPath string `mapstructure:"output-audio-path" json:"output-audio-path" yaml:"output-audio-path"`
	OutputVideoPath string `mapstructure:"output-video-path" json:"output-video-path" yaml:"output-video-path"`
	OutputLogPath   string `mapstructure:"output-log-path" json:"output-log-path" yaml:"output-log-path"`
}
