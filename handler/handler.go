package handler

import (
    "context"
    "fmt"
    "github.com/redis/go-redis/v9"
    "github.com/nats-io/nats.go"
)

var ctx = context.Background()
var redisClient *redis.Client
var natsConn *nats.Conn

func InitConnections() {
    redisClient = redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })
    var err error
    natsConn, err = nats.Connect(nats.DefaultURL)
    if err != nil {
        panic(err)
    }
    fmt.Println("Connected to Redis and NATS successfully!")
}

func RegisterFunction(name string, code string) {
    luaCode := fmt.Sprintf("return function(arg) %s return arg end", code)
    err := redisClient.Set(ctx, name, luaCode, 0).Err()
    if err != nil {
        fmt.Println("Error storing function in Redis:", err)
    } else {
        fmt.Println("Function", name, "stored in Redis with code:", luaCode)
    }
}

func PublishMessage(subject, msg string) {
    if err := natsConn.Publish(subject, []byte(msg)); err != nil {
        fmt.Println("Failed to publish message:", err)
    } else {
        fmt.Println("Message published to", subject)
    }
}

func DeregisterFunction(name string) {
    deleted, err := redisClient.Del(ctx, name).Result()
    if err != nil {
        fmt.Println("Error deleting function:", err)
    } else if deleted == 0 {
        fmt.Println("Function", name, "not found in Redis")
    } else {
        fmt.Println("Function", name, "deleted from Redis")
    }
}

func CheckFunctionExists(name string) bool {
    exists, err := redisClient.Exists(ctx, name).Result()
    if err != nil {
        fmt.Println("Error checking function in Redis:", err)
        return false
    }
    return exists > 0
}

func CallFunction(name string, arg string) {
    code, err := redisClient.Get(ctx, name).Result()
    if err == redis.Nil {
        fmt.Println("Function", name, "not found in Redis")
        return
    } else if err != nil {
        fmt.Println("Error retrieving function from Redis:", err)
        return
    }
    script := fmt.Sprintf("%s('%s')", code, arg)
    result, err := redisClient.Eval(ctx, script, []string{}).Result()
    if err != nil {
        fmt.Println("Error executing function:", err)
        return
    }
    fmt.Println("Execution result:", result)
}

func SubscribeInvoke() {
    _, err := natsConn.QueueSubscribe("invoke.>", "function_workers", func(msg *nats.Msg) {
        functionName := msg.Subject[len("invoke."):]
        arg := string(msg.Data)
        fmt.Println("Received request to invoke function:", functionName, "with arg:", arg)
        CallFunction(functionName, arg)
    })
    if err != nil {
        fmt.Println("NATS subscription error:", err)
        return
    }
    select {}
}

func ListFunctions() {
    keys, err := redisClient.Keys(ctx, "*").Result()
    if err != nil {
        fmt.Println("Error retrieving function list:", err)
        return
    }
    fmt.Println("Stored Functions:", keys)
}
