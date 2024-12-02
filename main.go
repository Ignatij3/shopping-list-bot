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
	mylog.Printf(mylogger.INFO+"chooseNewRecipe is called for %q recipe type\n", recipeType.String())
	// first, we go through recipes that have the closest ingredients to what we already have
	for dist := 1; dist < 2; dist++ {
		if goodRecipes, ok := remaining.FindMatchingRecipes(recips, dist); ok {
			for _, gr := range goodRecipes {
				if recipeType == gr.Typ {
					newRec = gr
					mylog.Printf(mylogger.INFO+"New recipe is chosen: %v\n", newRec)
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
		mylog.Printf(mylogger.WARN+"Recipe of type %q was not found\n", recipeType.String())
		return
	}

	// second, if nothing was found, we will choose one randomly
	newRec = recips[rand.Intn(len(recips))]
	for newRec.Typ != recipeType {
		newRec = recips[rand.Intn(len(recips))]
	}
	return
}

// TODO: // separate ingredient lists into separate files and place them next to pdfs.
// TODO: // add functionality to the bot so that it is possible to add additional recipes into consideration from interaction.
// TODO: // add functionality to the bot so that it automatically adds files (recipe.pdf and ingredients file) to correct folders and reads it.
// TODO: // add support for different categories.
// TODO: // add funct. to producre reports about daily consumption of calories/other stuff.
// TODO: // restrict updating aforementioned data to the next week, starting from Monday, 00:00. (meaning, it is not possible to update nutrition data for this week)
// TODO: // add differentiation between dishes that are cooked once and dishes that are cooked for the whole week (or variable amount of days).
// TODO: // introduce into a bot planner that will visualize what days are missing breakfast/lunch/evening meal/other....
// TODO: // real-time remaining product tracking (delete from database after a certain time in each day of the week has passed)

func main() {
	mylog = mylogger.NewLogger("MAIN: ")
	defer mylogger.CloseResources()

	mylog.Println(mylogger.INFO + "Beginning execution")
	// prods := recipes.GetProducts()
	recps := recipes.GetRecipes()
	mylog.Println(mylogger.INFO + "Reading kitchen state")
	remainingProducts := kitchen.ReadKitchenState()

	mylog.Println(mylogger.INFO + "Choosing recipes")
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

	mylog.Println(mylogger.INFO + "Summing ingredients")
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

	mylog.Println(mylogger.INFO + "Outputting result")
	fmt.Println("Here's the list of products that you need to buy for the next week:")
	i := 1
	for name, product := range productDifference {
		fmt.Printf("%d) %s (%d grams)\n", i, name, product)
		i++
	}
	fmt.Printf("Lunch recipes are:\nLunch - %s\nSalad - %s\n", lunch.Name, salad.Name)

	// imitate consuming all purchased foods and calculate remaining food
	// ** IMPORTANT ** - this is only available when minimal purchasing quantity will be added to product data
	// "Smallest package of XXX that you usually buy?"
}
