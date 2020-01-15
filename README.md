# Madprobe

This project has been made to provide company wanting simple liveness probe to be easily
installed on their systems.

## Getting Started

### Prerequisites

This project has be developed using `go version go1.13`

Golang can be installed from the official [go website](https://golang.org/dl/).

### Installing

It's quite simple to start the app. A basic CLI made with the standard `flag` package was made to specify common
options like the configuration file to create probes.

```
$ madprobe -yamlConfig /opt/probes/dummy.yaml
```

If you don't specify any options, by default, the application will search for a configuration file located
in `./test/sample.yml`

Each probe will run it's in own goroutine and will perform their checks independently.

An example of configuration can be found [here](configs/sample.yml).

## Contributing

I'll be more than happy to have feedbacks on the way I designed this application. Things can always be done better and
I m eager to learn what could be improved!

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details