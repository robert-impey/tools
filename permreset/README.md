# permreset

Small CLI to find executable scripts (files whose first line starts with `#!`) and set their permissions to `0755`.

Usage:

```sh
go run . -s ~/local-scripts
```

Tests:

```sh
go test ./...
```
