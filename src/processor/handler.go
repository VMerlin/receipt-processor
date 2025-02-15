package processor

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func Handle(service Service, w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			handleError(r, w)
		}
	}()
	HandleRequest(service, w, r)
}

func HandleRequest(service Service, w http.ResponseWriter, r *http.Request) {

	var responseBytes []byte
	var responseErr error
	pointsEndpoint := regexp.MustCompile(`^/receipts/(\S+)/points$`)
	matchGetPointsEP := pointsEndpoint.MatchString(r.URL.Path)

	switch {
	case r.URL.Path == "/receipts/process" && r.Method == http.MethodPost:
		request, err := buildAndValidateRequest(r)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error %v", err), http.StatusBadRequest)
			return
		}

		response := service.HandleReceiptProcessing(request)
		responseBytes, responseErr = json.Marshal(response)

	case matchGetPointsEP && r.Method == http.MethodGet:
		matches := pointsEndpoint.FindStringSubmatch(r.URL.Path)
		if len(matches) > 1 {
			id := matches[1]
			response := service.GetPoints(id)
			responseBytes, responseErr = json.Marshal(response)
		} else {
			http.Error(w, "Invalid ID format", http.StatusBadRequest)
			return
		}

	default:
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	if responseErr != nil {
		handleError(responseErr, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseBytes)
}

func buildAndValidateRequest(httpRequest *http.Request) (*Request, error) {
	rawRequest := &RawRequest{}
	if httpRequest.ContentLength > 0 {
		if err := json.NewDecoder(httpRequest.Body).Decode(rawRequest); err != nil {
			return nil, err
		}
	}

	request, err := validateRequest(rawRequest)

	if err != nil {
		return nil, err
	}

	return request, nil
}

func validateRequest(body *RawRequest) (*Request, error) {
	request := &Request{}
	const invalidReceiptError = "The receipt is invalid."

	if len(body.Items) == 0 {
		return nil, errors.New(invalidReceiptError)
	}

	retailerNameMatch, err := regexp.MatchString("^[\\w\\s\\-&]+$", body.Retailer)
	if err != nil || !retailerNameMatch {
		return nil, errors.New(invalidReceiptError)
	}
	request.Retailer = body.Retailer

	if matched, err := regexp.MatchString("^\\d+\\.\\d{2}$", body.Total); err != nil || !matched {
		return nil, errors.New(invalidReceiptError)
	}
	total, err := strconv.ParseFloat(body.Total, 64)
	if err != nil {
		return nil, errors.New(invalidReceiptError)
	}
	request.Total = float64(total)

	date, err := time.Parse("2006-01-02", body.PurchaseDate)
	if err != nil {
		return nil, errors.New(invalidReceiptError)
	}
	request.PurchaseDay = date.Day()

	hour, err := checkAndExtractHour(body.PurchaseTime)
	if err != nil {
		return nil, errors.New(invalidReceiptError)
	}
	request.PurchaseHour = hour

	var items []Item
	for _, rawItem := range body.Items {
		matchedPrice, err1 := regexp.MatchString("^\\d+\\.\\d{2}$", rawItem.Price)
		matchedDescription, err2 := regexp.MatchString("^[\\w\\s\\-]+$", rawItem.Description)
		if err1 != nil || err2 != nil || !matchedPrice || !matchedDescription {
			return nil, errors.New(invalidReceiptError)
		}
		floatValue, err := strconv.ParseFloat(rawItem.Price, 64)
		if err != nil {
			return nil, errors.New(invalidReceiptError)
		}

		item := Item{
			description: strings.TrimSpace(rawItem.Description),
			price:       floatValue,
		}
		items = append(items, item)
	}
	request.Items = items

	return request, nil
}

func checkAndExtractHour(timeStr string) (int, error) {
	pattern := `^([01]?[0-9]|2[0-3])\:([0-5][0-9])$`
	re := regexp.MustCompile(pattern)

	matches := re.FindStringSubmatch(timeStr)
	if matches == nil {
		return 0, errors.New("invalid time format")
	}

	hour, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, err
	}

	return hour, nil
}

func handleError(err interface{}, writer http.ResponseWriter) {
	http.Error(writer, fmt.Sprintf("Error %v", err), http.StatusBadRequest)
}
