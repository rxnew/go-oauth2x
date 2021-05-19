package main

import (
	"context"
	"os"

	"golang.org/x/oauth2"

	"github.com/rxnew/go-oauth2x"
)

func main() {
	c := &oauth2.Config{
		ClientID:     os.Getenv("OAUTH2_CLIENT_ID"),
		ClientSecret: os.Getenv("OAUTH2_CLIENT_SECRET"),
		Endpoint: oauth2.Endpoint{
			AuthURL:  os.Getenv("OAUTH2_AUTH_URL"),
			TokenURL: os.Getenv("OAUTH2_TOKEN_URL"),
		},
	}
	ctx := context.Background()
	cli := oauth2x.NewClient(ctx, c.TokenSource(ctx, nil))
	resp, err := cli.Get("https://example.com")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}
