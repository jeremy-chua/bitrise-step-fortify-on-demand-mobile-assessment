package fod

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

type MobileScanParams struct {
	ReleaseId                string
	AssessmentTypeId         int    `structs:"assessmentTypeId"`
	FrameworkType            string `structs:"frameworkType"`
	EntitlementId            string `structs:"entitlementId"`
	EntitlementFrequencyType string `structs:"entitlementFrequencyType"`
	PlatformType             string `structs:"platformType"`
	// IsRemediationScan        bool   `structs:"isRemediationScan"`
	FilePath string
}

type ScanResponse struct {
	ScanId int `json:"scanId"`
}

func isValidMobileParams(params MobileScanParams) error {

	if !isValidAssessmentTypeId(params.AssessmentTypeId) {
		return errors.New("invalid argument: AssessmentTypeId")
	}
	if !isValidMobileFramework(params.FrameworkType) {
		return errors.New("invalid argument: FrameworkType")
	}
	if !isValidMobilePlatform(params.PlatformType) {
		return errors.New("invalid argument: PlatformType")
	}
	if !isValidSubscriptionType(params.EntitlementFrequencyType) {
		return errors.New("invalid argument: EntitlementFrequencyType")
	}

	return nil
}

// Check release if there are any remedication scan available for Single Scan
func (c *Client) isMobileRemediationScanAvailable(accessToken string, releaseId string, assessmentTypeId int) bool {

	// Construct full api endpoint url with release id
	url := c.baseUrl + fmt.Sprintf(get_api_V3_releases_assessment_types, releaseId)

	type RemediationScansAvailableResponse struct {
		AssessmentTypeId          int    `json:"assessmentTypeId"`
		FrequencyType             string `json:"frequencyType"`
		IsRemediation             bool   `json:"isRemediation"`
		RemediationScansAvailable int    `json:"remediationScansAvailable"`
	}

	type Items struct {
		Items []RemediationScansAvailableResponse `json:"items"`
	}
	// Query prammaters to check if remdication scan exist
	var (
		resp      *resty.Response = nil
		items                     = Items{}
		err       error           = nil
		fodErrors                 = Errors{}
		scanType                  = "Mobile"
		filters                   = "isRemediation:true"
	)

	// Send request
	if resp, err = c.webClient.R().
		SetAuthScheme("bearer").
		SetAuthToken(accessToken).
		SetResult(&items).
		SetError(&fodErrors).
		SetQueryParam("scanType", scanType).
		SetQueryParam("filters", filters).
		Get(url); err == nil {

		// Parse return error
		switch resp.StatusCode() {
		case 200:
			for _, v := range items.Items {
				// Check for any existing remediation scan for single scan
				if v.IsRemediation && v.RemediationScansAvailable > 0 && v.AssessmentTypeId == assessmentTypeId && v.FrequencyType == "SingleScan" {
					return true
				}
			}

			return false
		case 400:
			err = fmt.Errorf("%+v", fodErrors)
		case 401:
			err = errors.New("[401] unauthorized")
		case 403:
			err = errors.New("[403] forbidden")
		case 404:
			err = errors.New("[404] not found")
		case 429:
			err = errors.New("[429] too many request")
		case 500:
			err = errors.New("[429] internal server error")
		default:
			err = errors.New("unknown error")
		}

		if err != nil {
			log.WithField("error", err).Error("server returned error")
		}
	} else {
		if err != nil {
			log.WithField("error", err).Error("resty client returned error")
		}
	}

	return false
}

// Submit a Mobile or Mobile+ scan.
// Time and date are set to current time. Timezone is set as Standard Time (GMT+08).
// Submit as remediation scan for Single Scan type if any remeidation scan is available.
func (c *Client) StartMobileScan(params MobileScanParams) (string, error) {

	if err := isValidMobileParams(params); err != nil {
		log.Error("invalid params, unable to submit mobile scan, " + err.Error())
		return "", err
	}

	log.WithFields(log.Fields{
		"releaseId":                params.ReleaseId,
		"assessmentTypeId":         params.AssessmentTypeId,
		"entitlementFrequencyType": params.EntitlementFrequencyType,
	}).Info("starting mobile scan")

	// "need to know" basis
	c.scope = []string{SCOPE_START_SCANS}

	var (
		authData  *AuthData
		resp      *resty.Response = nil
		err       error           = nil
		scanResp  *ScanResponse   = &ScanResponse{}
		scanError *Errors         = &Errors{}
		retScanId string          = ""
		data      []byte
		loc       *time.Location
	)

	// Perform authenitcation
	if authData, err = c.authenticate(); err == nil {

		// Construct full api endpoint url with release id
		url := c.baseUrl + fmt.Sprintf(post_api_V3_start_mobile_scan, params.ReleaseId)

		// Unmarshal struct to map[string]string
		querParams := structToFormData(params)

		// Read content of file
		if data, err = os.ReadFile(params.FilePath); err != nil {
			log.Error("file error, " + err.Error())
			return "", err
		}

		// check if remedication scan applies
		var isRescan bool
		if isRescan = c.isMobileRemediationScanAvailable(authData.AccessToken, params.ReleaseId, params.AssessmentTypeId); isRescan {
			querParams["isRemediationScan"] = fmt.Sprint(isRescan)
		}

		// Use current date time base on GMT Standard Time (UTC+00:00) timezone
		loc, err = time.LoadLocation("UTC")
		querParams["startDate"] = time.Now().In(loc).Format("2006-01-02 15:04")
		querParams["timeZone"] = "GMT Standard Time"

		if log.GetLevel() == log.DebugLevel {
			fields := log.Fields{}

			for k, v := range querParams {
				fields[k] = v
			}

			fields["base_url"] = c.baseUrl
			fields["block_size"] = block_size
			fields["file_size"] = fmt.Sprintf("%d bytes", len(data))

			log.WithFields(fields).Debug("mobile scan params")
		}

		var offset = 0

		splitDataToFragments(data, block_size, func(fragNo int, data []byte) error {

			var err error = nil

			querParams["fragNo"] = strconv.Itoa(fragNo)
			querParams["offset"] = strconv.Itoa(offset)

			logMsg := ""

			if isRescan {
				logMsg = "submitting remediation scan..."
			} else {
				logMsg = "submitting scan..."
			}

			log.WithFields(log.Fields{
				"fragNo": fragNo,
				"offset": offset,
			}).Info(logMsg)

			// Send request
			if resp, err = c.webClient.R().
				SetAuthScheme(authData.TokenType).
				SetAuthToken(authData.AccessToken).
				SetResult(&scanResp).
				SetError(scanError).
				SetQueryParams(querParams).
				SetBody(data).
				Post(url); err == nil {

				// Parse return error
				switch resp.StatusCode() {
				case 200:
					retScanId = strconv.Itoa(scanResp.ScanId)
				case 202:
					log.WithFields(log.Fields{
						"fragNo": fragNo,
						"offset": offset,
					}).Info("accepted")
				case 401:
					err = errors.New("[401] unauthorized")
				case 403:
					err = errors.New("[403] forbidden")
				case 404:
					err = errors.New("[404] not found")
				case 400, 422, 500:
					err = fmt.Errorf("%+v", scanError)
				default:
					err = errors.New("unknown error")
				}

				if err != nil {
					log.WithField("error", err).Error("error " + logMsg)
				}
			}

			offset += block_size
			time.Sleep(100 * time.Millisecond)

			return err
		})
	}

	return retScanId, err
}
