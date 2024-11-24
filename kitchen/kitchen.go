package kitchen

import (
	"bufio"
	"main/mylogger"
	"os"
	"strconv"
	"strings"
)

var mylog *mylogger.Mylogger

// init creates logger for the module
func init() {
	mylog = mylogger.NewLogger("KTCH: ")
}

// ReadKitchenState reads how much (in grams) I have of each ingredient and returns that mylogger.INFO in a map.
func ReadKitchenState() RemainingProducts {
	mylog.Println(mylogger.INFO + "ReadKitchenState is executing")

	remainingProd := make(RemainingProducts)
	for _, line := range readFromKitchenStateFile() {
		parsedLine := strings.Split(line, ",") // format: name,weight

		if weight, err := strconv.Atoi(parsedLine[1]); err != nil {
			mylog.Printf(mylogger.ERROR+"Couldn't parse weight in remaining products entry \"%s\", %v\n", parsedLine[1], err)
		} else {
			remainingProd[parsedLine[0]] = weight
		}
	}

	return remainingProd
}

// readFromKitchenStateFile returns lines of the "available products" file.
func readFromKitchenStateFile() []string {
	mylog.Println(mylogger.INFO + "readFromKitchenStateFile is executing")

	fin, err := os.OpenFile("data/available-products.dat", os.O_RDONLY, os.ModePerm)
	if err != nil {
		mylog.Printf(mylogger.WARN+"Could not open available products database: %v\n", err)
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

// SaveKitchenState converts products to string array and then writes product data to the file.
func SaveKitchenState(pr RemainingProducts) error {
	data := make([]string, 0)
	for name, weight := range pr {
		data = append(data, name+","+strconv.Itoa(weight)+"\n")
	}
	return SaveKitchenStateString(data)
}

// SaveKitchenStateString writes product data to the file.
func SaveKitchenStateString(pr []string) error {
	mylog.Println(mylogger.INFO + "SaveKitchenStateString is executing")
	fout, err := os.OpenFile("data/available-products.dat", os.O_TRUNC, 0644)
	if err != nil {
		mylog.Printf(mylogger.WARN+"Could not open products database for writing: %v\n", err)
		return err
	}
	defer fout.Close()

	buf := bufio.NewWriter(fout)
	for _, p := range pr {
		buf.WriteString(p)
	}
	buf.Flush()

	return nil
}
