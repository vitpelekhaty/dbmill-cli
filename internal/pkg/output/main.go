package output

import (
	"strings"
)

func init() {
	DefaultScriptsFolderOutput, _ = NewScriptsFolderOutput(strings.NewReader(defaultScriptsFolderOutput))
}
