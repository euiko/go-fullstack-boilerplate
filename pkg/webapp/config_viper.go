package webapp

import (
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

type (
	Viper struct {
		viper *viper.Viper

		// some options
		standalone bool
		paths      []string
	}

	ViperOptions interface {
		Configure(v *Viper)
	}

	ViperOptionsFunc func(v *Viper)

	viperValue struct {
		Value // ensure this implementation align with Value interface
		key   string
		viper *viper.Viper
	}
)

func (f ViperOptionsFunc) Configure(v *Viper) {
	f(v)
}

func (c *Viper) Sub(path string) Config {
	sub := c.viper.Sub(path)
	if sub == nil {
		sub = viper.New()
	}

	return &Viper{
		viper: sub,
	}
}

func (c *Viper) Get(path string) Value {
	return &viperValue{key: path, viper: c.viper}
}

func (c *Viper) Set(path string, val interface{}) error {
	c.viper.Set(path, val)
	return nil
}

func (c *Viper) Scan(out interface{}) error {
	v := new(viperValue)
	v.key = ""
	v.viper = c.viper

	return v.Scan(out)
}

func (c *Viper) Write() error {
	return c.viper.WriteConfig()
}

func (c *Viper) OnChange(callback OnChangedFunc) {
	c.viper.OnConfigChange(func(in fsnotify.Event) {
		callback()
	})
}

func (v *viperValue) Bool(defaults ...bool) bool {
	d := false
	if len(defaults) > 0 {
		d = defaults[0]
	}

	value := v.viper.Get(v.key)
	if val, ok := value.(bool); ok {
		return val
	}

	return d
}
func (v *viperValue) String(defaults ...string) string {
	d := ""
	if len(defaults) > 0 {
		d = defaults[0]
	}

	value := v.viper.Get(v.key)
	if val, ok := value.(string); ok {
		return val
	}

	return d

}
func (v *viperValue) Float64(defaults ...float64) float64 {
	d := 0.0
	if len(defaults) > 0 {
		d = defaults[0]
	}

	value := v.viper.Get(v.key)
	if val, ok := value.(float64); ok {
		return val
	}

	return d
}
func (v *viperValue) Duration(defaults ...time.Duration) time.Duration {
	d := time.Duration(0)
	if len(defaults) > 0 {
		d = defaults[0]
	}

	value := v.viper.Get(v.key)
	if val, ok := value.(time.Duration); ok {
		return val
	}

	return d
}
func (v *viperValue) StringSlice(defaults ...[]string) []string {
	d := []string{}
	if len(defaults) > 0 {
		d = defaults[0]
	}

	value := v.viper.Get(v.key)
	if val, ok := value.([]string); ok {
		return val
	}

	return d
}
func (v *viperValue) StringMap(defaults ...map[string]interface{}) map[string]interface{} {
	d := make(map[string]interface{})
	if len(defaults) > 0 {
		d = defaults[0]
	}

	if value, err := cast.ToStringMapE(v.viper.Get(v.key)); err == nil {
		return value
	}

	return d
}
func (v *viperValue) StringMapString(defaults ...map[string]string) map[string]string {
	d := make(map[string]string)
	if len(defaults) > 0 {
		d = defaults[0]
	}

	if value, err := cast.ToStringMapStringE(v.viper.Get(v.key)); err == nil {
		return value
	}

	return d
}

func (v *viperValue) Scan(val interface{}) error {
	if v.key == "" {
		return v.viper.Unmarshal(val)
	}

	return v.viper.UnmarshalKey(v.key, val)
}

func ViperStandalone() ViperOptions {
	return ViperOptionsFunc(func(v *Viper) {
		v.standalone = true
	})
}

func ViperPaths(paths ...string) ViperOptions {
	return ViperOptionsFunc(func(v *Viper) {
		v.paths = paths
	})
}

func NewViper(path string, opts ...ViperOptions) *Viper {
	v := viper.New()
	vpr := Viper{
		viper: v,
	}

	// decorate viper
	for _, o := range opts {
		o.Configure(&vpr)
	}

	// set the config name and path
	v.SetConfigName(path)
	v.AddConfigPath(".")

	// add additional config path
	for _, p := range vpr.paths {
		v.AddConfigPath(p)
	}

	// read the config
	if !vpr.standalone {
		if err := v.ReadInConfig(); err != nil {
			panic(err)
		}
	}

	return &vpr
}
