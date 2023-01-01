package config

type Config struct {
	Output   Output   `mapstructure:"output" json:"output" yaml:"output"`
	BasePath BasePath `mapstructure:"base-path" json:"base-path" yaml:"base-path"`
}
