package declaration

type Service struct {
	Name       string
	Executable string
	Arguments  []string
}

type Root struct {
	Services []*Service
}
