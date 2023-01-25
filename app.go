package main

import (
	"github.com/gin-gonic/gin"
	"os"
	"rasche-thalhofer.cloud/kv_celonis/membership"
	"rasche-thalhofer.cloud/kv_celonis/store"
	"sync"
)

type app struct {
	r  *gin.Engine
	mh membership.Handler
	s  store.Store
}

func newApp(e *gin.Engine, mh membership.Handler, s store.Store) *app {

	app := &app{r: e, mh: mh, s: s}

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

func (a app) addStoreRoutes() {
	sr := a.r.Group("values")

	// get a value
	sr.GET("/:key", func(ctx *gin.Context) {

		key := ctx.Param("key")
		replicatedNodes := a.mh.CalculateReplication(key)
		// we use a slice to collect results from the replicated values
		var returnValues []string

		// find out if we are currently on a node the key is replicated to
		for i, node := range replicatedNodes {
			if node.Id == os.Getenv("POD_IDENTITY") {
				//  we're currnetly running on one of the nodes the value is supposed to be readfrom
				// we remove the current node from the list of nodes for replicating
				replicatedNodes = append(replicatedNodes[:i], replicatedNodes[i+1:]...)
				value, found := a.s.Get(key)
				// if value not found everything is broken
				if found {
					returnValues = append(returnValues, value)
				} else {
					panic("should not happen")
				}
				// we can only run on a single node at once so we can break here
				break
			}
		}

		// mutex and waitgroup so we can write into result in parallel
		mut := sync.Mutex{}
		wg := sync.WaitGroup{}
		for _, node := range replicatedNodes {
			wg.Add(1)
			go func(key string, node membership.Node) {
				defer wg.Done()
				result, err := node.Get(key)
				if err != nil {
					println("not replicated yet")
				}
				mut.Lock()
				returnValues = append(returnValues, result)
				mut.Unlock()

			}(key, node)

		}
		wg.Wait()

		// check if returned values are the same
		var lastVal = ""
		for _, val := range returnValues {
			if lastVal != val {
				ctx.Status(500)
			}
			lastVal = val
		}

		if lastVal != "" {
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

		// find the 3 nodes which to replicate to
		replicatedNodes := a.mh.CalculateReplication(key)

		// find out if we are currently on that node
		for i, node := range replicatedNodes {
			if node.Id == os.Getenv("POD_IDENTITY") {
				//  we're currnetly running on one of the nodes the value is supposed to be replicated too
				// we remove the current node from the list of nodes for replicating
				replicatedNodes = append(replicatedNodes[:i], replicatedNodes[i+1:]...)

				// and put it in our local store instead
				a.s.Put(key, val)
				// we can only run on a single node at once so we can rbeak here
				break
			}
		}

		// replicate the value to the remaining nodes
		for _, node := range replicatedNodes {
			err := node.Replicate(key, val)
			if err != nil {
				// obviously here should be some more soffisticated error handling
				// probably even transactional behaviour
				panic(err.Error())
			}
		}
		ctx.Status(201)
	})

	// internal endpoint should only be used by the kv store itself for replication
	sr.PUT("/replicate/:key", func(ctx *gin.Context) {
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
