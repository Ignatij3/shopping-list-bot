package main

import (
	"bufio"
	"fmt"
	"log"
	"main/kitchen"
	"main/mylogger"
	"main/recipes"
	"math"
	"os"
	"time"

	"math/rand"
)

// TODO: // separate ingredient lists into separate files and place them next to pdfs.
// TODO: // add functionality to the bot so that it is possible to add additional recipes into consideration from interaction.
// TODO: // add functionality to the bot so that it automatically adds files (recipe.pdf and ingredients file) to correct folders and reads it.
// TODO: // add support for different categories.
// TODO: // add funct. to producre reports about daily consumption of calories/other stuff.
// TODO: // restrict updating aforementioned data to the next week, starting from Monday, 00:00. (meaning, it is not possible to update nutrition data for this week)
// TODO: // add differentiation between dishes that are cooked once and dishes that are cooked for the whole week (or variable amount of days).
// TODO: // introduce into a bot planner that will visualize what days are missing breakfast/lunch/evening meal/other....
// TODO: // real-time remaining product tracking (delete from database after a certain time in each day of the week has passed)

var mylog *mylogger.Mylogger

// chooseNewRecipeFromExistingProducts will randomly allocate a new recipe of specified type.
// false is returned if recipe of specified type was not found.
func chooseNewRecipeFromExistingProducts(recipeList recipes.Recipes, remainingProducts kitchen.RemainingProducts, recipeType recipes.DishType) (newRec recipes.Recipe, found bool) {
	mylog.Printf(mylogger.INFO+"chooseNewRecipe is called for %q recipe type\n", recipeType.String())
	for dist := 1; dist < 2; dist++ {
		if goodRecipes, ok := remainingProducts.FindMatchingRecipes(recipeList, dist); ok {
			for _, gr := range goodRecipes {
				if recipeType == gr.Typ {
					newRec = gr
					found = true
					mylog.Printf(mylogger.INFO+"New recipe is chosen: %v\n", newRec)
					return
				}
			}
		}
	}
	return
}

// chooseRandomRecipe chooses random recipe from list provided.
// If recipe does not exist, an empty recipe is returned.
func chooseRandomRecipe(recipeList recipes.Recipes, recipeType recipes.DishType) recipes.Recipe {
	// checking if recipe of specified type is present
	// TODO: make an array of recipe types so users can add their
	found := false
	for _, rec := range recipeList {
		if rec.Typ == recipeType {
			found = true
		}
	}
	if !found {
		mylog.Printf(mylogger.WARN+"Recipe of type %q was not found\n", recipeType.String())
		fmt.Printf("Recipe of type %q was not found\n", recipeType.String())
		return recipes.Recipe{}
	}

	newRecipe := recipeList[rand.Intn(len(recipeList))]
	for newRecipe.Typ != recipeType {
		newRecipe = recipeList[rand.Intn(len(recipeList))]
	}
	return newRecipe
}

