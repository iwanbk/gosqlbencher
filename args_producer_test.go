package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/iwanbk/gosqlbencher/query"
	"github.com/stretchr/testify/require"
)

func TestWorkProducerSequential(t *testing.T) {
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

func TestWorkProducerRandom(t *testing.T) {
	qps := []query.Arg{
		{
			DataType:       query.DataTypeInteger,
			GenType:        query.GenTypeRandom,
			RandomRangeMin: 10,
			RandomRangeMax: 20,
		},
		{
			DataType:       query.DataTypeString,
			GenType:        query.GenTypeRandom,
			Prefix:         "name_",
			RandomRangeMin: 20,
			RandomRangeMax: 30,
		},
	}
	var (
		ctx          = context.Background()
		wp           = newArgsProducer()
		numWorks     = 10000
		lenQueryArgs = len(qps)
		num          int
	)

	argsCh := wp.run(ctx, numWorks, qps)
	for args := range argsCh {
		num++
		require.Len(t, args, lenQueryArgs)

		// make sure the generated number is in the specified range
		randomNum := args[0].(int)
		require.True(t, randomNum >= qps[0].RandomRangeMin)
		require.True(t, randomNum <= qps[0].RandomRangeMax)

		// make sure the generated string has this format : prefix+randomNumber
		format := qps[1].Prefix + "%d"
		_, err := fmt.Sscanf(args[1].(string), format, &randomNum)
		require.NoError(t, err)
		require.True(t, randomNum >= qps[1].RandomRangeMin)
		require.True(t, randomNum <= qps[1].RandomRangeMax)
	}
}
