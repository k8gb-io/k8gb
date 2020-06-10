package controller

import (
	"github.com/AbsaOSS/ohmyglb/pkg/controller/gslb"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, gslb.Add)
}
