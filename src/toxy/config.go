package toxy

import (
	ini "gopkg.in/ini.v1"
)

const (
	METRIC_HANDLER_STATSD          = "statsd"
	METRIC_HANDLER_BUFFERED_STATSD = "buffered_statsd"
)

const (
	PROCESSOR_DEFAULT     = "default"
	PROCESSOR_MULTIPLEXED = "multiplexed"
)

type MetricConfig struct {
	Addr    string `ini:"addr"`
	Prefix  string `ini:"prefix"`
	Handler string `ini:"handler"`
}

type SentryConfig struct {
	Dsn string `ini:"dsn"`
}

type ProxyConfig struct {
	Addr      string `ini:"addr"`
	Processor string `ini:"processor"`
}

type ServiceConfig struct {
	Name        string
	Addr        string `ini:"addr"`
	Transport   string `ini:"transport"`
	Protocol    string `ini:"protocol"`
	Wrapper     string `ini:"wrapper"`
	Timeout     int    `ini:"timeout"`
	Multiplexed bool   `ini:"multiplexed"`
}

type Config struct {
	Metric   *MetricConfig
	Sentry   *SentryConfig
	Proxy    *ProxyConfig
	Services []*ServiceConfig
}

func (c *Config) init_metric(section *ini.Section) (err error) {
	c.Metric = &MetricConfig{
		Addr:    "0.0.0.0:8125",
		Prefix:  "",
		Handler: "buffered_statsd",
	}
	if err = section.MapTo(c.Metric); err != nil {
		return
	}
	return
}

func (c *Config) init_sentry(section *ini.Section) (err error) {
	c.Sentry = &SentryConfig{
		Dsn: "sentry dsn has not been specified",
	}
	if err = section.MapTo(c.Sentry); err != nil {
		return
	}
	return
}

func (c *Config) init_socket_server(section *ini.Section) (err error) {
	if err = section.MapTo(c.Proxy); err != nil {
		return
	}
	return
}

func (c *Config) add_backend_service(section *ini.Section) (err error) {
	var service = &ServiceConfig{
		Name:        section.Name()[8:],
		Addr:        "0.0.0.0:6001",
		Transport:   "socket",
		Wrapper:     "",
		Protocol:    "binary",
		Timeout:     5000,
		Multiplexed: false,
	}
	if err = section.MapTo(service); err != nil {
		return
	}
	c.Services = append(c.Services, service)
	return
}

func LoadConfig(filepath string) (config *Config, err error) {
	var f *ini.File
	var section *ini.Section

	config = &Config{
		Proxy: &ProxyConfig{
			Addr:      "0.0.0.0:6000",
			Processor: "default",
		},
	}

	// load config file
	if f, err = ini.Load(filepath); err != nil {
		return
	}

	if section, err = f.GetSection("metric"); err == nil {
		if err = config.init_metric(section); err != nil {
			return
		}
	}

	if section, err = f.GetSection("sentry"); err == nil {
		if err = config.init_sentry(section); err != nil {
			return
		}
	}

	if section, err = f.GetSection("socketserver"); err == nil {
		if err = config.init_socket_server(section); err != nil {
			return
		}
	}

	for _, section = range f.ChildSections("service") {
		if err = config.add_backend_service(section); err != nil {
			return
		}
	}
	return
}
