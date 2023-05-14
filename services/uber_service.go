package services

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/atomi-ai/atomi/models"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
)

const (
	AuthURL = "https://login.uber.com/oauth/v2/token"
	BaseURL = "https://api.uber.com/v1/customers"
)

type UberService interface {
	Quote(requestBody *models.QuoteRequest) (*models.QuoteResponse, error)
	CreateDelivery(requestBody *models.DeliveryData) (*models.DeliveryResponse, error)
	GetDelivery(deliveryID string) (*models.DeliveryResponse, error)
}

type UberServiceImpl struct {
	HTTPClient                  *resty.Client
	ClientID                    string
	ClientSecret                string
	DaasURL                     string
	Accessauthorization         string
	authorizationExpirationTime int64
}

func NewUberService() UberService {
	httpClient := resty.New()
	return &UberServiceImpl{
		HTTPClient:   httpClient,
		ClientID:     viper.GetString("uberClientId"),
		ClientSecret: viper.GetString("uberClientSecret"),
		DaasURL:      BaseURL + "/" + viper.GetString("uberCustomId"),
	}
}

func (u *UberServiceImpl) getAuthorization() (string, error) {
	if u.Accessauthorization == "" || (u.authorizationExpirationTime != 0 && time.Now().Unix() >= u.authorizationExpirationTime) {
		response := &models.TokenResponse{}
		errorResponse := &models.ErrorResponse{}
		resp, err := u.HTTPClient.R().
			SetFormData(map[string]string{
				"grant_type":    "client_credentials",
				"client_id":     u.ClientID,
				"client_secret": u.ClientSecret,
				"scope":         "eats.deliveries",
			}).
			SetResult(response).
			SetError(errorResponse).
			Post(AuthURL)

		if err != nil {
			return "", err
		}

		if resp.StatusCode() != http.StatusOK {
			return "", errors.New(errorResponse.ToString())
		}

		u.Accessauthorization = fmt.Sprintf("%s %s", response.TokenType, response.AccessToken)
		u.authorizationExpirationTime = time.Now().Unix() + response.ExpiresIn - 300
	}

	return u.Accessauthorization, nil
}

func (u *UberServiceImpl) Quote(requestBody *models.QuoteRequest) (*models.QuoteResponse, error) {
	authorization, err := u.getAuthorization()
	if err != nil {
		return nil, err
	}
	url := u.DaasURL + "/delivery_quotes"
	response := &models.QuoteResponse{}
	errorResponse := &models.ErrorResponse{}
	resp, err := u.HTTPClient.R().
		SetHeader("Authorization", authorization).
		SetBody(requestBody).
		SetResult(response).
		SetError(errorResponse).
		Post(url)
	fmt.Printf("Uber POST %s\n%v\n%v\n%v\n", url, *requestBody, resp, err)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, errors.New(errorResponse.ToString())
	}

	return response, nil
}

func (u *UberServiceImpl) CreateDelivery(requestBody *models.DeliveryData) (*models.DeliveryResponse, error) {
	authorization, err := u.getAuthorization()
	if err != nil {
		return nil, err
	}

	url := u.DaasURL + "/deliveries"
	response := &models.DeliveryResponse{}
	errorResponse := &models.ErrorResponse{}
	resp, err := u.HTTPClient.R().
		SetHeader("Authorization", authorization).
		SetBody(requestBody).
		SetResult(response).
		SetError(errorResponse).
		Post(url)
	fmt.Printf("Uber POST %s\n%v\n%v\n%v\n", url, *requestBody, resp, err)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, errors.New(errorResponse.ToString())
	}

	return response, nil
}

func (u *UberServiceImpl) GetDelivery(deliveryID string) (*models.DeliveryResponse, error) {
	authorization, err := u.getAuthorization()
	if err != nil {
		return nil, err
	}

	url := u.DaasURL + "/deliveries/" + deliveryID
	response := &models.DeliveryResponse{}
	errorResponse := &models.ErrorResponse{}
	resp, err := u.HTTPClient.R().
		SetHeader("Authorization", authorization).
		SetResult(response).
		SetError(errorResponse).
		Get(url)
	fmt.Printf("Uber GET %s\n%v\n%v\n", url, resp, err)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, errors.New(errorResponse.ToString())
	}

	return response, nil
}
