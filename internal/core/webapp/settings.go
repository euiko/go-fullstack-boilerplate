package webapp

import (
	"os"
	"path"
	"time"

	"github.com/spf13/viper"
)

type (
	Settings struct {
		Log          LogSettings          `mapstructure:"log"`
		Server       ServerSettings       `mapstructure:"server"`
		StaticServer StaticServerSettings `mapstructure:"static_server"`
		DB           DatabaseSettings     `mapstructure:"db"`

		extra *viper.Viper
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
		Enabled bool `mapstructure:"enabled"`
	}

	DatabaseSettings struct {
		// TODO: add support for multiple databases
		// TODO: support database other than sql (postgres)
		Uri             string        `mapstructure:"uri"`
		ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
		MaxIdleConns    int           `mapstructure:"max_idle_conns"`
		MaxOpenConns    int           `mapstructure:"max_open_conns"`
	}
)

func (s *Settings) GetExtra() *viper.Viper {
	return s.extra
}

func loadSettings(name string, shortName string) Settings {
	// default settings
	settings := Settings{
		Log: LogSettings{
			Level: "info",
		},
		Server: ServerSettings{
			Addr:         ":8080",
			ReadTimeout:  60 * time.Second,
			WriteTimeout: 60 * time.Second,
			IdleTimeout:  0,
			// TODO: add https support
		},
		StaticServer: StaticServerSettings{
			Enabled: true,
		},
		DB: DatabaseSettings{
			Uri:             "postgres://order-tracker:12345678@localhost:5432/order-tracker?sslmode=disable",
			ConnMaxLifetime: 60 * time.Second,
			MaxIdleConns:    10,
			MaxOpenConns:    10,
		},
		extra: nil,
	}

	// use viper for configuration
	v := viper.New()
	v.SetConfigName(name)
	v.AddConfigPath(".")

	// use short name as env prefix if it is defined
	if shortName != "" {
		v.SetEnvPrefix(shortName)
	}

	// add config in home directory is it is defined
	homeDir := os.Getenv("HOME")
	if homeDir != "" {
		v.AddConfigPath(homeDir)
		v.AddConfigPath(path.Join(homeDir, ".config", name))
	}

	// load settings
	v.Unmarshal(&settings)

	// set the setting's config
	settings.extra = v.Sub("extra")
	return settings
}
