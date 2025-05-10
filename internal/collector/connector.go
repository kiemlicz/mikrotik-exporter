package collector

import (
	"fmt"
	"github.com/go-routeros/routeros/v3"
	"sync"
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

type ConnectionPool struct {
	pool sync.Map // Target (IP:PORT) to *MikrotikDevice map
}

// GetConnector retrieves or creates a connector for the given target.
func (cp *ConnectionPool) GetConnector(host string, port int, username, password string) *MikrotikDevice {
	target := fmt.Sprintf("%s:%d", host, port)

	conn, _ := cp.pool.LoadOrStore(target, &MikrotikDevice{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	})

	return conn.(*MikrotikDevice)
}

func (cp *ConnectionPool) RemoveConnector(host string, port int) {
	target := fmt.Sprintf("%s:%d", host, port)
	cp.pool.Delete(target) //disconnect handling
}

type MikrotikCollector struct {
	device        *MikrotikDevice
	deviceMetrics sync.Map
}

func (c *MikrotikCollector) Collect() {

}
