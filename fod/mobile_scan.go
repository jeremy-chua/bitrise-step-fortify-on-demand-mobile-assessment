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

type ScanId struct {
	Value int `json:"scanId"`
}

// Submit a Mobile or Mobile+ scan.
func (c *Client) StartMobileScan(
	releaseId int,
	assessmentTypeId int,
	frameworkType string,
	entitlementId int,
	frequencyType string,
	platformType string,
	filePath string,
	isRemediationScan bool,
) (string, error) {

	log.WithFields(log.Fields{
		"releaseId":                releaseId,
		"assessmentTypeId":         assessmentTypeId,
		"frameworkType":            frameworkType,
		"entitlementId":            entitlementId,
		"entitlementFrequencyType": frequencyType,
		"platformType":             platformType,
		"filePath":                 filePath,
		"isRemediationScan":        isRemediationScan,
	}).Info("starting mobile scan")

	if !isValidMobileFramework(frameworkType) {
		return "", errors.New("invalid argument: frameworkType")
	}
	if !isValidMobilePlatform(platformType) {
		return "", errors.New("invalid argument: platformType")
	}
	if !isValidSubscriptionType(frequencyType) {
		return "", errors.New("invalid argument: entitlementFrequencyType")
	}

	var (
		authData  *AuthData       = nil
		resp      *resty.Response = nil
		err       error           = nil
		scanId    *ScanId         = &ScanId{}
		scanError *Errors         = &Errors{}
		retScanId string          = ""
		data      []byte
		loc       *time.Location
	)

	// Read content of file
	if data, err = os.ReadFile(filePath); err != nil {
		return "", err
	}

	// set scope "need to know" basis
	c.scope = []string{SCOPE_START_SCANS}

	// Perform authenitcation
	if authData, err = c.authenticate(); err == nil {

		var (
			querParams map[string]string = map[string]string{}
			offset     int               = 0
			logMsg     string            = ""
		)

		// Construct full api endpoint url with release id
		url := c.baseUrl + fmt.Sprintf(post_api_V3_start_mobile_scan, releaseId)

		querParams["releaseId"] = strconv.Itoa(releaseId)
		querParams["assessmentTypeId"] = strconv.Itoa(assessmentTypeId)
		querParams["frameworkType"] = frameworkType
		querParams["entitlementId"] = strconv.Itoa(entitlementId)
		querParams["entitlementFrequencyType"] = frequencyType
		querParams["platformType"] = platformType
		querParams["isRemediationScan"] = strconv.FormatBool(isRemediationScan)

		// 	// Use current date time base on GMT Standard Time (UTC+00:00) timezone
		loc, _ = time.LoadLocation("UTC")
		querParams["startDate"] = time.Now().In(loc).Format("2006-01-02 15:04")
		querParams["timeZone"] = "GMT Standard Time"

		if isRemediationScan {
			logMsg = "submitting remediation scan..."
		} else {
			logMsg = "submitting scan..."
		}

		splitDataToFragments(data, block_size, func(fragNo int, data []byte) error {

			var err error = nil

			querParams["fragNo"] = strconv.Itoa(fragNo)
			querParams["offset"] = strconv.Itoa(offset)

			log.WithFields(log.Fields{
				"fragNo": fragNo,
				"offset": offset,
			}).Info(logMsg)

			// Send request
			if resp, err = c.webClient.R().
				SetAuthScheme(authData.TokenType).
				SetAuthToken(authData.AccessToken).
				SetResult(&scanId).
				SetError(scanError).
				SetQueryParams(querParams).
				SetBody(data).
				Post(url); err == nil {

				// Parse return error
				switch resp.StatusCode() {
				case 200:
					retScanId = strconv.Itoa(scanId.Value)
				case 202:
					log.WithFields(log.Fields{
						"fragNo": fragNo,
						"offset": offset,
					}).Info("accepted")
				default:
					err = getErrorFromStatusCode(resp.StatusCode(), scanError)
				}
			}

			offset += block_size
			time.Sleep(100 * time.Millisecond)

			return err
		})
	}

	return retScanId, err
}
