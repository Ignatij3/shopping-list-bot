package recipes

import (
	"bufio"
	"fmt"
	"main/mylogger"
	"os"
	"strconv"
	"strings"
)

// GetRecipes reads file containing all recipes and returns them.
func GetRecipes() Recipes {
	mylog.Println(mylogger.INFO + "GetRecipes is executing")
	var (
		recipeN         int
		productN        int
		rtypeenum       DishType
		disregardRecipe bool
	)
	rcps := make(Recipes, 0)

	fin, err := os.OpenFile("data/recipes.dat", os.O_RDONLY, os.ModePerm)
	if err != nil {
		mylog.Fatalf(mylogger.ERROR+"Could not open products database: %v\n", err)
	}
	defer fin.Close()
	scanner := bufio.NewScanner(fin)

	// get recipe amount
	scanner.Scan()
	if num, err := strconv.Atoi(scanner.Text()); err == nil {
		recipeN = num
	} else {
		mylog.Fatalf(mylogger.ERROR+"Recipe number couldn't be read: %v\n", err)
	}
	mylog.Printf(mylogger.INFO+"Recipe number is %d\n", recipeN)

	recipeProgress := 0

	// read recipes, if their amount is not correct recipeProgress will not be reset at a correct moment
	for productProgress := 0; scanner.Scan(); productProgress++ {
		if productProgress == 0 {
			// parse recipe header
			recipeHeader := strings.Split(scanner.Text(), ",") // format: name,type,product_n
			if len(recipeHeader) != 3 {
				mylog.Printf(mylogger.ERROR+"Format of recipe header entry in recipes file is not correct: %v (len %d)\n", recipeHeader, len(recipeHeader))
				productProgress = -1
				continue
			}

			switch recipeHeader[1] {
			case "breakfast":
				rtypeenum = BREAKFAST
			case "lunch":
				rtypeenum = LUNCH
			case "salad":
				rtypeenum = SALAD
			case "soup":
				rtypeenum = SOUP
			default:
				mylog.Printf(mylogger.WARN+"Unknown recipe type: \"%s\"\n", recipeHeader[1])
			}

			if disregardRecipe {
				// replacing latest recipe with new one, since old is faulty
				mylog.Printf(mylogger.WARN+"Disregarding recipe \"%s\" due to faulty parsing\n", rcps[len(rcps)-1].Name)
				rcps[len(rcps)-1] = Recipe{Name: recipeHeader[0], Typ: rtypeenum, ing: Ingredients{}}
				disregardRecipe = false
			} else {
				rcps = append(rcps, Recipe{Name: recipeHeader[0], Typ: rtypeenum, ing: Ingredients{}})
			}

			if num, nok := strconv.Atoi(strings.TrimSuffix(recipeHeader[2], "\n")); nok == nil {
				productN = num
			} else {
				mylog.Printf(mylogger.ERROR+"Product amount in recipes file cannot be parsed: %s, %v\n", recipeHeader[2], err)
				disregardRecipe = true
			}

		} else if productProgress < productN {
			// if faulty ingredient is found, skip whole recipe
			if disregardRecipe {
				continue
			}

			// parse recipe ingredients
			productEntry := strings.Split(scanner.Text(), ",") // format: name,weight
			productProgress++

			if len(productEntry) != 2 {
				mylog.Printf(mylogger.ERROR+"Format of product entry in recipes file is not correct: %v (len %d)\n", productEntry, len(productEntry))
				disregardRecipe = true
				continue
			}

			if num, nok := strconv.Atoi(strings.TrimSuffix(productEntry[1], "\n")); nok == nil {
				rcps[len(rcps)-1].ing[productEntry[0]] = num
			} else {
				mylog.Printf(mylogger.ERROR+"product weight in recipes file cannot be parsed: %s, %v\n", productEntry[1], err)
				disregardRecipe = true
			}

		} else {
			productProgress = -1
			recipeProgress++
		}
	}

	if disregardRecipe {
		mylog.Printf(mylogger.ERROR+"Recipe ended sooner than anticipated, disregarding last recipe \"%s\"\n", rcps[len(rcps)-1].Name)
		rcps = rcps[:len(rcps)-1]
		disregardRecipe = false
	}

	if recipeN != recipeProgress {
		mylog.Printf(mylogger.ERROR+"Recipes amount that was declared is incorrect. Amount of recipes declared: %d, amount of recipes read: %d\n", recipeN, recipeProgress)
	}

	return rcps
}

