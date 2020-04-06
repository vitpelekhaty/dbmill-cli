package commands

import (
	"bufio"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/vitpelekhaty/dbmill-cli/cmd/engine"
	"github.com/vitpelekhaty/dbmill-cli/cmd/input"
	"github.com/vitpelekhaty/dbmill-cli/cmd/scriptsfolder"
	"github.com/vitpelekhaty/dbmill-cli/internal/pkg/dir"
	"github.com/vitpelekhaty/dbmill-cli/internal/pkg/filter"
	"github.com/vitpelekhaty/dbmill-cli/internal/pkg/log"
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

		options := make([]engine.DatabaseOption, 0)

		if logger != nil {
			options = append(options, engine.WithLogger(logger))
		}

		if include != nil {
			options = append(options, engine.WithIncludedObjects(include))
		}

		if exclude != nil {
			options = append(options, engine.WithExcludedObjects(exclude))
		}

		if outputDirStruct != nil {
			options = append(options, engine.WithOutputDirStructure(outputDirStruct))
		}

		return scriptsfolder.Run(Database, Path, IncludeData, Decrypt, options...)
	},
}

// OutputDirectoryStructure возвращает описание структуры директории, в которой будут создаваться скрипты объектов БД
func OutputDirectoryStructure(path string) (dir.IStructure, error) {
	if strings.Trim(path, " ") == "" {
		return dir.Default, nil
	}

	f, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer f.Close()

	return dir.NewStructure(f)
}

// ObjectFilter возвращает настроенный фильтр объектов БД
func ObjectFilter(path string, expressions []string) (filter.IFilter, error) {
	if len(expressions) > 0 {
		return filter.New(expressions)
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
