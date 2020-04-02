package dir

import (
	"fmt"
	"io"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Item элемент структуры каталога скриптов
type Item struct {
	// SubDirectory подкаталог
	SubDirectory string `yaml:"subdirectory"`
	// FilenameMask маска имени файла
	FilenameMask string `yaml:"mask"`
}

// IStructure интерфес описания структуры каталога скриптов
type IStructure interface {
	// ItemInfo возвращает целевой каталог и маску имени файла для указанного типа объекта itemType.
	// Если информация не найдена, то в параметре ok возвращается false, в противном случае - true
	ItemInfo(item StructItemType) (subdirectory, mask string, ok bool)
}

// Structure описание структуры каталога скриптов
type Structure struct {
	Items map[StructItemType]Item
}

// NewStructure конструктор описания структуры каталога скриптов
func NewStructure(in io.Reader) (*Structure, error) {
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

	return &Structure{Items: si}, nil
}

// ItemInfo возвращает целевой каталог и маску имени файла для указанного типа объекта itemType.
// Если информация не найдена, то в параметре ok возвращается false, в противном случае - true
func (self *Structure) ItemInfo(item StructItemType) (subdirectory, mask string, ok bool) {
	i, ok := self.Items[item]

	if ok {
		subdirectory = i.SubDirectory
		mask = i.FilenameMask
	}

	return
}

func parse(data []byte) (map[string]Item, error) {
	var s map[string]Item

	err := yaml.Unmarshal(data, &s)

	if err != nil {
		return s, err
	}

	return s, nil
}

func mapping(items map[string]Item) (map[StructItemType]Item, error) {
	si := make(map[StructItemType]Item)

	if len(items) == 0 {
		return si, nil
	}

	for key, value := range items {
		if k, ok := itemTypeMapping[key]; ok {
			si[k] = value
		} else {
			return si, fmt.Errorf("unknown object type %s", key)
		}
	}

	return si, nil
}
