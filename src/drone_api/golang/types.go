package drone

import (
	"fmt"
)

type DroneConfig struct {
	Host     string
	Port     int32
	Protocol string
	Token    string
}

func (c DroneConfig) apiUrl() string {
	return fmt.Sprintf("%s://%s:%d", c.Protocol, c.Host, c.Port)
}
