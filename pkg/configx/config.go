package configx

type Options struct {
	conf Config
}

type Option func(*Options)

type Config struct {
	Path  string
	Type  ConfigType
	Watch []func(string, error)
	c     interface{}
}

func WithLocal(path string, c interface{}) Option {
	return func(k *Options) {
		k.conf.Path = path
		k.conf.c = c
		k.conf.Type = Local
	}
}

func WithNacos(c interface{}) Option {
	return func(k *Options) {
		k.conf.c = c
		k.conf.Type = Nacos
	}
}

func WithWatch(watch ...func(string, error)) Option {
	return func(k *Options) {
		k.conf.Watch = watch
	}
}

func Must(opts ...Option) error {
	e := &Options{}
	for _, opt := range opts {
		opt(e)
	}
	if e.conf.Type == "" && e.conf.Path != "" {
		e.conf.Type = Local
	}
	return LoadConfig(e.conf.Path, e.conf.Type, e.conf.c, e.conf.Watch...)
}
