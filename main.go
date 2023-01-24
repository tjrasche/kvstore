package main

import (
	"github.com/gin-gonic/gin"
	"rasche-thalhofer.cloud/kv_celonis/membership"
	"rasche-thalhofer.cloud/kv_celonis/store"
)

func main() {
	app := newApp(gin.Default(), membership.NewK8sHandler(), store.NewStore())
	app.Run()
}
