package controllers

import "github.com/spaghettifunk/pinot-operator/controllers/tenant"

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, tenant.Add)
}
