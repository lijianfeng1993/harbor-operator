package v1

import (
	"context"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	RunningStatus = "running"
	FailedStatus  = "failed"
	FailedReason  = "fail to deploy harbor service"
)

func (cs *HarborServiceStatus) SetRunningStatus(url string) {
	c := Condition{
		Phase: "running",
	}
	cs.Condition = c
	cs.ExternalUrl = url
}

func (cs *HarborServiceStatus) SetFailedStatus(message string) {
	c := Condition{
		Phase:   FailedStatus,
		Reason:  FailedReason,
		Message: message,
	}
	cs.Condition = c
	cs.ExternalUrl = ""
}

func FlushInstanceStatus(c client.Client, instance *HarborService) error {
	c.Status().Update(context.TODO(), instance)
	return nil
}
