package declaration

import (
	"github.com/creichlin/gutil"
	"github.com/creichlin/gutil/testin"
	"gopkg.in/yaml.v2"
	"testing"
)

type testCase struct {
	name        string
	declaration map[string]interface{}
	errors      string
}

func TestDeclarationValidation(t *testing.T) {
	testin.RunMapTests(t, "testdata", func(source string, operation string, t *testing.T) string {
		declaration := map[string]interface{}{}
		err := yaml.Unmarshal([]byte(source), declaration)
		if err != nil {
			t.Error(err)
		}

		_, errs := Parse(gutil.ConvertToJSONTree(declaration))
		if errs != nil {
			return errs.Error()
		}
		return "nil"
	})
}
