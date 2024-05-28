# Masterclass - Golang

## Summary

- Basics
  - Create a project
  - Declare variables
  - Use File system
  - Handle errors
  - Build a project
- Folder structure
- gRPC
- OOP
  - Interfaces
  - Dependency Injection
- Parallelism
  - Go Routines
  - Channels
  - Mutex
- Extras
  - ZMQ
  - Tests

## Basics

### Create a project

Initialize a module to manage dependencies (like `package.json` or `requirements.txt`):

```bash
go mod init github.com/uandersonricardo/masterclass-go
```

It is a convention to use the repo URL as the module name (actually, it is the import path).

To manage dependencies, we use some commands:

- `go get [PACKAGE_NAME]` to download a specific package
- `go mod tidy` to scan the project, update `go.mod` and download all the dependencies automatically
- `go mod download` to download all the dependencies in `go.mod`

### Declare variables

First, we create a `main.go` file as our entrypoint.

```go
package main

import (
	"fmt"
)

func main() {
	message := "Hello, World!"
	fmt.Println(message)
}
```

Then, we run it executing `go run main.go` on terminal.

`:=` operator creates and assign the variable, so we don't need to:

```go
var message string
message = "Hello, World!"
```

### Use File system

Let's create a file with some string on the root of the project using the `os` package:

```go
package main

import (
	"fmt"
	"os"
)

func main() {
	os.WriteFile("file.txt", []byte("Some string"), 0777)
	fmt.Println("File created")
}
```

`WriteFile` receives path name, a slice of bytes and the permission bits for the file.

Now, let's read the file:

```go
package main

import (
	"fmt"
	"os"
)

func main() {
	content, _ := os.ReadFile("file.txt")
	fmt.Println(string(content))
}
```

### Handle errors

What if we pass a filename that does not exists?

```go
package main

import (
	"fmt"
	"os"
)

func main() {
	content, _ := os.ReadFile("test.txt")
	fmt.Println(string(content))
}
```

Nothing is passed, because the statement throws an error. Actually, the error was not "thrown". In Go, errors are values and it is idiomatic to return them:

```go
package main

import (
	"fmt"
	"os"
)

func main() {
	content, err := os.ReadFile("test.txt")

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(content))
}
```

### Build a project

Usually we create a Makefile for the tasks of the project, such as build, run, install system dependencies, compile protos...

```Makefile
run:
	@go run main.go

build:
	@go build -o bin/main main.go

start:
	@./bin/main
```

Then we can run `make build`.

## Folder structure

We created a `bin/` folder for the binaries. But how do we normally structure the application?

A common Go structure is as follows: https://gist.github.com/ayoubzulfiqar/9f1a34049332711fddd4d4b2bfd46096

(Move main.go to cmd/ and change Makefile tasks)

## gRPC

First, we need to install protoc: https://grpc.io/docs/protoc-installation/

Then, we download the packages:

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
```

I will create a protobuf here:

`protos/example.proto`

```proto
syntax = "proto3";

package masterclass.go;

option go_package = "github.com/uandersonricardo/masterclass-go/pkg/pb";

service FrameService {
    rpc GetFrame (GetFrameRequest) returns (Frame) {}
}

message GetFrameRequest {
    int32 id = 1;
}

message Frame {
    int32 id = 1;
}

```

Now, a script to compile the proto to Go files:

`scripts/compile-proto.sh`

```bash
protoc --go_out=./pkg/pb --go_opt=module=github.com/uandersonricardo/masterclass-go/pkg/pb \
       --go-grpc_out=./pkg/pb --go-grpc_opt=module=github.com/uandersonricardo/masterclass-go/pkg/pb \
       ./protos/example.proto
```

Finally, let's add a task on Makefile:

```Makefile
compile-proto:
	@bash ./scripts/compile-proto.sh
```

Now, we run:

```bash
make compile-proto
go mod tidy
```

## OOP

### Interfaces

Now, we can declare a interface for our gRPC server:

`internal/server.go`

```go
package internal

type Server interface {
	Start() error
}
```

Interfaces only contain method signatures, not properties.

### Dependency Injection

We do not have the concept of constructor for objects, so we create a function that builds the struct (passing the dependencies) and return it:

`internal/grpc_server.go`

```go
package internal

