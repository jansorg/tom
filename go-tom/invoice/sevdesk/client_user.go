package sevdesk

import (
	"fmt"
	"time"
)

type User struct {
	ID                     string      `json:"id"`
	ObjectName             string      `json:"objectName"`
	AdditionalInformation  interface{} `json:"additionalInformation"`
	Create                 time.Time   `json:"create"`
	Update                 time.Time   `json:"update"`
	Fullname               string      `json:"fullname"`
	FirstName              string      `json:"firstName"`
	LastName               string      `json:"lastName"`
	Username               string      `json:"username"`
	Status                 string      `json:"status"`
	Email                  string      `json:"email"`
	Gender                 string      `json:"gender"`
	Role                   string      `json:"role"`
	MemberCode             string      `json:"memberCode"`
	SevClient              IDWithType  `json:"sevClient"`
	LastLogin              time.Time   `json:"lastLogin"`
	LastLoginIP            string      `json:"lastLoginIp"`
	WelcomeScreenSeen      string      `json:"welcomeScreenSeen"`
	SMTPName               string      `json:"smtpName"`
	SMTPMail               string      `json:"smtpMail"`
	SMTPUser               interface{} `json:"smtpUser"`
	SMTPPort               interface{} `json:"smtpPort"`
	SMTPSsl                interface{} `json:"smtpSsl"`
	SMTPHost               interface{} `json:"smtpHost"`
	LanguageCode           string      `json:"languageCode"`
	TwoFactorAuth          string      `json:"twoFactorAuth"`
	ForcePasswordChange    string      `json:"forcePasswordChange"`
	ClientOwner            string      `json:"clientOwner"`
	DefaultReceiveMailCopy string      `json:"defaultReceiveMailCopy"`
	HideMapsDirections     string      `json:"hideMapsDirections"`
	StartDate              interface{} `json:"startDate"`
	EndDate                interface{} `json:"endDate"`
	LastPasswordChange     time.Time   `json:"lastPasswordChange"`
}

func (api *Client) GetCurrentUser() (*User, error) {
	resp, err := api.doRequest("GET", "/SevUser", nil, nil)
	if err != nil {
		return nil, err
	}

	usersDef := []User{}
	err = api.unwrapJSONResponse(resp, &usersDef)
	if err != nil {
		return nil, err
	}

	if len(usersDef) != 1 {
		return nil, fmt.Errorf("unexpected number of sevdesk users: %d", len(usersDef))
	}
	return &usersDef[0], nil
}
