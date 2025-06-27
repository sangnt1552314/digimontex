package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/sangnt1552314/digimontex/internal/models"
)

const (
	baseURL = "https://digi-api.com/api/v1/digimon"
)

func GetDigimonList(params models.DigimonSearchQueryParams) ([]models.Digimon, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %v", err)
	}

	q := u.Query()
	if params.Name != "" {
		q.Add("name", params.Name)
	}
	if params.Level != "" {
		q.Add("level", params.Level)
	}
	if params.Page > 0 {
		q.Add("page", strconv.Itoa(params.Page))
	}
	if params.PageSize > 0 {
		q.Add("pageSize", strconv.Itoa(params.PageSize))
	}
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, fmt.Errorf("failed to fetch digimon data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned non-200 status code: %d", resp.StatusCode)
	}

	var apiResp models.DigimonResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	var digimonList []models.Digimon
	for _, d := range apiResp.Content {
		digimonList = append(digimonList, models.Digimon{
			ID:    d.ID,
			Name:  d.Name,
			Href:  d.Href,
			Image: d.Image,
		})
	}

	return digimonList, nil
}

func GetDigimonByName(name string) (*models.DigimonDetail, error) {
	url := fmt.Sprintf("%s/%s", baseURL, url.QueryEscape(name))

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch digimon by name: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned non-200 status code: %d", resp.StatusCode)
	}

	var digimon models.DigimonDetail
	if err := json.NewDecoder(resp.Body).Decode(&digimon); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &digimon, nil
}

func GetDigimonByID(id int) (*models.DigimonDetail, error) {
	url := fmt.Sprintf("%s/%d", baseURL, id)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch digimon by ID: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned non-200 status code: %d", resp.StatusCode)
	}

	var digimon models.DigimonDetail
	if err := json.NewDecoder(resp.Body).Decode(&digimon); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &digimon, nil
}
