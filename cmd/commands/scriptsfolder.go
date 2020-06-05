package commands

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/vitpelekhaty/dbmill-cli/cmd/engine"
	"github.com/vitpelekhaty/dbmill-cli/cmd/engine/commands"
	"github.com/vitpelekhaty/dbmill-cli/cmd/input"
	"github.com/vitpelekhaty/dbmill-cli/internal/pkg/filter"
	"github.com/vitpelekhaty/dbmill-cli/internal/pkg/log"
	"github.com/vitpelekhaty/dbmill-cli/internal/pkg/output"
)

// cmdScriptsFolder команда создания скриптов на основе схемы
var cmdScriptsFolder = &cobra.Command{
	Use:   "scriptsfolder",
	Short: "creates scripts based on the schema",
	RunE: func(cmd *cobra.Command, args []string) error {
		if strings.Trim(Username, " ") == "" || strings.Trim(Password, " ") == "" {
			user, pwd, err := engine.Credentials(Database)

			if err != nil {
				return err
			}

			if strings.Trim(Username, " ") == "" {
				Username = user
			}

			if strings.Trim(Password, " ") == "" {
				Password = pwd
			}
		}

		if strings.Trim(Username, " ") == "" {
			Username = input.Username()
		}

		if strings.Trim(Password, " ") == "" {
			Password = input.Password()
		}

		Database, err := engine.SetCredentials(Database, Username, Password)

		if err != nil {
			return err
		}

		var logger log.ILogger

		if strings.Trim(LogFilename, " ") != "" {
			var logLevel = log.InfoLevel

			if strings.Trim(LogLevel, " ") != "" {
				logLevel, err = ParseLogLevel(LogLevel)

				if err != nil {
					return err
				}
			}

			fl, err := os.Create(LogFilename)

			if err != nil {
				return err
			}

			defer fl.Close()

			logger = log.New(log.WithLevel(logLevel), log.WithOutput(fl))
		}

		include, err := ObjectFilter(FilterPath, Filter)

		if err != nil {
			return err
		}

		exclude, err := ObjectFilter(ExcludePath, Exclude)

		if err != nil {
			return err
		}

		outputDirStruct, err := OutputDirectoryStructure(DirStructFilename)

		if err != nil {
			return err
		}

		engineOptions := make([]engine.Option, 0)
		commandOptions := make([]commands.ScriptsFolderOption, 0)

		commandOptions = append(commandOptions, commands.WithObjectDefinitionCallback(
			func(objectCatalog, objectSchema, objectName string, objectType output.DatabaseObjectType,
				objectDefinition []byte) error {
				return SaveDefinition(Path, objectCatalog, objectSchema, objectName, objectType, objectDefinition,
					outputDirStruct)
			}))

		commandOptions = append(commandOptions, commands.WithDatabaseObjectTypes(outputDirStruct.DatabaseObjects()))

		if logger != nil {
			engineOptions = append(engineOptions, engine.WithLogger(logger))
		}

		if include != nil {
			commandOptions = append(commandOptions, commands.WithIncludedObjects(include))
		}

		if exclude != nil {
			commandOptions = append(commandOptions, commands.WithExcludedObjects(exclude))
		}

		if Decrypt {
			commandOptions = append(commandOptions, commands.WithDecrypt())
		}

		if IncludeData {
			commandOptions = append(commandOptions, commands.WithStaticData())
		}

		engn, err := engine.New(Database, engineOptions...)

		if err != nil {
			return err
		}

		return engn.ScriptsFolder(commandOptions...).Run()
	},
}

// OutputDirectoryStructure возвращает описание структуры директории, в которой будут создаваться скрипты объектов БД
func OutputDirectoryStructure(path string) (output.IScriptsFolderOutput, error) {
	if strings.Trim(path, " ") == "" {
		return output.DefaultScriptsFolderOutput, nil
	}

	f, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer f.Close()

	return output.NewScriptsFolderOutput(f)
}

// ObjectFilter возвращает настроенный фильтр объектов БД
func ObjectFilter(path string, expressions []string) (filter.IFilter, error) {
	if len(expressions) > 0 {
		return filter.New(expressions)
	}

	if strings.Trim(path, " ") == "" {
		return nil, nil
	}

	f, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	var exp []string

	for scanner.Scan() {
		exp = append(exp, scanner.Text())
	}

	return filter.New(exp)
}

// SaveDefinition сохраняет определение объекта БД в скрипт
func SaveDefinition(path string, objectCatalog, objectSchema, objectName string, objectType output.DatabaseObjectType,
	objectDefinition []byte, rules output.IScriptsFolderOutput) error {
	subdirectory, filename, ok := rules.Rules(objectType)

	if !ok {
		return fmt.Errorf("no ouput rules for %v", objectType)
	}

	filename = strings.ReplaceAll(filename, "$schema$", objectSchema)
	filename = strings.ReplaceAll(filename, "$object$", objectName)
	filename = strings.ReplaceAll(filename, "$database$", objectCatalog)
	filename = strings.ReplaceAll(filename, "$type$", objectType.String())

	err := os.MkdirAll(filepath.Join(path, subdirectory), os.ModePerm)

	if err != nil {
		return err
	}

	return ioutil.WriteFile(filepath.Join(filepath.Join(path, subdirectory), filename), objectDefinition, 0664)
}
