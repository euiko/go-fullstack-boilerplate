package webapp

import "time"

type (

	// Value represent value of config fields
	Value interface {
		// Bool returns the value as a boolean
		Bool(def ...bool) bool

		// String returns the value as a string
		String(def ...string) string

		// Float64 returns the value as a float64
		Float64(def ...float64) float64

		// Duration returns the value as a time.Duration
		Duration(def ...time.Duration) time.Duration

		// StringSlice returns the value as a []string
		StringSlice(def ...[]string) []string

		// StringMap returns the value as a map[string]interface{}
		StringMap(def ...map[string]interface{}) map[string]interface{}

		// StringMapString returns the value as a map[string]string
		StringMapString(def ...map[string]string) map[string]string

		// Scan scans the value into a struct
		Scan(val interface{}) error
	}

	// OnChangedFunc for calling on config change callback
	OnChangedFunc func()

	Config interface {
		// Sub returns a new config that is a sub-config of the current config
		Sub(path string) Config

		// Get returns a value from the config
		Get(path string) Value

		// Set sets a value in the config
		Set(path string, val interface{}) error

		// Scan scans the config into a struct
		Scan(out interface{}) error

		// Write writes the config to the config file
		Write() error

		// OnChange registers a callback that is called when the config changes
		OnChange(callback OnChangedFunc)
	}
)
