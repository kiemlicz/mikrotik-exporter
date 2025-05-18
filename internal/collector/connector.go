package collector

import (
	"fmt"
	"net/http"
	"path"
	"time"
)

type MetricRequest struct {
	MetricName string
	Query      string
}
type MetricResponse struct {
	MetricName string
	Data       map[string]string
	Err        error
}

type MikrotikConnector struct {
	url      string
	username string
	password string

	client *http.Client

	responseChannel chan *MetricResponse
}

func NewMikrotikConnector(
	host string,
	port int,
	username string,
	password string,
	tls bool,
	connectTimeout time.Duration,
	connectionKeepalive time.Duration,
) *MikrotikConnector {
	u := fmt.Sprintf("https://%s:%d", host, port)
	if !tls {
		u = fmt.Sprintf("http://%s:%d", host, port)
	}

	client := &http.Client{
		Timeout: connectTimeout,
		Transport: &http.Transport{
			IdleConnTimeout: connectionKeepalive,
		},
	}

	return &MikrotikConnector{
		url:           u,
		username:      username,
		password:      password,
		client:        client,
		resultChannel: make(chan *MetricResponse),
	}
}

func (mc *MikrotikConnector) Request(request MetricRequest) {
	go func() {
		url := path.Join(mc.url, request.Query)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			mc.responseChannel <- &MetricResponse{
				MetricName: request.MetricName,
				Data:       map[string]string{},
				Err:        err,
			}
			return
		}

		resp, err := mc.client.Do(req)
		defer resp.Body.Close()
		var data []map[string]string
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			mc.responseChannel <- &MetricResponse{
				MetricName: request.MetricName,
				Data:       data,
				Err:        err,
			}
			return
		}
		//fixme cont here
		// For simplicity, send the first object or empty if none
		result := map[string]string{}
		if len(data) > 0 {
			result = data[0]
		}
		mc.responseChannel <- &MetricResponse{
			MetricName: request.MetricName,
			Data:       result,
			Err:        nil,
		}

		if err != nil {
			mc.responseChannel <- &MetricResponse{
				MetricName: request.MetricName,
				Data:       map[string]string{},
				Err:        err,
			}
		} else {

			mc.responseChannel <- &MetricResponse{
				MetricName: request.MetricName,
				Data:       map[string]string{},
				Err:        nil,
			}
		}
	}()
}
