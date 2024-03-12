package config

import (
	"fmt"
	"github.com/cometbft/cometbft/version"
	"net/http"
	"path/filepath"
	"sync"
	"time"
)

// import (
//
//	"sync"
//
// )
//
// // TODO Change this to MUTEX for better performance and Struct Types
const (
	StationRPC           = "http://localhost:8545" //Station RPC
	StationID            = "1"                     //Station Unique ID
	StationName          = "Station1"              // NAME OF THE STATION
	StationType          = "1"                     //EVM, SVM ,COSMWASM
	PODSize              = 25                      // P0D Size
	StationBlockDuration = 5                       // In Seconds
	JunctionRPC          = "1"                     // Junction RPC
	DAType               = "mock"                  // Data Availability Type  : -Eigen , Avail , Celestia,Mock
	DARpc                = "localhost:8080"        // Data Availability RPC
)

//
//type Config struct {
//	mu              sync.RWMutex // Mutex for safe concurrent access
//	LatestPodNumber uint64
//	Peers           []string
//
//	LatestProof          []byte
//	PreviousProof        []byte
//	StationRPC           string
//	StationID            string
//	StationName          string
//	StationType          string
//	PODSize              int
//	StationBlockDuration int
//	JunctionRPC          string
//	DAType               string
//	DARpc                string
//}
//
//// NewConfig creates a new Config instance with default values
//func NewConfig() *Config {
//	return &Config{
//		LatestPodNumber:      0,          // Example default value
//		Peers:                []string{}, // No default peers
//		LatestProof:          nil,        // No default latest proof
//		PreviousProof:        nil,        // No default previous proof
//		StationRPC:           "http://localhost:8545",
//		StationID:            "1",
//		StationName:          "Station1",
//		StationType:          "1",
//		PODSize:              25,
//		StationBlockDuration: 5,
//		JunctionRPC:          "1",
//		DAType:               "mock",
//		DARpc:                "localhost:8080",
//	}
//}
//
////// TODO: SyncConfig should be a method of Config
////func SyncConfig(config *Config) *Config {
////	// sync config detaile from database and nodes, then return right config data
////}
//
////type LatestUnverifiedPodData struct {
////	mu    sync.Mutex
////	count uint64
////}
////
////func (pod *LatestUnverifiedPodData) IncrementUnverifiedPod() {
////	pod.mu.Lock()   // Lock the mutex before accessing count
////	pod.count++     // Critical section: modify count
////	pod.mu.Unlock() // Unlock the mutex after accessing count
////}
////
////func (pod *LatestUnverifiedPodData) ValueUnverifiedPod() uint64 {
////	pod.mu.Lock()         // Lock the mutex before accessing count
////	defer pod.mu.Unlock() // Unlock the mutex after accessing count using defer
////	return pod.count      // Critical section: read count
////}
////
////type LatestUnverifiedProofData struct {
////	mu   sync.Mutex
////	data []byte
////}
////
////func (proof *LatestUnverifiedProofData) UpdateUnverifiedProof(p []byte) {
////	proof.mu.Lock()   // Lock the mutex before accessing count
////	proof.data = p    // Critical section: modify count
////	proof.mu.Unlock() // Unlock the mutex after accessing count
////}
////
////func (proof *LatestUnverifiedProofData) ValueUnverifiedProof() []byte {
////	proof.mu.Lock()         // Lock the mutex before accessing count
////	defer proof.mu.Unlock() // Unlock the mutex after accessing count using defer
////	return proof.data       // Critical section: read count
////}
////
////// call from main with default values... e.g. pod := NewLatestUnverifiedPodData(10) // Start count at 10
////func NewLatestUnverifiedPodData(initialCount uint64) *LatestUnverifiedPodData {
////	return &LatestUnverifiedPodData{
////		count: initialCount,
////	}
////}
////func NewLatestVerified() *LatestUnverifiedProofData {
////	return &LatestUnverifiedProofData{
////		data: nil,
////	}
////}
//
//// Combined struct for unverified pod data and proof
//
//type LatestUnverifiedData struct {
//	Mtx       sync.Mutex
//	Count     uint64
//	ProofData []byte
//}
//
//func NewLatestUnverifiedData(initialCount uint64, initialProofData []byte) *LatestUnverifiedData {
//	return &LatestUnverifiedData{
//		Count:     initialCount,
//		ProofData: initialProofData,
//	}
//}
//
//type ProofUpdater interface {
//	UpdateUnverifiedProof([]byte)
//}
//
//type CountIncrementer interface {
//	IncrementUnverifiedPod()
//}
//
//func (data *LatestUnverifiedData) IncrementUnverifiedPod() {
//	data.Mtx.Lock()
//	data.Count++
//	data.Mtx.Unlock()
//}
//
//func (data *LatestUnverifiedData) UpdateUnverifiedProof(proof []byte) {
//	data.Mtx.Lock()
//	data.ProofData = proof
//	data.Mtx.Unlock()
//}

