package http

import (
	"context"
	"sync/atomic"
	"testing"

	"github.com/elazarl/goproxy"
	fixtures "github.com/go-git/go-git-fixtures/v4"
	"github.com/jesseduffield/go-git/v5/internal/transport/http/test"
	"github.com/jesseduffield/go-git/v5/plumbing/transport"
	"github.com/stretchr/testify/suite"
)

func TestProxySuite(t *testing.T) {
	suite.Run(t, new(ProxySuite))
}

type ProxySuite struct {
	UploadPackSuite
}

func (s *ProxySuite) TestAdvertisedReferences() {
	var proxiedRequests int32

	s.SetupTest()
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true
	test.SetupHTTPProxy(proxy, &proxiedRequests)

	httpProxyAddr, proxyServer, httpListener := test.SetupProxyServer(s.T(), proxy, false, true)
	defer httpListener.Close()
	defer proxyServer.Close()

	endpoint := s.prepareRepository(fixtures.Basic().One(), "basic.git")
	endpoint.Proxy = transport.ProxyOptions{
		URL:      httpProxyAddr,
		Username: "user",
		Password: "pass",
	}

	s.ups.Client = NewClient(nil)
	session, err := s.ups.Client.NewUploadPackSession(endpoint, nil)
	s.Nil(err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	info, err := session.AdvertisedReferencesContext(ctx)
	s.Nil(err)
	s.NotNil(info)
	proxyUsed := atomic.LoadInt32(&proxiedRequests) > 0
	s.Equal(true, proxyUsed)

	atomic.StoreInt32(&proxiedRequests, 0)
	test.SetupHTTPSProxy(proxy, &proxiedRequests)

	httpsProxyAddr, tlsProxyServer, httpsListener := test.SetupProxyServer(s.T(), proxy, true, true)
	defer httpsListener.Close()
	defer tlsProxyServer.Close()

	endpoint, err = transport.NewEndpoint("https://github.com/git-fixtures/basic.git")
	s.Nil(err)
	endpoint.Proxy = transport.ProxyOptions{
		URL:      httpsProxyAddr,
		Username: "user",
		Password: "pass",
	}
	endpoint.InsecureSkipTLS = true

	session, err = s.ups.Client.NewUploadPackSession(endpoint, nil)
	s.Nil(err)

	info, err = session.AdvertisedReferencesContext(ctx)
	s.Nil(err)
	s.NotNil(info)
	proxyUsed = atomic.LoadInt32(&proxiedRequests) > 0
	s.Equal(true, proxyUsed)
}
