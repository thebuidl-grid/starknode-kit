package types

type Process struct {
	Processes []process `yaml:"processes"`
}

type process struct {
	Name ClientType `yaml:"name"`
	Pid  int        `yaml:"pid"`
	Path string     `yaml:"path"`
}