/////////////////////////////////////////////////////////////////////////

const (
	defaultMoniker                = "tracks"
	DefaultTracksDir              = ".tracks"
	DefaultConfigDir              = "config"
	DefaultDataDir                = "data"
	DefaultConfigFileName         = "sequencer.toml"
	defaultSubscriptionBufferSize = 200
)

var (
	defaultConfigFilePath = filepath.Join(DefaultConfigDir, DefaultConfigFileName)
)

type Config struct {
	BaseConfig `mapstructure:",squash"`
	RPC        *RPCConfig
	P2P        *P2PConfig
	StateSync  *StateSyncConfig
	Consensus  *ConsensusConfig
	DA         *DAConfig
	Station    *StationConfig
	Junction   *JunctionConfig
}

func DefaultConfig() *Config {
	return &Config{
		BaseConfig: DefaultBaseConfig(),
		RPC:        DefaultRPCConfig(),
		P2P:        DefaultP2PConfig(),
		StateSync:  NewStateSyncConfig(),
		Consensus:  DefaultConsensusConfig(),
		DA:         DefaultDAConfig(),
		Station:    DefaultStationConfig(),
		Junction:   DefaultJunctionConfig(),
	}
}

// SetRoot sets the RootDir for all Config structs
func (cfg *Config) SetRoot(root string) *Config {
	fmt.Println("sddsd", root)
	cfg.BaseConfig.RootDir = root
	cfg.RPC.RootDir = root
	cfg.P2P.RootDir = root
	cfg.Consensus.RootDir = root
	fmt.Println(cfg)
	return cfg
}

type BaseConfig struct {
	Version     string `mapstructure:"version"`
	RootDir     string `mapstructure:"home"`
	ProxyApp    string `mapstructure:"proxy_app"`
	Moniker     string `mapstructure:"moniker"`
	DBBackend   string `mapstructure:"db_backend"`
	DBPath      string `mapstructure:"db_dir"`
	FilterPeers bool   `mapstructure:"filter_peers"`
}

func DefaultBaseConfig() BaseConfig {
	return BaseConfig{
		Version:     version.TMCoreSemVer,
		Moniker:     defaultMoniker,
		ProxyApp:    "tcp://127.0.0.1:26658",
		FilterPeers: false,
		DBBackend:   "goleveldb",
		DBPath:      DefaultDataDir,
	}
}

type RPCConfig struct {
	mu                        sync.RWMutex
	RootDir                   string        `mapstructure:"home"`
	ListenAddress             string        `mapstructure:"laddr"`
	CORSAllowedOrigins        []string      `mapstructure:"cors_allowed_origins"`
	CORSAllowedMethods        []string      `mapstructure:"cors_allowed_methods"`
	CORSAllowedHeaders        []string      `mapstructure:"cors_allowed_headers"`
	GRPCListenAddress         string        `mapstructure:"grpc_laddr"`
	GRPCMaxOpenConnections    int           `mapstructure:"grpc_max_open_connections"`
	Unsafe                    bool          `mapstructure:"unsafe"`
	MaxOpenConnections        int           `mapstructure:"max_open_connections"`
	MaxSubscriptionClients    int           `mapstructure:"max_subscription_clients"`
	MaxSubscriptionsPerClient int           `mapstructure:"max_subscriptions_per_client"`
	SubscriptionBufferSize    int           `mapstructure:"experimental_subscription_buffer_size"`
	WebSocketWriteBufferSize  int           `mapstructure:"experimental_websocket_write_buffer_size"`
	CloseOnSlowClient         bool          `mapstructure:"experimental_close_on_slow_client"`
	TimeoutBroadcastTxCommit  time.Duration `mapstructure:"timeout_broadcast_tx_commit"`
	MaxBodyBytes              int64         `mapstructure:"max_body_bytes"`
	MaxHeaderBytes            int           `mapstructure:"max_header_bytes"`
	TLSCertFile               string        `mapstructure:"tls_cert_file"`
	TLSKeyFile                string        `mapstructure:"tls_key_file"`
	PprofListenAddress        string        `mapstructure:"pprof_laddr"`
}

