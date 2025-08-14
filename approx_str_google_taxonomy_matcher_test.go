package fastcategorymatcher_test

import (
	"testing"

	fastcategorymatcher "github.com/ndx-technologies/fast-category-matcher"
	"github.com/ndx-technologies/fast-category-matcher/googleproducttaxonomy"
)

func TestApproxStrGoogleTaxonomyMatcher(t *testing.T) {
	tests := map[string]googleproducttaxonomy.ProductCategory{
		"asdf asdf": googleproducttaxonomy.Unknown,
		"Food & Grocery>Produce>Citrus Fruit>Oranges": mustProductCategory(t, "Food, Beverages & Tobacco > Food Items > Fruits & Vegetables > Fresh & Frozen Fruits > Citrus Fruits > Oranges"),
		"Oranges":                       mustProductCategory(t, "Food, Beverages & Tobacco > Food Items > Fruits & Vegetables > Fresh & Frozen Fruits > Citrus Fruits > Oranges"),
		"Citrus":                        googleproducttaxonomy.Unknown,
		"Citrus Fruits":                 mustProductCategory(t, "Food, Beverages & Tobacco > Food Items > Fruits & Vegetables > Fresh & Frozen Fruits > Citrus Fruits"),
		"Fresh & Frozen Fruits":         mustProductCategory(t, "Food, Beverages & Tobacco > Food Items > Fruits & Vegetables > Fresh & Frozen Fruits"),
		"Shorts":                        mustProductCategory(t, "Apparel & Accessories > Clothing > Shorts"),
		"Food & Grocery>Snacks>Popcorn": mustProductCategory(t, "Food, Beverages & Tobacco > Food Items > Snack Foods > Popcorn"),
		"Food Items":                    mustProductCategory(t, "Food, Beverages & Tobacco > Food Items"),
		"Snack Foods":                   mustProductCategory(t, "Food, Beverages & Tobacco > Food Items > Snack Foods"),
		"Popcorn":                       mustProductCategory(t, "Food, Beverages & Tobacco > Food Items > Snack Foods > Popcorn"),
		"Potatoes":                      mustProductCategory(t, "Food, Beverages & Tobacco > Food Items > Fruits & Vegetables > Fresh & Frozen Vegetables > Potatoes"),
		"Potato":                        mustProductCategory(t, "Food, Beverages & Tobacco > Food Items > Fruits & Vegetables > Fresh & Frozen Vegetables > Potatoes"),
		"BabyFood":                      mustProductCategory(t, "Baby & Toddler > Nursing & Feeding > Baby & Toddler Food > Baby Food"),
		"Beef":                          googleproducttaxonomy.Unknown, // make sure it does not match to Beer
		"Milk":                          mustProductCategory(t, "Food, Beverages & Tobacco > Beverages > Milk"),
		"Eggs":                          mustProductCategory(t, "Food, Beverages & Tobacco > Food Items > Meat, Seafood & Eggs > Eggs"),
		"Meat":                          mustProductCategory(t, "Food, Beverages & Tobacco > Food Items > Meat, Seafood & Eggs > Meat"),
		"Seafood":                       mustProductCategory(t, "Food, Beverages & Tobacco > Food Items > Meat, Seafood & Eggs > Seafood"),
		"Food, Beverages & Tobacco > Beverages > Juice":  mustProductCategory(t, "Food, Beverages & Tobacco > Beverages > Juice"),
		"Food, Beverages & Tobacco > Beverages > Juices": mustProductCategory(t, "Food, Beverages & Tobacco > Beverages > Juice"),
	}
	for str, exp := range tests {
		t.Run(str, func(t *testing.T) {
			var config fastcategorymatcher.ApproxStrGoogleTaxonomyMatcherConfig
			s := fastcategorymatcher.NewApproxStrGoogleTaxonomyMatcher(config.WithDefaults())

			category, err := s.MatchGoogleProductCategory(str)
			if err != nil {
				t.Error(err)
			}

			if category != exp {
				t.Error(googleproducttaxonomy.Categories[category], category, exp)
			}
		})
	}
}

func BenchmarkApproxStrGoogleTaxonomyMatcher(b *testing.B) {
	var config fastcategorymatcher.ApproxStrGoogleTaxonomyMatcherConfig
	s := fastcategorymatcher.NewApproxStrGoogleTaxonomyMatcher(config.WithDefaults())

	v := "Food, Beverages & Tobacco>Beverages>Juices"

	for b.Loop() {
		s.MatchGoogleProductCategory(v)
	}
}

func mustProductCategory(t *testing.T, s string) googleproducttaxonomy.ProductCategory {
	v, err := googleproducttaxonomy.ProductCategoryFromString(s)
	if err != nil {
		t.Error(err)
	}
	return v
}
