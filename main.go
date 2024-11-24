package main

import (
	"fmt"
	kitchen "main/kitchen"
	mylogger "main/mylogger"
	recipes "main/recipes"
	"math/rand"
)

var mylog *mylogger.Mylogger

// chooseNewRecipe will randomly allocate a new recipe of specified type.
// Changes must be saved to the file manually.
func chooseNewRecipe(recips recipes.Recipes, remaining *kitchen.RemainingProducts, recipeType recipes.DishType) (newRec recipes.Recipe) {
	mylog.Printf(mylogger.INFO+"chooseNewRecipe is called for %s recipe type\n", recipeType.String())
	// first, we go through recipes that have the closest ingredients to what we already have
	for dist := 1; dist < 2; dist++ {
		if goodRecipes, ok := remaining.FindMatchingRecipes(recips, dist); ok {
			for _, gr := range goodRecipes {
				if recipeType == gr.Typ {
					newRec = gr
					mylog.Printf(mylogger.INFO+"new recipe is chosen: %v\n", newRec)
					return
				}
			}
		}
	}

	// checking if recipe of that type is present
	found := false
	for _, rec := range recips {
		if rec.Typ == recipeType {
			found = true
		}
	}
	if !found {
		mylog.Printf(mylogger.WARN+"Recipe of type \"%s\" was not found\n", recipeType.String())
		return
	}

	// second, if nothing was found, we will choose one randomly
	newRec = recips[rand.Intn(len(recips))]
	for newRec.Typ != recipeType {
		newRec = recips[rand.Intn(len(recips))]
	}
	return
}

func main() {
	mylog = mylogger.NewLogger("MAIN: ")
	defer mylogger.CloseResources()

	mylog.Println(mylogger.INFO + "beginning execution")
	// prods := recipes.GetProducts()
	recps := recipes.GetRecipes()
	mylog.Println(mylogger.INFO + "reading kitchen state")
	remainingProducts := kitchen.ReadKitchenState()

	mylog.Println(mylogger.INFO + "choosing recipes")
	// get all recipes to further calculate what is needed to buy

	breakfasts := make(recipes.Recipes, 0, 7)
	for _, recipeBreakfast := range recps {
		if recipeBreakfast.Typ == recipes.BREAKFAST {
			breakfasts = append(breakfasts, recipeBreakfast)
		}
	}

	lunch := chooseNewRecipe(recps, &remainingProducts, recipes.LUNCH)
	salad := chooseNewRecipe(recps, &remainingProducts, recipes.SALAD)

	totalNeeded := make(recipes.Ingredients)

	mylog.Println(mylogger.INFO + "summing ingredients")
	// this function assures that when same ingredient comes up, it is correctly accounted for
	correctIngredientSum := func(name string, ingWeight int) int {
		if wght, ok := totalNeeded[name]; ok {
			ingWeight += wght
		}
		return ingWeight
	}

	// all ingredients are summed up and collected in one place to determine shopping cart (what will we buy)
	for _, brkfast := range breakfasts {
		for name, ing := range brkfast.GetIngredients() {
			totalNeeded[name] = correctIngredientSum(name, ing)
		}
	}
	for name, ing := range lunch.GetIngredients() {
		totalNeeded[name] = correctIngredientSum(name, ing)
	}
	for name, ing := range salad.GetIngredients() {
		totalNeeded[name] = correctIngredientSum(name, ing)
	}

	// difference between what we have and what we need is calculated
	productDifference := make(recipes.Ingredients)
	for name := range totalNeeded {
		if totalNeeded[name] > remainingProducts[name] {
			productDifference[name] = totalNeeded[name] - remainingProducts[name]
		}
	}

	mylog.Println(mylogger.INFO + "outputting result")
	// presenting the user with output and calculate how much we bought (need minimal purchasing quantity for this) and how much is remaining after this
	fmt.Println("Here's the list of products that you need to buy for the next week:")
	i := 1
	for name, product := range productDifference {
		fmt.Printf("%d) %s (%d grams)\n", i, name, product)
		i++
	}
}
