package collector

import (
	"fmt"
	"github.com/go-routeros/routeros/v3"
)

type MikrotikDevice struct {
	Host     string
	Port     int
	Username string
	Password string
}

// TODO make this method lazy, call as late as possible
func (c *MikrotikDevice) connect() (*routeros.Client, error) {
	return routeros.Dial(fmt.Sprintf("%s:%d", c.Host, c.Port), c.Username, c.Password)
}
