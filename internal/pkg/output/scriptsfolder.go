package output

import (
	"fmt"
	"io"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// ScriptsFolderOutputRule элемент структуры каталога скриптов
type ScriptsFolderOutputRule struct {
	// SubDirectory подкаталог
	SubDirectory string `yaml:"subdirectory"`
	// FilenameMask маска имени файла
	FilenameMask string `yaml:"mask"`
}

// IScriptsFolderOutput интерфес описания структуры каталога скриптов
type IScriptsFolderOutput interface {
	// Rules возвращает целевой каталог и маску имени файла для указанного типа объекта itemType.
	// Если информация не найдена, то в параметре ok возвращается false, в противном случае - true
	Rules(objectType DatabaseObjectType) (subdirectory, mask string, ok bool)
}

// ScriptsFolderOutput описание структуры каталога скриптов
type ScriptsFolderOutput struct {
	rules map[DatabaseObjectType]ScriptsFolderOutputRule
}

// DefaultScriptsFolderOutput структура каталога скриптов по умолчанию
var DefaultScriptsFolderOutput *ScriptsFolderOutput

// NewScriptsFolderOutput конструктор описания структуры каталога скриптов
func NewScriptsFolderOutput(in io.Reader) (*ScriptsFolderOutput, error) {
	data, err := ioutil.ReadAll(in)

	if err != nil {
		return nil, err
	}

	ss, err := parse(data)

	if err != nil {
		return nil, err
	}

	si, err := mapping(ss)

	if err != nil {
		return nil, err
	}

	return &ScriptsFolderOutput{rules: si}, nil
}

// Rules возвращает целевой каталог и маску имени файла для указанного типа объекта itemType.
// Если информация не найдена, то в параметре ok возвращается false, в противном случае - true
func (output *ScriptsFolderOutput) Rules(item DatabaseObjectType) (subdirectory, mask string, ok bool) {
	i, ok := output.rules[item]

	if ok {
		subdirectory = i.SubDirectory
		mask = i.FilenameMask
	}

	return
}

func parse(data []byte) (map[string]ScriptsFolderOutputRule, error) {
	var s map[string]ScriptsFolderOutputRule

	err := yaml.Unmarshal(data, &s)

	if err != nil {
		return s, err
	}

	return s, nil
}

func mapping(rules map[string]ScriptsFolderOutputRule) (map[DatabaseObjectType]ScriptsFolderOutputRule, error) {
	si := make(map[DatabaseObjectType]ScriptsFolderOutputRule)

	if len(rules) == 0 {
		return si, nil
	}

	for key, value := range rules {
		if k, ok := databaseObjectTypeMappingReverse[key]; ok {
			si[k] = value
		} else {
			return si, fmt.Errorf("unknown object type %s", key)
		}
	}

	return si, nil
}

const defaultScriptsFolderOutput = `
database:
  subdirectory: Database/
  mask: "$object$.sql"

table:
  subdirectory: Tables/
  mask: $schema$.$object$.sql

staticData:
  subdirectory: Tables/StaticData
  mask: $schema$.$object$.Data.sql

view:
  subdirectory: Views
  mask: $schema$.$object$.sql

procedure:
  subdirectory: Programmability/Procedures
  mask: $schema$.$object$.sql

function:
  subdirectory: Programmability/Functions
  mask: $schema$.$object$.sql

trigger:
  subdirectory: Programmability/Database/Triggers
  mask: $schema$.$object$.sql

domain:
  subdirectory: Programmability/User Types/Data Types
  mask: $schema$.$object$.sql

schema:
  subdirectory: Security/Schemas
  mask: $object$.sql
`
