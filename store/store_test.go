package store_test

import (
	"github.com/google/uuid"
	"rasche-thalhofer.cloud/kv_celonis/store"
	"sync"
	"testing"
	"time"
)

func TestParallelWrite(t *testing.T) {
	numberOfVals := 100_000

	// setup test store
	s := store.NewStore()
	// generate random testvalues (using uuid as simple random strings)
	var keyValuePairs []struct{ key, value string }
	for i := 0; i < numberOfVals; i++ {
		uuid1, _ := uuid.NewRandom()
		uuid2, _ := uuid.NewRandom()
		keyValuePairs = append(keyValuePairs, struct{ key, value string }{uuid1.String(), uuid2.String()})
	}
	wg := sync.WaitGroup{}

	// time writes
	ts := time.Now()
	q := 0
	for _, kvPair := range keyValuePairs {
		wg.Add(1)
		q++
		kvPair := kvPair
		go func() {
			s.Put(kvPair.key, kvPair.value)
			wg.Done()
		}()
	}
	wg.Wait()
	tm := time.Now().Sub(ts)

	t.Logf("Took %s ", tm)
	if s.CountElements() != uint64(numberOfVals) {
		t.Fail()

	}

}
