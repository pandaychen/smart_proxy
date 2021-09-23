package config

//smartproxy配置定义

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

//全局配置，使用viper实现配置文件hot-reload
var g_Config *SmartProxyConfig

func GetSmartproxyConf() *SmartProxyConfig {
	return g_Config
}

// Config is application global config
type SmartProxyConfig struct {
	ProjectName    string           `mapstructure:"project"`
	MysqlConf      MysqlConfig      `mapstructure:"mysql"`
	RedisConf      RedisConfig      `mapstructure:"redis"`
	LogConfig      LoggerConfig     `mapstructure:"loger"`
	ControllerConf ControllerConfig `mapstructure:"manager"`
	MetricsConf    MetricsConfig    `mapstructure:"metrics"`
	//DiscoveryListConf    []DiscoveryConfig    `mapstructure:"discovery"`
	DiscoveryConf        DiscoveryConfig      `mapstructure:"discovery"`
	ReverseProxyListConf []ReverseProxyConfig `mapstructure:"reverseproxy_group"` //support multi proxys
}

type PoolConfig struct {
	Address string `mapstructure:"address"`
	Weight  int    `mapstructure:"weight"`
}

type ReverseProxyConfig struct {
	ProxyName    string       `mapstructure:"name"`
	BindAddr     string       `mapstructure:"bind_addr"`
	TlsOn        bool         `mapstructure:"tls"`
	Key          string       `mapstructure:"key"`
	Cert         string       `mapstructure:"cert"`
	LbType       string       `mapstructure:"lbtype"`
	SingnatureOn bool         `mapstructure:"singnature"`
	DnsName      string       `mapstructure:"dns_name"`
	PoolConfList []PoolConfig `mapstructure:"pool"`
}

//采用dns时，服务名字定义在reverseproxy_group中
type DiscoveryConfig struct {
	DiscoveryType string `mapstructure:"type"`
	ClusterAddr   string `mapstructure:"cluster"`
	RootPrefix    string `mapstructure:"prefix"` //watcher root prefix
	Name          string `mapstructure:"name"`
}

type ControllerConfig struct {
	Host           string `mapstructure:"host"`
	Port           int    `mapstructure:"port"`
	ControllerType string `mapstructure:"type"`
	Username       string `mapstructure:"username"`
	Password       string `mapstructure:"password"`
}

type MetricsConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type MysqlConfig struct {
	Dbname          string `mapstructure:"dbname"`
	Host            string `mapstructure:"host"`
	Port            string `mapstructure:"port"`
	Username        string `mapstructure:"username"`
	Password        string `mapstructure:"password"`
	MaximumPoolSize int    `mapstructure:"maximum-pool-size"`
	MaximumIdleSize int    `mapstructure:"maximum-idle-size"`
	//LogMode         bool   `mapstructure:"log-mode"`
}

type RedisConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	Password     string `mapstructure:"password"`
	Db           int    `mapstructure:"db"`
	PoolSize     int    `mapstructure:"pool-size"`
	MinIdleConns int    `mapstructure:"min-idle-conns"`
	IdleTimeout  int    `mapstructure:"idle-timeout"`
}

//ZAP日志配置
type LoggerConfig struct {
	Level      string `mapstructure:"level"`
	FileName   string `mapstructure:"file-name"`
	TimeFormat string `mapstructure:"time-format"`
	MaxSize    int    `mapstructure:"max-size"`
	MaxBackups int    `mapstructure:"max-backups"`
	MaxAge     int    `mapstructure:"max-age"`
	Compress   bool   `mapstructure:"compress"`
	LocalTime  bool   `mapstructure:"local-time"`
	Console    bool   `mapstructure:"console"`
}

// config loader
func LoadSmartproxyConfig(configFilePath string) *SmartProxyConfig {
	//set config path
	setConfigPath(configFilePath)

	//init config
	if err := initConfig(); err != nil {
		panic(err)
	}

	//watcher
	watchConfig()

	return g_Config
}

func initConfig() error {
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	viper.SetEnvPrefix("APPLICATION")
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	// parser config file
	g_Config = &SmartProxyConfig{}
	if err := viper.Unmarshal(g_Config); err != nil {
		panic(err)
	}
	return nil
}

// monitor file modify event
func watchConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		log.Printf("Config file changed: %s, will reload it", in.Name)
		// reload config files
		LoadSmartproxyConfig(in.Name)
	})
}

func setConfigPath(path string) {
	if path != "" {
		viper.SetConfigFile(path)
	} else {
		// set default config path
		viper.AddConfigPath("./conf")
		viper.SetConfigName("smartproxy.yaml")
	}
}

func (c SmartProxyConfig) CheckReverseproxyValid() bool {

	return true
}

func main() {
	g_Config = LoadSmartproxyConfig("smartproxy.yaml")
	fmt.Println(GetSmartproxyConf())
	for {
		//wait to watch config changed
		time.Sleep(1 * time.Second)
		fmt.Println(GetSmartproxyConf())
	}
}
