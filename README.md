# NixLight

Simple fishtank light controller.

## Running it

First of all... don't.  I make no guarantees that this won't blow things up.

### Vendoring

I'm experimenting with the 'to-be-official' golang `dep` tool... but you don't have to.

```
go get -u github.com/golang/dep/...
dep init
dep ensure -update
```

Vendor somehow or you can't build from the included `Dockerfile`.

### Build

```
docker build . -t nixlight
```

### Run

Run the container with something like this...

```
docker run -v `pwd`/example:/app/config --rm -it nixlight -config /app/config/nixlight.toml
```
