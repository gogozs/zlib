package dbgenerator

import "github.com/gogozs/zlib/tools"

type KeyString string

func (n KeyString) String() string {
	return string(n)
}

func (n KeyString) UpperCamelCase() string {
	return tools.ToUpperCamelCase(n.String())
}

func (n KeyString) LowerCamelCase() string {
	return tools.ToLowerCamelCase(n.String())
}

func (n KeyString) LowerSnake() string {
	return tools.ToLowerSnakeString(n.String())
}

func (n KeyString) UpperSnake() string {
	return tools.ToUpperSnakeString(n.String())
}
