package models

type MongoCfg struct {
	Mongo struct {
		URI string `yaml:"uri"`
		DB  string `yaml:"db"`
	} `yaml:"mongo"`
}
type JwtCfg struct {
	Jwt struct {
		SecKey     string `yaml:"secret_key"`
		RefreshKey string `yaml:"refresh_key"`
	}
}
