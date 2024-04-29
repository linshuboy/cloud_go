package jobs

import "goravel/app/models"

type PinyiAddIpPool struct {
}

// Signature The name and signature of the job.
func (receiver *PinyiAddIpPool) Signature() string {
	return "PinyiAddIpPool"
}

// Handle Execute the job.
func (receiver *PinyiAddIpPool) Handle(args ...any) error {
	ProxyPinYi := models.ProxyPinYi{}
	for _, arg := range args {
		ProxyPinYi = arg.(models.ProxyPinYi)
	}
	ProxyPinYi.AddProxyToIpPool()
	return nil
}
