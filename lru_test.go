package cache_test

import (
	"errors"
	"testing"

	"github.com/berquerant/cache"
)

func TestLRU(t *testing.T) {
	t.Run("invalid size", func(t *testing.T) {
		_, err := cache.NewLRU(1, func(_ string) (int, error) { return 1, nil })
		if !errors.Is(err, cache.ErrInvalidSize) {
			t.Fail()
		}
	})

	t.Run("scenario size 2", func(t *testing.T) {
		runner := &testRunner[int, string]{
			source: &testStringIntSource{},
			newCache: func(x cache.Source[int, string]) (cache.Cache[int, string], error) {
				return cache.NewLRU(2, x)
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
				// latest <= 2 1
				{
					title:     "3rd miss(3)",
					arg:       3,
					wantArgs:  []int{1, 2, 3},
					wantValue: "3",
				},
				// latest <= 3 2
				{
					title:     "2nd hit(2)",
					arg:       2,
					wantArgs:  []int{1, 2, 3},
					wantValue: "2",
				},
				// latest <= 2 3
				{
					title:     "4th miss(4)",
					arg:       4,
					wantArgs:  []int{1, 2, 3, 4},
					wantValue: "4",
				},
				//latest <= 4 2
				{
					title:     "5th miss(3)",
					arg:       3,
					wantArgs:  []int{1, 2, 3, 4, 3},
					wantValue: "3",
				},
			},
		}

		runner.test(t)
	})

	t.Run("scenario size 3", func(t *testing.T) {
		runner := &testRunner[int, string]{
			source: &testStringIntSource{},
			newCache: func(x cache.Source[int, string]) (cache.Cache[int, string], error) {
				return cache.NewLRU(3, x)
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
				// latest <= 2 1
				{
					title:     "3rd miss(3)",
					arg:       3,
					wantArgs:  []int{1, 2, 3},
					wantValue: "3",
				},
				// latest <= 3 2 1
				{
					title:     "2nd hit(2)",
					arg:       2,
					wantArgs:  []int{1, 2, 3},
					wantValue: "2",
				},
				// latest <= 2 3 1
				{
					title:     "4th miss(4)",
					arg:       4,
					wantArgs:  []int{1, 2, 3, 4},
					wantValue: "4",
				},
				//latest <= 4 2 3
				{
					title:     "5th hit(3)",
					arg:       3,
					wantArgs:  []int{1, 2, 3, 4},
					wantValue: "3",
				},
				// latest <= 3 4 2
				{
					title:     "6th miss(6)",
					arg:       6,
					wantArgs:  []int{1, 2, 3, 4, 6},
					wantValue: "6",
				},
			},
		}

		runner.test(t)
	})
}
