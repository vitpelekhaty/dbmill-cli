package filter

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
)

// ErrorNotMatched ошибка "Значение не соответствует фильтру"
var ErrorNotMatched = errors.New("not matched")

// IFilter интерфейс фильтра объектов
type IFilter interface {
	// Match проверяет соответствие значения value фильтру
	Match(value string) error
}

// Filter фильтр объектов базы данных
type Filter struct {
	expressions []string
}

// New возвращает новый экземпляр фильтра объектов
func New(expressions []string) (*Filter, error) {
	for _, exp := range expressions {
		if _, err := regexp.Compile(exp); err != nil {
			return nil, fmt.Errorf("%s: %v", exp, err)
		}
	}

	return &Filter{
		expressions: expressions,
	}, nil
}

// FromFile возвращает новый экземпляр фильтра объектов.
// Список выражений загружается из указанного файла
func FromFile(path string) (*Filter, error) {
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

	return New(exp)
}

// Match проверяет соответствие значения value фильтру
func (self *Filter) Match(value string) error {
	if len(self.expressions) == 0 {
		return nil
	}

	for _, exp := range self.expressions {
		if err := self.match(exp, value); err == nil {
			return nil
		} else {
			if err != ErrorNotMatched {
				return err
			}
		}
	}

	return ErrorNotMatched
}

func (self *Filter) match(pattern, value string) error {
	re, err := regexp.Compile(pattern)

	if err != nil {
		return err
	}

	ok := re.MatchString(value)

	if ok {
		return nil
	}

	return ErrorNotMatched
}
