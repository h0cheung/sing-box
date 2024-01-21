package adapter

import (
	"context"
	"net/http"
	"net/netip"
	"time"

	"github.com/sagernet/sing-box/common/geoip"
	"github.com/sagernet/sing-dns"
	"github.com/sagernet/sing-tun"
	"github.com/sagernet/sing/common/control"
	N "github.com/sagernet/sing/common/network"
	"github.com/sagernet/sing/service"

	mdns "github.com/miekg/dns"
)

type Router interface {
	Service
	PreStarter
	PostStarter

	Outbounds() []Outbound
	Outbound(tag string) (Outbound, bool)
	OutboundsWithProvider() []Outbound
	OutboundWithProvider(tag string) (Outbound, bool)
	DefaultOutbound(network string) (Outbound, error)

	OutboundProviders() []OutboundProvider
	OutboundProvider(tag string) (OutboundProvider, bool)

	FakeIPStore() FakeIPStore

	ConnectionRouter

	GeoIPReader() *geoip.Reader
	LoadGeosite(code string) (Rule, error)

	RuleSets() []RuleSet
	RuleSet(tag string) (RuleSet, bool)

	NeedWIFIState() bool

	Exchange(ctx context.Context, message *mdns.Msg) (*mdns.Msg, error)
	Lookup(ctx context.Context, domain string, strategy dns.DomainStrategy) ([]netip.Addr, error)
	LookupDefault(ctx context.Context, domain string) ([]netip.Addr, error)
	ClearDNSCache()

	InterfaceFinder() control.InterfaceFinder
	UpdateInterfaces() error
	DefaultInterface() string
	AutoDetectInterface() bool
	AutoDetectInterfaceFunc() control.Func
	DefaultMark() int
	NetworkMonitor() tun.NetworkUpdateMonitor
	InterfaceMonitor() tun.DefaultInterfaceMonitor
	PackageManager() tun.PackageManager
	WIFIState() WIFIState
	Rules() []Rule

	ClashServer() ClashServer
	SetClashServer(server ClashServer)

	V2RayServer() V2RayServer
	SetV2RayServer(server V2RayServer)

	ResetNetwork() error
}

func ContextWithRouter(ctx context.Context, router Router) context.Context {
	return service.ContextWith(ctx, router)
}

func RouterFromContext(ctx context.Context) Router {
	return service.FromContext[Router](ctx)
}

type HeadlessRule interface {
	Match(metadata *InboundContext) bool
	RuleCount() int
	UseIPRule() bool
}

type Rule interface {
	HeadlessRule
	Service
	Type() string
	UpdateGeosite() error
	SkipResolve() bool
	UseIPRule() bool
	Outbound() string
	String() string
}

type DNSRule interface {
	Rule
	DisableCache() bool
	RewriteTTL() *uint32
	ClientSubnet() *netip.Addr
	WithAddressLimit() bool
	MatchAddressLimit(metadata *InboundContext) bool
}

type RuleSet interface {
	Tag() string
	Type() string
	Format() string
	UpdatedTime() time.Time
	Update(router Router) error
	StartContext(ctx context.Context, startContext RuleSetStartContext) error
	PostStart() error
	Metadata() RuleSetMetadata
	Close() error
	HeadlessRule
}

type RuleSetMetadata struct {
	ContainsProcessRule bool
	ContainsWIFIRule    bool
	ContainsIPCIDRRule  bool
}

type RuleSetStartContext interface {
	HTTPClient(detour string, dialer N.Dialer) *http.Client
	Close()
}

type InterfaceUpdateListener interface {
	InterfaceUpdated()
}

type WIFIState struct {
	SSID  string
	BSSID string
}
