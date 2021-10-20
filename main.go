package main

import (
	"bitrise-step-fortify-on-demand-mobile-assessment/fod"
	"fmt"
	"os"

	"github.com/bitrise-io/go-steputils/stepconf"
	"github.com/bitrise-io/go-utils/env"
	"github.com/bitrise-io/go-utils/log"
)

// Config ...
type Config struct {
	ClientId       string `env:"client_id,required"`
	ClientSecret   string `env:"client_secret,required"`
	Datacenter     string `env:"datacenter,required"`
	EntitlementId  int    `env:"entitlement_id,required"`
	ReleaseId      int    `env:"release_id,required"`
	AssessmentType string `env:"assessment_type,required"`
	FrameworkType  string `env:"framework_type,required"`
	PlatformType   string `env:"platform_type,required"`
	FilePath       string `env:"file_path,required"`
}

func main() {

	// Initialize new logger
	logger := log.NewLogger()
	logger.Infof("Start Fortify on Demand mobile assessment step")

	var (
		cfg               Config         = Config{}
		envGetter         env.Repository = env.NewRepository()
		client            *fod.Client    = nil
		scanId            string         = ""
		err               error          = nil
		assessmentTypeId  int            = -1
		frequencyType     string         = ""
		isRemediationScan bool           = false
		assessmentTypes   []fod.AssessmentType
	)

	// Parse configuration from environment
	if err = stepconf.NewInputParser(envGetter).Parse(&cfg); err != nil {
		logger.Errorf("%v\n", err)
		os.Exit(1)
	}

	stepconf.Print(cfg)

	// Create a new FoD client with Client Credentials
	if client = fod.NewWithClientCredentials(cfg.ClientId, cfg.ClientSecret, cfg.Datacenter); client == nil {
		logger.Errorf("invalid parameters")
		stepconf.Print(cfg)
		os.Exit(1)
	}

	// Get assessment types for the release
	if assessmentTypes, err = client.SetDebug(true).GetAssessmentTypes(cfg.ReleaseId, fod.SCAN_TYPE_MOBILE); err != nil {
		logger.Errorf("err: %v\n", err)
		os.Exit(1)
	}

	for _, at := range assessmentTypes {
		if cfg.AssessmentType == fmt.Sprintf("%s (%s)", at.Name, "Single Scan") && at.FrequencyType == fod.SINGLE_SCAN {
			assessmentTypeId = at.AssessmentTypeId
			frequencyType = at.FrequencyType

			// Check if remediation scan is available for single scan
			if at.IsRemediation {
				isRemediationScan = true
			}

		} else if cfg.AssessmentType == fmt.Sprintf("%s (%s)", at.Name, "Subscription") && at.FrequencyType == fod.SUBSCRIPTION {
			assessmentTypeId = at.AssessmentTypeId
			frequencyType = at.FrequencyType
		}
	}

	// Check if assessment type exist in entitlement
	if assessmentTypeId < 0 {
		logger.Errorf("assessment type not found in entitlement\n")
		os.Exit(1)
	}

	// submit mobile scan
	if scanId, err = client.SetDebug(true).StartMobileScan(
		cfg.ReleaseId,
		assessmentTypeId,
		cfg.FrameworkType,
		cfg.EntitlementId,
		frequencyType,
		cfg.PlatformType,
		cfg.FilePath,
		isRemediationScan,
	); err == nil {
		logger.Infof("mobile scan submitted successfully, Scan ID: %s\n", scanId)
	} else {
		logger.Errorf("err: %v\n", err)
		os.Exit(1)
	}

}
