# Worker
> ⛏ Mineway's worker is an easy-to-install worker for manage your rigs with Mineway interface

## Summary

1. Development Environment

## Development Environment

First of all, you need to clone this repository in your workspace :

```go
git clone git@github.com:mineway/worker.git
```

Before running worker, you need to generate all routes from annotation,
for that, you need to run worker with `build` parameter like that :

```go
go run -v ./cmd/worker build
```

`routes.json` file is now fill, you can run your worker :

```go
go run -v ./cmd/worker
```

That it !

You can also use Makefile :

```
make worker-dev
```