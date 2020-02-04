# conduit

An implementation of UNIX pipes over HTTP implementing MPMC / PubSub inspired
by [PatchBay](https://patchbay.pub/).

More to come...

## Examples

### Simple Queues / Multi-Producer=Multi-Consumer

```sh
$ curl http://localhost:8000/queue/hello -d 'Hello World'
```

```sh
$ curl http://localhost:8000/queue/hello
Hello World
```

### Simple PubSub

```sh
$ curl http://localhost:8000/topic/foo -d 'Hello World'
```

```sh
$ curl http://localhost:8000/topic/foo
Hello World
```

```sh
$ curl http://localhost:8000/topic/foo
Hello World
```

### Poor man's Web Server

```sh
$ while true; do curl http://localhost:8000/queue/index.html --data-binary @./index.html; done
```

...
