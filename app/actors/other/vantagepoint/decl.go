package vantagepoint

import "github.com/ottemo/foundation/env"

const (
	ConstErrorModule = "vantagepoint"
	ConstErrorLevel  = env.ConstErrorLevelActor

	ConstConfigPathVantagePoint = "general.vantagepoint"
	ConstConfigPathVantagePointEnabled = "general.vantagepoint.enabled"
	ConstConfigPathVantagePointUploadPath = "general.vantagepoint.upload.path"
	ConstConfigPathVantagePointUploadFileMask = "general.vantagepoint.upload.filemask"

	ConstSchedulerTaskName = "vantagePointCheckNewUploads"
)

type UploadProcessorInterface interface {
	Process() error
}

