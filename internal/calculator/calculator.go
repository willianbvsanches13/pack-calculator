// Package calculator implements the pack calculation algorithm using dynamic programming.
package calculator

import (
	"errors"
	"math"
	"sort"
)

var (
	ErrNoPackSizes     = errors.New("no pack sizes configured")
	ErrInvalidAmount   = errors.New("amount must be greater than zero")
	ErrInvalidPackSize = errors.New("pack size must be greater than zero")
)

type Calculator struct {
	packSizes []int // sorted descending
}

func New(packSizes []int) (*Calculator, error) {
	if len(packSizes) == 0 {
		return nil, ErrNoPackSizes
	}

	sizes := make([]int, len(packSizes))
	for i, size := range packSizes {
		if size <= 0 {
			return nil, ErrInvalidPackSize
		}
		sizes[i] = size
	}

	sort.Sort(sort.Reverse(sort.IntSlice(sizes)))

	return &Calculator{packSizes: sizes}, nil
}

func (c *Calculator) GetPackSizes() []int {
	result := make([]int, len(c.packSizes))
	copy(result, c.packSizes)
	return result
}

func (c *Calculator) SetPackSizes(packSizes []int) error {
	if len(packSizes) == 0 {
		return ErrNoPackSizes
	}

	sizes := make([]int, len(packSizes))
	for i, size := range packSizes {
		if size <= 0 {
			return ErrInvalidPackSize
		}
		sizes[i] = size
	}

	sort.Sort(sort.Reverse(sort.IntSlice(sizes)))
	c.packSizes = sizes
	return nil
}

// Calculate finds the optimal pack combination using DP (similar to coin change).
func (c *Calculator) Calculate(amount int) (map[int]int, error) {
	if amount <= 0 {
		return nil, ErrInvalidAmount
	}

	if len(c.packSizes) == 0 {
		return nil, ErrNoPackSizes
	}

	smallestPack := c.packSizes[len(c.packSizes)-1]
	largestPack := c.packSizes[0]

	// upper bound for DP - no valid solution exceeds this
	maxTarget := amount + largestPack

	const impossible = math.MaxInt32
	dp := make([]int, maxTarget+1)     // dp[i] = min packs to get i items
	parent := make([]int, maxTarget+1) // tracks which pack was used

	for i := range dp {
		dp[i] = impossible
		parent[i] = -1
	}
	dp[0] = 0

	// build DP table
	for i := 0; i <= maxTarget; i++ {
		if dp[i] == impossible {
			continue
		}

		for _, packSize := range c.packSizes {
			next := i + packSize
			if next > maxTarget {
				continue
			}

			if dp[i]+1 < dp[next] {
				dp[next] = dp[i] + 1
				parent[next] = packSize
			}
		}
	}

	// find smallest total >= amount
	target := -1
	for i := amount; i <= maxTarget; i++ {
		if dp[i] != impossible {
			target = i
			break
		}
	}

	if target == -1 {
		// Fallback: shouldnt happen with valid pack sizes
		packsNeeded := (amount + smallestPack - 1) / smallestPack
		return map[int]int{smallestPack: packsNeeded}, nil
	}

	// backtrack to find which packs were used
	result := make(map[int]int)
	current := target
	for current > 0 {
		pack := parent[current]
		if pack == -1 {
			break
		}
		result[pack]++
		current -= pack
	}

	return result, nil
}

type CalculationResult struct {
	Packs       map[int]int `json:"packs"`
	TotalItems  int         `json:"total_items"`
	TotalPacks  int         `json:"total_packs"`
	OrderAmount int         `json:"order_amount"`
}

func (c *Calculator) CalculateWithDetails(amount int) (*CalculationResult, error) {
	packs, err := c.Calculate(amount)
	if err != nil {
		return nil, err
	}

	var totalItems, totalPacks int
	for size, qty := range packs {
		totalItems += size * qty
		totalPacks += qty
	}

	return &CalculationResult{
		Packs:       packs,
		TotalItems:  totalItems,
		TotalPacks:  totalPacks,
		OrderAmount: amount,
	}, nil
}
