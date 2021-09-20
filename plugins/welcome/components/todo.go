package components

import (
	"github.com/iagapie/go-spring/pkg/manager"
	"reflect"
)

type todo struct {
	props manager.ComponentProps
}

func NewTodo(props manager.ComponentProps) manager.Component {
	return &todo{
		props: props,
	}
}

func (c *todo) Details() manager.ComponentDetails {
	return manager.ComponentDetails{
		Code:     "todo",
		Name:     "Todo Component",
		ViewFile: "default",
	}
}

func (c *todo) CfgProps() manager.ComponentCfgProps {
	return map[string]manager.ComponentCfgProp{
		"max": {
			Title: "Max",
			Type:  reflect.Int,
		},
	}
}

func (c *todo) Props() manager.ComponentProps {
	return c.props
}
