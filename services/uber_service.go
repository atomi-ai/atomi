package services

import (
	"fmt"
	"time"

	"github.com/atomi-ai/atomi/models"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
)

const (
	AUTH_URL = "https://login.uber.com/oauth/v2/token"
	BASE_URL = "https://api.uber.com/v1/customers"
)

type UberService interface {
	Quote(requestBody *models.QuoteRequest) (*models.QuoteResponse, error)
	CreateDelivery(requestBody *models.DeliveryData) (*models.DeliveryResponse, error)
	GetDelivery(deliveryId string) (*models.DeliveryResponse, error)
}

type UberServiceImpl struct {
	HttpClient                  *resty.Client
	ClientId                    string
	ClientSecret                string
	DAAS_URL                    string
	Accessauthorization         string
	authorizationExpirationTime int64
}

func NewUberService() UberService {
	httpClient := resty.New()
	return &UberServiceImpl{
		HttpClient:   httpClient,
		ClientId:     viper.GetString("uberClientId"),
		ClientSecret: viper.GetString("uberClientSecret"),
		DAAS_URL:     BASE_URL + "/" + viper.GetString("uberCustomId"),
	}
}

func (u *UberServiceImpl) getAuthorization() (string, error) {
	if u.Accessauthorization == "" || (u.authorizationExpirationTime != 0 && time.Now().Unix() >= u.authorizationExpirationTime) {
		response := &models.TokenResponse{}
		_, err := u.HttpClient.R().
			SetFormData(map[string]string{
				"grant_type":    "client_credentials",
				"client_id":     u.ClientId,
				"client_secret": u.ClientSecret,
				"scope":         "eats.deliveries",
			}).
			SetResult(response).
			Post(AUTH_URL)

		if err != nil {
			return "", err
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
	url := u.DAAS_URL + "/delivery_quotes"
	response := &models.QuoteResponse{}
	resp, err := u.HttpClient.R().
		SetHeader("Authorization", authorization).
		SetBody(requestBody).
		SetResult(response).
		Post(url)
	fmt.Printf("Uber POST %s\n%v\n%v\n%v\n", url, *requestBody, resp, err)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (u *UberServiceImpl) CreateDelivery(requestBody *models.DeliveryData) (*models.DeliveryResponse, error) {
	authorization, err := u.getAuthorization()
	if err != nil {
		return nil, err
	}

	url := u.DAAS_URL + "/deliveries"
	response := &models.DeliveryResponse{}
	resp, err := u.HttpClient.R().
		SetHeader("Authorization", authorization).
		SetBody(requestBody).
		SetResult(response).
		Post(url)
	fmt.Printf("Uber POST %s\n%v\n%v\n%v\n", url, *requestBody, resp, err)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (u *UberServiceImpl) GetDelivery(deliveryId string) (*models.DeliveryResponse, error) {
	authorization, err := u.getAuthorization()
	if err != nil {
		return nil, err
	}

	url := u.DAAS_URL + "/deliveries/" + deliveryId
	response := &models.DeliveryResponse{}
	resp, err := u.HttpClient.R().
		SetHeader("Authorization", authorization).
		SetResult(response).
		Get(url)
	fmt.Printf("Uber GET %s\n%v\n%v\n", url, resp, err)
	if err != nil {
		return nil, err
	}

	return response, nil
}
