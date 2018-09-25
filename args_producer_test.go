package main

import (
	"context"
	"fmt"
	"testing"
	"time"

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
			Min:      10,
		},
		{
			DataType: query.DataTypeString,
			GenType:  query.GenTypeSequential,
			Prefix:   prefix,
			Min:      20,
		},
		{
			DataType: query.DataTypeTime,
			GenType:  query.GenTypeSequential,
		},
	}
	var (
		ctx          = context.Background()
		wp           = newArgsProducer()
		numWorks     = 100
		lenQueryArgs = len(qps)
		num          int
	)

	start := time.Now()
	argsCh := wp.run(ctx, numWorks, qps)
	for args := range argsCh {
		num++
		require.Len(t, args, lenQueryArgs)
		require.Equal(t, num+qps[0].Min, args[0])
		require.Equal(t, fmt.Sprintf("%s%d", prefix, num+qps[1].Min), args[1])

		generatedTime := args[2].(time.Time)
		require.True(t, generatedTime.After(start))
		require.True(t, generatedTime.Before(time.Now()))
	}
}

func TestWorkProducerRandom(t *testing.T) {
	qps := []query.Arg{
		{
			DataType: query.DataTypeInteger,
			GenType:  query.GenTypeRandom,
			Min:      10,
			Max:      20,
		},
		{
			DataType: query.DataTypeString,
			GenType:  query.GenTypeRandom,
			Prefix:   "name_",
			Min:      20,
			Max:      30,
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
		require.True(t, randomNum >= qps[0].Min)
		require.True(t, randomNum <= qps[0].Max)

		// make sure the generated string has this format : prefix+randomNumber
		format := qps[1].Prefix + "%d"
		_, err := fmt.Sscanf(args[1].(string), format, &randomNum)
		require.NoError(t, err)
		require.True(t, randomNum >= qps[1].Min)
		require.True(t, randomNum <= qps[1].Max)
	}
}
