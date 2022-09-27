package cache_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/berquerant/cache"
)

func equalSlice[T comparable](t *testing.T, want, got []T, msg string) {
	t.Helper()
	if len(want) != len(got) {
		t.Errorf("%s len want %d but got %d", msg, len(want), len(got))
		return
	}
	for i, w := range want {
		g := got[i]
		if w != g {
			t.Errorf("%s at index %d, want %v but got %v", msg, i, w, g)
		}
	}
}

type testcase[K comparable, V comparable] struct {
	title     string
	arg       K
	wantArgs  []K
	wantValue V
	wantErr   error
}

type testSource[K comparable, V comparable] interface {
	Call(K) (V, error) // Source[K, V]
	GetArgs() []K      // accumulated arguments of Call
}

type testRunner[K comparable, V comparable] struct {
	cases    []*testcase[K, V]
	source   testSource[K, V]
	newCache func(cache.Source[K, V]) (cache.Cache[K, V], error)
}

func (r *testRunner[K, V]) test(t *testing.T) {
	c, err := r.newCache(r.source.Call)
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range r.cases {
		tc := tc
		t.Run(tc.title, func(t *testing.T) {
			got, err := c.Get(tc.arg)
			t.Log(c)
			equalSlice(t, tc.wantArgs, r.source.GetArgs(), "args")
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Errorf("err %v", err)
				}
			} else if tc.wantValue != got {
				t.Errorf("value %v", got)
			}
		})
	}
}

type testStringIntSource struct {
	args []int
}

var errTestStringIntSourceNegative = errors.New("StringIntSourceNegative")

func (s *testStringIntSource) Call(key int) (string, error) {
	if key < 0 {
		return "", errTestStringIntSourceNegative
	}
	s.args = append(s.args, key)
	return fmt.Sprint(key), nil
}

func (s *testStringIntSource) GetArgs() []int { return s.args }
