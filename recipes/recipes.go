package recipes

import "main/mylogger"

// enumeration of dish type by time of consumption.
const (
	BREAKFAST DishType = iota
	LUNCH     DishType = iota
	SALAD     DishType = iota
	SOUP      DishType = iota
)

type (
	DishType uint8
	// Ingredients contains information on how much of product we have in a recipe.
	Ingredients map[string]int

	// represents a recipe, distinguished by type and ingredient list.
	Recipe struct {
		Name string
		Typ  DishType
		ing  Ingredients
	}
	Recipes []Recipe
)

var mylog *mylogger.Mylogger

// init creates logger for the module
func init() {
	mylog = mylogger.NewLogger("RECP: ")
}

// GetIngredients returns ingredients of the recipe.
func (r Recipe) GetIngredients() Ingredients {
	return r.ing
}

// GetIngredientNames returns a list of all ingredients' names in a recipe
func (r Recipe) GetIngredientNames() []string {
	keys := make([]string, len(r.ing))

	i := 0
	for k := range r.ing {
		keys[i] = k
		i++
	}
	return keys
}

// AddIngredient adds a new ingredient to the recipe. Collisions are not checked, so if a name of new ingredient collides with the old one, old one is replaced.
func (r *Recipe) AddIngredient(name string, newProductWeight int) {
	(*r).ing[name] = newProductWeight
}

// ReplaceIngredient replaces an ingredient if found. Returned value is true when replacement has been successful.
func (r *Recipe) ReplaceIngredient(name string, newProductWeight int) bool {
	if _, ok := (*r).ing[name]; ok {
		(*r).ing[name] = newProductWeight
		return true
	}
	return false
}

// DeleteIngredient deletes an ingredient if it is found.
func (r *Recipe) DeleteIngredient(name string, newPrd Product) {
	delete((*r).ing, name)
}

func (d DishType) String() string {
	switch d {
	case BREAKFAST:
		return "BREAKFAST"
	case LUNCH:
		return "LUNCH"
	case SALAD:
		return "SALAD"
	case SOUP:
		return "SOUP"
	default:
		return "__UNKNOWN__"
	}
}
