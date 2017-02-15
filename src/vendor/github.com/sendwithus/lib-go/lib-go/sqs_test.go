package swu

import (
	"testing"
)

// We should build a queue for this smoke test that always has items in it or something.
func TestQueueSize(t *testing.T) {
	service := NewSQSService()
	count, err := service.GetQueueSize("1_prd_CreateOrUpdateCustomerWorker", true)
	if err != nil {
		t.Error(err.Error())
	}
	if count == 0 {
		t.Error("Invalid queue size, this queue should have stuff in it.")
	}
	shared := NewSharedSQSService()
	count, err = shared.GetQueueSize("1_prd_CreateOrUpdateCustomerWorker", true)
	if err != nil {
		t.Error(err.Error())
	}
	if count == 0 {
		t.Error("Invalid queue size, this queue should have stuff in it.")
	}
}