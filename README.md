# From TCP to HTTP

Implementation from TCP of a simple HTTP server prototype

## To test it

- git clone the project

- ```bash
  go mod tidy
  ```

- ```bash
  go run cmd/httpserver/main.go
  ```

- go to any page for 200 [http://localhost:42069/](http://localhost:42069/)
- go to [http://localhost:42069/yourproblem](http://localhost:42069/yourproblem) for 400
- go to [http://localhost:42069/myproblem](http://localhost:42069/myproblem) for 500

## To modify the default behavior

Here is the implementation of the default behavior

```go
func main() {
 server, err := server.Serve(port,
func()
 {
  if req.RequestLine.RequestTarget == "/yourproblem" {
   w.WriteStatusLine(response.BadRequest)
   w.WriteBody(response.BadRequestBody)
   return
  }

  if req.RequestLine.RequestTarget == "/myproblem" {
   w.WriteStatusLine(response.InternalServerError)
   w.WriteBody(response.InternalErrorBody)
   return
  }

  w.WriteBody(response.OkBody)
 })
 if err != nil {
  log.Fatalf("Error starting server: %v", err)
 }
 defer server.Close()
 log.Println("Server started on port", port)

 sigChan := make(chan os.Signal, 1)
 signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
 <-sigChan
 log.Println("Server gracefully stopped")
}
```

Simply modify the function passed in paramater of Serve with a signature of

```go
func(w *response.Writer, req *request.Request)
```
