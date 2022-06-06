package client

import "context"

type ProfileManager struct {
}

func (p *ProfileManager) Verify(ctx context.Context, userId string) (string, error) {
	return "qkw测试", nil
}
