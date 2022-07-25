package data

import (
	log "billing-api/logging"
	models "billing-api/models"
	"errors"
	"fmt"
	"os"
)

var (
	StarterPrice float64
	ProPrice     float64
	ElitePrice   float64
)

var ResourcePricing models.ResourcePricing

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

func parseResourcePricing() models.ResourcePricing {
	resourcepricing := models.ResourcePricing{}

	file1, err1 := os.Open("./config/pricing/CPUCoreHourPrice.txt")
	if err1 != nil {
		log.Fatal("error parsing MemoryByteHourPrice ", err1)
	}
	defer file1.Close()
	fmt.Fscan(file1, &resourcepricing.MemoryByteHourPrice)

	file2, err2 := os.Open("./config/pricing/CPUCoreHourPrice.txt")
	if err2 != nil {
		log.Fatal("error parsing CPUCoreHourPrice ", err2)
	}
	defer file2.Close()
	fmt.Fscan(file2, &resourcepricing.CPUCoreHourPrice)

	file3, err3 := os.Open("./config/pricing/networkReceiveBytesPrice.txt")
	if err3 != nil {
		log.Fatal("error parsing networkReceiveBytesPrice ", err3)
	}
	defer file3.Close()
	fmt.Fscan(file3, &resourcepricing.NetworkReceiveBytesPrice)

	file4, err4 := os.Open("./config/pricing/networkTransferBytesPrice.txt")
	if err4 != nil {
		log.Fatal("error parsing networkTransferBytesPrice ", err4)
	}
	defer file4.Close()
	fmt.Fscan(file4, &resourcepricing.NetworkTransferBytesPrice)

	file5, err5 := os.Open("./config/pricing/PVByteHourPrice.txt")
	if err5 != nil {
		log.Fatal("error parsing PVByteHourPrice ", err5)
	}
	defer file5.Close()
	fmt.Fscan(file5, &resourcepricing.PVByteHourPrice)

	ResourcePricing = resourcepricing
	return resourcepricing
}

func GetPlanPrice(planName models.Plan) (float64, error) {

	switch planName {
	case "Starter":
		return StarterPrice, nil
	case "Pro":
		return ProPrice, nil
	case "ElitePrice":
		return ElitePrice, nil
	}
	return 0.0, errors.New("plan not found!")
}

func GetResourcesPrices() models.ResourcePricing {
	return ResourcePricing
}

func ParseFiles() {
	parsePlansPricing()
	parseResourcePricing()
}
