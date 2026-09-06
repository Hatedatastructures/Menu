package config

type Develop struct {
  IsCasbinMode bool `mapstructure:"is_casbin_mode" json:"is_casbin_mode" yaml:"is_casbin_mode"`
}