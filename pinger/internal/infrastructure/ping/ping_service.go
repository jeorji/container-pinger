package ping

import (
	"time"

	"github.com/prometheus-community/pro-bing"

	"pinger/internal/infrastructure"
)

type ProBingService struct{}

func NewProBingService() infrastructure.PingService {
	return &ProBingService{}
}

func (p *ProBingService) Ping(ip string) (string, float32, error) {
	pinger, err := probing.NewPinger(ip)
	if err != nil {
		return "fail", 0, err
	}

	pinger.Count = 3
	pinger.Timeout = 3 * time.Second
	pinger.SetPrivileged(true)

	err = pinger.Run()
	if err != nil {
		return "fail", 0, err
	}

	stats := pinger.Statistics()
	if stats.PacketsRecv == 0 {
		return "fail", 0, nil
	}

	avgMs := float32(stats.AvgRtt.Microseconds()) / 1000.0
	return "success", avgMs, nil
}

