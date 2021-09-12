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
	ClientID       string `env:"client_id,required"`
	ClientSecret   string `env:"client_secret,required"`
	Datacenter     string `env:"datacenter,required"`
	EntitlementID  string `env:"entitlement_id,required"`
	ReleaseID      string `env:"release_id,required"`
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
		cfg       Config         = Config{}
		envGetter env.Repository = env.NewRepository()
		client    *fod.Client    = nil
	)

	// Parse configuration from environment
	if err := stepconf.NewInputParser(envGetter).Parse(&cfg); err != nil {
		logger.Errorf("%v\n", err)
		os.Exit(1)
	}

	// Create a new FoD client with Client Credentials
	if client = fod.NewWithClientCredentials(cfg.ClientID, cfg.ClientSecret, cfg.Datacenter); client == nil {
		logger.Errorf("invalid parameters")
		stepconf.Print(cfg)
		os.Exit(1)
	}

	// Parse config inputs to mobie scan parameters
	params := fod.MobileScanParams{
		ReleaseId:     cfg.ReleaseID,
		FrameworkType: cfg.FrameworkType,
		EntitlementId: cfg.EntitlementID,
		PlatformType:  cfg.PlatformType,
		FilePath:      cfg.FilePath,
	}

	// Parse config assessment type to mobile scan parameters
	switch cfg.AssessmentType {
	case "Mobile Assessment (Single Scan)":
		params.AssessmentTypeId = fod.MOBILE_ASSESSMENT
		params.EntitlementFrequencyType = fod.SINGLE_SCAN
	case "Mobile+ Assessment (Single Scan)":
		params.AssessmentTypeId = fod.MOBILE_PLUS_ASSESSMENT
		params.EntitlementFrequencyType = fod.SINGLE_SCAN
	case "Mobile Assessment (Subscription)":
		params.AssessmentTypeId = fod.MOBILE_ASSESSMENT
		params.EntitlementFrequencyType = fod.SUBSCRIPTION
	case "Mobile+ Assessment (Subscription)":
		params.AssessmentTypeId = fod.MOBILE_PLUS_ASSESSMENT
		params.EntitlementFrequencyType = fod.SUBSCRIPTION
	}

	stepconf.Print(params)
	fmt.Println()

	var (
		scanID string
		err    error
	)

	// submit mobile scan
	if scanID, err = client.SetDebug(true).StartMobileScan(params); err == nil {
		fmt.Println()
		logger.Infof("mobile scan submitted successfully, Scan ID: %s\n", scanID)
	} else {
		logger.Errorf("err: %v\n", err)
		os.Exit(1)
	}
}
