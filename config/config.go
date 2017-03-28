package config

import "github.com/go-ini/ini"

type Config struct {
	WebManager WebManager
	File       *ini.File
}

type WebManager struct {
	Enable bool
	Path   string
	User   string
	Pwd    string
}

func NewConfig(fileName string) *Config {
	cfg := new(Config)

	f, err := ini.Load(fileName)
	if err != nil {
		return cfg
	}

	cfg.WebManager.Enable, _ = f.Section("WebManager").Key("enable").Bool()
	cfg.WebManager.Path = f.Section("WebManager").Key("path").String()
	cfg.WebManager.User = f.Section("WebManager").Key("user").String()
	cfg.WebManager.Pwd = f.Section("WebManager").Key("pwd").String()

	return cfg
}
