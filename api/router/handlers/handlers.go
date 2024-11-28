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

// Struct defining the receipt as a whole
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

	// Associate generated UUID with processed receipt
	receiptId := uuid.New().String()
	receipt.ReceiptId = receiptId

	receipts = append(receipts, receipt)
	receiptIdJson, err := json.Marshal(map[string]string{ "id": receiptId })
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// json.Marshal(receiptId)
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

	// Loop through receipts array to find a matching receipt based on receiptId
	for _, receipt := range receipts {
		// Validate that receiptId exists on the receipt
		if receipt.ReceiptId != "" {
			if receipt.ReceiptId == receiptId {
				// Calculate points by rule
				points += getAlphanumPoints(receipt.Retailer)
				// Check if dollar amount is round by seeing if it matches truncated value
				floatTotal, err := strconv.ParseFloat(receipt.Total, 64)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				if floatTotal == math.Trunc(floatTotal) {
					points += 50
				}
				// Check if total is a multiple of 0.25
				if math.Mod(floatTotal, 0.25) == 0 {
					points += 25
				}
				// Award 5 points for every 2 items on the receipt
				itemCount := len(receipt.Items)
				points += (itemCount / 2) * 5
				// Award points based on trimmed description length
				items := receipt.Items
				for _, item := range items {

					// Trim description of all whitespace
					description := item.ShortDescription
					description = strings.TrimSpace(description)

					// Getting length of trimmed description
					descriptionLength := len(strings.TrimSpace(description))
					// Check if description length is a multiple of 3
					if math.Mod(float64(descriptionLength), 3) == 0 {
						// Convert price to float
						priceFloat, err := strconv.ParseFloat(item.Price, 64)
						if err != nil {
							http.Error(w, err.Error(), http.StatusBadRequest)
							return
						}
						// Calculate rounded and converted points value
						points += int(math.Ceil(priceFloat * 0.2))
					}
				}

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
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}

				if purchaseTimeInt > 1400 && purchaseTimeInt < 1600 {
					points += 10
				}

			}
		}

	}
	pointsJson, err := json.Marshal(map[string]int{ "points": points })
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte(pointsJson))
	return
}

// Helper function to calculate points based on the alphanumeric rule
func getAlphanumPoints(retailerName string) int {
	regex := regexp.MustCompile("^[a-zA-Z0-9]*$")
	alphaNumPoints := 0

	// Loop through provided string as an array
	for _, c := range retailerName {
		if regex.MatchString(string(c)) {
			alphaNumPoints += 1
		}
	}

	return alphaNumPoints
}
