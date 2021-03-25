package memory

import (
	"github.com/mhchlib/register/plugin"
)

func init() {
	plugin.RegisterRegisterPlugin("memory", newMemoryRegister)
}
