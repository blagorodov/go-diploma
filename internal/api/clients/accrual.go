package clients

import (
	"context"
	"diploma/internal/config"
	"diploma/internal/errs"
	"diploma/internal/logger"
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/http"
	"time"
)

type AccrualClient struct {
	client  *resty.Client
	retryAt *time.Time
}

type Order struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}

func NewAccrualClient(url string) AccrualClient {
	return AccrualClient{
		client: resty.New().SetBaseURL(url),
	}
}

func (client *AccrualClient) GetOrderInfo(ctx context.Context, orderNumber string) (*Order, error) {
	logger.Log("Poll accrual system")
	order := new(Order)

	response, err := client.client.R().
		SetContext(ctx).
		SetResult(order).
		Get(fmt.Sprintf("%s/api/orders/%s", config.Options.AccrualAddress, orderNumber))

	logger.Log("- error: ")
	logger.Log(err)

	if err != nil {
		return nil, errs.ErrNoAccrual
	}

	logger.Log("- step2")

	if err = client.isBlocked(response); err != nil {
		return nil, err
	}
	logger.Log("- step3")

	logger.Log(fmt.Sprintf("got from accrual [%s] order: %s, message: %s",
		response.Status(),
		orderNumber,
		string(response.Body())))

	if response.StatusCode() == http.StatusNoContent {
		return nil, nil
	}

	return order, nil
}

func (client *AccrualClient) CanMakeRequest() error {
	if client.retryAt == nil {
		return nil
	}

	if time.Now().After(*client.retryAt) {
		client.retryAt = nil
		return nil
	}

	return fmt.Errorf("accrual client will be unblocked in %s", time.Until(*client.retryAt))
}

func (client *AccrualClient) isBlocked(response *resty.Response) error {
	if !response.IsError() && response.StatusCode() != http.StatusTooManyRequests {
		return nil
	}

	retryAfter, err := time.ParseDuration(fmt.Sprintf("%ss", response.Header().Get("Retry-After")))
	if err != nil {
		return err
	}

	body := string(response.Body())
	logger.Log("Accrual system is blocked")
	*client.retryAt = time.Now().Add(retryAfter)

	return fmt.Errorf(body)
}
