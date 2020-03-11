# Madprobe

This project has been made to provide companies wanting simple aliveness probe to be easily installed on their systems.

## Getting Started

### Prerequisites

This project has be developed using `go version go1.13`

It uses the well known `gorilla/mux` to serve our API.

Golang can be installed from the official [go website](https://golang.org/dl/).

### Installing

It's quite simple to start the app. A basic CLI using with the standard `flag` package was made to specify common
options like the port on which the server should listen.

```
$ ./madprobe -port 6666 -graceful-timeout 5
```

If you don't specify any options, by default, the application will start a webserver on port `3000`
and provide a graceful timeout of `15 seconds`.

Each probe will run it's in own goroutine and will perform their checks independently.

### API

The API is accessible through HTTP. It implements basic CRUD operations to manage the
state of probes.

For now, there is no probe persistence so stopping `madprobe` will result in losing all
probes states.

> :warning: **HTTP(s) is not supported yet**: Be very careful before exposing anything!

Endpoints: 
  - POST /api/v1/probe/create
````   
{
    "Name": "simple-service-http",
    "URL": "http://localhost:8080/actuator/health",
    "Delay": 5
}
````
  - GET /api/v1/probe/{name}
  - GET /api/v1/probe
  - PUT /api/v1/probe/{name}
````   
{
    "Name": "simple-service-http",
    "URL": "http://localhost:8080/actuator/health",
    "Delay": 5
}
````
  - DELETE /api/v1/probe/{name}

## Contributing

I'll be more than happy to have feedback on the way I designed this application. Things can always be done better and
I m eager to learn what could be improved!

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details