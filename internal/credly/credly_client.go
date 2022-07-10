package credly

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type CredlyClient interface {
	GetBadgesByUsername(username string) []Badge
	GetMostRecentBadgeByUsername(username string) (Badge, error)
}

type credlyClient struct {
	client http.Client
}

func NewCredlyClient(client http.Client) CredlyClient {

	return credlyClient{
		client,
	}
}

type CredlyBadgeResponse struct {
	Data []Badge `json:"data"`
	//TODO get metadata for paging
}

type Badge struct {
	ID            string        `json:"id"`
	IssuedAt      time.Time     `json:"issued_at"`
	IssuedTo      string        `json:"issued_to"`
	BadgeTemplate BadgeTemplate `json:"badge_template"`
}

type BadgeTemplate struct {
	ImageURL string `json:"image_url"`
}

func (c credlyClient) GetBadgesByUsername(username string) []Badge {
	req, err := http.NewRequest("GET",
		fmt.Sprintf("https://www.credly.com/users/%s/badges?sort=-state_updated_at&page=1", username),
		nil,
	)

	if err != nil {
		panic(err)
	}

	req.Header.Set("Accept", "application/json")

	resp, err := c.client.Do(req)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}

	var credlyResponse CredlyBadgeResponse

	err = json.Unmarshal(responseBody, &credlyResponse)
	if err != nil {
		panic(err)
	}

	return credlyResponse.Data
}

func (c credlyClient) GetMostRecentBadgeByUsername(username string) (Badge, error) {
	badges := c.GetBadgesByUsername(username)

	if len(badges) > 0 {
		return badges[0], nil
	}

	return Badge{}, errors.New("no badges for user")
}

//https://api.credly.com/v1/<endpoint_path>

///Accept: application/json
//Authorization: Basic SFRkYXVTanFYeWVzNUxieExPNUdadzo=
//Content-Type: application/json

//curl -X GET 'https://www.credly.com/users/robert-banaszak/badges?sort=most_popular&page=1' -H 'Accept: application/json'
