package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/iwanbk/gosqlbencher/query"
)

type argsProducer struct {
}

func newArgsProducer() *argsProducer {
	rand.Seed(time.Now().UnixNano())

	return &argsProducer{}
}

func (ap *argsProducer) run(ctx context.Context, num int, qas []query.Arg) <-chan []interface{} {
	resCh := make(chan []interface{}, 1000)
	go func() {
		defer close(resCh)
		for i := 1; i <= num; i++ {
			// generate the params
			var params []interface{}
			for _, qa := range qas {
				var param interface{}
				switch qa.DataType {
				case query.DataTypeInteger:
					param = i
				case query.DataTypeString:
					param = fmt.Sprintf("%s%d", qa.Prefix, i)
				}
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
	switch qa.DataType {
	case query.DataTypeInteger:
		param = i
	case query.DataTypeString:
		param = fmt.Sprintf("%s%v", qa.Prefix, i)
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
