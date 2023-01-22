package main

import (
	"github.com/gin-gonic/gin"
	"rasche-thalhofer.cloud/kv_celonis/membership"
	"rasche-thalhofer.cloud/kv_celonis/store"
)

type app struct {
	r  *gin.Engine
	mh membership.Handler
	s  store.Store
}

func newApp(e *gin.Engine, mh membership.Handler, s store.Store) *app {

	app := &app{r: e, mh: mh, s: s}

	app.addMembershipRoutes()
	app.addStoreRoutes()
	return app
}

func (a app) Run() {
	defer a.r.Run(":8080")
	a.r.GET("ping", func(ctx *gin.Context) {
		ctx.Status(200)
	})

	a.r.PUT("values")
}

func (a app) addMembershipRoutes() {
	// @Todo implement

}

func (a app) addStoreRoutes() {
	sr := a.r.Group("values")

	// get a value
	sr.GET("/:key", func(ctx *gin.Context) {
		key, found := a.s.Get(ctx.Param("key"))
		if found {
			ctx.String(200, key)
		} else {
			ctx.Status(404)
		}
	})

	sr.DELETE("/:key", func(ctx *gin.Context) {
		a.s.Delete(ctx.Param("key"))
	})

	sr.PUT("/:key", func(ctx *gin.Context) {
		raw, err := ctx.GetRawData()
		if err != nil {
			ctx.Error(err)
			ctx.Status(422)
		}
		key := ctx.Param("key")
		val := string(raw)

		a.s.Put(key, val)
		ctx.Status(201)
	})
}
