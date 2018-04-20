package main

// The configuration options for kindlepush.
type appConfig struct {
	Verbose       bool        `yaml:"verbose"`
	KindleAddress string      `yaml:"kindleAddress"`
	CacheDir      string      `yaml:"cacheDir"`
	ResizeImage   string      `yaml:"resizeImage"`
	MaxFileSize   int         `yaml:"maxFileSize"`
	Feeds         []string    `yaml:"feeds"`
	Smtp          *smtpConfig `yaml:"smtp"`
	Kindlegen     string      `yaml:"kindlegen"`
	Proxy         string      `yaml:"proxy"`
}

type smtpConfig struct {
	SenderAddress string `yaml:"senderAddress"`
	HostAndPort   string `yaml:"hostAndPort"`
	SSL           bool   `yaml:"ssl"`
	Account       string `yaml:"account"`
	Password      string `yaml:"password"`
}
