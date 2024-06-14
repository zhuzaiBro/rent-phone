package sys_config

type SysConfig struct {
	ServiceList `yaml:"service_list" json:"service_list"`
}

type ServiceList []string
