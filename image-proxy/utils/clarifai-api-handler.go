package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
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
				"user_id": %s,
				"app_id": "%s"
			},
			"inputs": [
				{
					"data": {
						"image": {
							"url": "https://samples.clarifai.com/metro-north.jpg"
						}
					}
				}
			]
		}
	`, env.ClarifaiUserID, env.ClarifaiAppID))

	req, err := http.NewRequest("POST", URL, requestOptions)
	if err != nil {
		panic(err)
	}
	req.Header.Set("user-agent", "golang application")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Key 0f90fc9066624052ae143cd96bf5e9f0")
	response, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body) // response body is []byte
	return string(body)
}

// get the coordinates the detected face
// takes clarifai api response as param
// return the coordinates as a box(BoundinBox)
func GetBoundingBox(jsonData string) BoundingBox {
	data := clarifaiResponse{}
	json.Unmarshal([]byte(jsonData), &data)
	box := data.Outputs[0].Data.Regions[0].RegionInfo.BoundingBox

	return BoundingBox{
		TopRow:    box.TopRow,
		RightCol:  box.RightCol,
		BottomRow: box.BottomRow,
		LeftCol:   box.LeftCol,
	}
}
