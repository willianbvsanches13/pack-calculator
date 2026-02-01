package calculator

import (
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name      string
		packSizes []int
		wantErr   error
	}{
		{"valid", []int{250, 500, 1000}, nil},
		{"empty", []int{}, ErrNoPackSizes},
		{"nil", nil, ErrNoPackSizes},
		{"zero value", []int{250, 0, 1000}, ErrInvalidPackSize},
		{"negative", []int{250, -500, 1000}, ErrInvalidPackSize},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(tt.packSizes)
			if err != tt.wantErr {
				t.Errorf("got error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestCalculate(t *testing.T) {
	tests := []struct {
		name      string
		packSizes []int
		amount    int
		wantItems int
		wantCount int
		wantErr   error
	}{
		{
			name:      "1 item",
			packSizes: []int{250, 500, 1000, 2000, 5000},
			amount:    1,
			wantItems: 250,
			wantCount: 1,
		},
		{
			name:      "exact 250",
			packSizes: []int{250, 500, 1000, 2000, 5000},
			amount:    250,
			wantItems: 250,
			wantCount: 1,
		},
		{
			name:      "251 should use 1x500",
			packSizes: []int{250, 500, 1000, 2000, 5000},
			amount:    251,
			wantItems: 500,
			wantCount: 1,
		},
		{
			name:      "501 items",
			packSizes: []int{250, 500, 1000, 2000, 5000},
			amount:    501,
			wantItems: 750,
			wantCount: 2,
		},
		{
			name:      "12001 items",
			packSizes: []int{250, 500, 1000, 2000, 5000},
			amount:    12001,
			wantItems: 12250,
			wantCount: 4,
		},
		{
			name:      "zero amount",
			packSizes: []int{250, 500, 1000},
			amount:    0,
			wantErr:   ErrInvalidAmount,
		},
		{
			name:      "negative amount",
			packSizes: []int{250, 500, 1000},
			amount:    -100,
			wantErr:   ErrInvalidAmount,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calc, err := New(tt.packSizes)
			if err != nil {
				t.Fatalf("failed to create calculator: %v", err)
			}

			result, err := calc.CalculateWithDetails(tt.amount)

			if err != tt.wantErr {
				t.Errorf("error = %v, want %v", err, tt.wantErr)
				return
			}

			if tt.wantErr != nil {
				return
			}

			if result.TotalItems != tt.wantItems {
				t.Errorf("totalItems = %v, want %v", result.TotalItems, tt.wantItems)
			}

			if result.TotalPacks != tt.wantCount {
				t.Errorf("totalPacks = %v, want %v", result.TotalPacks, tt.wantCount)
			}
		})
	}
}

// edge case from the challenge email: pack sizes 23, 31, 53 with amount 500000
func TestEdgeCaseFromEmail(t *testing.T) {
	calc, err := New([]int{23, 31, 53})
	if err != nil {
		t.Fatalf("failed to create calculator: %v", err)
	}

	result, err := calc.CalculateWithDetails(500000)
	if err != nil {
		t.Fatalf("Calculate() error = %v", err)
	}

	// expected: 23*2 + 31*7 + 53*9429 = 46 + 217 + 499737 = 500000
	expectedTotal := 500000
	expectedPacks := 2 + 7 + 9429

	if result.TotalItems != expectedTotal {
		t.Errorf("TotalItems = %v, want %v", result.TotalItems, expectedTotal)
	}

	if result.TotalPacks != expectedPacks {
		t.Errorf("TotalPacks = %v, want %v", result.TotalPacks, expectedPacks)
	}

	// verify sum matches
	sum := 0
	for pack, qty := range result.Packs {
		sum += pack * qty
	}

	if sum != result.TotalItems {
		t.Errorf("sum of packs (%v) != TotalItems (%v)", sum, result.TotalItems)
	}

	t.Logf("result: %v (total: %d items, %d packs)", result.Packs, result.TotalItems, result.TotalPacks)
}

// rule 2 should take priority (minimize items first)
func TestMinimizesItems(t *testing.T) {
	calc, err := New([]int{250, 500})
	if err != nil {
		t.Fatalf("failed to create calculator: %v", err)
	}

	result, err := calc.CalculateWithDetails(251)
	if err != nil {
		t.Fatalf("Calculate() error = %v", err)
	}

	// should use 1x500, not 2x250
	if result.Packs[500] != 1 {
		t.Errorf("expected 1x500, got %v", result.Packs)
	}

	if result.TotalPacks != 1 {
		t.Errorf("expected 1 pack, got %d", result.TotalPacks)
	}
}

func TestSinglePackSize(t *testing.T) {
	calc, err := New([]int{100})
	if err != nil {
		t.Fatalf("failed to create calculator: %v", err)
	}

	tests := []struct {
		amount    int
		wantPacks int
		wantItems int
	}{
		{1, 1, 100},
		{100, 1, 100},
		{101, 2, 200},
		{250, 3, 300},
		{999, 10, 1000},
	}

	for _, tt := range tests {
		result, err := calc.CalculateWithDetails(tt.amount)
		if err != nil {
			t.Errorf("Calculate(%d) error = %v", tt.amount, err)
			continue
		}

		if result.TotalPacks != tt.wantPacks {
			t.Errorf("Calculate(%d) packs = %d, want %d", tt.amount, result.TotalPacks, tt.wantPacks)
		}

		if result.TotalItems != tt.wantItems {
			t.Errorf("Calculate(%d) items = %d, want %d", tt.amount, result.TotalItems, tt.wantItems)
		}
	}
}

func TestSetPackSizes(t *testing.T) {
	calc, _ := New([]int{250, 500})

	result1, _ := calc.CalculateWithDetails(251)
	if result1.TotalItems != 500 {
		t.Errorf("initial: expected 500 items, got %d", result1.TotalItems)
	}

	err := calc.SetPackSizes([]int{100, 200})
	if err != nil {
		t.Fatalf("SetPackSizes() error = %v", err)
	}

	result2, _ := calc.CalculateWithDetails(251)
	if result2.TotalItems != 300 {
		t.Errorf("after update: expected 300 items, got %d", result2.TotalItems)
	}
}

func TestGetPackSizes(t *testing.T) {
	calc, _ := New([]int{100, 500, 250})

	sizes := calc.GetPackSizes()

	// should be sorted descending
	if sizes[0] != 500 || sizes[1] != 250 || sizes[2] != 100 {
		t.Errorf("GetPackSizes() = %v, expected [500 250 100]", sizes)
	}

	// modifying returned slice shouldnt affect calculator
	sizes[0] = 9999
	original := calc.GetPackSizes()
	if original[0] != 500 {
		t.Error("GetPackSizes() should return a copy, not the original slice")
	}
}

func BenchmarkCalculate(b *testing.B) {
	calc, _ := New([]int{23, 31, 53})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		calc.Calculate(500000)
	}
}

func BenchmarkCalculateStandardPacks(b *testing.B) {
	calc, _ := New([]int{250, 500, 1000, 2000, 5000})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		calc.Calculate(12001)
	}
}
