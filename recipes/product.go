package recipes

import (
	"fmt"
	"main/mylogger"
	"os"
	"strings"
)

type (
	// Product structure has nutritional information and weight, in grams (used for recipes).
	Product struct {
		// Weight           int
		Calr             int
		Prot, Fats, Carb float32
	}
	Products map[string]Product
)

// AddNewProductToList adds a new product to database. If error is encountered, it is returned.
func AddNewProductToList(nameProd string, newProd Product) error {
	mylog.Printf(mylogger.INFO+"AddNewProductToList is called with arguments: %s, %v\n", nameProd, newProd)
	fout, err := os.OpenFile("data/products.dat", os.O_WRONLY, os.ModePerm)
	if err != nil {
		mylog.Printf(mylogger.WARN+"Could not open products database for writing: %v\n", err)
		return err
	}
	fout.WriteString(fmt.Sprintf("\n%s,%d,%f,%f,%f", nameProd, newProd.Calr, newProd.Prot, newProd.Fats, newProd.Carb))
	fout.Close()
	return nil
}

// DeleteProductFromList deletes product with specified database. If product is not found, nothing will change.
// Returned value is true when deletion was successful.
func DeleteProductFromList(name string) bool {
	mylog.Printf(mylogger.INFO+"DeleteProductFromList is called with arguments: %s\n", name)
	data := readProductsFromFile()
	if data == nil {
		return false
	}

	for i, product := range data {
		if strings.Contains(product, name) {
			data = append(data[:i], data[i+1:]...)
		}
	}

	if err := SaveProductsToFileString(data); err != nil {
		return false
	}
	return true
}

// GetNutrition sums up a nutritious value of a product list and returns it in a new product.
func (p Products) GetNutrition() (totalNut Product) {
	for _, pr := range p {
		totalNut.Calr += pr.Calr
		totalNut.Carb += pr.Carb
		totalNut.Fats += pr.Fats
		totalNut.Prot += pr.Prot
	}
	return
}
