# Example gRPC + HTTP Go service with nginx

This is a companion to the blog post [A misleading error when using gRPC with
Go and nginx](https://kennethjenkins.net/posts/go-nginx-grpc/), illustrating
a problem I ran into when using nginx to terminate TLS for a Go service that
served both HTTP and gRPC on the same port.

To run:

    $ ./run.sh

The script requires Docker and docker-compose. You'll also want to install
[grpcurl](https://github.com/fullstorydev/grpcurl) to make gRPC requests from
the command line.

First, the script runs openssl within an Alpine Linux container to generate a
self-signed TLS certificate for nginx, and then uses docker-compose to build
and start the Go service as well as an nginx reverse proxy container.

## The nginx error message

Once running, the buggy behavior is demonstrated on port 8080. Using `curl`
to make an HTTP request succeeds:

    $ curl --cacert ssl.cert https://localhost:8080
    Hello world!

But using `grpcurl` to make a gRPC request fails:

    $ grpcurl -d '{"message": "hello"}' -cacert ssl.cert -proto service.proto localhost:8080 Echo.Echo
    ERROR:
      Code: Unavailable
      Message: unexpected HTTP status code received from server: 502 (Bad Gateway); transport: received unexpected content-type "text/html"

The nginx logs shown by docker-compose should contain the error that inspired
the blog post:

> [error] 29#29: *1 upstream sent too large http2 frame: 4740180 while reading response header from upstream, client: 172.18.0.1, server: , request: "POST /Echo/Echo HTTP/2.0", upstream: "grpc://172.18.0.2:8080", host: "localhost:8080"

## The Go service fix

The fixed behavior is demonstrated on port 8081:

    $ curl --cacert ssl.cert https://localhost:8081
    Hello world!
    $ grpcurl -d '{"message": "hello"}' -cacert ssl.cert -proto service.proto localhost:8081 Echo.Echo
    {
      "message": "hello"
    }
