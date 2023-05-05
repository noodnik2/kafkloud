# `producer`

`producer` is intended to serve as an adapter between the front-end application
(e.g., `portal`) and the message-based "back-end" (e.g., `streamer`).

## Design

The `producer` will provide REST APIs that `portal` can use to control its
behavior, and will interact directly with `streamer` using messages.

## Implementation

It was thought to continue developing experience with [Typescript](https://www.typescriptlang.org/)
technologies for implementing `producer`.  Accordingly, some
[popular REST API frameworks](https://externlabs.com/blogs/extensive-comparison-of-nodejs-framework/)
[were considered](https://npmtrends.com/express-vs-fastify-vs-hapi-vs-koa-vs-nest.js-vs-next-vs-restify),
such as:

- [ExpressJS](https://expressjs.com/)
- [Fastify](https://www.fastify.io/)
- [Hapi](https://hapi.dev/)
- [Koa](https://koajs.com/)
- [NestJS](https://nestjs.com/)
- [Restify](http://restify.com/)

After consideration, `ExpressJS` and `Typescript` are being chosen.

## Running

### Kubernetes

To create and validate the service in Kubernetes, do something like:

```shell
$ kubectl create -f k8s
deployment.apps/producer-dpl created
service/producer-svc created
$ curl http://localhost:30000/_health
Healthy!
```

### Docker

The target run environment is Kubernetes; however, running it
in docker is also possible:

```shell
$ docker run --env-file docker.env -p 8000:8000 mdross510/kafkloud-producer
$ curl localhost:8000/_health
Healthy!
```

### Locally
Run it locally with:

```shell
$ npm run dev:server
```

#### Notes

The `nodemon` restart loop often produces an `EADDRINUSE` error.
There's no known solution for this.  
Just kill the process that's still running and try again, e.g.:

```shell
$ kill $(lsof -t -i:8000)
$ : or, via script:
$ npm run dev:kill-server
```

## Testing

The following `curl` may be relevant:

```shell
$ curl -H "Content-Type: application/json" \
  -d '{"session": "123", "type":"produce", "message":{"text":"word1 word2"}}' \
  http://localhost:30000/produce/stream/key
```

- _Tip: use the correct port; e.g. `30000` (K8s) or `8000` (local)_
