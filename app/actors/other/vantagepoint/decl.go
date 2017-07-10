package vantagepoint

import "github.com/ottemo/foundation/env"

const (
	ConstErrorModule = "vantagepoint"
	ConstErrorLevel  = env.ConstErrorLevelActor
	ConstLogStorage  = "vantagepoint.log"

	ConstConfigPathVantagePoint                = "general.vantagepoint"
	ConstConfigPathVantagePointScheduleEnabled = "general.vantagepoint.schedule.enabled"
	ConstConfigPathVantagePointScheduleHour    = "general.vantagepoint.schedule.hour"
	ConstConfigPathVantagePointUploadPath      = "general.vantagepoint.upload.path"
	ConstConfigPathVantagePointUploadFileMask  = "general.vantagepoint.upload.filemask"

	ConstSchedulerTaskName = "vantagePointCheckNewUploads"
)

type UploadProcessorInterface interface {
	Process() error
}
