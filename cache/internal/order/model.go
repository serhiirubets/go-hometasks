package order

import (
	"math/rand/v2"
	"strconv"
	"time"
)

type Order struct {
	UUID      string
	Value     any
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewOrder(value any) *Order {
	return &Order{

		// Emulate creating UUID. Not in the scope of current task
		UUID:      "vfdf3ref" + strconv.Itoa(rand.IntN(1000)),
		Value:     value,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
