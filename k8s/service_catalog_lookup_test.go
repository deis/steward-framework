package k8s

import (
	"fmt"
	"testing"

	"github.com/arschles/assert"
	"github.com/deis/steward-framework"
	"github.com/pborman/uuid"
)

func TestServiceCatalogLookupCatalogKey(t *testing.T) {
	svcID := uuid.New()
	planID := uuid.New()

	key := catalogKey(svcID, planID)
	assert.Equal(t, key, fmt.Sprintf("%s-%s", svcID, planID), "catalog key")
}

func TestServiceCatalogLookupGetSet(t *testing.T) {
	initial := []*ServiceCatalogEntry{
		&ServiceCatalogEntry{
			Info: framework.ServiceInfo{ID: "testsvc1"},
			Plan: framework.ServicePlan{ID: "testplan1"},
		},
		&ServiceCatalogEntry{
			Info: framework.ServiceInfo{ID: "testsvc2"},
			Plan: framework.ServicePlan{ID: "testplan2"},
		},
	}
	lookup := NewServiceCatalogLookup(initial)
	for _, entry := range initial {
		fetched := lookup.Get(entry.Info.ID, entry.Plan.ID)
		if fetched == nil {
			t.Fatalf("service expected but not found for service ID %s, plan ID %s", entry.Info.ID, entry.Plan.ID)
		}
		assert.Equal(t, fetched.Info.ID, entry.Info.ID, "service ID")
		assert.Equal(t, fetched.Plan.ID, entry.Plan.ID, "plan ID")
	}
	newEntry := &ServiceCatalogEntry{
		Info: framework.ServiceInfo{ID: "testsvc3"},
		Plan: framework.ServicePlan{ID: "testplan3"},
	}
	lookup.Set(newEntry)
	fetched := lookup.Get(newEntry.Info.ID, newEntry.Plan.ID)
	if fetched == nil {
		t.Fatalf("entry with service ID %s, plan ID %s not found after it was set", newEntry.Info.ID, newEntry.Plan.ID)
	}
	assert.Equal(t, fetched.Info.ID, newEntry.Info.ID, "service ID")
	assert.Equal(t, fetched.Plan.ID, newEntry.Plan.ID, "plan ID")
}
