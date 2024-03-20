package elasticex

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/xm-chentl/goresource"
)

type Config struct {
	Addrs    []string
	User     string
	Password string
}

type resource struct {
	config *Config
	client *elasticsearch.Client // 连接
}

func (f resource) Db(args ...interface{}) goresource.IRepository {
	repo := &repository{
		client: f.client,
	}
	for _, arg := range args {
		if ctx, ok := arg.(context.Context); ok {
			repo.ctx = ctx
		}
	}

	return repo
}

func (f resource) Uow() goresource.IUnitOfWork {
	return nil
}

func New(config *Config) goresource.IResource {
	if config == nil {
		panic("elastic config is nil")
	}

	cfg := elasticsearch.Config{
		Addresses: config.Addrs,
		Username:  config.User,
		Password:  config.Password,
		Transport: &http.Transport{
			//MaxIdleConnsPerHost 如果非零，控制每个主机保持的最大空闲(keep-alive)连接。如果为零，则使用默认配置2。
			MaxIdleConnsPerHost: 10,
			//ResponseHeaderTimeout 如果非零，则指定在写完请求(包括请求体，如果有)后等待服务器响应头的时间。
			ResponseHeaderTimeout: time.Second,
			//DialContext 指定拨号功能，用于创建不加密的TCP连接。如果DialContext为nil(下面已弃用的Dial也为nil)，那么传输拨号使用包网络。
			DialContext: (&net.Dialer{Timeout: time.Second}).DialContext,
			// TLSClientConfig指定TLS.client使用的TLS配置。
			//如果为空，则使用默认配置。
			//如果非nil，默认情况下可能不启用HTTP/2支持。
			TLSClientConfig: &tls.Config{
				MaxVersion: tls.VersionTLS11,
				//InsecureSkipVerify 控制客户端是否验证服务器的证书链和主机名。
				InsecureSkipVerify: true,
			},
		},
	}
	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		panic(fmt.Errorf("elastic init is failed: %v", err))
	}

	return &resource{
		config: config,
		client: client,
	}
}
