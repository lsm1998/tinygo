package httpServer

import "github.com/gin-gonic/gin"

const defaultRouterName = "defaultRouterName"

func AddRouter(register Register, routerName ...string) {
	if len(routerName) > 0 {
		registers[routerName[0]] = append(registers[routerName[0]], register)
	}
	registers[defaultRouterName] = append(registers[defaultRouterName], register)
}

var registers = make(map[string][]Register)

type Register func(engine *gin.Engine)