// GetProducts returns parsed product list.
func GetProducts() Products {
	mylog.Println(mylogger.INFO + "GetProducts is executing")
	ings := make(Products)

	var (
		calr             int
		prot, fats, carb float64
		err              error
	)
	for _, productLine := range readProductsFromFile() {
		parsedLine := strings.Split(productLine, ",") // format: name,calories,proteins,fats,carbs

		if len(parsedLine) != 5 {
			mylog.Printf(mylogger.ERROR+"Product entry \"%v\" (len %d) has incorrect format, skipping", parsedLine, len(parsedLine))
			continue
		}

		if calr, err = strconv.Atoi(parsedLine[1]); err != nil {
			mylog.Printf(mylogger.ERROR+"Couldn't parse calories in product entry \"%s\", %v\n", parsedLine[1], err)
			continue
		}

		if prot, err = strconv.ParseFloat(parsedLine[2], 32); err != nil {
			mylog.Printf(mylogger.ERROR+"Couldn't parse proteins in product entry \"%s\", %v\n", parsedLine[2], err)
			continue
		}

		if fats, err = strconv.ParseFloat(parsedLine[3], 32); err != nil {
			mylog.Printf(mylogger.ERROR+"Couldn't parse fats in product entry \"%s\", %v\n", parsedLine[3], err)
			continue
		}

		if carb, err = strconv.ParseFloat(parsedLine[4], 32); err != nil {
			mylog.Printf(mylogger.ERROR+"Couldn't parse carbs in product entry \"%s\", %v\n", parsedLine[4], err)
			continue
		}

		ings[parsedLine[0]] = Product{
			Calr: calr,
			Prot: float32(prot),
			Fats: float32(fats),
			Carb: float32(carb),
		}
	}

	return ings
}

// readProductsFromFile reads products database and returns an array of lines. If an error occurs, function returns nil.
func readProductsFromFile() []string {
	mylog.Println(mylogger.INFO + "readProductsFromFile is executing")

	fin, err := os.OpenFile("data/products.dat", os.O_RDONLY, os.ModePerm)
	if err != nil {
		mylog.Printf(mylogger.WARN+"Could not open products database: %v", err)
		return nil
	}
	defer fin.Close()
	scanner := bufio.NewScanner(fin)

	data := make([]string, 0)
	for scanner.Scan() {
		data = append(data, scanner.Text())
	}

	return data
}

// SaveProductsToFile converts data to string and flushes product data to the file.
func SaveProductsToFile(products Products) error {
	data := make([]string, 0)

	for name, pr := range products {
		data = append(data, fmt.Sprintf("%s,%d,%f,%f,%f\n", name, pr.Calr, pr.Prot, pr.Fats, pr.Carb))
	}
	return SaveProductsToFileString(data)
}

// SaveProductsToFileString flushes product data to the file.
func SaveProductsToFileString(products []string) error {
	mylog.Println(mylogger.INFO + "SaveProductsToFileString is executing")
	fout, err := os.OpenFile("data/products.dat", os.O_TRUNC, 0644)
	if err != nil {
		mylog.Printf(mylogger.WARN+"Could not open products database for writing: %v\n", err)
		return err
	}
	defer fout.Close()

	buf := bufio.NewWriter(fout)
	for _, pr := range products {
		buf.WriteString(pr)
	}
	buf.Flush()

	return nil
}
