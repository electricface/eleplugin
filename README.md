# eleplugin

example

file hello.go, pkg path is `ele-test/goplugin/plugin_hello`:
```go
package plugin_hello

import "fmt"

func PluginStart() {
	fmt.Println("hello")
}
```

file `use_eleplugin.go`:

```go
package main

import (
	"log"

	"github.com/electricface/eleplugin"
)

func main() {
	plug, err := eleplugin.Open("/home/del0/gocode/pkg/linux_amd64_dynlink/libele-test-goplugin-plugin_hello.so",
		"ele-test/goplugin/plugin_hello")
	if err != nil {
		log.Fatal(err)
	}
	err = plug.Start("")
	if err != nil {
		log.Fatal(err)
	}
}

```

my GOPATH is `/home/del0/gocode/`。

compile and run:
```sh
$ go install -buildmode=shared std
$ go install -v -buildmode=shared -linkshared ele-test/goplugin/plugin_hello
$ go build -v -linkshared
$ ./use_eleplugin
hello
```