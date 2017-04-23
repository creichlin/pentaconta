package declaration

import (
	"github.com/creichlin/goschema"
	"github.com/mitchellh/mapstructure"
)

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

func Parse(data interface{}) (*Root, error) {
	schema := goschema.NewObjectType("Pentaconta service declaration", func(o goschema.ObjectType) {
		o.Map("services", "Service definitions", func(m goschema.MapType) {
			m.Object("Services start executables and restart them when terminated/crashed", func(o goschema.ObjectType) {
				o.String("executable", "Path to the executable, if not absolute will use PATH env var")
				o.List("arguments", "List of arguments", func(l goschema.ListType) {
					l.String("Single argument")
				})
			})
		})
		o.Map("fs-triggers", "Filesystem triggers", func(m goschema.MapType) {
			m.Object("Filesystem watchers trigger termination (and implicit restart) of services", func(o goschema.ObjectType) {
				o.String("path", "Folder or file to watch. Folder must exists. If it's a file, parent folder must exist.")
				o.List("services", "List of service names to restart", func(l goschema.ListType) {
					l.String("Service name")
				})
			})
		}).Optional()
	})

	errors := goschema.ValidateGO(schema, data)
	if errors.Has() {
		return nil, errors
	}

	root := &Root{}
	err := mapstructure.Decode(data, root)
	return root, err
}
