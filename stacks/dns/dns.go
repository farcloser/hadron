package dns

import (
	"time"

	"github.com/farcloser/hadron/sdk"
	qx "github.com/farcloser/quark/sdk"
)

type Config struct {
	Name     string
	LogLevel string
	Alias    string
}

func DNS(con *sdk.ContainerBuilder, img *qx.Image) *sdk.Container {
	// Security
	con = con.
		ReadOnly().
		CapDrop("ALL").
		Memory("512m").
		MemoryReservation("128m").
		CPUShares(256).
		CPUs("0.25").
		PIDsLimit(50).
		Restart("unless-stopped").
		SecurityOpt("no-new-privileges")

	// Metrics and monitoring
	con = con.Label("prometheus.scrape", "true").
		Label("prometheus.port", "2019").
		HealthCheck(sdk.UDPCheck(53).
			WithTimeout(30 * time.Second).
			WithInterval(30 * time.Second).
			WithRetries(3))

	return con.Image(img.Domain()+"/"+img.Name()+":"+img.Version()+"@"+img.Digest()).
		Port("53:53").
		CapAdd("NET_BIND_SERVICE").
		Env("HEALTHCHECK_URL", "127.0.0.1:53").
		Env("HEALTHCHECK_QUESTION", "dns.autonomous.healthcheck.duncan.st").
		Env("HEALTHCHECK_TYPE", "udp").
		Env("METRICS_LISTEN", "0.0.0.0:4242").
		Env("DNS_PORT", "53").
		Env("DNS_STUFF_MDNS", "false").
		Env("DNS_FORWARD_ENABLED", "true").
		Env("DNS_FORWARD_UPSTREAM_NAME", "cloudflare-dns.com").
		Env("DNS_FORWARD_UPSTREAM_IP_1", "tls://1.1.1.1").
		Env("DNS_FORWARD_UPSTREAM_IP_2", "tls://1.0.0.1").
		Env("DNS_OVER_TLS_ENABLED", "false").
		Env("DNS_OVER_TLS_DOMAIN", "").
		Env("DNS_OVER_TLS_PORT", "").
		Env("DNS_OVER_TLS_LEGO_PORT", "").
		Env("DNS_OVER_TLS_LEGO_EMAIL", "").
		Env("DNS_OVER_TLS_LE_USE_STAGING", "false").
		Build()
}
