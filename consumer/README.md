# Consumer

Consumer focuses on processing messages received from the project's
message-based "back-end" (e.g., `streamer`).

## Design

The `consumer` will:
- interact directly with `streamer` as a consumer of messages
- process those messages in a way that measurably changes its state 
- provide REST APIs that `portal` can use to access its state and control
  its behavior

## Implementation

The `go` language is being chosen to implement `consumer` in part because
of previous good experience implementing `REST` services in `go`, and also
because `go` was initially targeted for inclusion in this project.

### Frameworks

Given that `go` will be used:

- [gorilla/mux](https://github.com/gorilla/mux) was chosen as the foundation
  of the "middleware" layer on top of which to implement the REST services.
  - _NOTE: Because `mux` is no longer maintained, alternatives were considered,
    (such as [gin](https://gin-gonic.com/) and [chi](https://github.com/go-chi/chi));
    however, `mux` prevailed due to its maturity and my familiarity with it.
    Also, another significant and relevant factor was `mux`'s "modularity":
    whereas `gin` and `chi` both provide many components for a full web framework,
    the `mux` library (within the `gorilla` toolkit) provides what seems to be
    a more focused set of services relevant to the needs of implementing basic
    REST API within `portal`._

#### Source Code Structure

After reviewing the following literature:

- [Golang Standards / Project Layout](https://github.com/golang-standards/project-layout)
- [Project Layout Discussion](https://www.developer.com/languages/go-project-layout/)
- [Kamaleshwar's Opinions](https://github.com/bnkamalesh/goapp)
- [Developer20](https://developer20.com/how-to-structure-go-code/)
- [Ankit Kumar Article](https://dev.to/jinxankit/go-project-structure-and-guidelines-4ccm)
- [Ben Johnson Article](https://medium.com/sellerapp/golang-project-structuring-ben-johnson-way-2a11035f94bc)
- [To `/pkg` or not?](https://travisjeffery.com/b/2019/11/i-ll-take-pkg-over-internal/)

The initial source code structure is targeted at:

```text
- main.go
- internal
  - configs
    \- config.go
  - <other packages specific to implementing consumer>
```

### Docs & Tips on Accessing Kafka using `go`

- [Working with Kafka](https://www.sohamkamani.com/golang/working-with-kafka/)
- [Confluent's Client](https://github.com/confluentinc/confluent-kafka-go)
- [How to use Kafka withi Go](https://www.educative.io/answers/how-to-use-kafka-with-go)
- [Kafka Go Client](https://docs.confluent.io/kafka-clients/go/current/overview.html)
- [sarama](https://github.com/Shopify/sarama)

### Configuration

The `configs` package focuses on providing application-wide
configuration values.  There are hard-coded configuration 
values which can be overridden by values found in the
environment, or read from the `.env` file (which will be )

#### See
- [Loading environment variables properly in Go with env and godotenv](https://www.thedevelopercafe.com/articles/loading-environment-variables-properly-in-go-with-env-and-godotenv-7ec94d4101a7)

## Running

### Kubernetes

To create and validate the service in Kubernetes, do something like:

```shell
$ kubectl create -f k8s
deployment.apps/consumer-dpl created
service/consumer-svc created
$ curl http://localhost:30072/_health
OK!
$ curl http://localhost:30072/status
{"n_consumed":0}
```

### Docker

The target run environment is Kubernetes; however, running it
in docker is also possible:

```shell
$ docker run -e docker.env -p 8072:8072 mdross510/kafkloud-consumer
```