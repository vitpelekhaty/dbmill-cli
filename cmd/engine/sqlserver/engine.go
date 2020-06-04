package sqlserver

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/denisenkom/go-mssqldb"

	"github.com/vitpelekhaty/dbmill-cli/cmd/engine/commands"
	"github.com/vitpelekhaty/dbmill-cli/internal/pkg/log"
	"github.com/vitpelekhaty/dbmill-cli/internal/pkg/output"
)

// Engine реализация функциональности утилиты dbmill-cli для MS SQL Server
type Engine struct {
	db     *sql.DB
	logger log.ILogger
	output output.IScriptsFolderOutput
}

const timeout = time.Second * 30

// NewEngine возвращает экземпляр Engine
func NewEngine(connection string) (*Engine, error) {
	db, err := sql.Open("sqlserver", connection)

	if err != nil {
		return nil, err
	}

	err = db.PingContext(context.Background())

	if err != nil {
		return nil, err
	}

	return &Engine{
		db:     db,
		logger: nil,
		output: output.DefaultScriptsFolderOutput,
	}, nil
}

// SetLogger устанавливает логгер событий
func (engine *Engine) SetLogger(logger log.ILogger) {
	engine.logger = logger
}

// SetOutputDirectoryStructure устанавливает описание структуры каталога, где будут созданы скрипты
func (engine *Engine) SetOutputDirectoryStructure(dirStruct output.IScriptsFolderOutput) {
	if dirStruct == nil {
		engine.output = output.DefaultScriptsFolderOutput
	} else {
		engine.output = dirStruct
	}
}

// ScriptsFolder создает скрипты объектов БД по указанному пути path
func (engine *Engine) ScriptsFolder(options ...commands.ScriptsFolderOption) commands.IScriptsFolderCommand {
	return NewScriptsFolderCommand(engine, options...)
}

// Log создает запись в логе, если указан логгер
func (engine *Engine) Log(level log.Level, args ...interface{}) {
	if engine.logger == nil {
		return
	}

	engine.logger.Print(level, args...)
}

// Logf создает фомрматированную запись в логе, если указан логгер
func (engine *Engine) Logf(level log.Level, format string, args ...interface{}) {
	if engine.logger == nil {
		return
	}

	engine.logger.Printf(level, format, args...)
}