// DefaultRPCConfig returns a default configuration for the RPC server
func DefaultRPCConfig() *RPCConfig {
	return &RPCConfig{
		ListenAddress:          "tcp://127.0.0.1:26657",
		CORSAllowedOrigins:     []string{},
		CORSAllowedMethods:     []string{http.MethodHead, http.MethodGet, http.MethodPost},
		CORSAllowedHeaders:     []string{"Origin", "Accept", "Content-Type", "X-Requested-With", "X-Server-Time"},
		GRPCListenAddress:      "",
		GRPCMaxOpenConnections: 900,

		Unsafe:             false,
		MaxOpenConnections: 900,

		MaxSubscriptionClients:    100,
		MaxSubscriptionsPerClient: 5,
		SubscriptionBufferSize:    defaultSubscriptionBufferSize,
		TimeoutBroadcastTxCommit:  10 * time.Second,
		WebSocketWriteBufferSize:  defaultSubscriptionBufferSize,

		MaxBodyBytes:   int64(1000000), // 1MB
		MaxHeaderBytes: 1 << 20,        // same as the net/http default

		TLSCertFile: "",
		TLSKeyFile:  "",
	}
}

// P2P COnfiguration
type P2PConfig struct {
	RootDir                      string        `mapstructure:"home"`
	ListenAddress                string        `mapstructure:"laddr"`
	ExternalAddress              string        `mapstructure:"external_address"`
	Seeds                        string        `mapstructure:"seeds"`
	PersistentPeers              string        `mapstructure:"persistent_peers"`
	MaxNumInboundPeers           int           `mapstructure:"max_num_inbound_peers"`
	MaxNumOutboundPeers          int           `mapstructure:"max_num_outbound_peers"`
	UnconditionalPeerIDs         string        `mapstructure:"unconditional_peer_ids"`
	PersistentPeersMaxDialPeriod time.Duration `mapstructure:"persistent_peers_max_dial_period"`
	FlushThrottleTimeout         time.Duration `mapstructure:"flush_throttle_timeout"`
	MaxPacketMsgPayloadSize      int           `mapstructure:"max_packet_msg_payload_size"`
	SendRate                     int64         `mapstructure:"send_rate"`
	RecvRate                     int64         `mapstructure:"recv_rate"`
	PexReactor                   bool          `mapstructure:"pex"`
	SeedMode                     bool          `mapstructure:"seed_mode"`
	PrivatePeerIDs               string        `mapstructure:"private_peer_ids"`
	AllowDuplicateIP             bool          `mapstructure:"allow_duplicate_ip"`
	HandshakeTimeout             time.Duration `mapstructure:"handshake_timeout"`
	DialTimeout                  time.Duration `mapstructure:"dial_timeout"`
}

func DefaultP2PConfig() *P2PConfig {
	return &P2PConfig{
		ListenAddress:                "tcp://0.0.0.0:2300",
		ExternalAddress:              "", // Empty means it will be set automatically.
		MaxNumInboundPeers:           40,
		MaxNumOutboundPeers:          10,
		PersistentPeersMaxDialPeriod: 0, // No delay by default, can use exponential backoff.
		FlushThrottleTimeout:         100 * time.Millisecond,
		MaxPacketMsgPayloadSize:      1024,      // 1 KB
		SendRate:                     5_120_000, // 5 MB/s
		RecvRate:                     5_120_000, // 5 MB/s
		PexReactor:                   true,
		SeedMode:                     false,
		AllowDuplicateIP:             false,
		HandshakeTimeout:             20 * time.Second,
		DialTimeout:                  3 * time.Second,
	}
}

