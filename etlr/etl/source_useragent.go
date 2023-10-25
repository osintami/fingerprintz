package etl

type UserAgent struct {
}

func NewUserAgent() *UserAgent {
	return &UserAgent{}
}

func (x *UserAgent) Transform(job IETLJob) error {
	err := job.Tools().FileSystem.Copy(job.Info().inputFile, job.Info().snapshotFile)

	job.Tools().Items["browser.useragent.family"] = Item{
		Item:        "browser/useragent/family",
		Enabled:     true,
		GJSON:       "UserAgent.Family",
		Description: "User agent family.",
		Type:        "String"}

	job.Tools().Items["browser.useragent.major"] = Item{
		Item:        "browser/useragent/major",
		Enabled:     true,
		GJSON:       "UserAgent.Major",
		Description: "User agent major version",
		Type:        "String"}
	job.Tools().Items["browser.useragent.minor"] = Item{
		Item:        "browser/useragent/minor",
		Enabled:     true,
		GJSON:       "UserAgent.Minor",
		Description: "User agent minor version",
		Type:        "String"}
	job.Tools().Items["browser.useragent.patch"] = Item{
		Item:        "browser/useragent/patch",
		Enabled:     true,
		GJSON:       "UserAgent.Patch",
		Description: "User agent patch",
		Type:        "String"}
	job.Tools().Items["browser.useragent.osFamily"] = Item{
		Item:        "browser/useragent/osFamily",
		Enabled:     true,
		GJSON:       "Os.Family",
		Description: "OS family",
		Type:        "String"}
	job.Tools().Items["browser.useragent.osMajor"] = Item{
		Item:        "browser/useragent/osMajor",
		Enabled:     true,
		GJSON:       "Os.Major",
		Description: "OS major version",
		Type:        "String"}
	job.Tools().Items["browser.useragent.osMinor"] = Item{
		Item:        "browser/useragent/osMinor",
		Enabled:     true,
		GJSON:       "Os.Minor",
		Description: "OS minor version",
		Type:        "String"}
	job.Tools().Items["browser.useragent.osPatch"] = Item{
		Item:        "browser/useragent/osPatch",
		Enabled:     true,
		GJSON:       "Os.Patch",
		Description: "OS patch",
		Type:        "String"}
	job.Tools().Items["browser.useragent.osPatchMinor"] = Item{
		Item:        "browser/useragent/osPatchMinor",
		Enabled:     true,
		GJSON:       "Os.PatchMinor",
		Description: "OS minor patch",
		Type:        "String"}
	job.Tools().Items["browser.useragent.deviceFamily"] = Item{
		Item:        "browser/useragent/deviceFamily",
		Enabled:     true,
		GJSON:       "Device.Family",
		Description: "Device family",
		Type:        "String"}
	job.Tools().Items["browser.useragent.deviceBrand"] = Item{
		Item:        "browser/useragent/deviceBrand",
		Enabled:     true,
		GJSON:       "Device.Brand",
		Description: "Device brand",
		Type:        "String"}
	job.Tools().Items["browser.useragent.deviceModel"] = Item{
		Item:        "browser/useragent/deviceModel",
		Enabled:     true,
		GJSON:       "Device.Model",
		Description: "Device model",
		Type:        "String"}

	return err
}
