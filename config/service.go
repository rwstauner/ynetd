package config

// Service holds string representations of Service attributes.
type Service struct {
	AutoStart      bool `yaml:"auto_start"`
	Proxy          map[string]string
	Command        []string
	StopAfter      string `yaml:"stop_after"`
	StopSignal     string `yaml:"stop_signal"`
	Timeout        string
	WaitAfterStart string `yaml:"wait_after_start"`
}
