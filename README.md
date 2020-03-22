# Eurek8s

This is a PoC of a controller that listens for ingresses and syncs them
to Eureka.  We used it internally at [Adevinta](https://adevinta.com) as
the basis to create a solution to migrate a fleet of microservices from
AWS to Kubernetes.  The final version we used in production remains
private to Adevinta, but [I wrote in depth about what we implemented
here](https://srvaroa.github.io/kubernetes/eureka/paas/microservices/loadbalancing/2020/02/12/eureka-kubernetes.html)..

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


