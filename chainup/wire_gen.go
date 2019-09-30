// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package chainup

import (
	"chainup.dev/chainup/ansible"
	"chainup.dev/chainup/database"
	"chainup.dev/chainup/database/transaction"
	"chainup.dev/chainup/infrastructure"
	"chainup.dev/chainup/provision"
	"chainup.dev/chainup/statemachine/middleware"
	"chainup.dev/chainup/terraform"
	"chainup.dev/lib/log"
	"testing"
)

// Injectors from inject_database.go:

func SetupDatabaseApp() (*App, func(), error) {
	provider := ProvideFileConfigProvider()
	config := ProvideConfig(provider)
	databaseConfig := config.Database
	logConfig := config.Log
	db, cleanup, err := database.ProvideDB(databaseConfig, logConfig)
	if err != nil {
		return nil, nil, err
	}
	providerSettingsRepository := database.NewProviderSettingsRepository(db)
	serverRepository := database.NewServerRepository(db)
	jobRepository := database.NewJobRepository(db)
	deploymentRepository := database.NewDeploymentRepository(db)
	jobScheduler := provision.NewJobScheduler(db, jobRepository, serverRepository, deploymentRepository)
	terraformConfig := config.Terraform
	terraformTerraform := terraform.ConfigureTerraform(terraformConfig)
	serverProvisioner := provision.NewServerProvisioner(terraformTerraform, serverRepository)
	stepProvisionServer := provision.NewStepProvisionServer(serverProvisioner, jobRepository)
	ansibleConfig := config.Ansible
	ansibleAnsible := ansible.ConfigureAnsible(ansibleConfig)
	deploymentProvisioner := provision.NewDeploymentProvisioner(ansibleAnsible, deploymentRepository)
	stepProvisionDeployment := provision.NewStepProvisionDeployment(deploymentProvisioner, jobRepository)
	transactional := middleware.NewTransactional(db)
	jobStateMachine := provision.ConfigureJobStateMachine(stepProvisionServer, stepProvisionDeployment, transactional)
	serverDestroyer := provision.NewServerDestroyer(terraformTerraform, db, serverRepository, deploymentRepository)
	provisioner := provision.NewProvisioner(jobStateMachine, jobScheduler, terraformTerraform, serverDestroyer)
	consoleLogger := log.NewConsoleLogger(logConfig)
	app := NewApp(config, providerSettingsRepository, serverRepository, jobScheduler, provisioner, consoleLogger)
	return app, func() {
		cleanup()
	}, nil
}

// Injectors from inject_memory.go:

func SetupInMemoryApp() *App {
	provider := ProvideFileConfigProvider()
	config := ProvideConfig(provider)
	inMemoryProviderSettingsRepository := infrastructure.NewInMemoryProviderSettingsRepository()
	inMemoryServerRepository := infrastructure.NewInMemoryServerRepository()
	inMemoryTxContext := transaction.NewInMemoryTransactionContext()
	inMemoryJobRepository := provision.NewInMemoryJobRepository()
	inMemoryDeploymentRepository := infrastructure.NewInMemoryDeploymentRepository()
	jobScheduler := provision.NewJobScheduler(inMemoryTxContext, inMemoryJobRepository, inMemoryServerRepository, inMemoryDeploymentRepository)
	terraformConfig := config.Terraform
	terraformTerraform := terraform.ConfigureTerraform(terraformConfig)
	serverProvisioner := provision.NewServerProvisioner(terraformTerraform, inMemoryServerRepository)
	stepProvisionServer := provision.NewStepProvisionServer(serverProvisioner, inMemoryJobRepository)
	ansibleConfig := config.Ansible
	ansibleAnsible := ansible.ConfigureAnsible(ansibleConfig)
	deploymentProvisioner := provision.NewDeploymentProvisioner(ansibleAnsible, inMemoryDeploymentRepository)
	stepProvisionDeployment := provision.NewStepProvisionDeployment(deploymentProvisioner, inMemoryJobRepository)
	transactional := middleware.NewTransactional(inMemoryTxContext)
	jobStateMachine := provision.ConfigureJobStateMachine(stepProvisionServer, stepProvisionDeployment, transactional)
	serverDestroyer := provision.NewServerDestroyer(terraformTerraform, inMemoryTxContext, inMemoryServerRepository, inMemoryDeploymentRepository)
	provisioner := provision.NewProvisioner(jobStateMachine, jobScheduler, terraformTerraform, serverDestroyer)
	logConfig := config.Log
	consoleLogger := log.NewConsoleLogger(logConfig)
	app := NewApp(config, inMemoryProviderSettingsRepository, inMemoryServerRepository, jobScheduler, provisioner, consoleLogger)
	return app
}

// Injectors from inject_testing.go:

func SetupTestApp(t *testing.T) *App {
	provider := ProvideTestConfigProvider()
	config := ProvideConfig(provider)
	inMemoryProviderSettingsRepository := infrastructure.NewInMemoryProviderSettingsRepository()
	inMemoryServerRepository := infrastructure.NewInMemoryServerRepository()
	inMemoryTxContext := transaction.NewInMemoryTransactionContext()
	inMemoryJobRepository := provision.NewInMemoryJobRepository()
	inMemoryDeploymentRepository := infrastructure.NewInMemoryDeploymentRepository()
	jobScheduler := provision.NewJobScheduler(inMemoryTxContext, inMemoryJobRepository, inMemoryServerRepository, inMemoryDeploymentRepository)
	terraformConfig := config.Terraform
	terraformTerraform := terraform.ConfigureTerraform(terraformConfig)
	serverProvisioner := provision.NewServerProvisioner(terraformTerraform, inMemoryServerRepository)
	stepProvisionServer := provision.NewStepProvisionServer(serverProvisioner, inMemoryJobRepository)
	ansibleConfig := config.Ansible
	ansibleAnsible := ansible.ConfigureAnsible(ansibleConfig)
	deploymentProvisioner := provision.NewDeploymentProvisioner(ansibleAnsible, inMemoryDeploymentRepository)
	stepProvisionDeployment := provision.NewStepProvisionDeployment(deploymentProvisioner, inMemoryJobRepository)
	transactional := middleware.NewTransactional(inMemoryTxContext)
	jobStateMachine := provision.ConfigureJobStateMachine(stepProvisionServer, stepProvisionDeployment, transactional)
	serverDestroyer := provision.NewServerDestroyer(terraformTerraform, inMemoryTxContext, inMemoryServerRepository, inMemoryDeploymentRepository)
	provisioner := provision.NewProvisioner(jobStateMachine, jobScheduler, terraformTerraform, serverDestroyer)
	testingLogger := log.NewTestingLogger(t)
	app := NewApp(config, inMemoryProviderSettingsRepository, inMemoryServerRepository, jobScheduler, provisioner, testingLogger)
	return app
}
