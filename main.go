package main

import (
    "fmt"
    "time"
    "faas-proyecto/handler"
)

func main() {
    handler.InitConnections()
    handler.ListFunctions()

    handler.DeregisterFunction("exampleFunc")

    exists := handler.CheckFunctionExists("exampleFunc")
    if exists {
        fmt.Println("Function still exists in Redis. Deletion failed!")
    } else {
        fmt.Println("Function deleted successfully!")
    }

    handler.RegisterFunction("exampleFunc", "redis.call('SET', 'output', arg)")
    fmt.Println("Function stored in Redis")

    handler.SubscribeInvoke()

    time.Sleep(2 * time.Second)

    handler.PublishMessage("invoke.exampleFunc", "This is a test message")

    select {}
}
