package cache_test

import (
	"testing"

	"github.com/berquerant/cache"
)

func TestInfinity(t *testing.T) {
	t.Run("scenario", func(t *testing.T) {
		runner := &testRunner[int, string]{
			source: &testStringIntSource{},
			newCache: func(x cache.Source[int, string]) (cache.Cache[int, string], error) {
				return cache.NewInfinity(x), nil
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
				{
					title:     "3rd miss(3)",
					arg:       3,
					wantArgs:  []int{1, 2, 3},
					wantValue: "3",
				},
				{
					title:     "4th hit(1)",
					arg:       1,
					wantArgs:  []int{1, 2, 3},
					wantValue: "1",
				},
				{
					title:     "5th hit(3)",
					arg:       3,
					wantArgs:  []int{1, 2, 3},
					wantValue: "3",
				},
			},
		}

		runner.test(t)
	})

	t.Run("random", func(t *testing.T) {
		randomRunner := &randomTestRunner{
			n:        256,
			minValue: 0,
			maxValue: 10,
			newCache: func(f cache.Source[int, int]) (cache.Cache[int, int], error) {
				return cache.NewFIFO(5, f)
			},
		}

		randomRunner.test(t)
	})
}
