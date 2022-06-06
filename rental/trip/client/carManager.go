package client

import (
	"context"
	rentalpb "server2/rental/api/gen/v1"
)

type CarManager struct {
}

func (c *CarManager) Verify(ctx context.Context, carId string, location *rentalpb.Location) error {
	return nil
}

func (c *CarManager) Unlock(ctx context.Context, carId string) error {
	return nil
}
