package logger

import (
	"github.com/FactomProject/logrustash"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"strings"
	"time"
)

const production = "production"

//Alias for map used in withfields methods
type Fields map[string]interface{}

type logger struct {
	logger      *logrus.Logger
	serviceName string
	dFields     map[string]interface{}
}

type entry struct {
	entry       *logrus.Entry
	serviceName string
	dFields     map[string]interface{}
}

type Logger interface {
	WithFields(map[string]interface{}) Logger
	Debug(...interface{})
	Info(...interface{})
	Warn(...interface{})
	Error(...interface{})
	Fatal(...interface{})
	Panic(...interface{})
}

type Config struct {
	ServiceName     string
	EnvironmentName string
	LogstashServer  string
	LogstashPort    string
	LogLevel        string
	DefaultFields   map[string]interface{}
}

func New(config *Config) Logger {
	validate(config)
	newLogger := &logger{
		logger:      logrus.StandardLogger(),
		serviceName: config.ServiceName,
		dFields:     config.DefaultFields,
	}
	configure(config)
	return newLogger

}
func validate(config *Config) {
	if strings.TrimSpace(config.ServiceName) == "" {
		log.Fatal("Required attribute 'ServiceName' missing on logger config")
	}
}

func (l *logger) WithFields(fields map[string]interface{}) Logger {
	return &entry{
		entry:       l.logger.WithFields(collectFields(l.dFields, fields, l.serviceName)),
		serviceName: l.serviceName,
		dFields:     l.dFields,
	}
}

func (l *logger) Debug(message ...interface{}) {
	l.logger.WithFields(collectFields(l.dFields, map[string]interface{}{}, l.serviceName)).Debug(message)
}

func (l *logger) Info(message ...interface{}) {
	l.logger.WithFields(collectFields(l.dFields, map[string]interface{}{}, l.serviceName)).Info(message)
}

func (l *logger) Warn(message ...interface{}) {
	l.logger.WithFields(collectFields(l.dFields, map[string]interface{}{}, l.serviceName)).Warn(message)
}

func (l *logger) Error(message ...interface{}) {
	l.logger.WithFields(collectFields(l.dFields, map[string]interface{}{}, l.serviceName)).Error(message)
}

func (l *logger) Fatal(message ...interface{}) {
	l.logger.WithFields(collectFields(l.dFields, map[string]interface{}{}, l.serviceName)).Fatal(message)
}

func (l *logger) Panic(message ...interface{}) {
	l.logger.WithFields(collectFields(l.dFields, map[string]interface{}{}, l.serviceName)).Panic(message)
}

func (e *entry) WithFields(fields map[string]interface{}) Logger {
	e.entry.WithFields(collectFields(e.dFields, fields, e.serviceName))
	return e
}

func (e *entry) Debug(message ...interface{}) {
	e.entry.WithFields(collectFields(e.dFields, map[string]interface{}{}, e.serviceName)).Debug(message)
}

func (e *entry) Info(message ...interface{}) {
	e.entry.WithFields(collectFields(e.dFields, map[string]interface{}{}, e.serviceName)).Info(message)
}

func (e *entry) Warn(message ...interface{}) {
	e.entry.WithFields(collectFields(e.dFields, map[string]interface{}{}, e.serviceName)).Warn(message)
}

func (e *entry) Error(message ...interface{}) {
	e.entry.WithFields(collectFields(e.dFields, map[string]interface{}{}, e.serviceName)).Error(message)
}

func (e *entry) Fatal(message ...interface{}) {
	e.entry.WithFields(collectFields(e.dFields, map[string]interface{}{}, e.serviceName)).Fatal(message)
}

func (e *entry) Panic(message ...interface{}) {
	e.entry.WithFields(collectFields(e.dFields, map[string]interface{}{}, e.serviceName)).Panic(message)
}

func collectFields(a map[string]interface{}, b map[string]interface{}, serviceName string) map[string]interface{} {
	var allFields = make(map[string]interface{}, len(a)+len(b)+1)
	for k, v := range a {
		allFields[k] = v
	}
	for k, v := range b {
		allFields[k] = v
	}
	allFields["service_name"] = serviceName
	return allFields
}

func configure(configuration *Config) {

	logrus.SetLevel(getLevel(configuration.LogLevel))
	logrus.SetFormatter(getFormatter(configuration.EnvironmentName))

	if hook := createHook(configuration); hook != nil {
		logrus.AddHook(hook)
	}
}

func getLevel(logLevel string) logrus.Level {
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		return logrus.InfoLevel
	}
	return level
}

func getFormatter(environment string) logrus.Formatter {
	envType := valueOrDefault(environment, "development")

	if envType == production {
		return &logrus.JSONFormatter{}
	}
	return &logrus.TextFormatter{}
}

func createHook(configuration *Config) logrus.Hook {
	logstashServer := configuration.LogstashServer
	if strings.TrimSpace(logstashServer) == "" {
		return nil
	}

	logstashPort := valueOrDefault(configuration.LogstashPort, "5000")
	address := logstashServer + ":" + logstashPort

	hook, err := logrustash.NewAsyncHook("tcp", address, configuration.ServiceName)

	if err != nil {
		return nil
	}

	hook.ReconnectBaseDelay = time.Second
	hook.ReconnectDelayMultiplier = 2
	hook.MaxReconnectRetries = 10

	return hook
}

func valueOrDefault(name string, defValue string) string {
	v := strings.TrimSpace(os.Getenv(name))

	if v == "" {
		return defValue
	} else {
		return v
	}
}
