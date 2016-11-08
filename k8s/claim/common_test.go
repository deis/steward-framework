package claim

import (
	"github.com/deis/steward-framework/k8s"
	"github.com/pborman/uuid"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/watch"
)

func getEvent(claim k8s.ServicePlanClaim) *Event {
	return &Event{
		claim: &k8s.ServicePlanClaimWrapper{
			Claim: &claim,
			ObjectMeta: v1.ObjectMeta{
				ResourceVersion: "1",
				Name:            "testclaim",
				Namespace:       "testns",
				Labels:          map[string]string{"label-1": "label1"},
			},
		},
		operation: watch.Added,
	}
}

func getClaim(action k8s.ServicePlanClaimAction) k8s.ServicePlanClaim {
	return k8s.ServicePlanClaim{
		TargetName: "target1",
		ServiceID:  "svc1",
		PlanID:     "plan1",
		ClaimID:    uuid.New(),
		Action:     action.String(),
	}
}

func getClaimWithStatus(action k8s.ServicePlanClaimAction, status k8s.ServicePlanClaimStatus) k8s.ServicePlanClaim {
	cl := getClaim(action)
	cl.Status = status.String()
	return cl
}
