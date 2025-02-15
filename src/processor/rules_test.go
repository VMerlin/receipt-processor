package processor

import (
	"testing"
)

func TestCalculatePointsForRetailerName(t *testing.T) {
	name := "M&M Corner Market"
	expected := 14
	result := calculatePointsForRetailerName(name)
	if result != expected {
		t.Errorf("Expected %d points, got %d", expected, result)
	}
}

func TestCalculatePointsForTotal(t *testing.T) {
	tests := []struct {
		total    float64
		expected int64
	}{
		{9.00, 75},
		{10.15, 0},
		{23.25, 25},
	}

	for _, test := range tests {
		result := calculatePointsForTotal(test.total)
		if result != test.expected {
			t.Errorf("For total %.2f, expected %d points, got %d", test.total, test.expected, result)
		}
	}
}

func TestCalculatePointsForItems(t *testing.T) {
	items := []Item{
		{description: "Gatorade", price: 2.25},
		{description: "Gatorade", price: 2.25},
		{description: "Gatorade", price: 2.25},
		{description: "Gatorade", price: 2.25},
	}
	expected := int64(10)
	result := calculatePointsForItems(items)
	if result != expected {
		t.Errorf("Expected %d points, got %d", expected, result)
	}

	items = append(items, Item{description: "Gatorade", price: 2.25})
	result = calculatePointsForItems(items)
	if result != expected {
		t.Errorf("Expected %d points, got %d", expected, result)
	}

	items = append(items, Item{description: "Gatorade3", price: 2.25})
	expected = int64(16)
	result = calculatePointsForItems(items)
	if result != expected {
		t.Errorf("Expected %d points, got %d", expected, result)
	}
}

func TestCalculatePointsForDay(t *testing.T) {
	tests := []struct {
		day      int
		expected int
	}{
		{11, 6},
		{2, 0},
		{3, 6},
		{24, 0},
	}

	for _, test := range tests {
		result := calculatePointsForDay(test.day)
		if result != test.expected {
			t.Errorf("For day %d, expected %d points, got %d", test.day, test.expected, result)
		}
	}
}

func TestCalculatePointsForTime(t *testing.T) {
	tests := []struct {
		hour     int
		expected int
	}{
		{13, 0},
		{14, 10},
		{15, 10},
		{16, 0},
	}

	for _, test := range tests {
		result := calculatePointsForTime(test.hour)
		if result != test.expected {
			t.Errorf("For hour %d, expected %d points, got %d", test.hour, test.expected, result)
		}
	}
}
