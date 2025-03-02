package profile_test

import (
	"cache/internal/order"
	"cache/internal/profile"
	"sync"
	"testing"
	"time"
)

func TestCacheSetAndGet(t *testing.T) {
	cache := profile.NewCache(2 * time.Second)

	order1 := order.NewOrder("Order 1")
	order2 := order.NewOrder("Order 2")
	prof := NewProfile("Test User", []*order.Order{order1, order2})

	cache.Set(prof.UUID, prof)

	result := cache.Get(prof.UUID)
	if result == nil {
		t.Fatal("Expected profile, got nil")
	}
	if result.Name != prof.Name || len(result.Orders) != len(prof.Orders) {
		t.Errorf("Expected profile %v, got %v", prof, result)
	}
}

func TestCacheTTLExpiration(t *testing.T) {
	cache := profile.NewCache(100 * time.Millisecond)

	order1 := order.NewOrder("Order 1")
	prof := NewProfile("Test User", []*order.Order{order1})
	cache.Set(prof.UUID, prof)

	if result := cache.Get(prof.UUID); result == nil {
		t.Fatal("Expected profile before TTL expiration, got nil")
	}

	time.Sleep(150 * time.Millisecond)

	if result := cache.Get(prof.UUID); result != nil {
		t.Errorf("Expected nil after TTL expiration, got %v", result)
	}
}

func TestCacheThreadSafety(t *testing.T) {
	cache := profile.NewCache(2 * time.Second)
	uuid := "test-uuid"

	var wg sync.WaitGroup
	const goroutines = 50
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func(i int) {
			defer wg.Done()

			prof := NewProfile("User "+string(rune(i)), nil)
			cache.Set(uuid, prof)

			if result := cache.Get(uuid); result == nil {
				t.Errorf("Goroutine %d: Expected profile, got nil", i)
			}
		}(i)
	}

	wg.Wait()

	if result := cache.Get(uuid); result == nil {
		t.Error("Expected profile after concurrent access, got nil")
	}
}

func TestCacheDataIsolation(t *testing.T) {
	cache := profile.NewCache(2 * time.Second)

	order1 := order.NewOrder("Order 1")
	prof := NewProfile("Test User", []*order.Order{order1})
	cache.Set(prof.UUID, prof)

	prof.Name = "Modified User"
	prof.Orders[0].Value = "Modified Order"

	result := cache.Get(prof.UUID)
	if result == nil {
		t.Fatal("Expected profile, got nil")
	}
	if result.Name != "Test User" || result.Orders[0].Value != "Order 1" {
		t.Errorf("Expected unchanged data %v, got %v", prof, result)
	}
}

func TestCacheUpdateResetsTTL(t *testing.T) {
	cache := profile.NewCache(100 * time.Millisecond) // Короткий TTL для теста

	prof := NewProfile("Test User", nil)
	cache.Set(prof.UUID, prof)

	time.Sleep(80 * time.Millisecond)

	updatedProf := NewProfile("Updated User", nil)
	cache.Set(prof.UUID, updatedProf)

	time.Sleep(50 * time.Millisecond)

	result := cache.Get(prof.UUID)
	if result == nil {
		t.Fatal("Expected updated profile, got nil")
	}
	if result.Name != "Updated User" {
		t.Errorf("Expected updated profile name 'Updated User', got %v", result.Name)
	}
}
