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
type boxClarifaiResponse struct {
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

type personInfoClarifaiResponse struct {
	Status struct {
		Code        int    `json:"code"`
		Description string `json:"description"`
		ReqID       string `json:"req_id"`
	} `json:"status"`
	Outputs []struct {
		Data struct {
			Concepts []struct {
				ID    string  `json:"id"`
				Name  string  `json:"name"`
				Value float64 `json:"value"`
			} `json:"concepts"`
		} `json:"data"`
	} `json:"outputs"`
}

var SURE float64 = 0.8
var UNCERTAIN float64 = 0.6

// make a call to clarifai api in order to make prediction on an image
// takes api url as param
// return response as json string
func Fetch(URL string, modelType string) string {
	modelID, modelVersionID := "", ""

	switch modelType {
	case "face":
		modelID = "f76196b43bbd45c99b4f3cd8e8b40a8a"
		modelVersionID = "45fb9a671625463fa646c3523a3087d5"
	case "celebrity":
		modelID = "cfbb105cb8f54907bb8d553d68d9fe20"
		modelVersionID = "0676ebddd5d6413ebdaa101570295a39"
	case "gender":
		modelID = "af40a692dfe6040f23ca656f4e144fc2"
		modelVersionID = "ff83d5baac004aafbe6b372ffa6f8227"
	default:
		return ""
	}

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

	clarifaiApiURL := "https://api.clarifai.com/v2/models/" + modelID + "/versions/" + modelVersionID + "/outputs?"
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
	data := boxClarifaiResponse{}
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

// get the name and confidence level of the prediction
// takes clarifai api response as param
// return name and cofidence level
func GetCelebrity(jsonData string) (string, int, error) {
	data := personInfoClarifaiResponse{}
	json.Unmarshal([]byte(jsonData), &data)

	if len(data.Outputs) == 0 || len(data.Outputs[0].Data.Concepts) == 0 {
		return "", -1, errors.New("Empty data")
	}

	concepts := data.Outputs[0].Data.Concepts

	max := concepts[0].Value
	index := 0
	for i := 1; i < len(concepts); i++ {
		if max < concepts[i].Value {
			max = concepts[i].Value
			index = i
		}
	}

	name := concepts[index].Name
	confidenceLevel := 0
	switch value := concepts[index].Value; {
	case value >= SURE:
		confidenceLevel = 2
	case value > UNCERTAIN:
		confidenceLevel = 1
	default:
		name = ""
	}

	return name, confidenceLevel, nil
}

// get the gender of the face supplied
// takes clarifai api response as param
// return gender
func GetGender(jsonData string) (bool, error) {
	data := personInfoClarifaiResponse{}
	json.Unmarshal([]byte(jsonData), &data)

	if len(data.Outputs) == 0 || len(data.Outputs[0].Data.Concepts) == 0 {
		return false, errors.New("Empty data")
	}

	concepts := data.Outputs[0].Data.Concepts

	gender := false
	if concepts[0].Value > concepts[1].Value {
		if concepts[0].Name[0] == 'M' {
			gender = true
		}
	}

	return gender, nil
}
