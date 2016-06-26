// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

type Config struct {
	Scdbeat ScdbeatConfig
}

type ScdbeatConfig struct {
	Period string `config:"period"`
	Path   string `config:"path"` //path 환경변수 값 추가.
}
