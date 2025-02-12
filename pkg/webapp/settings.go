package webapp

import "time"

type (
	Settings struct {
		Log          LogSettings          `mapstructure:"log"`
		Server       ServerSettings       `mapstructure:"server"`
		StaticServer StaticServerSettings `mapstructure:"static_server"`

		config Config
	}

	LogSettings struct {
		Level string `mapstructure:"level"`
	}

	ServerSettings struct {
		Addr         string        `mapstructure:"addr"`
		ReadTimeout  time.Duration `mapstructure:"read_timeout"`
		WriteTimeout time.Duration `mapstructure:"write_timeout"`
		IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
		// TODO: add https support
	}

	StaticServerSettings struct {
		Enabled    bool   `mapstructure:"enabled"`
		Path       string `mapstructure:"path"`
		IndexPath  string `mapstructure:"index_path"`
		AssetsPath string `mapstructure:"assets_path"`
	}
)

func (s *Settings) GetConfig() Config {
	return s.config
}
