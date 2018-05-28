# ptail
Read the file from the end with golang and execute the hook

# usage
```golang
package main

import (
	"fmt"

	"github.com/pyama86/ptail/ptail"
)

func main() {
	p := ptail.NewPtail("/path/to/example.log", 100)
	cnt := 0

	p.Use(func(l []byte) error {
		fmt.Println(string(l))
		return nil
	})

	p.Use(func(l []byte) error {
		cnt++
		return nil
	})

	p.Execute()
	fmt.Println(cnt)
}
```

## Middleware Function
```golang
type middlewareFunc func([]byte) error
```

# author
@pyama86
