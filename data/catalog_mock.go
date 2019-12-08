package data

/*
 Copyright 2019 Crunchy Data Solutions, Inc.
 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at
      http://www.apache.org/licenses/LICENSE-2.0
 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

import (
	"bytes"
	"fmt"
	"strconv"
	"text/template"
)

type catalogMock struct {
	layers    []*Layer
	layerData map[string][]string
}

var instance Catalog

// CatMockInstance tbd
func CatMockInstance() Catalog {
	// TODO: make a singleton
	instance = newCatalogMock()
	return instance
}

func newCatalogMock() Catalog {
	var layers []*Layer
	layers = append(layers, &Layer{
		ID:          "mock_a",
		Title:       "Mock A",
		Description: "This dataset contains mock data about A",
		Extent:      Extent{Minx: 0, Miny: 0, Maxx: 80, Maxy: 90},
		Srid:        999,
	})
	layers = append(layers, &Layer{
		ID:          "mock_b",
		Title:       "Mock B",
		Description: "This dataset contains mock data about B (100 points)",
		Extent:      Extent{Minx: -130, Miny: 40, Maxx: -120, Maxy: 60},
		Srid:        999,
	})

	layerData := map[string][]string{}
	layerData["mock_a"] = featuresMock
	//layerData["mock_b"] = features
	layerData["mock_b"] = makeFeatures(10, 10)

	catMock := catalogMock{
		layers:    layers,
		layerData: layerData,
	}

	return &catMock
}

func (cat *catalogMock) Layers() ([]*Layer, error) {
	return cat.layers, nil
}

func (cat *catalogMock) LayerByName(name string) (*Layer, error) {
	for _, lyr := range cat.layers {
		if lyr.ID == name {
			return lyr, nil
		}
	}
	// not found
	return nil, fmt.Errorf(errMsgBadLayerName, name)
}

func (cat *catalogMock) LayerFeatures(name string) ([]string, error) {
	features, ok := cat.layerData[name]
	if !ok {
		return []string{}, fmt.Errorf(errMsgBadLayerName, name)
	}
	//fmt.Println("LayerFeatures: " + name)
	//fmt.Println(layerData)
	return features, nil
}

func (cat *catalogMock) LayerFeature(name string, id string) (string, error) {
	features, ok := cat.layerData[name]
	if !ok {
		return "", fmt.Errorf(errMsgBadLayerName, name)
	}
	index, err := strconv.Atoi(id)
	if err != nil {
		return "", fmt.Errorf(errMsgNoBadFeatureID, id)
	}

	//fmt.Println("LayerFeatures: " + name)
	//fmt.Println(layerData)
	return features[index], nil
}

var featuresMock = []string{
	`{ "type": "Feature", "id": 1,  "geometry": {"type": "Point","coordinates": [  -75,	  45]  },
	  "properties": { "value": "89.9"  } }`,
	`{ "type": "Feature", "id": 2,  "geometry": {"type": "Point","coordinates": [  -75,	  40]  },
	  "properties": { "value": "89.9"  } }`,
	`{ "type": "Feature", "id": 3,  "geometry": {"type": "Point","coordinates": [  -75,	  35]  },
	  "properties": { "value": "89.9"  } }`,
}

type featurePointMock struct {
	ID  int
	X   float64
	Y   float64
	Val string
}

var templateFeaturePoint = `{ "type": "Feature", "id": {{ .ID }},
"geometry": {"type": "Point","coordinates": [  {{ .X }}, {{ .Y }} ]  },
"properties": { "value": "{{ .Val }}"  } }`

func makeFeatures(nx int, ny int) []string {
	tmpl, err := template.New("feature").Parse(templateFeaturePoint)
	if err != nil {
		panic(err)
	}
	n := nx * ny
	features := make([]string, n)
	var tempOut bytes.Buffer
	index := 0
	for ix := 0; ix < nx; ix++ {
		for iy := 0; iy < ny; iy++ {
			x := -75 + 0.01*float64(ix)
			y := 45 + 0.01*float64(iy)
			val := fmt.Sprintf("data value %v", index)
			feat := featurePointMock{index, x, y, val}
			tempOut.Reset()
			tmpl.Execute(&tempOut, feat)
			features[index] = tempOut.String()
			//fmt.Println(features[index])

			index++
		}
	}
	return features
}