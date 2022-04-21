package html

import (
	_ "embed"
	"github.com/stretchr/testify/assert"
	"strconv"
	"sync"
	"testing"
)

//go:embed testdata/template006/index.html
var poolTemplate string

//go:embed testdata/template006/expect.html
var poolExpected string

func TestNewPool(t *testing.T) {
	testcases := []struct {
		description string
		poolSize    int
		bufferSize  int
	}{
		{
			description: "pool size 1",
			poolSize:    1,
			bufferSize:  8192,
		},
		{
			description: "pool size 10",
			poolSize:    10,
			bufferSize:  8192,
		},
		{
			description: "pool size 100",
			poolSize:    100,
			bufferSize:  8192,
		},
	}

	for _, testcase := range testcases {
		vdom, err := New(poolTemplate, BufferSize(testcase.bufferSize))
		if !assert.Nil(t, err, testcase.description) {
			continue
		}

		pool := NewPool(testcase.poolSize, vdom)

		goroutines := 100
		wg := sync.WaitGroup{}
		wg.Add(goroutines)
		for i := 0; i < goroutines; i++ {
			go func() {
				dom := pool.New()
				defer func() {
					pool.Put(dom)
					wg.Done()
				}()

				attrIt := dom.SelectAttributes("img", "src")
				counter := 0
				for attrIt.Has() {
					attribute, _ := attrIt.Next()
					attribute.Set("abcdef" + strconv.Itoa(counter))
					counter++
				}
				html := dom.Render()
				assert.Equal(t, poolExpected, html, testcase.description)
			}()
		}
		wg.Wait()
	}
}
