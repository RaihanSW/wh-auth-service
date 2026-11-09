package config

import "time"

type ApplicationConfig struct {
	SecretKey   string                    `json:"secretkey"`
	ServiceName string                    `json:"service_name"`
	Port        PortConfig                `json:"port"`
	Database    map[string]DatabaseConfig `json:"database"`
	RedisCfg    RedisConfig               `json:"redis"`
	State       StateConfig               `json:"state"`
	Jwt         JwtConfig                 `json:"jwt"`
	Log         Log                       `json:"log"`
}

type DatabaseConfig struct {
	Driver          string        `json:"driver"`
	Username        string        `json:"username"`
	Password        string        `json:"password"`
	Hostname        string        `json:"hostname"`
	DBname          string        `json:"dbname"`
	Timeout         time.Duration `json:"timeout"`
	MaxOpenConns    int           `json:"maxopenconns"`
	MaxIdleConns    int           `json:"maxidleconns"`
	ConnMaxLifetime time.Duration `json:"connmaxlifetime"`
}

type MovieConfig struct {
	ImageLink   string `json:"imagelink"`
	TrailerLink string `json:"trailerlink"`
}

type PortConfig struct {
	Service         int           `json:"service"`
	Profiler        int           `json:"profiler"`
	ServiceTimeout  time.Duration `json:"servicetimeout"`
	ProfilerTimeout time.Duration `json:"profilertimeout"`
}

type RedisConfig struct {
	Address              string        `json:"address"`  //address - ip
	Password             string        `json:"password"` //password
	Port                 string        `json:"port"`     // port
	DB                   int           `json:"db"`
	MaxIdleConnections   int           `json:"maxidle"`
	IdleTimeout          time.Duration `json:"idletimeout"`
	MaxActiveConnections int           `json:"maxactive"`
}

type StateConfig struct {
	RewriteSID bool `json:"rewrite_sid"`
}

type JwtConfig struct {
	TokenExp    int `json:"tokenexp"`
	RefTokenExp int `json:"reftokenexp"`
}

type Log struct {
	Request          bool `json:"request"`
	SuccessResponse  bool `json:"success_response"`
	FailedResponse   bool `json:"failed_response"`
	ExternalRequest  bool `json:"external_request"`
	ExternalResponse bool `json:"external_response"`
}

type DynamicURLConfig struct {
	URLCafeList       string `json:"url_cafe_list"`
	URLPromoList      string `json:"url_promo_list"`
	URLPromoDetail    string `json:"url_promo_detail"`
	URLCafeMenu       string `json:"url_cafe_menu"`
	URLCafeMenuBundle string `json:"url_cafe_menu_bundle"`
}
