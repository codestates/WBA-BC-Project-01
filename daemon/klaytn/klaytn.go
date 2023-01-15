package klaytn

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/airbloc/airbloc-go/pkg/kas"
	klayClient "github.com/klaytn/klaytn/client"
	"github.com/klaytn/klaytn/networks/rpc"
	"github.com/pkg/errors"
)

type clientData struct{ *klayClient.Client }
type Client struct {
	clientData
}

// Dial connects to Klaytn API Service and returns a JSON-RPC client.
func Dial(cfg Config) (*klayClient.Client, error) {
	if _, ok := networkNameToChainIDs[cfg.Network]; !ok {
		return nil, fmt.Errorf("validate KAS config: unknown network: \"%s\"", cfg.Network)
	}
	if cfg.AccessKeyID == "" {
		return nil, fmt.Errorf("validate KAS config: access key ID is empty")
	}
	if cfg.SecretAccessKey == "" {
		return nil, fmt.Errorf("validate KAS config: secret access key is empty")
	}
	cli := new(http.Client)
	cli.Transport = kasAuthTransport{cfg}
	rpcCli, err := rpc.DialHTTPWithClient(cfg.Endpoint, cli)
	if err != nil {
		return nil, err
	}
	return klayClient.NewClient(rpcCli), nil
}

// kasAuthTransport wraps http.RoundTripper and adds authentication header to HTTP request
// for communicating with Klaytn API Service.
type kasAuthTransport struct {
	cfg Config
}

// RoundTrip adds authentication header to JSONRPC request.
// For details, please refer to KAS documentation (https://console.klaytnapi.com/ko/service/node)
func (k kasAuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.SetBasicAuth(k.cfg.AccessKeyID, k.cfg.SecretAccessKey)
	req.Header.Set("x-chain-id", fmt.Sprintf("%d", networkNameToChainIDs[k.cfg.Network]))

	return http.DefaultTransport.RoundTrip(req)
}

func NewClientWithKAS(ctx context.Context, cfg kas.Config) (*Client, error) {
	client, err := kas.Dial(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "dial klaytn api")
	}
	cid, err := client.NetworkID(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "fetch network id")
	}
	log.Println("Using {} network", cid.String())
	return &Client{
		clientData: clientData{client},
	}, nil
}

/* 클레이튼 네트워크 소켓 연결*/
// var cfg kas.Config
// cfg.AccessKeyID = cf.Klaytn.AccessKeyID
// cfg.Endpoint = cf.Klaytn.Endpoint
// cfg.SecretAccessKey = cf.Klaytn.SecretKeyID
// cfg.Network = cf.Klaytn.Network

// client, err := klaytn.NewClientWithKAS(context.Background(), cfg)
// if err != nil {
// 	log.Fatal(err)
// }
// headers := make(chan *types.Header)

// sub, err := client.SubscribeNewHead(context.Background(), headers)
// if err != nil {
// 	log.Fatal(err)
// }
// for {
// 	select {
// 	case err := <-sub.Err():
// 		log.Fatal(err)

// 	case header := <-headers:
// 		block, err := client.BlockByNumber(context.Background(), header.Number)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		fmt.Printf("%+v", block)
// 	}
// }
