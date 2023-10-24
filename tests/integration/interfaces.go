package integration

import "testing"

type ServerInstance interface {
	Port(*testing.T) int
	Close(*testing.T)
	Address(*testing.T) string
	Prep(*testing.T)
}
