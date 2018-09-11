package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/iwanbk/gosqlbencher/query"
	"github.com/stretchr/testify/require"
)

func TestWorkProducerBasic(t *testing.T) {
	const (
		prefix = "name_"
	)
	qps := []query.Arg{
		{
			DataType: query.DataTypeInteger,
			GenType:  query.GenTypeSequential,
		},
		{
			DataType: query.DataTypeString,
			GenType:  query.GenTypeSequential,
			Prefix:   prefix,
		},
	}
	var (
		ctx          = context.Background()
		wp           = newArgsProducer()
		numWorks     = 100
		lenQueryArgs = len(qps)
		num          int
	)

	argsCh := wp.run(ctx, numWorks, qps)
	for args := range argsCh {
		num++
		require.Len(t, args, lenQueryArgs)
		require.Equal(t, num, args[0])
		require.Equal(t, fmt.Sprintf("%s%d", prefix, num), args[1])
	}
}
