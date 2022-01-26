package logical

import (
	"fmt"
	"github.com/go-various/consul"
	"github.com/go-various/micro"
	"time"
)

func init() {
	micro.InitializeCircuit(5, time.Minute * 5)
}

var _ micro.Service = (*microService)(nil)

type microService struct {
	cli consul.Client
}

func NewMicroServiceClient(cli consul.Client) *microService {
	return &microService{cli}
}

//GetServers 微服务实例获取接口
//此处获取的服务是根据微服务注册时候的 consul.Service 对象处理
func (m *microService) GetServers(name, tags string) ([]micro.Server, error) {
	ss, err := m.cli.GetServices(name, tags)
	if err != nil {
		return nil, fmt.Errorf("fetch service[%s]: %v", name , err)
	}

	var servers []micro.Server
	for _, s := range ss {
		host := s.TaggedAddresses[consul.WanAddrKey]
		addr := fmt.Sprintf("http://%s:%d",host.Address, host.Port)
		server := micro.Server{
			ID:          s.ID,
			Address:     addr,
			Weight:      s.Weights.Passing,
			TPSDelay:    0,
		}
		servers = append(servers, server)
	}
	return servers, nil
}
