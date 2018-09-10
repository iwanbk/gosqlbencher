package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWorkProducerBasic(t *testing.T) {
	const (
		prefix = "name_"
	)
	qps := []queryParam{
		queryParam{
			DataType: integerDataType,
			GenType:  sequentialGenType,
		},
		queryParam{
			DataType: stringDataType,
			GenType:  sequentialGenType,
			Prefix:   prefix,
		},
	}
	var (
		ctx            = context.Background()
		wp             = &workProducer{}
		numWorks       = 100
		lenQueryParams = len(qps)
		num            int
	)

	argsCh := wp.run(ctx, numWorks, qps)
	for args := range argsCh {
		num++
		require.Len(t, args, lenQueryParams)
		require.Equal(t, num, args[0])
		require.Equal(t, fmt.Sprintf("%s%d", prefix, num), args[1])
	}
}
