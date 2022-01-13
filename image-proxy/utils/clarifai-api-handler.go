package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	ib "github.com/neeejm/image-box/utils"
)

// part of the structure(json) of the clarifai api response
type clarifaiResponse struct {
	Status struct {
		Code        int    `json:"code"`
		Description string `json:"description"`
		ReqID       string `json:"req_id"`
	} `json:"status"`
	Outputs []struct {
		Data struct {
			Regions []struct {
				ID         string `json:"id"`
				RegionInfo struct {
					BoundingBox struct {
						TopRow    float64 `json:"top_row"`
						LeftCol   float64 `json:"left_col"`
						BottomRow float64 `json:"bottom_row"`
						RightCol  float64 `json:"right_col"`
					} `json:"bounding_box"`
				} `json:"region_info"`
			} `json:"regions"`
		} `json:"data"`
	} `json:"outputs"`
}

// make a call to clarifai api in order to make prediction on an image
// takes api url as param
// return response as json string
func Fetch(URL string) string {
	env := GetENV()

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	requestOptions := strings.NewReader(fmt.Sprintf(`
		{
			"user_app_id": {
				"user_id": "%s",
				"app_id": "%s"
			},
			"inputs": [
				{
					"data": {
						"image": {
							"url": "%s"
						}
					}
				}
			]
		}
	`, env.ClarifaiUserID, env.ClarifaiAppID, URL))

	clarifaiApiURL := "https://api.clarifai.com/v2/models/f76196b43bbd45c99b4f3cd8e8b40a8a/versions/45fb9a671625463fa646c3523a3087d5/outputs?"
	req, err := http.NewRequest("POST", clarifaiApiURL, requestOptions)
	if err != nil {
		panic(err)
	}
	req.Header.Set("user-agent", "golang application")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Key "+env.ClarifaiApiKey)
	response, err := client.Do(req)
	if err != nil {
		log.Println("me")
		panic(err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body) // response body is []byte
	return string(body)
}

// get the coordinates the detected face
// takes clarifai api response as param
// return the coordinates as a box(BoundinBox)
func GetBoundingBox(jsonData string) ([]ib.Box, error) {
	data := clarifaiResponse{}
	json.Unmarshal([]byte(jsonData), &data)

	if len(data.Outputs) == 0 || len(data.Outputs[0].Data.Regions) == 0 {
		return nil, errors.New("Empty data")
	}

	regions := data.Outputs[0].Data.Regions
	box := []ib.Box{}

	for _, r := range regions {
		box = append(box, ib.Box{
			TopRow:    r.RegionInfo.BoundingBox.TopRow,
			RightCol:  r.RegionInfo.BoundingBox.RightCol,
			BottomRow: r.RegionInfo.BoundingBox.BottomRow,
			LeftCol:   r.RegionInfo.BoundingBox.LeftCol,
		})
	}

	return box, nil
}
