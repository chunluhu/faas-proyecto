package main

import (
    "faas-proyecto/handler"
    "github.com/gin-gonic/gin"
    "net/http"
)

func main() {
    handler.InitConnections()
    r := gin.Default()

    r.POST("/func/register", func(c *gin.Context) {
        var req struct {
            Name string `json:"name"`
            Code string `json:"code"`
        }
        if err := c.BindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
            return
        }
        handler.RegisterFunction(req.Name, req.Code)
        c.JSON(http.StatusOK, gin.H{"message": "Function registered successfully"})
    })

    r.POST("/func/invoke", func(c *gin.Context) {
        var req struct {
            Name string   `json:"name"`
            Args []string `json:"args"`
        }
        if err := c.BindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
            return
        }
        result := handler.CallFunction(req.Name, req.Args...)
        c.JSON(http.StatusOK, gin.H{"result": result})
    })

    r.DELETE("/func/delete", func(c *gin.Context) {
        var req struct {
            Name string `json:"name"`
        }
        if err := c.BindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
            return
        }
        handler.DeregisterFunction(req.Name)
        c.JSON(http.StatusOK, gin.H{"message": "Function deleted successfully"})
    })

    r.GET("/func/list", func(c *gin.Context) {
        functions := handler.ListFunctions()
        c.JSON(http.StatusOK, gin.H{"functions": functions})
    })

    r.Run(":8080")
}
