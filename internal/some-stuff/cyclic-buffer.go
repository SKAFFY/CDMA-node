package some_stuff

import (
	"errors"
)

var EmptyError = errors.New("CyclicBuffer is empty")

type CyclicBuffer[Type any] struct {
	tail int // to pull
	head int // to push

	body []Type
}

func NewCyclicBuffer[Type any](n int) *CyclicBuffer[Type] {
	return &CyclicBuffer[Type]{
		tail: 0,
		head: 0,
		body: make([]Type, n),
	}
}

func (c *CyclicBuffer[Type]) Push(v Type) {
	c.head = (c.head + 1) % len(c.body)

	c.body[c.head] = v

	if c.head == c.tail {
		c.tail = (c.tail + 1) % len(c.body)
	}
}

func (c *CyclicBuffer[Type]) Pull() (Type, error) {
	if c.head == c.tail {
		var zero Type

		return zero, EmptyError
	}

	v := c.body[c.tail]

	c.tail = (c.tail + 1) % len(c.body)

	return v, nil
}

func (c *CyclicBuffer[Type]) Body() []Type {
	result := make([]Type, len(c.body))

	start := c.tail
	i := 0

	for {
		result[i] = c.body[start]

		start = (start + 1) % len(c.body)
		if start == c.head {
			break
		}

		i++
	}

	return result
}
