package config

import (
	"os"
	"strings"
	"sync"
	"time"

	"github.com/namsral/flag"
	"github.com/pkg/errors"
)

// Flag defaults.
const (
	defaultHost            = "0.0.0.0"
	defaultLogLevel        = "info"
	defaultPort            = 8765
	defaultAppName         = "url-shortener"
	defaultShutdownTimeout = 10 * time.Second
	defaultDBDriver        = "postgres"
	defaultDSN             = "user=.. password=.. host=.. dbname=.. sslmode=disable"
)

var flagErrorHandling = flag.ContinueOnError

type Config struct {
	AppName  string        `json:"appName"`
	LogLevel string        `json:"logLevel"`
	HTTP     HTTPConfig    `json:"http"`
	Storage  StorageConfig `json:"storage"`
}

func (c Config) Validate() []error {
	errs := c.HTTP.Validate()
	errs = append(errs, c.Storage.Validate()...)
	if strings.TrimSpace(c.LogLevel) == "" {
		errs = append(errs, errors.New("Config requires a non-empty LogLevel config value"))
	}
	return errs
}

type HTTPConfig struct {
	Host            string        `json:"host"`
	Port            int           `json:"port"`
	ShutdownTimeout time.Duration `json:"shutdownTimeout"`
}

func (c HTTPConfig) Validate() []error {
	var errs []error

	if len(c.Host) == 0 {
		errs = append(errs, errors.Errorf("HTTPConfig requires a non-empty Host config value"))
	}

	if c.Port <= 0 {
		errs = append(errs, errors.Errorf("HTTPConfig requires a postive Port config value"))
	}

	return errs
}

type StorageConfig struct {
	Driver string `json:"driver"`
	DSN    string `json:"dsn"`
}

func (c StorageConfig) Validate() []error {
	var errs []error

	if len(c.Driver) == 0 {
		errs = append(errs, errors.New("db driver cannot be empty"))
	}

	if len(c.DSN) == 0 {
		errs = append(errs, errors.New("db DSN cannot be empty"))
	}

	return errs
}

func initConfig() (*Config, error) {
	var config Config
	flagset := flag.NewFlagSetWithEnvPrefix(defaultAppName, "", flagErrorHandling)

	// App
	flagset.StringVar(&config.AppName, "app-name", defaultAppName, "Service name.")
	flagset.StringVar(&config.LogLevel, "log-level", defaultLogLevel, "Log level (debug, info, warn, error)")

	// HTTP
	flagset.StringVar(&config.HTTP.Host, "host", defaultHost, "Host part of listening address.")
	flagset.IntVar(&config.HTTP.Port, "port", defaultPort, "Listening port.")
	flagset.DurationVar(&config.HTTP.ShutdownTimeout, "shutdown-timeout", defaultShutdownTimeout, "Shutdown timeout for http http.")

	//DB
	flagset.StringVar(&config.Storage.Driver, "db-driver", defaultDBDriver, "Data http driver.")
	flagset.StringVar(&config.Storage.DSN, "db-dsn", defaultDSN, "Data http data source name.")

	if err := flagset.Parse(os.Args[1:]); err != nil {
		return nil, errors.Wrap(err, "parsing flags")
	}

	// Validate the config.
	if errs := config.Validate(); len(errs) > 0 {
		return nil, errors.Errorf("invalid flag(s): %s", errs)
	}
	return &config, nil
}

var cfg *Config
var once sync.Once

func GetConfig() (*Config, error) {
	err := error(nil)
	if cfg == nil {
		once.Do(func() {
			cfg, err = initConfig()
		})
	}
	return cfg, err
}
