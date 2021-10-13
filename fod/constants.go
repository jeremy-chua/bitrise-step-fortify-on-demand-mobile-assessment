package fod

const (
	DATACENTER_AMS  = "ams"
	DATACENTER_APAC = "apac"
	DATACENTER_EMEA = "emea"
	DATACENTER_FED  = "fed"

	SCOPE_API_TENANT           = "api-tenant"           // Grants access to all endpoints
	SCOPE_START_SCANS          = "start-scans"          // Configure and start static, dynamic, and mobile scans; import static and dynamic scans
	SCOPE_MANAGE_APPS          = "manage-apps"          // // Manage applications
	SCOPE_VIEW_APPS            = "view-apps"            // View applications
	SCOPE_MANAGE_ISSUES        = "manage-issues"        // Manage issues
	SCOPE_VIEW_ISSUES          = "view-issues"          // View issues
	SCOPE_MANAGE_REPORTS       = "manage-reports"       // Manage reports
	SCOPE_VIEW_REPORTS         = "view-reports"         // View reports
	SCOPE_MANAGE_USERS         = "manage-users"         // Manage users
	SCOPE_VIEW_USERS           = "view-users"           // View users
	SCOPE_MANAGE_NOTIFICATIONS = "manage-notifications" // Manage notifications
	SCOPE_VIEW_TENANT_DATA     = "view-tenant-data"     // View data at the tenant level

	MOBILE_ASSESSMENT      = 271 // Mobile Assessment
	MOBILE_PLUS_ASSESSMENT = 272 // Mobile+ Assessment

	SINGLE_SCAN  = "SingleScan"   // Single Scan
	SUBSCRIPTION = "Subscription" // Subscription

	MMOBILE_PLATFORM_PHONE  = "Phone"  // Phone mobile platform
	MMOBILE_PLATFORM_TABLET = "Tablet" // Tablet mobile platform
	MMOBILE_PLATFORM_BOTH   = "Both"   // Mobile Platform for both phone and tablet

	MOBILE_FRAMEWORK_IOS     = "iOS"     // iOS mobile framework
	MOBILE_FRAMEWORK_ANDROID = "Android" // Android mobile framework
)

const (

	// Block size of 1MB (1024 x 1024) for data being uploaded
	block_size = 1048576

	grant_type_client_credentials = "client_credentials"
	grant_type_password           = "password"

	post_api_authenticate                = "/oauth/token"
	post_api_V3_start_mobile_scan        = "/api/v3/releases/%s/mobile-scans/start-scan"
	get_api_V3_releases_assessment_types = "/api/v3/releases/%s/assessment-types"
)
