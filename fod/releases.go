package fod

import (
	"errors"
	"fmt"

	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

type AssessmentType struct {
	AssessmentTypeId               int    `json:"assessmentTypeId"`
	Name                           string `json:"name"`
	ScanType                       string `json:"scanType"`
	ScanTypeId                     int    `json:"scanTypeId"`
	EntitlementId                  int    `json:"entitlementId"`
	FrequencyType                  string `json:"frequencyType"`
	FrequencyTypeId                int    `json:"frequencyTypeId"`
	Units                          int    `json:"units"`
	UnitsAvailable                 int    `json:"unitsAvailable"`
	SubscriptionEndDate            string `json:"subscriptionEndDate"`
	IsRemediation                  bool   `json:"isRemediation"`
	RemediationScansAvailable      int    `json:"remediationScansAvailable"`
	IsBundledAssessment            bool   `json:"isBundledAssessment"`
	ParentAssessmentTypeId         int    `json:"parentAssessmentTypeId"`
	ParentAssessmentTypeName       string `json:"parentAssessmentTypeName"`
	ParentAssessmentTypeScanType   string `json:"parentAssessmentTypeScanType"`
	ParentAssessmentTypeScanTypeId int    `json:"parentAssessmentTypeScanTypeId"`
	EntitlementDescription         string `json:"entitlementDescription"`
}

// Get list of assessment types for a release
func (c *Client) GetAssessmentTypes(releaseId int, scanType string) ([]AssessmentType, error) {

	if !isValidScanType(scanType) {
		return nil, errors.New("invalid scan type")
	}

	log.WithFields(log.Fields{
		"releaseId": releaseId,
		"scanType":  scanType,
	}).Info("get assessment types")

	// "need to know" basis
	c.scope = []string{SCOPE_START_SCANS}

	type Items struct {
		Items []AssessmentType `json:"items"`
	}

	var (
		authData  *AuthData
		resp      *resty.Response = nil
		err       error           = nil
		items     *Items          = &Items{}
		fodErrors *Errors         = &Errors{}
	)

	// Perform authenitcation
	if authData, err = c.authenticate(); err == nil {

		// Construct full api endpoint url with release id
		url := c.baseUrl + fmt.Sprintf(get_api_V3_releases_assessment_types, releaseId)

		log.WithFields(log.Fields{
			"url": url,
		}).Info("get assessment types")

		// Send request
		if resp, err = c.webClient.R().
			SetAuthScheme(authData.TokenType).
			SetAuthToken(authData.AccessToken).
			SetQueryParam("scanType", scanType).
			SetResult(&items).
			SetError(&fodErrors).
			Get(url); err == nil {

			if err = getErrorFromStatusCode(resp.StatusCode(), fodErrors); err != nil {
				return nil, err
			}

		} else {
			return nil, err
		}
	}

	return items.Items, nil
}
