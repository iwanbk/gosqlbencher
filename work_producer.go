package main

import (
	"context"
	"fmt"
)

type workProducer struct {
}

const (
	integerDataType = "integer"
	stringDataType  = "string"
)

const (
	sequentialGenType = "sequential"
)

func (wp *workProducer) run(ctx context.Context, num int, qps []queryParam) <-chan []interface{} {
	resCh := make(chan []interface{}, 1000)
	go func() {
		defer close(resCh)
		for i := 1; i <= num; i++ {
			// generate the params
			var params []interface{}
			for _, qp := range qps {
				var param interface{}
				switch qp.DataType {
				case integerDataType:
					param = i
				case stringDataType:
					param = fmt.Sprintf("%s%d", qp.Prefix, i)
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
