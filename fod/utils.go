package fod

import (
	"fmt"
	"strings"

	"github.com/fatih/structs"
)

type Error struct {
	ErrorCode int    `json:"errorCode"`
	Message   string `json:"message"`
}

type Errors struct {
	Errors []Error `json:"errors"`
}

// Get base FoD RESTFul endpoint for datacenter
func getBaseUrl(datacenter string) string {

	url := ""

	if dc := strings.ToLower(strings.TrimSpace(datacenter)); dc != "" {

		switch dc {
		case DATACENTER_AMS, DATACENTER_APAC, DATACENTER_EMEA:
			url = "https://api." + dc + ".fortify.com"
		case DATACENTER_FED:
			url = "https://api." + dc + ".fortifygov.com"
		default:
			url = ""
		}
	}

	return url
}

// Check for valid datacenter
func isValidDatacenter(dc string) bool {
	return getBaseUrl(dc) != ""
}

// Get error from http status code and additional error message
func getErrorFromStatusCode(statusCode int, i interface{}) error {

	var (
		msg string = ""
	)

	switch statusCode {
	case 200, 202:
		return nil
	case 400: // bad request
		msg = fmt.Sprintf("[%d] bad request", statusCode)
	case 401: // unauthorized
		msg = fmt.Sprintf("[%d] unauthorized", statusCode)
	case 403: // forbidden
		msg = fmt.Sprintf("[%d] forbidden", statusCode)
	case 404: // not found
		msg = fmt.Sprintf("[%d] not found", statusCode)
	case 422: // not found
		msg = fmt.Sprintf("[%d] unprocessable entity", statusCode)
	case 429: // too many request
		msg = fmt.Sprintf("[%d] too many request", statusCode)
	case 500: // internal server error
		msg = fmt.Sprintf("[%d] internal server error", statusCode)
	default:
		msg = fmt.Sprintf("[%d] unknown error", statusCode)
	}

	if i != nil {
		return fmt.Errorf("%s, %+v", msg, i)
	}

	return fmt.Errorf(msg)
}

// Check for valid scope
// func isValidScope(s string) bool {
// 	switch s {
// 	case SCOPE_API_TENANT, SCOPE_START_SCANS, SCOPE_MANAGE_APPS, SCOPE_VIEW_APPS, SCOPE_MANAGE_ISSUES, SCOPE_VIEW_ISSUES, SCOPE_MANAGE_REPORTS, SCOPE_VIEW_REPORTS, SCOPE_MANAGE_USERS, SCOPE_VIEW_USERS, SCOPE_MANAGE_NOTIFICATIONS, SCOPE_VIEW_TENANT_DATA:
// 		return true
// 	}

// 	return false
// }

// Check for valid scan type
func isValidScanType(st string) bool {
	switch st {
	case SCAN_TYPE_STATIC, SCAN_TYPE_DYNAMIC, SCAN_TYPE_MOBILE:
		return true
	}
	return false
}

// Check for valid assessment type id
func isValidAssessmentType(at string) bool {
	switch at {
	case MOBILE_ASSESSMENT, MOBILE_PLUS_ASSESSMENT:
		return true
	}
	return false
}

// Check for valid subscription type
func isValidSubscriptionType(t string) bool {
	switch t {
	case SINGLE_SCAN, SUBSCRIPTION:
		return true
	}

	return false
}

// Check for valid mobile platform type
func isValidMobilePlatform(t string) bool {
	switch t {
	case MMOBILE_PLATFORM_PHONE, MMOBILE_PLATFORM_TABLET, MMOBILE_PLATFORM_BOTH:
		return true
	}

	return false
}

// Check for valid mobile framework
func isValidMobileFramework(t string) bool {
	switch t {
	case MOBILE_FRAMEWORK_IOS, MOBILE_FRAMEWORK_ANDROID:
		return true
	}

	return false
}

// Generate a string from an array of scope string
func scopeToString(scopeArr []string) (scope string) {

	for _, s := range scopeArr {
		scope += s + " "
	}

	return scope
}

// Convert structs to map[string]string format
func structToFormData(s interface{}) map[string]string {

	var (
		out = map[string]string{}
		m   = structs.Map(s)
	)

	for k, v := range m {
		out[k] = fmt.Sprint(v)
	}

	return out
}

// Split data into fragments, each fragment is pass back through callback function.
// fragNo indicates the id of the fragment extracted. -1 indicates the last fragment
// Array of fragments is returned.
func splitDataToFragments(indaata []byte, block_size int, callback func(fragNo int, data []byte) error) [][]byte {

	if indaata == nil || block_size <= 0 {
		return nil
	}

	var (
		indata_len = len(indaata)
		outdata    [][]byte
		fragNo     int   = 0
		rem        int   = 0
		numFrag    int   = 0
		offset     int   = 0
		err        error = nil
	)

	// Compute the quotient amd remainder
	numFrag, rem = indata_len/block_size, indata_len%block_size

	for fragNo = 0; fragNo < numFrag; fragNo++ {

		outdata = append(outdata, indaata[offset:offset+block_size])

		// Move to the next fragment
		offset += block_size

		if callback != nil {

			// check for last fragment
			if rem == 0 {
				err = callback(-1, outdata[fragNo])
			} else {
				err = callback(fragNo, outdata[fragNo])
			}

			if err != nil {
				return nil
			}
		}
	}

	// Remainder as the last fragment
	if rem > 0 {
		outdata = append(outdata, indaata[offset:offset+rem])

		if callback != nil {
			err = callback(-1, outdata[fragNo])

			if err != nil {
				return nil
			}
		}
	}

	return outdata
}
