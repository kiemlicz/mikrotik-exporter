package collector

import (
	"fmt"
	"github.com/go-routeros/routeros/v3"
	"strings"
	"sync"
	"time"
)

type MikrotikRequest = string
type MikrotikResponse = string

type MikrotikDevice struct {
	Host                string
	Port                int
	Username            string
	Password            string
	connectTimeout      time.Duration
	connectionKeepalive time.Duration
}

type MikrotikConnector struct {
	device *MikrotikDevice
	conn   *routeros.Client
	timer  *time.Timer
	mutex  sync.Mutex
}

func (mc *MikrotikConnector) login() {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	if mc.conn == nil {
		c, err := routeros.DialTimeout(
			fmt.Sprintf("%s:%d", mc.device.Host, mc.device.Port),
			mc.device.Username,
			mc.device.Password,
			mc.device.connectTimeout,
		)
		if err != nil {
			//?
		}
		mc.conn = c
		mc.conn.Async()
		if mc.timer != nil {
			mc.timer.Stop()
		}
		mc.timer = time.AfterFunc(mc.device.connectionKeepalive, func() {
			mc.mutex.Lock()
			defer mc.mutex.Unlock()
			if mc.conn != nil {
				mc.conn.Close()
				mc.conn = nil
			}
		})
	}
}

func (mc *MikrotikConnector) Run(request MikrotikRequest) (MikrotikResponse, error) {
	mc.login()
	mc.timer.Reset(mc.device.connectionKeepalive)
	reply, err := mc.conn.RunArgs(strings.Split(string(request), " "))
	//fixme very bad, check how to read from channel
	if err != nil {
		return "", err
	} else {
		return reply.String(), nil
	}
}
