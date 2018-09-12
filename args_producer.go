package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/iwanbk/gosqlbencher/query"
)

// argsProducer defines the query arguments generator
type argsProducer struct {
}

// newArgsProducer creates new argsProducer object
func newArgsProducer() *argsProducer {
	rand.Seed(time.Now().UnixNano())

	return &argsProducer{}
}

// run runs this producer in a goroutine and send the generated arguments
// in the returned channel.
// the channel will be closed after it finished generated the specified number of arguments
func (ap *argsProducer) run(ctx context.Context, num int, args []query.Arg) <-chan []interface{} {
	resCh := make(chan []interface{}, 1000)
	go func() {
		defer close(resCh)

		// generate the params
		for i := 1; i <= num; i++ {
			var params []interface{}
			for _, arg := range args {
				param := ap.generateArg(i, arg)
				params = append(params, param)
			}

			// publish it
			select {
			case <-ctx.Done():
				return
			default:
				resCh <- params
			}
		}
	}()
	return resCh
}

func (ap *argsProducer) generateArg(i int, qa query.Arg) interface{} {
	if qa.GenType == query.GenTypeRandom {
		return ap.generateRandomArg(i, qa)
	}
	return ap.generateSeqArg(i, qa)
}

// generate random argument
func (ap *argsProducer) generateRandomArg(i int, qa query.Arg) interface{} {
	var param interface{}
	randNum := rand.Intn(qa.RandomRangeMax-qa.RandomRangeMin+1) + qa.RandomRangeMin
	switch qa.DataType {
	case query.DataTypeInteger:
		param = randNum
	case query.DataTypeString:
		param = fmt.Sprintf("%s%v", qa.Prefix, randNum)
	}
	return param
}

// generateSeqArg generates sequential arg
func (ap *argsProducer) generateSeqArg(i int, qa query.Arg) interface{} {
	var param interface{}
	switch qa.DataType {
	case query.DataTypeInteger:
		param = i
	case query.DataTypeString:
		param = fmt.Sprintf("%s%v", qa.Prefix, i)
	}
	return param
}
