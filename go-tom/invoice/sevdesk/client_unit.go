package sevdesk

import "fmt"

type Unit struct {
	ID         string `json:"id"`
	ObjectName string `json:"objectName"`
	Name       string `json:"name"`
	// AdditionalInformation interface{} `json:"additionalInformation"`
	// Create                interface{} `json:"create"`
	// TranslationCode       string      `json:"translationCode"`
	// UnitySystem           string      `json:"unitySystem"`
}

// locates the unit in the cached data
func (api *Client) FindUnit(name string) (*Unit, error) {
	if api.units == nil {
		return nil, fmt.Errorf("no units loaded")
	}

	for _, def := range *api.units {
		if def.Name == name {
			return &def, nil
		}
	}
	return nil, fmt.Errorf("no unit found for name %s", name)
}

func (api *Client) GetUnits() ([]Unit, error) {
	resp, err := api.doRequest("GET", "/Unity", nil, nil)
	if err != nil {
		return nil, err
	}

	var units []Unit
	err = api.unwrapJSONResponse(resp, &units)
	if err != nil {
		return nil, err
	}
	return units, nil
}
