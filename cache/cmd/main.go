package main

import (
	"cache/internal/order"
	"cache/internal/profile"
	"fmt"
	"time"
)

func main() {
	cache := profile.NewCache(2 * time.Second)

	order1 := order.NewOrder("Order 1")
	order2 := order.NewOrder("Order 2")

	profile1 := profile.NewProfile("User 1", []*order.Order{order1, order2})

	cache.Set(profile1.UUID, profile1)

	fmt.Print(cache.Get(profile1.UUID))
}
