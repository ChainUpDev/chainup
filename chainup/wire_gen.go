// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package chainup

import (
	"chainup.dev/chainup/infrastructure"
)

// Injectors from inject_memory.go:

func SetupInMemoryApp() *App {
	stateMachine := infrastructure.ConfigureStateMachine()
	provisioner := NewProvisioner(stateMachine)
	inMemoryServerRepository := infrastructure.NewInMemoryServerRepository()
	app := NewApp(provisioner, inMemoryServerRepository)
	return app
}
