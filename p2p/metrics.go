// Contains the meters and timers used by the networking layer.

package p2p

import (
	"net"

	"github.com/rcrowley/go-metrics"
)

var (
	ingressConnectMeter = metrics.GetOrRegisterMeter("p2p/InboundConnects", metrics.DefaultRegistry)
	ingressTrafficMeter = metrics.GetOrRegisterMeter("p2p/InboundTraffic", metrics.DefaultRegistry)
	egressConnectMeter  = metrics.GetOrRegisterMeter("p2p/OutboundConnects", metrics.DefaultRegistry)
	egressTrafficMeter  = metrics.GetOrRegisterMeter("p2p/OutboundTraffic", metrics.DefaultRegistry)
)

// meteredConn is a wrapper around a network TCP connection that meters both the
// inbound and outbound network traffic.
type meteredConn struct {
	*net.TCPConn // Network connection to wrap with metering
}

// newMeteredConn creates a new metered connection, also bumping the ingress or
// egress connection meter.
func newMeteredConn(conn net.Conn, ingress bool) net.Conn {
	if ingress {
		ingressConnectMeter.Mark(1)
	} else {
		egressConnectMeter.Mark(1)
	}
	return &meteredConn{conn.(*net.TCPConn)}
}

// Read delegates a network read to the underlying connection, bumping the ingress
// traffic meter along the way.
func (c *meteredConn) Read(b []byte) (n int, err error) {
	n, err = c.TCPConn.Read(b)
	ingressTrafficMeter.Mark(int64(n))
	return
}

// Write delegates a network write to the underlying connection, bumping the
// egress traffic meter along the way.
func (c *meteredConn) Write(b []byte) (n int, err error) {
	n, err = c.TCPConn.Write(b)
	egressTrafficMeter.Mark(int64(n))
	return
}
