package serverMsg

import (
	"github.com/asynkron/protoactor-go/actor"
)

//#############################3messages for Selector###################################################3

//Server -----> Selector (server pings Selector to Start FL Round by sending the avavilable device count to Coordianator)
type SelectorSendCountToCoordinator struct {
	CoordinatorPid *actor.PID
}

// Selector -----> Coordinator(Selector sends count of devices available)
type SelectorInfoMsg struct {
	Devices_available int
	SelectorPid       *actor.PID
}

//Selector ---> Coordinator
type DevicesAssignmentSuccess struct{}

// Coordinator -----> Selector (if not enough devices available)
type NotEnoughDevicesAvailable struct{}

//Coordinator ------> Selector(once Aggregator are spawned, Selector Assign Trainers to Each Aggregator)
type AssignDevicesToAggr struct {
	NumAggr        int
	DevicesPerAggr int
	AggrPidList    []*actor.PID
}

//#############################messages for Coordinator###################################################3

//Coordianator -------> Master Aggregator
type SpawnNewAggr struct {
	NumAggr int
}

//Coordinator ---> MAster Aggr ---> Aggr ---> Trainer(pass model)
type ModelCarryingMessage struct {
	Model string
}

//#############################messages for Master Aggregator###################################################3

//Master Aggregator ----> Coordinator(aggr spawned successfully)
type AggrSpawnSuccessful struct {
	AggrPidList []*actor.PID
}

//#############################messages for Aggregator###################################################3

type ReceiveDevices struct {
	DevicesPidList []*actor.PID
}

//#############################messages for Trainer###################################################3

//Trainer ----> Selector(send availability status)
type TrainerPingtoSelector struct {
}

//Server ----> Trainer(to start pinging to Selector)
type ServerPingtoTrainer struct {
	SelectorPid *actor.PID
}

//Trainer ---> Aggr ---> Master Aggr ---> Coordinator
type ModelTrainingCompletion struct {
}

//just a demo msg
type Demo struct {
}
