package common

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type WbClient struct {
	Token string
}

func (w *WbClient) WbSearch(ctx context.Context, query string, page int) ([]SearchResponse, error) {
	params := url.Values{
		"ab_testing":         {"false"},
		"appType":            {"1"},
		"curr":               {"rub"},
		"dest":               {"1259570991"},
		"hide_dtype":         {"15"},
		"hide_vflags":        {"4294967296"},
		"inheritFilters":     {"true"},
		"lang":               {"ru"},
		"locale":             {"ru"},
		"page":               {strconv.Itoa(page)},
		"query":              {query},
		"resultset":          {"catalog"},
		"sort":               {"popular"},
		"spp":                {"30"},
		"suppressSpellcheck": {"false"},
	}

	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		"https://www.wildberries.ru/__internal/u-search/exactmatch/ru/common/v18/search?"+params.Encode(),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("new search request: %w", err)
	}
	req.Header.Set("deviceid", "site_7f2ffa244cbb49599e3678e498a4e726")
	req.Header.Set("cookie", w.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do search request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code %d", resp.StatusCode)
	}

	var body SearchBody
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return nil, fmt.Errorf("failed to decode body: %w", err)
	}

	searchResp := make([]SearchResponse, 0, len(body.Products))
	for _, product := range body.Products {
		searchResp = append(searchResp, SearchResponse{
			ID:             product.ID,
			Name:           product.Name,
			Brand:          product.Brand,
			PriceRub:       product.Sizes[0].Price.Product / 100,
			ReviewRating:   product.ReviewRating,
			SupplierRating: product.SupplierRating,
			Feedbacks:      product.Feedbacks,
			RootID:         int(product.Root),
		})
	}

	return searchResp, nil
}

func (w *WbClient) GetCard(ctx context.Context, id int) (GetCardBody, error) {
	sID := strconv.Itoa(id)
	l := len(sID)

	client := &http.Client{}
	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		"https://rst-basket-cdn-12.geobasket.ru/vol"+sID[:l-5]+"/part"+sID[:l-3]+"/"+sID+"/info/ru/card.json",
		nil,
	)
	if err != nil {
		return GetCardBody{}, fmt.Errorf("new card request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return GetCardBody{}, fmt.Errorf("do card request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return GetCardBody{}, fmt.Errorf("status code %d", resp.StatusCode)
	}

	var body GetCardBody
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return body, fmt.Errorf("failed to decode body: %w", err)
	}

	return body, nil
}

func (w *WbClient) Feedbacks(ctx context.Context, rootID int) ([]FeedbackResponse, error) {
	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, "GET", "https://feedbacks2.wb.ru/feedbacks/v1/"+strconv.Itoa(rootID), nil)
	if err != nil {
		return nil, fmt.Errorf("new feedbacks request: %w", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do feedbacks request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code %d", resp.StatusCode)
	}

	var body FeedbacksBody
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return nil, fmt.Errorf("failed to decode body: %w", err)
	}

	feedbacks := make([]FeedbackResponse, 0, len(body.Feedbacks))
	for _, feedback := range body.Feedbacks {
		if feedback.Text != "" {
			feedbacks = append(feedbacks, FeedbackResponse{
				Text:   feedback.Text,
				Pros:   feedback.Pros,
				Cons:   feedback.Cons,
				Rating: feedback.ProductValuation,
				Date:   feedback.CreatedDate.Format("2006-01-02 15:04"),
			})
		}
	}

	return feedbacks, nil
}

type SearchResponse struct {
	ID             int     `csv:"id"`
	Name           string  `csv:"name"`
	Brand          string  `csv:"brand"`
	PriceRub       int     `csv:"price_rub"`
	ReviewRating   float64 `csv:"review_rating"`
	SupplierRating float64 `csv:"supplier_rating"`
	Feedbacks      int     `csv:"feedbacks"`
	RootID         int     `csv:"root_id"`
}

