package wechat

import "fmt"

type Service struct {
}

func (s Service) Resolve(code string) (string, error) {
	return fmt.Sprintf("111"), nil
}
