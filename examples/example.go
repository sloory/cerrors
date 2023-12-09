package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/sloory/cerrors"
)

func handler(ctx context.Context, request map[string]any) error {
	ctx = cerrors.InComponent(ctx, "handler")

	// Place in context data from request for get more information related to error when it happened
	ctx = cerrors.WithCtxField(ctx, "userId", request["userId"])
	ctx = cerrors.WithCtxField(ctx, "requestId", request["requestId"])

	err := service(ctx)
	if err != nil {
		// Hide internal error message from clients.
		// Because we wrap real error in Sentry we will have stacktrace and real reason od error 
		return cerrors.Opaque("Internal error", err)
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
		// Place additional information related to error
		cerrors.WithField(
			errors.New("some error"),
			"db", "postgres",
		),
	)
}

func main() {
	// It is like http/rpc handler example
	err := handler(
		context.Background(),
		map[string]any{
			"userId":    11,
			"requestId": "0f1e6fb0-964f-11ee-b9d1-0242ac120002",
		},
	)
	// It can be placed in middleware as common place for log errors
	if err != nil {
		fmt.Println(cerrors.Components(err))
		fmt.Println(cerrors.Fields(err))
		fmt.Printf("%+v", err)
	}
}
