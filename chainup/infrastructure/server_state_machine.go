package infrastructure

import (
	"chainup.dev/chainup/statemachine"
	"chainup.dev/chainup/statemachine/step"
)

var (
	// StateRequested describes the requested server specification.
	StateRequested = statemachine.NewState("requested")

	// StateProvisioning indicates that the server is being provisioned by the infrastructure provider.
	StateProvisioning = statemachine.NewState("provisioning")

	// StateReady is the terminating state representing a successful server provisioning.
	StateReady = statemachine.NewState("completed").Successful()

	// StateFailed is the terminating state representing provisioning server failure.
	// @TODO: Add failure message to job somewhere.
	StateFailed = statemachine.NewState("failed").Failure()

	// ValidStates of an infrastructure.Server.
	ValidStates = []statemachine.State{StateRequested, StateProvisioning, StateReady, StateFailed}
)

// ConfigureStateMachine returns a preconfigured StateMachine
// for running server provisioning.
func ConfigureStateMachine() *statemachine.StateMachine {
	return statemachine.Builder(ValidStates).
		Step(StateRequested, step.Simple(StateFailed)).
		Build()
}
