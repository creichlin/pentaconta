package declaration

type Service struct {
	Executable string
	Arguments  []string
}

type FSTrigger struct {
	Path     string
	Services []string
}

type Root struct {
	Services   map[string]*Service
	FSTriggers map[string]*FSTrigger `mapstructure:"fs-triggers"`
}
