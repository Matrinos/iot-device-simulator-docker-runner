package main

type SimulatorStatus string

const (
	Running       SimulatorStatus = "running"
	Idle          SimulatorStatus = "idle"
	Uninitialized SimulatorStatus = "uninitialized"
	Error         SimulatorStatus = "error"
)

type WorkflowAlias string
