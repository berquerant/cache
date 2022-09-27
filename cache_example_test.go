package cache_test

import (
	"fmt"

	"github.com/berquerant/cache"
)

func ExampleCache() {
	c, err := cache.NewLRU(3, func(x int) (int, error) {
		fmt.Printf("src: %d\n", x)
		return x * x, nil
	})
	if err != nil {
		panic(err)
	}
	get := func(x int) {
		y, _ := c.Get(x)
		fmt.Println(y)
	}

	for _, x := range []int{
		1, 1, 2, 3, 4, 3, 1,
	} {
		get(x)
	}

	fmt.Println(c.Hit(), c.Miss(), c.Size())
	// Output:
	// src: 1
	// 1
	// 1
	// src: 2
	// 4
	// src: 3
	// 9
	// src: 4
	// 16
	// 9
	// src: 1
	// 1
	// 2 5 3
}
