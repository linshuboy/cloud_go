package jobs

import "goravel/app/models"

type PinyiRemoveInvalidIp struct {
}

// Signature The name and signature of the job.
func (receiver *PinyiRemoveInvalidIp) Signature() string {
	return "PinyiRemoveInvalidIp"
}

// Handle Execute the job.
func (receiver *PinyiRemoveInvalidIp) Handle(args ...any) error {
	ProxyPinYi := models.IpPool{}
	ProxyPinYi.RemoveInvalidIp()
	return nil
}
