package data

import (
	log "billing-api/logging"
	"fmt"
	"os"
)

var (
	StarterPrice float64
	ProPrice     float64
	ElitePrice   float64
)

func parsePlansPricing() {
	//TODO: make ths function scan all files inside config/plans/ and make env variable by "filename:contentOfFile"

	file1, err1 := os.Open("data/pricing/plans/StarterPrice.txt")
	if err1 != nil {
		log.Error("error parsing Starter plan price ", err1)
	}
	defer file1.Close()
	fmt.Fscan(file1, &StarterPrice)

	file2, err2 := os.Open("./data/pricing/plans/ProPrice.txt")
	if err2 != nil {
		log.Error("error parsing Starter plan price ", err2)
	}
	defer file2.Close()
	fmt.Fscan(file2, &ProPrice)

	file3, err3 := os.Open("./data/pricing/plans/ElitePrice.txt")
	if err3 != nil {
		log.Error("error parsing Starter plan price ", err3)
	}
	defer file3.Close()
	fmt.Fscan(file3, &ElitePrice)
}

func ParseFiles() {
	parsePlansPricing()
}
