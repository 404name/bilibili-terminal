package config

type BasePath struct {
	VideoCoverImg string `mapstructure:"video-cover-img" json:"video-cover-img" yaml:"video-cover-img"`
	VideoStopImg  string `mapstructure:"video-stop-img" json:video-stop-img" yaml:"video-stop-img"`
}