// chooseRecipes will return in a single array recipes that were chosen by user.
func chooseRecipes(recipeList recipes.Recipes, remainingProducts kitchen.RemainingProducts) recipes.Recipes {
	var lunch, salad recipes.Recipe
	allRecipes := make(recipes.Recipes, 0)

	mylog.Println(mylogger.INFO + "Choosing recipes")
	breakfasts := make(recipes.Recipes, 0, 7)
	for _, recipeBreakfast := range recipeList {
		if recipeBreakfast.Typ == recipes.BREAKFAST {
			breakfasts = append(breakfasts, recipeBreakfast)
		}
	}

	for repeat := true; repeat; {
		answer := ""
		for answer != "y" && answer != "n" {
			fmt.Print("Do you wish to take products that you already have in the kitchen into account when choosing a new recipe?(y/n):")
			fmt.Scan(&answer)
		}

		if answer == "y" {
			if newRecipe, ok := chooseNewRecipeFromExistingProducts(recipeList, remainingProducts, recipes.LUNCH); !ok {
				lunch = chooseRandomRecipe(recipeList, recipes.LUNCH)
			} else {
				lunch = newRecipe
			}
			if newRecipe, ok := chooseNewRecipeFromExistingProducts(recipeList, remainingProducts, recipes.SALAD); !ok {
				salad = chooseRandomRecipe(recipeList, recipes.SALAD)
			} else {
				salad = newRecipe
			}
		} else {
			lunch = chooseRandomRecipe(recipeList, recipes.LUNCH)
			salad = chooseRandomRecipe(recipeList, recipes.SALAD)
		}
		fmt.Printf("Here are the recipes which were chosen for you:\nLunch - %s\nSalad - %s\n\n", lunch.Name, salad.Name)

		satisfied := ""
		for satisfied != "y" && satisfied != "n" {
			fmt.Print("\nAre you satisfied with the choice?(y/n):")
			fmt.Scan(&satisfied)
		}

		if satisfied == "n" {
			mylog.Println(mylogger.INFO + "Repeating procedure")
			repeat = true
		} else {
			repeat = false
		}
	}

	allRecipes = append(allRecipes, breakfasts...)
	allRecipes = append(allRecipes, lunch)
	allRecipes = append(allRecipes, salad)

	return allRecipes
}

// getShoppingCartList will derive shopping list from the recipes that will be cooked.
func getShoppingCartList(allRecipes recipes.Recipes, productList recipes.Products, remainingProducts kitchen.RemainingProducts) recipes.Ingredients {
	shoppingCartList := make(recipes.Ingredients)
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
	for _, recp := range allRecipes {
		for name, ing := range recp.GetIngredients() {
			totalNeeded[name] = correctIngredientSum(name, ing)
		}
	}

	// first, shortage of products is calculated, then, based on calculated shortage of products we
	// calculate amount of products to be purchased.
	// if shortage is more than minimal purchasing quantity and the product is packaged (only discreete
	// steps of weight are allowed), then we buy it n times (n being ceil(diff/qty))
	var productDiff int
	for name := range totalNeeded {
		if totalNeeded[name] > remainingProducts[name] {
			productDiff = totalNeeded[name] - remainingProducts[name]

			if productDiff <= productList[name].ShopQty {
				shoppingCartList[name] = productList[name].ShopQty
			} else if productList[name].Discrete {
				shoppingCartList[name] = productList[name].ShopQty * int(math.Ceil(float64(productDiff)/float64(productList[name].ShopQty)))
			} else if !productList[name].Discrete {
				shoppingCartList[name] = productDiff
			}
		}
	}

	return shoppingCartList
}

