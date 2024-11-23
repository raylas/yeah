package vendors

import (
	"slices"
	"testing"
)

var testVendors = []VendorEntry{
	{Oui: "001122", Name: "Shoe Sparkle", Address: "123 Soles St"},
	{Oui: "001123", Name: "Pulp Pirates", Address: "456 Paper Ave"},
	{Oui: "001134", Name: "Data Dudes", Address: "789 Bytes Blvd"},
}

func TestInsertAndSearch(t *testing.T) {
	cases := []struct {
		prefix   string
		vendors  []VendorEntry
		expected []VendorEntry
	}{
		{
			prefix:  "001122",
			vendors: testVendors,
			expected: []VendorEntry{
				{Oui: "001122", Name: "Shoe Sparkle", Address: "123 Soles St"},
			},
		},
		{
			prefix:  "00:11:23",
			vendors: testVendors,
			expected: []VendorEntry{
				{Oui: "001123", Name: "Pulp Pirates", Address: "456 Paper Ave"},
			},
		},
		{
			prefix:  "00-11-2",
			vendors: testVendors,
			expected: []VendorEntry{
				{Oui: "001122", Name: "Shoe Sparkle", Address: "123 Soles St"},
				{Oui: "001123", Name: "Pulp Pirates", Address: "456 Paper Ave"},
			},
		},
		{
			prefix:   "00.11.28",
			vendors:  testVendors,
			expected: []VendorEntry{},
		},
		{
			prefix:   "00:xx:xx",
			vendors:  testVendors,
			expected: []VendorEntry{},
		},
	}

	for _, c := range cases {
		v := New()
		for _, vendor := range c.vendors {
			v.Insert(vendor.Oui, &vendor)
		}
		results := v.Search(c.prefix)

		if len(results) != len(c.expected) {
			t.Errorf("expected %v, got %v", c.expected, results)
		}
		for _, result := range results {
			if !slices.Contains(c.expected, *result) {
				t.Errorf("expected %+v, got %+v", c.expected, results)
			}
		}
	}
}

func TestNormalize(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"A1:B2:C3", "A1B2C3"},
		{"00:1A:2B:3C:4D:5E", "001A2B3C4D5E"},
		{"00-04-F3-2E-B1-A2", "0004F32EB1A2"},
		{"00.23.45.67.89.AB", "0023456789AB"},
		{"001122AABBCC", "001122AABBCC"},
	}

	for _, c := range cases {
		if got := normalize(c.input); got != c.expected {
			t.Errorf("expected %v, got %v", c.expected, got)
		}
	}
}
