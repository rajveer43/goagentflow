package goagentflow

import "log"

type Logger interface {
	Printf(format string, args ...any)
}

type Tracer interface {
	StartSpan(name string) func()
}

type Config struct {
	MaxSteps    int
	Logger      Logger
	Tracer      Tracer
	Memory      Memory
	RetryPolicy RetryPolicy
	Observers   []Observer
}

type Option func(*Config)

func WithMaxSteps(maxSteps int) Option {
	return func(cfg *Config) { cfg.MaxSteps = maxSteps }
}

func WithLogger(logger Logger) Option {
	return func(cfg *Config) { cfg.Logger = logger }
}

func WithTracer(tracer Tracer) Option {
	return func(cfg *Config) { cfg.Tracer = tracer }
}

func WithMemory(memory Memory) Option {
	return func(cfg *Config) { cfg.Memory = memory }
}

func WithRetryPolicy(policy RetryPolicy) Option {
	return func(cfg *Config) { cfg.RetryPolicy = policy }
}

func WithObserver(observer Observer) Option {
	return func(cfg *Config) { cfg.Observers = append(cfg.Observers, observer) }
}

func DefaultConfig() Config {
	return Config{
		MaxSteps:    8,
		Logger:      log.Default(),
		RetryPolicy: DefaultRetryPolicy(),
	}
}
