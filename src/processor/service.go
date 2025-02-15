package processor

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"strings"
	"sync"
)

type Service interface {
	HandleReceiptProcessing(request *Request) *Response1
	GetPoints(id string) *Response2
}

type service struct {
	context  context.Context
	receipts map[string]int64
	mutex    sync.Mutex
}

func New() Service {
	return &service{
		context:  context.Background(),
		receipts: make(map[string]int64),
	}
}

func (s *service) HandleReceiptProcessing(request *Request) *Response1 {
	response := &Response1{}

	id := uuid.New().String()
	uniqueId := strings.ReplaceAll(id, "-", "")

	var points int64
	points += int64(calculatePointsForRetailerName(request.Retailer))
	points += int64(calculatePointsForDay(request.PurchaseDay))
	points += int64(calculatePointsForTime(request.PurchaseHour))
	points += calculatePointsForItems(request.Items)
	points += calculatePointsForTotal(request.Total)

	response.Id = uniqueId
	s.mutex.Lock()
	s.receipts[uniqueId] = points
	s.mutex.Unlock()
	fmt.Println("Receipt ID: ", uniqueId, "   Points = ", points)
	return response
}

func (s *service) GetPoints(id string) *Response2 {
	response := &Response2{}
	if len(s.receipts) > 0 {
		response.Points = s.receipts[id]
	}
	return response
}
