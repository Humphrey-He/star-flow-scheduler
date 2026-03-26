package redisx

type Config struct {
	Addr           string `json:",optional"`
	Password       string `json:",optional"`
	DB             int    `json:",optional"`
	PoolSize       int    `json:",optional"`
	DialTimeoutMs  int64  `json:",optional"`
	ReadTimeoutMs  int64  `json:",optional"`
	WriteTimeoutMs int64  `json:",optional"`
}
