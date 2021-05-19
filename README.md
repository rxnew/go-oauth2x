# Extension of the OAuth2 library for Go

A library extends [OAuth2 for Go](https://github.com/golang/oauth2).

- Asynchronous preloading of a token

## Installation

```shell
go get github.com/rxnew/go-oauth2x
```

## Examples

```go
package main

import (
	"context"

	"github.com/rxnew/go-oauth2x"
	"golang.org/x/oauth2"
)

func main() {
	c := &oauth2.Config{ /* Configure */ }
	ctx := context.Background()
	cli := oauth2x.NewClient(ctx, c.TokenSource(ctx, nil))
	resp, err := cli.Get("https://example.com")
	// Process the response
}
```
