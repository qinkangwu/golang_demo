package client

import rentalpb "server2/rental/api/gen/v1"

type PoiManager struct {
}

func (p PoiManager) GetPoiName(l *rentalpb.Location) (string, error) {
	return "", nil
}

