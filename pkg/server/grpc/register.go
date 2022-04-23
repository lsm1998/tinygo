package grpcServer

import (
	"google.golang.org/grpc"
)

const defaultRouterName = "defaultRouterName"

func AddComponent(register Register) {
	registers[defaultRouterName] = append(registers[defaultRouterName], register)
}

var registers = make(map[string][]Register)

type Register func(server *grpc.Server)
