package cache_test

import (
	"errors"
	"fmt"
	"math/rand"
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

// scenario test
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

type randomTestRunner struct {
	n                  int
	minValue, maxValue int
	newCache           func(cache.Source[int, int]) (cache.Cache[int, int], error)
}

func (r *randomTestRunner) rand() int { return r.minValue + rand.Int()%(r.maxValue+1) }

func (r *randomTestRunner) test(t *testing.T) {
	c, err := r.newCache(func(v int) (int, error) { return v, nil }) // identity
	if err != nil {
		t.Fatal(err)
	}

	rand.Seed(42)
	for i := 0; i < r.n; i++ {
		val := r.rand()
		t.Logf("n = %d val = %d\n%v", i, val, c)
		got, err := c.Get(val)
		t.Logf("got = %d err = %v\n%v", got, err, c)
		if err != nil {
			t.Errorf("err %v", err)
		}
		if got != val {
			t.Errorf("got %d", got)
		}
	}
}
