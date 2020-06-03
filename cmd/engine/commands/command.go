package commands

// IEngineCommand интерфейс команды "движка"
type IEngineCommand interface {
	Run() error
}
