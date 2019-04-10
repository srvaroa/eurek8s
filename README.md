# Eurek8s

PoC of a controller that listens for ingresses and syncs corresponding
ingresses with Eureka.

Started from <https://github.com/trstringer/k8s-controller-core-resource>.

## Running

```
$ go run cmd/eurek8s/main.go
```

or

```
make build && ./bin/eurek8s
```

## Running Eureka, locally

For testing, which really, is the only purpose this should be used for
at this point:

```
docker run -p 8080:8080 netflixoss/eureka:1.3.1
```


