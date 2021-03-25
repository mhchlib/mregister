package memory

import (
	"github.com/mhchlib/mregister/plugin"
)

func init() {
	plugin.RegisterRegisterPlugin("memory", newMemoryRegister)
}
