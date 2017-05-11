package declaration

import (
	"fmt"
	"github.com/creichlin/goschema"
	"github.com/creichlin/gutil"
	"github.com/mitchellh/mapstructure"
	"sort"
)

type Service struct {
	Executable string
	WorkingDir string `mapstructure:"working-dir"`
	Arguments  []string
}

type FSTrigger struct {
	Path     string
	Services []string
}

type Stats struct {
	File    string
	Seconds int
}

type Root struct {
	Stats      *Stats
	Services   map[string]*Service
	FSTriggers map[string]*FSTrigger `mapstructure:"fs-triggers"`
}

func Doc() string {
	return goschema.Doc(buildSchema())
}

func Parse(data interface{}) (*Root, error) {
	schema := buildSchema()
	errors := goschema.ValidateGO(schema, data)
	if errors.Has() {
		return nil, errors
	}

	root := &Root{}
	err := mapstructure.Decode(data, root)
	errors = validate(root)
	if errors.Has() {
		return nil, errors
	}
	return root, err
}

func validate(r *Root) *gutil.ErrorCollector {
	ec := gutil.NewErrorCollector()

	keys := []string{}
	for name, _ := range r.FSTriggers {
		keys = append(keys, name)
	}
	sort.Strings(keys)

	for _, name := range keys {
		for _, serviceName := range r.FSTriggers[name].Services {
			if _, contains := r.Services[serviceName]; !contains {
				ec.Add(fmt.Errorf("fs-trigger %v has unknown service %v as target", name, serviceName))
			}
		}
	}

	return ec
}

func buildSchema() goschema.Type {
	return goschema.NewObjectType("Pentaconta service declaration", func(o goschema.ObjectType) {
		o.Map("services", "Service definitions", func(m goschema.MapType) {
			m.Object("Services start executables and restart them when terminated/crashed", func(o goschema.ObjectType) {
				o.String("executable", "Path to the executable, if not absolute will use PATH env var")
				o.String("working-dir", "Path to the working dir, absolute or will use path to current working dir").Optional()
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
		o.Object("stats", "Statistics writer", func(o goschema.ObjectType) {
			o.Int("seconds", "Over how many seconds stats are printed").Min(5).Max(60)
			o.String("file", "Where to write the stats to")
		}).Optional()
	})
}
