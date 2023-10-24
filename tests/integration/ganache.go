package integration

import (
	"github.com/testcontainers/testcontainers-go"
	"testing"
)

type GanacheServer struct {
	instance testcontainers.Container
}

func (g GanacheServer) Port(t *testing.T) int {
	//TODO implement me
	panic("implement me")
}

func (g GanacheServer) Close(t *testing.T) {
	//TODO implement me
	panic("implement me")
}

func (g GanacheServer) Address(t *testing.T) string {
	//TODO implement me
	panic("implement me")
}

func (g GanacheServer) Prep(t *testing.T) {
	//TODO implement me
	panic("implement me")
}
