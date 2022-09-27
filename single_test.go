package cache_test

import (
	"testing"

	"github.com/berquerant/cache"
)

func TestSingle(t *testing.T) {
	runner := &testRunner[int, string]{
		source: &testStringIntSource{},
		newCache: func(x cache.Source[int, string]) (cache.Cache[int, string], error) {
			return cache.NewSingle(x)
		},
		cases: []*testcase[int, string]{
			{
				title:     "1st miss(1)",
				arg:       1,
				wantArgs:  []int{1},
				wantValue: "1",
			},
			{
				title:     "1st hit(1)",
				arg:       1,
				wantArgs:  []int{1},
				wantValue: "1",
			},
			{
				title:    "source error",
				arg:      -1,
				wantArgs: []int{1},
				wantErr:  errTestStringIntSourceNegative,
			},
			{
				title:     "2nd miss(2)",
				arg:       2,
				wantArgs:  []int{1, 2},
				wantValue: "2",
			},
		},
	}

	t.Run("scenario", runner.test)

	randomRunner := &randomTestRunner{
		n:        256,
		minValue: 0,
		maxValue: 10,
		newCache: func(f cache.Source[int, int]) (cache.Cache[int, int], error) {
			return cache.NewSingle(f)
		},
	}

	t.Run("random", randomRunner.test)
}