// StateSyncConfig holds configuration settings related to syncing pods.
type StateSyncConfig struct {
	Enable            bool          `mapstructure:"enable"`              // Enable or disable pod syncing
	TempDir           string        `mapstructure:"temp_dir"`            // Directory for temporary storage during pod syncing
	RPCServers        []string      `mapstructure:"rpc_servers"`         // List of RPC servers for fetching pods
	PodTrustPeriod    time.Duration `mapstructure:"pod_trust_period"`    // Period for which a pod is considered trusted
	PodTrustHeight    int64         `mapstructure:"pod_trust_height"`    // Height at which the pod's trust starts
	PodTrustHash      string        `mapstructure:"pod_trust_hash"`      // Hash of a trusted pod to start syncing from
	PodDiscoveryTime  time.Duration `mapstructure:"pod_discovery_time"`  // Time for discovering new pods
	PodRequestTimeout time.Duration `mapstructure:"pod_request_timeout"` // Timeout for pod requests
	PodChunkFetchers  int32         `mapstructure:"pod_chunk_fetchers"`  // Number of concurrent fetchers for pod chunks
}

// NewStateSyncConfig creates a new instance of StateSyncConfig with default values.
func NewStateSyncConfig() *StateSyncConfig {
	return &StateSyncConfig{
		Enable:            false,
		TempDir:           "./podsync_temp",
		RPCServers:        []string{},
		PodTrustPeriod:    168 * time.Hour, // 7 days
		PodTrustHeight:    0,
		PodTrustHash:      "",
		PodDiscoveryTime:  15 * time.Minute,
		PodRequestTimeout: 5 * time.Second,
		PodChunkFetchers:  4,
	}
}

// ConsensusConfig holds the configuration options for the consensus layer.
type ConsensusConfig struct {
	RootDir string `mapstructure:"home"`
	WalPath string `mapstructure:"wal_file"`
	walFile string // unexported, overrides WalPath if set

	// Timing configurations for the pod consensus process
	TimeoutPropose        time.Duration `mapstructure:"timeout_propose"`
	TimeoutProposeDelta   time.Duration `mapstructure:"timeout_propose_delta"`
	TimeoutPrevote        time.Duration `mapstructure:"timeout_prevote"`
	TimeoutPrevoteDelta   time.Duration `mapstructure:"timeout_prevote_delta"`
	TimeoutPrecommit      time.Duration `mapstructure:"timeout_precommit"`
	TimeoutPrecommitDelta time.Duration `mapstructure:"timeout_precommit_delta"`
	TimeoutCommit         time.Duration `mapstructure:"timeout_commit"`

	// Configuration to skip the commit timeout for faster consensus on pods
	SkipTimeoutCommit bool `mapstructure:"skip_timeout_commit"`

	// Pod-specific configurations
	ValidatePods               bool          `mapstructure:"validate_pods"`                 // Whether to validate pods before accepting them
	PodValidationSleepDuration time.Duration `mapstructure:"pod_validation_sleep_duration"` // Sleep duration between pod validations

	DoubleSignCheckHeight int64 `mapstructure:"double_sign_check_height"` // Height to check for double signing
}

// DefaultConsensusConfig returns a default configuration for the consensus service.
func DefaultConsensusConfig() *ConsensusConfig {

	return &ConsensusConfig{
		WalPath:                    filepath.Join(DefaultDataDir, "pods.wal", "wal"),
		TimeoutPropose:             3 * time.Second,
		TimeoutProposeDelta:        500 * time.Millisecond,
		TimeoutPrevote:             1 * time.Second,
		TimeoutPrevoteDelta:        500 * time.Millisecond,
		TimeoutPrecommit:           1 * time.Second,
		TimeoutPrecommitDelta:      500 * time.Millisecond,
		TimeoutCommit:              1 * time.Second,
		SkipTimeoutCommit:          false,
		ValidatePods:               true,
		PodValidationSleepDuration: 100 * time.Millisecond,
		DoubleSignCheckHeight:      0,
	}
}

type DAConfig struct {
}

// DefaultDAConfig returns a default configuration for the data availability layer.
func DefaultDAConfig() *DAConfig {
	return &DAConfig{}
}

type StationConfig struct {
}

// DefaultStationConfig returns a default configuration for the station.
func DefaultStationConfig() *StationConfig {
	return &StationConfig{}
}

type JunctionConfig struct {
}

// DefaultJunctionConfig returns a default configuration for the junction.
func DefaultJunctionConfig() *JunctionConfig {
	return &JunctionConfig{}
}
