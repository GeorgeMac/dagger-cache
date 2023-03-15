package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"dagger.io/dagger"
)

func daggerClient(ctx context.Context) (*dagger.Client, error) {
	return dagger.Connect(ctx,
		dagger.WithLogOutput(os.Stdout),
	)
}

func main() {
	ctx := context.Background()

	client, err := daggerClient(ctx)
	if err != nil {
		panic(err)
	}

	src := client.Host().Directory(".")

	id, err := src.File("go.Sum").ID(ctx)
	if err != nil {
		panic(err)
	}

	fooCache := client.CacheVolume(fmt.Sprintf("foo-cache-%s", id))

	client.Container().
		From("golang:1.18-alpine3.16").
		WithMountedDirectory("/src", src).
		WithMountedCache("/foo", fooCache).
		WithEnvVariable("CACHE_BUST", time.Now().String()).
		WithExec([]string{"touch", "/foo/bar.txt"}).
		WithExec([]string{"ls", "/foo"}).
		ExitCode(ctx)
}