// saveResultsInAFile writes shopping list to the file.
func saveResultsInAFile(allRecipes recipes.Recipes, shoppingCartList recipes.Ingredients) {
	mylog.Println(mylogger.INFO + "Saving shopping list in a file")

	if _, err := os.Stat("shopping-lists"); os.IsNotExist(err) {
		if err := os.Mkdir("shopping-lists", 0755); err != nil {
			log.Fatalf(mylogger.ERROR+"Could not create \"shopping-lists\" directory: %v\n", err)
		}
	}

	fout, err := os.OpenFile("shopping-lists/shopping-list_"+time.Now().Format(time.DateOnly)+".txt", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		mylog.Printf(mylogger.WARN+"Could not create shopping list for writing: %v\n", err)
		return
	}
	defer fout.Close()

	writer := bufio.NewWriter(fout)
	i := 1
	for name, weight := range shoppingCartList {
		fmt.Fprintf(writer, "%d) %s (%dg)\n", i, name, weight)
		i++
	}

	breakfasts := make(recipes.Recipes, 0)
	lunches := make(recipes.Recipes, 0)
	soups := make(recipes.Recipes, 0)
	salads := make(recipes.Recipes, 0)
	for _, recipe := range allRecipes {
		switch recipe.Typ {
		case recipes.BREAKFAST:
			breakfasts = append(breakfasts, recipe)
		case recipes.LUNCH:
			lunches = append(lunches, recipe)
		case recipes.SALAD:
			salads = append(salads, recipe)
		case recipes.SOUP:
			soups = append(soups, recipe)
		}
	}

	if len(breakfasts) > 0 {
		fmt.Fprintln(writer, "\nBreakfasts")
		i := 1
		for _, breakf := range breakfasts {
			fmt.Fprintf(writer, "%d) %s\n", i, breakf.Name)
			i++
		}
	}

	if len(lunches) > 0 {
		fmt.Fprintln(writer, "\nLunches")
		i := 1
		for _, lunch := range lunches {
			fmt.Fprintf(writer, "%d) %s\n", i, lunch.Name)
			i++
		}
	}

	if len(soups) > 0 {
		fmt.Fprintln(writer, "\nSoups")
		i := 1
		for _, soup := range soups {
			fmt.Fprintf(writer, "%d) %s\n", i, soup.Name)
			i++
		}
	}

	if len(salads) > 0 {
		fmt.Fprintln(writer, "\nSoups")
		i := 1
		for _, salad := range salads {
			fmt.Fprintf(writer, "%d) %s\n", i, salad.Name)
			i++
		}
	}

	writer.Flush()
}

func main() {
	mylog = mylogger.NewLogger("MAIN: ")
	defer mylogger.CloseResources()
	mylog.Println(mylogger.INFO + "Beginning execution")

	productList := recipes.GetProducts()
	recipeList := recipes.GetRecipes()
	remainingProducts := kitchen.ReadKitchenState()

	// TODO: restructure this when scheduled recipes will be implemented
	allRecipes := chooseRecipes(recipeList, remainingProducts)
	shoppingCartList := getShoppingCartList(allRecipes, productList, remainingProducts)

	mylog.Println(mylogger.INFO + "Outputting result")
	fmt.Println("Here's the list of products that you need to buy for the next week:")
	i := 1
	for name, product := range shoppingCartList {
		fmt.Printf("%d) %s (%dg)\n", i, name, product)
		i++
	}
	saveResultsInAFile(allRecipes, shoppingCartList)

	// imitate purchasing and consuming all products
	mylog.Println(mylogger.INFO + "Calculating how much food will remain")
	for name, productWeight := range shoppingCartList {
		remainingProducts[name] += productWeight
	}

	for _, recps := range allRecipes {
		remainingProducts.CookRecipe(recps)
	}

	mylog.Println(mylogger.INFO + "Asking if the following remaining food will be consumed")

	// additionalNutrition is used to store amount of nutrition consumed with products
	// It will be evenly distributed among workdays
	var additionalNutrition recipes.Product
	const CUTOFF_WEIGHT int = 100
	for name, weight := range remainingProducts {
		if weight <= CUTOFF_WEIGHT {
			answer := ""
			for answer != "y" && answer != "n" {
				fmt.Printf("\nDo you wish to consume %s (%dg) entirely?(y/n):", name, weight)
				fmt.Scan(&answer)
			}

			if answer == "y" {
				mylog.Printf(mylogger.INFO+"Consuming product %s\n", name)
				additionalNutrition.Calr += productList[name].Calr
				additionalNutrition.Carb += productList[name].Carb
				additionalNutrition.Prot += productList[name].Prot
				additionalNutrition.Fats += productList[name].Fats
				remainingProducts.DeleteProduct(name)
			}
		}
	}

	kitchen.SaveKitchenState(remainingProducts)

	// mylog.Println(mylogger.INFO + "Calculating total energy input for every day")
	// --IMPORTANT-- this will be available when generalization of recipe consumption will be added (with flexible schedules which user will set)

	mylog.Println(mylogger.INFO + "Exiting...")
}
