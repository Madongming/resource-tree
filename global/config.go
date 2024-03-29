package global

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

var Configs *Config

type ConfigFileAll struct {
	TestFile    *ConfigFile `yaml:"test"`
	DebugFile   *ConfigFile `yaml:"debug"`
	ReleaseFile *ConfigFile `yaml:"release"`
}

//Parse Yaml File
type ConfigFile struct {
	Log               string     `yaml:"log"`
	Mysql             *MysqlFile `yaml:"mysql"`
	UserCacheSize     int64      `yaml:"user_cache_size"`
	ResourceCacheSize int64      `yaml:"resource_cache_size"`
	GraphCacheSize    int64      `yaml:"graph_cache_size"`
}

type MysqlFile struct {
	User         string `yaml:"user"`
	Pass         string `yaml:"pass"`
	Db           string `yaml:"db"`
	Protocol     string `yaml:"protocol"`
	Address      string `yaml:"address"`
	Params       string `yaml:"params"`
	MaxOpenConns int64  `yaml:"max_open_conns"`
	MaxIdleConns int64  `yaml:"max_idle_conns"`
}

//Userd for Program
type Config struct {
	Log               string
	Mysql             *Mysql
	UserCacheSize     int64
	ResourceCacheSize int64
	GraphCacheSize    int64
}

type Mysql struct {
	User         string
	Pass         string
	Db           string
	Protocol     string
	Address      string
	Params       string
	MaxOpenConns int64
	MaxIdleConns int64
}

func initConfig() {
	cfgFile := os.Getenv("CONFIG_FILE")
	if cfgFile == "" {
		cfgFile = path.Join(filepath.Dir(os.Args[0]), "../etc/config.yml")
	}
	if _, err := os.Stat(cfgFile); err != nil {
		fmt.Fprintf(os.Stderr, "No suitable config file %s, %v", cfgFile, err)
		os.Exit(1)
	}
	yamlConfigAll, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Read config file from %s failed: %s\n", cfgFile, err)
		os.Exit(1)
	}

	//读出所有配置
	configFileAll := &ConfigFileAll{}
	err = yaml.Unmarshal(yamlConfigAll, configFileAll)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse config file from %s failed: %s\n", cfgFile, err)
		os.Exit(1)
	}

	//用于存储根据环境变量确定的配置
	var configFile *ConfigFile
	switch os.Getenv("RESOURCE_TREE_MODE") {
	case "test", "debug":
		configFile = configFileAll.TestFile
	case "release":
		configFile = configFileAll.ReleaseFile
	default:
		configFile = configFileAll.TestFile
	}

	Configs = &Config{
		Log: configFile.Log,
		Mysql: &Mysql{
			User:         configFile.Mysql.User,
			Pass:         configFile.Mysql.Pass,
			Db:           configFile.Mysql.Db,
			Protocol:     configFile.Mysql.Protocol,
			Address:      configFile.Mysql.Address,
			Params:       configFile.Mysql.Params,
			MaxOpenConns: configFile.Mysql.MaxOpenConns,
			MaxIdleConns: configFile.Mysql.MaxIdleConns,
		},
		UserCacheSize:     configFile.UserCacheSize,
		ResourceCacheSize: configFile.ResourceCacheSize,
		GraphCacheSize:    configFile.GraphCacheSize,
	}
}
