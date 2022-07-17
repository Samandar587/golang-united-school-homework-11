package batch

import (
	"fmt"
	"sync"
	"time"
)

type user struct {
	ID int64
}

func getOne(id int64) user {
	time.Sleep(time.Millisecond * 100)
	return user{
		ID: id,
	}
}

func getBatch(n int64, pool int64) (res []user) {

	var wg sync.WaitGroup

	res = make([]user, 0, n)
	workSize := n / pool

	counter := SafeCounter{v: make(map[string]int)}
	counter.v["userId"] = 0

	for i := 1; i <= int(pool); i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			for j := 1; j <= int(workSize); j++ {

				if counter.Value("userId") == int(n) {
					return
				}
				res = append(res, getOne(int64(counter.GetNewUserId("userId"))))
			}
		}()
	}

	wg.Wait()

	return
}

type SafeCounter struct {
	mu sync.Mutex
	v  map[string]int
}

func (c *SafeCounter) GetNewUserId(key string) int {
	c.mu.Lock()
	defer c.mu.Unlock()
	temp := c.v[key]
	c.v[key]++
	return temp
}

func (c *SafeCounter) Inc(key string) {
	c.mu.Lock()
	fmt.Println(c.v[key])
	c.v[key]++
	c.mu.Unlock()
}
func (c *SafeCounter) Value(key string) int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.v[key]
}
