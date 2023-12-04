package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/sloory/cerrors"
)

func handler(ctx context.Context) error {
	ctx = cerrors.InComponent(ctx, "handler")

	err := service(ctx)
	if err != nil {
		return err
	}

	return nil
}

func service(ctx context.Context) error {
	ctx = cerrors.InComponent(ctx, "service")

	err := repository(ctx)
	if err != nil {
		return err
	}

	return nil
}

func repository(ctx context.Context) error {
	ctx = cerrors.InComponent(ctx, "repository")

	return cerrors.Enrich(
		ctx,
		cerrors.WithField(
			errors.New("some error"),
			"db", "postgres",
		),
	)
}

func main() {
	err := handler(context.Background())
	if err != nil {
		fmt.Println(cerrors.Components(err))
		fmt.Println(cerrors.Fields(err))
		fmt.Printf("%+v", err)
	}
}
