package ancore

type InitOptions struct {
	LogPath    string
	ConfigPath string
	Debug      *bool
}

type Option func(*InitOptions)

func WithLogPath(p string) Option {
	return func(o *InitOptions) { o.LogPath = p }
}

func WithConfigPath(p string) Option {
	return func(o *InitOptions) { o.ConfigPath = p }
}

func WithDebug(b bool) Option {
	return func(o *InitOptions) { o.Debug = &b }
}

func ptrBool(b bool) *bool {
	return &b
}
