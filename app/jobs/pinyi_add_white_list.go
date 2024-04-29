package jobs

import "goravel/app/models"

type PinyiAddWhiteList struct {
}

// Signature The name and signature of the job.
func (receiver *PinyiAddWhiteList) Signature() string {
	return "PinyiAddWhiteList"
}

// Handle Execute the job.
func (receiver *PinyiAddWhiteList) Handle(args ...any) error {
	ProxyPinYi := models.ProxyPinYi{}
	for _, arg := range args {
		ProxyPinYi = arg.(models.ProxyPinYi)
	}
	if ProxyPinYi.HasFlowBalance() {
		ProxyPinYi.AddThisIpToWhiteList()
	}
	return nil
}
