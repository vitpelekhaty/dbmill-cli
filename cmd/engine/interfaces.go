package engine

type DatabaseOption func(engine IDatabase)

// IDatabase интерфейс базы данных
type IDatabase interface {
	SetLogger(logger interface{})
	ScriptsFolder(path string, includeData, decrypt bool) error
}

// WithLogger устанавливает
func WithLogger(logger interface{}) DatabaseOption {
	return func(engine IDatabase) {
		engine.SetLogger(logger)
	}
}