type SearchBody struct {
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

type GetCardBody struct {
	ImtID        int    `json:"imt_id"`
	NmID         int    `json:"nm_id"`
	ImtName      string `json:"imt_name"`
	Slug         string `json:"slug"`
	SubjName     string `json:"subj_name"`
	SubjRootName string `json:"subj_root_name"`
	VendorCode   string `json:"vendor_code"`
	Description  string `json:"description"`
	Options      []struct {
		Name             string   `json:"name"`
		Value            string   `json:"value"`
		CharcType        int      `json:"charc_type"`
		IsVariable       bool     `json:"is_variable,omitempty"`
		VariableValues   []string `json:"variable_values,omitempty"`
		VariableValueIDs []int    `json:"variable_value_IDs,omitempty"`
	} `json:"options"`
	Compositions []struct {
		Name string `json:"name"`
	} `json:"compositions"`
	SizesTable struct {
		DetailsProps []string `json:"details_props"`
		Values       []struct {
			TechSize string   `json:"tech_size"`
			ChrtID   int      `json:"chrt_id"`
			Details  []string `json:"details"`
		} `json:"values"`
	} `json:"sizes_table"`
	Certificate struct {
		Verified bool `json:"verified"`
	} `json:"certificate"`
	NmColorsNames string `json:"nm_colors_names"`
	Colors        []int  `json:"colors"`
	Contents      string `json:"contents"`
	FullColors    []struct {
		NmID int `json:"nm_id"`
	} `json:"full_colors"`
	Selling struct {
		BrandName   string `json:"brand_name"`
		BrandHash   string `json:"brand_hash"`
		SupplierID  int    `json:"supplier_id"`
		NoReturnMap int    `json:"no_return_map"`
	} `json:"selling"`
	Media struct {
		HasVideo   bool `json:"has_video"`
		PhotoCount int  `json:"photo_count"`
	} `json:"media"`
	Data struct {
		SubjectID     int   `json:"subject_id"`
		SubjectRootID int   `json:"subject_root_id"`
		ChrtIDs       []int `json:"chrt_ids"`
	} `json:"data"`
	GroupedOptions []struct {
		GroupName string `json:"group_name"`
		Options   []struct {
			Name             string   `json:"name"`
			Value            string   `json:"value"`
			CharcType        int      `json:"charc_type"`
			IsVariable       bool     `json:"is_variable,omitempty"`
			VariableValues   []string `json:"variable_values,omitempty"`
			VariableValueIDs []int    `json:"variable_value_IDs,omitempty"`
		} `json:"options"`
	} `json:"grouped_options"`
	HasSellerRecommendations bool      `json:"has_seller_recommendations"`
	NeedKiz                  bool      `json:"need_kiz"`
	UserFlags                int       `json:"user_flags"`
	Properties               int       `json:"properties"`
	UpdateDate               time.Time `json:"update_date"`
	CreateDate               time.Time `json:"create_date"`
}

type FeedbacksBody struct {
	PhotosUris []any `json:"photosUris"`
	Photo      []int `json:"photo"`
	Photos     []struct {
		ID        int    `json:"id"`
		Key       string `json:"key"`
		IsBlurred bool   `json:"isBlurred"`
		IsReady   bool   `json:"isReady"`
	} `json:"photos"`
	Valuation             string `json:"valuation"`
	ValuationSum          int    `json:"valuationSum"`
	ValuationDistribution struct {
		Num1 int `json:"1"`
		Num2 int `json:"2"`
		Num3 int `json:"3"`
		Num4 int `json:"4"`
		Num5 int `json:"5"`
	} `json:"valuationDistribution"`
	ValuationDistributionPercent struct {
		Num1 int `json:"1"`
		Num2 int `json:"2"`
		Num3 int `json:"3"`
		Num4 int `json:"4"`
		Num5 int `json:"5"`
	} `json:"valuationDistributionPercent"`
	NmValuationDistribution []struct {
		Nm                    int `json:"nm"`
		ValuationDistribution struct {
			Num1 int `json:"1"`
			Num2 int `json:"2"`
			Num3 int `json:"3"`
			Num4 int `json:"4"`
			Num5 int `json:"5"`
		} `json:"valuationDistribution"`
		ValuationDistributionPercent struct {
			Num1 int `json:"1"`
			Num2 int `json:"2"`
			Num3 int `json:"3"`
			Num4 int `json:"4"`
			Num5 int `json:"5"`
		} `json:"valuationDistributionPercent"`
	} `json:"nmValuationDistribution"`
	MatchingSizePercentages   any `json:"matchingSizePercentages"`
	NmMatchingSizePercentages any `json:"nmMatchingSizePercentages"`
	FeedbackCount             int `json:"feedbackCount"`
	FeedbackCountWithPhoto    int `json:"feedbackCountWithPhoto"`
	FeedbackCountWithText     int `json:"feedbackCountWithText"`
	FeedbackCountWithVideo    int `json:"feedbackCountWithVideo"`
	Feedbacks                 []struct {
		ID            string `json:"id"`
		GlobalUserID  string `json:"globalUserId"`
		WbUserID      int    `json:"wbUserId"`
		WbUserDetails struct {
			Country  string `json:"country"`
			Name     string `json:"name"`
			HasPhoto bool   `json:"hasPhoto"`
		} `json:"wbUserDetails"`
		NmID                int       `json:"nmId"`
		Text                string    `json:"text"`
		Pros                string    `json:"pros"`
		Cons                string    `json:"cons"`
		MatchingSize        string    `json:"matchingSize"`
		MatchingPhoto       string    `json:"matchingPhoto"`
		MatchingDescription string    `json:"matchingDescription"`
		ProductValuation    int       `json:"productValuation"`
		Color               string    `json:"color"`
		Size                string    `json:"size"`
		CreatedDate         time.Time `json:"createdDate"`
		UpdatedDate         time.Time `json:"updatedDate"`
		Answer              struct {
			Text         string    `json:"text"`
			State        string    `json:"state"`
			LastUpdate   any       `json:"lastUpdate"`
			CreateDate   time.Time `json:"createDate"`
			RejectReason any       `json:"rejectReason"`
			Metadata     struct {
				EditText         string `json:"editText"`
				EditRejectReason int    `json:"editRejectReason"`
			} `json:"metadata"`
		} `json:"answer"`
		FeedbackHelpfulness any `json:"feedbackHelpfulness"`
		Video               any `json:"video"`
		Votes               struct {
			Pluses  int `json:"pluses"`
			Minuses int `json:"minuses"`
		} `json:"votes"`
		Rank     float64 `json:"rank"`
		StatusID int     `json:"statusId,omitempty"`
		Reasons  struct {
			Good []any `json:"good"`
			Bad  []any `json:"bad"`
		} `json:"reasons"`
		ExcludedFromRating struct {
			IsExcluded bool  `json:"isExcluded"`
			Reasons    []any `json:"reasons"`
		} `json:"excludedFromRating"`
		Tags []struct {
			ID     int    `json:"id"`
			Status string `json:"status"`
		} `json:"tags,omitempty"`
		Bables []string `json:"bables,omitempty"`
		Photos []struct {
			ID        int    `json:"id"`
			Key       string `json:"key"`
			IsBlurred bool   `json:"isBlurred"`
			IsReady   bool   `json:"isReady"`
		} `json:"photos,omitempty"`
		Photo            []int  `json:"photo,omitempty"`
		ParentFeedbackID string `json:"parentFeedbackId,omitempty"`
		ChildFeedbackID  string `json:"childFeedbackId,omitempty"`
	} `json:"feedbacks"`
}

type FeedbackResponse struct {
	Text   string `csv:"text"`
	Pros   string `csv:"pros"`
	Cons   string `csv:"cons"`
	Rating int    `csv:"rating"`
	Date   string `csv:"date"`
}
