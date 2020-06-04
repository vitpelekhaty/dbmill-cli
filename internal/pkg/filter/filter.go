package filter

import (
	"errors"
	"fmt"
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

// Match проверяет соответствие значения value фильтру
func (filter *Filter) Match(value string) error {
	if len(filter.expressions) == 0 {
		return nil
	}

	for _, exp := range filter.expressions {
		if err := filter.match(exp, value); err == nil {
			return nil
		} else {
			if err != ErrorNotMatched {
				return err
			}
		}
	}

	return ErrorNotMatched
}

func (filter *Filter) match(pattern, value string) error {
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
