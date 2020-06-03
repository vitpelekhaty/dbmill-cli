package commands

import (
	"github.com/vitpelekhaty/dbmill-cli/internal/pkg/filter"
)

// ScriptFoldersOption тип параметра выполнения команды ScriptsFolder
type ScriptFoldersOption func(command IScriptsFolderCommand)

// WithIncludedObjects указывает "движку" при выполнении команды использовать фильтр filter, чтобы ограничить список
// объектов БД указанным пользователем
func WithIncludedObjects(filter filter.IFilter) ScriptFoldersOption {
	return func(command IScriptsFolderCommand) {
		command.SetIncludedObjects(filter)
	}
}

// WithExcludedObjects указывает "движку" при выполнении команды использовать фильтр filter, чтобы исключить из
// обработки указанные пользователем объекты БД
func WithExcludedObjects(filter filter.IFilter) ScriptFoldersOption {
	return func(command IScriptsFolderCommand) {
		command.SetExcludedObjects(filter)
	}
}

// ObjectDefinitionCallback тип callback-функции, вызываемой при чтении определения объекта БД
type ObjectDefinitionCallback func(objectCatalog, objectSchema, objectName, objectType string,
	objectDefinition []byte) error

// WithObjectDefinitionCallback устанавливает callback для чтения определений объектов БД
func WithObjectDefinitionCallback(fn ObjectDefinitionCallback) ScriptFoldersOption {
	return func(command IScriptsFolderCommand) {
		command.SetObjectDefinitionCallback(fn)
	}
}

// WithStaticData опция выгрузки скриптов вставки данных в таблицы
func WithStaticData() ScriptFoldersOption {
	return func(command IScriptsFolderCommand) {
		command.StaticData(true)
	}
}

// WithDecrypt по возможности расшифровывать определения объектов БД
func WithDecrypt() ScriptFoldersOption {
	return func(command IScriptsFolderCommand) {
		command.Decrypt(true)
	}
}

// IScriptsFolderCommand интерфейс команды ScriptsFolder
type IScriptsFolderCommand interface {
	IEngineCommand

	// SetIncludedObjects устанавливает фильтр, позволяющий выбирать только те объекты БД, которые должны быть
	// обработаны
	SetIncludedObjects(filter filter.IFilter)
	// SetExcludedObjects устанавливает фильтр, позволяющий игнорировать объекты БД, которые должны быть заторонуты
	// обработкой
	SetExcludedObjects(filter filter.IFilter)
	// SetObjectDefinitionCallback устанавливает callback для чтения определений объектов БД
	SetObjectDefinitionCallback(callback ObjectDefinitionCallback)
	// StaticData опция выгрузки скриптов вставки данных
	StaticData(on bool)
	// Decrypt по возможности расшифровывать определения объектов БД
	Decrypt(on bool)
}
