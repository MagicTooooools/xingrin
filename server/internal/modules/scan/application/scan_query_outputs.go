package application

type ScanTargetName struct {
	Name string
	Type string
}

type ScanStatistics struct {
	Total           int64
	Running         int64
	Completed       int64
	Failed          int64
	TotalVulns      int64
	TotalSubdomains int64
	TotalEndpoints  int64
	TotalWebsites   int64
	TotalAssets     int64
}
