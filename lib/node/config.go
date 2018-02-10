package node

import (
	"errors"
	"regexp"
	"time"

	"github.com/centrifugal/centrifugo/lib/channel"
)

// Config contains Application configuration options.
type Config struct {
	// Name of this server node - must be unique, used as human readable
	// and meaningful node identificator.
	Name string

	// Secret is a secret key, used to generate signatures.
	Secret string

	// channel.Options embedded to config.
	channel.Options

	// Namespaces - list of namespaces for custom channel options.
	Namespaces []channel.Namespace

	// NodePingInterval is an interval how often node must send ping
	// control message.
	NodePingInterval time.Duration
	// NodeInfoCleanInterval is an interval in seconds, how often node must
	// clean information about other running nodes.
	NodeInfoCleanInterval time.Duration
	// NodeInfoMaxDelay is an interval in seconds – how many seconds node
	// info considered actual.
	NodeInfoMaxDelay time.Duration
	// NodeMetricsInterval detects interval node will use to aggregate metrics.
	NodeMetricsInterval time.Duration

	// PresencePingInterval is an interval how often connected clients
	// must update presence info.
	PresencePingInterval time.Duration
	// PresenceExpireInterval is an interval how long to consider
	// presence info valid after receiving presence ping.
	PresenceExpireInterval time.Duration

	// PingInterval sets interval server will send ping messages to clients.
	ClientPingInterval time.Duration
	// ClientInsecure turns on insecure mode for client connections - when it's
	// turned on then no authentication required at all when connecting to Centrifugo,
	// anonymous access and publish allowed for all channels, no connection expire
	// performed. This can be suitable for demonstration or personal usage.
	ClientInsecure bool
	// ClientExpire turns on client connection expire mechanism so Centrifugo
	// will close expired connections (if not refreshed).
	ClientExpire bool
	// ExpiredConnectionCloseDelay is an interval given to client to
	// refresh its connection in the end of connection lifetime.
	ClientExpiredCloseDelay time.Duration
	// ClientStaleCloseDelay is an interval in seconds after which
	// connection will be closed if still not authenticated.
	ClientStaleCloseDelay time.Duration
	// MessageWriteTimeout is maximum time of write message operation.
	// Slow client will be disconnected. By default we don't use this option (i.e. it's 0)
	// and slow client connections will be closed when there queue size exceeds
	// ClientQueueMaxSize. In case of SockJS transport we don't have control over it so
	// it only affects raw websocket.
	ClientMessageWriteTimeout time.Duration
	// ClientRequestMaxSize sets maximum size in bytes of allowed client request.
	ClientRequestMaxSize int
	// ClientQueueMaxSize is a maximum size of client's message queue in bytes.
	// After this queue size exceeded Centrifugo closes client's connection.
	ClientQueueMaxSize int
	// ClientChannelLimit sets upper limit of channels each client can subscribe to.
	ClientChannelLimit int

	// UserConnectionLimit limits number of connections from user with the
	// same ID. 0 - unlimited.
	UserConnectionLimit int

	// PrivateChannelPrefix is a prefix in channel name which indicates that
	// channel is private.
	ChannelPrivatePrefix string
	// NamespaceChannelBoundary is a string separator which must be put after
	// namespace part in channel name.
	ChannelNamespaceBoundary string
	// UserChannelBoundary is a string separator which must be set before allowed
	// users part in channel name.
	ChannelUserBoundary string
	// UserChannelSeparator separates allowed users in user part of channel name.
	ChannelUserSeparator string
	// ClientChannelBoundary is a string separator which must be set before client
	// connection ID in channel name so only client with this ID can subscribe on
	// that channel.
	ChannelClientBoundary string
	// ChannelMaxLength is a maximum length of channel name.
	ChannelMaxLength int
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// Validate validates config and returns error if problems found
func (c *Config) Validate() error {
	errPrefix := "config error: "
	pattern := "^[-a-zA-Z0-9_]{2,}$"

	var nss []string
	for _, n := range c.Namespaces {
		name := string(n.Name)
		match, _ := regexp.MatchString(pattern, name)
		if !match {
			return errors.New(errPrefix + "wrong namespace name – " + name)
		}
		if stringInSlice(name, nss) {
			return errors.New(errPrefix + "namespace name must be unique")
		}
		nss = append(nss, name)
	}
	return nil
}

// channelOpts searches for channel options for specified namespace key.
func (c *Config) channelOpts(namespaceName string) (channel.Options, bool) {
	if namespaceName == "" {
		return c.Options, true
	}
	for _, n := range c.Namespaces {
		if n.Name == namespaceName {
			return n.Options, true
		}
	}
	return channel.Options{}, false
}

const (
	// DefaultName of node.
	DefaultName = "centrifugo"
	// DefaultNodePingInterval used in default config.
	DefaultNodePingInterval = 3
)

// DefaultConfig is Config initialized with default values for all fields.
var DefaultConfig = &Config{
	Name: DefaultName,

	NodePingInterval:      DefaultNodePingInterval * time.Second,
	NodeInfoCleanInterval: DefaultNodePingInterval * 3 * time.Second,
	NodeInfoMaxDelay:      DefaultNodePingInterval*2*time.Second + 1*time.Second,
	NodeMetricsInterval:   60 * time.Second,

	PresencePingInterval:   25 * time.Second,
	PresenceExpireInterval: 60 * time.Second,

	ChannelMaxLength:         255,
	ChannelPrivatePrefix:     "$", // so private channel will look like "$gossips"
	ChannelNamespaceBoundary: ":", // so namespace "public" can be used "public:news"
	ChannelUserBoundary:      "#", // so user limited channel is "user#2694" where "2696" is user ID
	ChannelUserSeparator:     ",", // so several users limited channel is "dialog#2694,3019"
	ChannelClientBoundary:    "&", // so client channel is sth like "client&7a37e561-c720-4608-52a8-a964a9db7a8a"

	ClientInsecure:            false,
	ClientMessageWriteTimeout: 0,
	ClientPingInterval:        25 * time.Second,
	ClientExpiredCloseDelay:   25 * time.Second,
	ClientStaleCloseDelay:     25 * time.Second,
	ClientRequestMaxSize:      65536,    // 64KB by default
	ClientQueueMaxSize:        10485760, // 10MB by default
	ClientChannelLimit:        128,
}