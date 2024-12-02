package kitchen

import (
	"main/mylogger"
	recipes "main/recipes"
)

// map in which a key is a name of product and number is amount in grams of how much products are left.
type RemainingProducts map[string]int

// FindMatchingRecipe will return all recipes that it is possible to produce using "rmProd".
// Dist - is amount of products that are allowed to be in quantities not high enough to cook a recipe.
// Second returned value is true if a recipe is found.
// Beware that if one recipe is made, it is not guaranteed that other recipes still can be made
func (rp RemainingProducts) FindMatchingRecipes(rec recipes.Recipes, dist int) (recipes.Recipes, bool) {
	mylog.Printf(mylogger.INFO+"FindMatchingRecipes is called with arguments: %v, %d\n", rec, dist)
	if dist < 0 {
		return recipes.Recipes{}, false
	}

	result := make(recipes.Recipes, 0)
	for _, recipe := range rec {
		invalidCounter := 0
		for name, ingrdWeight := range recipe.GetIngredients() {
			if wght, ok := rp[name]; !ok || wght < ingrdWeight {
				invalidCounter++
			}
		}

		if invalidCounter <= dist {
			result = append(result, recipe)
		}
	}

	return result, false
}

// AddToProduct adds weight to product.
// If the product does not exist, it is added.
func (p *RemainingProducts) AddToProduct(name string, weight int) {
	(*p)[name] += weight
}

// ReduceProduct reduces weight of product.
func (p *RemainingProducts) ReduceProduct(name string, weight int) {
	if _, ok := (*p)[name]; ok {
		if (*p)[name] < weight {
			delete((*p), name)
		} else {
			(*p)[name] -= weight
		}
	}
}

// DeleteProduct deleted the product's weight.
func (p *RemainingProducts) DeleteProduct(name string) {
	delete((*p), name)
}

// CookRecipe simulates using up all ingredients that are specified in the recipe.
func (p *RemainingProducts) CookRecipe(rcp recipes.Recipe) {
	mylog.Printf(mylogger.INFO+"CookRecipe is called with arguments: %v\n", rcp)
	for name, prod := range rcp.GetIngredients() {
		(*p).ReduceProduct(name, prod)
	}
}

// UnCookRecipe is the same as CookRecipe, but it adds ingredients to the remaining products.
func (p *RemainingProducts) UnCookRecipe(rcp recipes.Recipe) {
	mylog.Printf(mylogger.INFO+"UnCookRecipe is called with arguments: %v\n", rcp)
	for name, prod := range rcp.GetIngredients() {
		(*p).AddToProduct(name, prod)
	}
}
