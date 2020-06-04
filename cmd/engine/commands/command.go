package commands

// IEngineCommand интерфейс команды "движка"
type IEngineCommand interface {
	// Run запускает выполнение команды
	Run() error
}
