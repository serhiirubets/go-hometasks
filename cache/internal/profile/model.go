package profile

import (
	"cache/internal/order"
	"math/rand/v2"
	"strconv"
)

type Profile struct {
	UUID   string
	Name   string
	Orders []*order.Order
}

func NewProfile(name string, orders []*order.Order) *Profile {
	return &Profile{
		// Emulate creating UUID. Not in the scope of current task
		UUID:   "fdfwe33fd" + strconv.Itoa(rand.IntN(1000)),
		Name:   name,
		Orders: orders,
	}
}
