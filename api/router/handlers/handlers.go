package api

import (
	"encoding/json"
	"io"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var receipts []Receipt

// Struct defining each individual item in a receipt
type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

type Receipt struct {
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Items        []Item `json:"items"`
	Total        string `json:"total"`
	ReceiptId    string `json:"receiptId"`
}

func ProcessReceipts(w http.ResponseWriter, r *http.Request) {
	reqBody := r.Body
	defer reqBody.Close()

	body, err := io.ReadAll((reqBody))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var receipt Receipt
	err = json.Unmarshal(body, &receipt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	receiptId := uuid.New().String()
	receipt.ReceiptId = receiptId

	receipts = append(receipts, receipt)
	receiptIdJson, err := json.Marshal(map[string]string{"id": receiptId})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write([]byte(receiptIdJson))
}

func CalculatePoints(w http.ResponseWriter, r *http.Request) {
	// Return an error if no receipts have been processed yet.
	if len(receipts) == 0 {
		http.Error(w, "No receipts available for points calculation.", http.StatusBadRequest)
		return
	}

	points := 0

	receiptId, ok := mux.Vars(r)["id"]
	if !ok {
		http.Error(w, "No receipt ID available.", http.StatusBadRequest)
	}

	for _, receipt := range receipts {
		if receipt.ReceiptId != "" && receipt.ReceiptId == receiptId {
			// Increment points based on the # of alphanumeric characters in the retailer name
			points += getAlphanumPoints(receipt.Retailer)

			floatTotal, err := strconv.ParseFloat(receipt.Total, 64)
			if err != nil {
				http.Error(w, "Error converting total to float.", http.StatusBadRequest)
				return
			}

			// Check if total is a whole number; if so, award 50 points
			points += getRoundPoints(floatTotal)

			if math.Mod(floatTotal, 0.25) == 0 {
				points += 25
			}

			itemCount := len(receipt.Items)
			points += (itemCount / 2) * 5

			// Award points based on trimmed description length of items
			items := receipt.Items
			descPoints, err := getDescriptionPoints(items)
			if err != nil {
				http.Error(w, "Error calculating description length points.", http.StatusBadRequest)
				return
			}
			points += descPoints

			// Award points if purchase day is odd
			purchaseDate := receipt.PurchaseDate
			purchaseDateArray := strings.Split(purchaseDate, "-")
			purchaseDay, err := strconv.ParseFloat(purchaseDateArray[2], 64)
			if math.Mod(float64(purchaseDay), 2.0) != 0 {
				points += 6
			}

			// Award points if purchase time is between 2pm and 4pm
			purchaseTime := receipt.PurchaseTime
			purchaseTimeArray := strings.Split(purchaseTime, ":")
			purchaseTimeInt, err := strconv.Atoi(strings.Join(purchaseTimeArray, ""))
			if err != nil {
				http.Error(w, "Error converting purchase time to int.", http.StatusBadRequest)
				return
			}

			if purchaseTimeInt > 1400 && purchaseTimeInt < 1600 {
				points += 10
			}

		}

	}
	pointsJson, err := json.Marshal(map[string]int{"points": points})
	if err != nil {
		http.Error(w, "Error marshalling points into JSON.", http.StatusBadRequest)
		return
	}

	w.Write([]byte(pointsJson))
	return
}

// Helper function to calculate points based on the alphanumeric rule
func getAlphanumPoints(retailerName string) int {
	regex := regexp.MustCompile("^[a-zA-Z0-9]*$")
	alphaNumPoints := 0

	for _, c := range retailerName {
		if regex.MatchString(string(c)) {
			alphaNumPoints += 1
		}
	}
	return alphaNumPoints
}

func getRoundPoints(total float64) int {
	if total == math.Trunc(total) {
		return 50
	}
	return 0
}

func getDescriptionPoints(items []Item) (int, error) {
	pointsToReturn := 0
	for _, item := range items {

		description := item.ShortDescription
		descriptionLength := len(strings.TrimSpace(description))

		if math.Mod(float64(descriptionLength), 3) == 0 {
			priceFloat, err := strconv.ParseFloat(item.Price, 64)
			if err != nil {
				return 0, err
			}
			pointsToReturn += int(math.Ceil(priceFloat * 0.2))
		}
	}
	return pointsToReturn, nil
}
