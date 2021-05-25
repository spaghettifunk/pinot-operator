package controllers

import "github.com/spaghettifunk/pinot-operator/controllers/pinot"

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, pinot.Add)
}
