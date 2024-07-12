package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/MatthieuCoder/OrionV3/internal"
	"github.com/MatthieuCoder/OrionV3/internal/proto"
	"github.com/vishvananda/netlink"
	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"google.golang.org/grpc"
)

func main() {

	// Get TLS credentials
	cred := internal.NewClientTLS()

	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", "reg.orionet.re", 6443), grpc.WithTransportCredentials(cred), grpc.WithIdleTimeout(time.Second*120))
	if err != nil {
		log.Fatalf("Unable to connect gRPC channel %v", err)
	}

	// Create the gRPC client
	_ = proto.NewRegistryClient(conn)
	holepunch := proto.NewHolePunchingServiceClient(conn)

	privatek, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		panic(err)
	}
	publickey := privatek.PublicKey()
	ctx := context.Background()
	events, err := holepunch.Session(ctx, &proto.HolePunchingInitialize{
		PublicKey: publickey[:],
	})
	if err != nil {
		panic(err)
	}
	wgclient, err := wgctrl.New()
	if err != nil {
		panic(err)
	}

	for {
		event, err := events.Recv()
		if err == io.EOF {
			return
		}
		if err != nil {
			return
		}
		if event.GetComplete() != nil {
			fmt.Println("Hole punching is complete!")
			fmt.Println(event.GetComplete())
		}
		if event.GetInitializationResponse() != nil {
			fmt.Println("Connecting to wg to hp;..")
			// Needs to initialize a new wg tunnel
			ir := event.GetInitializationResponse()
			fmt.Println(ir)
			// Create interface
			wglink := internal.WireguardNetLink{
				Id:     200,
				Prefix: "clu",
			}

			err = netlink.LinkAdd(wglink)
			if err != nil {
				panic(err)
			}

			port := 40000 + 200
			defer func() {
				netlink.LinkDel(wglink)
			}()
			ip, err := net.LookupIP(ir.EndpointAddr)
			if err != nil {
				panic(err)
			}
			keepalive := time.Second * 5
			preshared := wgtypes.Key(ir.PresharedKey)
			thirtyonemask := net.CIDRMask(31, 32)
			device := wgtypes.Config{
				ListenPort: &port,
				PrivateKey: &privatek,
				Peers: []wgtypes.PeerConfig{
					{
						PublicKey: wgtypes.Key(ir.PublicKey),
						Endpoint: &net.UDPAddr{
							IP:   ip[0],
							Port: int(ir.EndpointPort),
						},
						AllowedIPs: []net.IPNet{
							{
								IP:   net.ParseIP(ir.ClientAddress).Mask(thirtyonemask),
								Mask: thirtyonemask,
							},
						},
						PersistentKeepaliveInterval: &keepalive,
						PresharedKey:                &preshared,
					},
				},
			}

			err = wgclient.ConfigureDevice("clu200", device)
			if err != nil {
				panic(err)
			}
			err = netlink.LinkSetUp(wglink)
			if err != nil {
				panic(err)
			}

		}
	}
}
