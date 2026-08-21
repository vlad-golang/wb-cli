package common

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

type WbClient struct {
	Token string
}

func (w *WbClient) WbSearch(ctx context.Context, query string) (SearchResponse, error) {
	client := &http.Client{}
	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		"https://www.wildberries.ru/__internal/u-search/exactmatch/ru/common/v18/search?ab_testing=false&appType=1&curr=rub&dest=1259570991&hide_dtype=15&hide_vflags=4294967296&inheritFilters=true&lang=ru&locale=ru&query="+url.QueryEscape(
			query,
		)+"&resultset=catalog&sort=popular&spp=30&suppressSpellcheck=false",
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("deviceid", "site_7f2ffa244cbb49599e3678e498a4e726")
	req.Header.Set("cookie", w.Token)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return SearchResponse{}, fmt.Errorf("status code %d", resp.StatusCode)
	}

	var body SearchResponse
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return SearchResponse{}, fmt.Errorf("failed to decode body: %w", err)
	}

	return body, nil
}

type SearchResponse struct {
	Metadata struct {
		CatalogType        string   `json:"catalog_type"`
		CatalogValue       string   `json:"catalog_value"`
		Normquery          string   `json:"normquery"`
		SearchResult       struct{} `json:"search_result"`
		Name               string   `json:"name"`
		Rmi                string   `json:"rmi"`
		Title              string   `json:"title"`
		Rs                 int      `json:"rs"`
		Context            []string `json:"context"`
		Qv                 string   `json:"qv"`
		Kcl                string   `json:"kcl"`
		PresetNormqueryMap struct {
			Num203083314 string `json:"203083314"`
		} `json:"preset_normquery_map"`
	} `json:"metadata"`
	Products []struct {
		ID          int    `json:"id"`
		Root        int64  `json:"root"`
		KindID      int    `json:"kindId"`
		Brand       string `json:"brand"`
		BrandID     int    `json:"brandId"`
		SiteBrandID int    `json:"siteBrandId"`
		Colors      []struct {
			Name string `json:"name"`
			ID   int    `json:"id"`
		} `json:"colors"`
		SubjectID       int     `json:"subjectId"`
		SubjectParentID int     `json:"subjectParentId"`
		SemanticID      []int   `json:"semanticId"`
		Name            string  `json:"name"`
		Entity          string  `json:"entity"`
		MatchID         int     `json:"matchId"`
		Supplier        string  `json:"supplier"`
		SupplierID      int     `json:"supplierId"`
		SupplierRating  float64 `json:"supplierRating"`
		SupplierFlags   int     `json:"supplierFlags"`
		Pics            int     `json:"pics"`
		PicsUpdate      int     `json:"picsUpdate"`
		Rating          int     `json:"rating"`
		ReviewRating    float64 `json:"reviewRating"`
		Feedbacks       int     `json:"feedbacks"`
		PanelPromoID    int     `json:"panelPromoId,omitempty"`
		Volume          int     `json:"volume"`
		Weight          float64 `json:"weight"`
		ViewFlags       int     `json:"viewFlags"`
		Sizes           []struct {
			Name     string `json:"name"`
			OrigName string `json:"origName"`
			Rank     int    `json:"rank"`
			OptionID int    `json:"optionId"`
			Wh       int    `json:"wh"`
			Time1    int    `json:"time1"`
			Time2    int    `json:"time2"`
			Dtype    int64  `json:"dtype"`
			Price    struct {
				Basic     int `json:"basic"`
				Product   int `json:"product"`
				Logistics int `json:"logistics"`
				Return    int `json:"return"`
				Cashback  int `json:"cashback"`
			} `json:"price"`
			SaleConditions int    `json:"saleConditions"`
			Payload        string `json:"payload"`
		} `json:"sizes"`
		TotalQuantity int    `json:"totalQuantity"`
		Time1         int    `json:"time1"`
		Time2         int    `json:"time2"`
		Wh            int    `json:"wh"`
		Dtype         int64  `json:"dtype"`
		Dist          int    `json:"dist"`
		Logs          string `json:"logs,omitempty"`
		Meta          struct {
			Tokens   []any `json:"tokens"`
			PresetID int   `json:"presetId"`
		} `json:"meta"`
		FeedbackPoints int  `json:"feedbackPoints,omitempty"`
		IsNew          bool `json:"isNew,omitempty"`
	} `json:"products"`
	Total int `json:"total"`
}
