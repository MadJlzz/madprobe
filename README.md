# Madprobe
[![Go report](https://goreportcard.com/badge/github.com/MadJlzz/madprobe)](https://goreportcard.com/report/github.com/MadJlzz/madprobe)

This project has been made to provide companies wanting simple aliveness probe to be easily installed on their systems.

## Getting Started

### Prerequisites

This project has be developed using `go version go1.13`

It uses the well known `gorilla/mux` to serve our API.

Golang can be installed from the official [go website](https://golang.org/dl/).

### Installing

It's quite simple to start the app. A basic CLI using the `pflag` package to specify common
options like the port on which the server should listen. `viper` is used to manage configuration context.

```shell script
./madprobe --port 6666 --graceful-timeout 5
```

If you don't specify any options, by default, the application will start an HTTP webserver on port `3000`
and provide a graceful timeout of `15 seconds`.

If you want to run the server in HTTPs, you need to pass more arguments to the command line:
```shell script
./madprobe --cert configs/certs/public.pem --key configs/certs/key.pem
```

Do not forget to set the `--ca-cert` flag if you desire your probes to do HTTPS requests.
```shell script
./madprobe --cert configs/certs/public.pem --key configs/certs/key.pem --ca-cert configs/certs/cacert.pem
```

Also, be aware you can configure `madprobe` using a `yaml configuration` file. Here's an example:
```yaml
port: 3000
cert: configs/certs/public.pem
key: configs/certs/key.pem
ca-cert: configs/certs/cacert.pem
```

Environment variables are a way to configure `madprobe` too!
```shell script
export PORT=3000
export CERT=configs/certs/public.pem
export KEY=configs/certs/key.pem
export CA-CERT=configs/certs/cacert.pem
```

> :warning: **Pay attention to the override direction**: defaults, config file, env. variables, flags

If you want to generate basic certificates, please look in the configs/certs directory.
`gencert.sh` is based on `cfssl` and `cfssljson` which are easier to use than `openssl`.

Each probe will run it's in own goroutine and will perform their checks independently.

### API

The API is accessible through HTTP. It implements basic CRUD operations to manage the
state of probes.

For now, there is no probe persistence so stopping `madprobe` will result in losing all
probes states.

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
  - DELETE /api/v1/probe/{name}

## Contributing

I'll be more than happy to have feedback on the way I designed this application. Things can always be done better and
I m eager to learn what could be improved!

## License

This project license is the MIT License - see the [LICENSE](LICENSE) file for details