import (
	"context"
	"net"

	"github.com/uandersonricardo/masterclass-go/pkg/pb"
	"google.golang.org/grpc"
)

type GrpcServer struct {
	server  *grpc.Server
	address string

	pb.UnimplementedFrameServiceServer
}

func NewGrpcServer(address string) *GrpcServer {
	server := grpc.NewServer()

	return &GrpcServer{
		server:  server,
		address: address,
	}
}

func (s *GrpcServer) Start() error {
	pb.RegisterFrameServiceServer(s.server, s)
	lis, err := net.Listen("tcp", s.address)

	if err != nil {
		return err
	}

	s.server.Serve(lis)
	return nil
}

func (s *GrpcServer) GetFrame(ctx context.Context, in *pb.GetFrameRequest) (*pb.Frame, error) {
	return &pb.Frame{
		Id: 1,
	}, nil
}

```

(It is important to know that when a function starts with a capital letter, it is like a public method, that is exported and visible to another files and packages)

Now, we can update `main.go`:

```go
package main

import (
	"fmt"
	"os"

	"github.com/uandersonricardo/masterclass-go/internal"
)

func main() {
	fmt.Println("Starting server...")
	server := internal.NewGrpcServer(":8080")
	err := server.Start()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
```

Everything OK, we can test with Insomnia now.

## Parallelism

### Go Routines

What if we wanted to run another server in our application? Let's say a HTTP server:

`internal/http_server`

```go
package internal

import "net/http"

type HTTPServer struct {
	address string
}

func NewHTTPServer(address string) *HTTPServer {
	http.HandleFunc("/health", HealthHandler)
	return &HTTPServer{
		address: address,
	}
}

func (s *HTTPServer) Start() error {
	return http.ListenAndServe(s.address, nil)
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
```

Now, we add this to our main file:

```go
package main

import (
	"fmt"
	"os"

	"github.com/uandersonricardo/masterclass-go/internal"
)

func main() {
	fmt.Println("Starting server...")
	server := internal.NewGrpcServer(":8080")
	err := server.Start()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Starting HTTP server...")
	httpServer := internal.NewHTTPServer(":8081")
	err = httpServer.Start()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
```

If we run this, we will notice that the 'Stating HTTP Server' message will never be printed. This is obvious, since `server.Start()` is blocking our main thread. We can refactor our code to use Go routines to deal with them:

```go
package main

import (
	"fmt"
	"os"

	"github.com/uandersonricardo/masterclass-go/internal"
)

func startGrpcServer() {
	fmt.Println("Starting server...")
	server := internal.NewGrpcServer(":8080")
	err := server.Start()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func startHTTPServer() {
	fmt.Println("Starting HTTP server...")
	httpServer := internal.NewHTTPServer(":8081")
	err := httpServer.Start()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	go startGrpcServer()
	go startHTTPServer()
}
```

Now, when we run this, our code exits without printing anything. This happens because we are not waiting the threads:

```go
package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/uandersonricardo/masterclass-go/internal"
)

func startGrpcServer(wg *sync.WaitGroup) {
	fmt.Println("Starting server...")
	server := internal.NewGrpcServer(":8080")
	err := server.Start()

	if err != nil {
		fmt.Println(err)
		wg.Done()
	}
}

func startHTTPServer(wg *sync.WaitGroup) {
	fmt.Println("Starting HTTP server...")
	httpServer := internal.NewHTTPServer(":8081")
	err := httpServer.Start()

	if err != nil {
		fmt.Println(err)
		wg.Done()
	}
}

func main() {
	wg := &sync.WaitGroup{} // or var wg *sync.WaitGroup
	wg.Add(2)

	go startGrpcServer(wg)
	go startHTTPServer(wg)

	wg.Wait()
	os.Exit(1)
}
```

### Channels

Channels are a thread-safe way to communicate between threads, so we can send data through the channel without needing to worry about mutexes:

```go
ch := make(chan string)
ch <- "Shakespeare"
name := <-ch
```

### Mutex

```go
var mu *sync.Mutex

mu.Lock()
// anything
mu.Unlock()
```

## Extras

### ZMQ

https://zeromq.org/languages/go/

### Tests

https://go.dev/doc/tutorial/add-a-test

### BFFs

https://github.com/robocin/ssl-core/pull/89
