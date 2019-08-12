# shutdown

package `github.com/iamken1204/shutdown` offers a simple way to watch OS termination signals and execute cleanup function.

## Install

`go get -u github.com/iamken1204/shutdown`

## Usage

```go
package main

import (
	"fmt"

	"github.com/iamken1204/shutdown"
)

func main() {
	fmt.Println("Starting the app.")

	srv := myServer{}
	s := shutdown.New()

	h := shutdown.NewHook(func() error {
		fmt.Println("executing my hook...")
		return nil
	})
	s.AddHook(
		h,
		shutdown.NewHook(func() error {
			fmt.Println("hook the 2nd is running")
			return nil
		}),
		srv,
	)

	s.Listen()
	fmt.Println("The app has been shut down.")
}

type myServer struct{}

func (myServer) Cleanup() error {
	fmt.Println("terminating myServer")
	return nil
}
```