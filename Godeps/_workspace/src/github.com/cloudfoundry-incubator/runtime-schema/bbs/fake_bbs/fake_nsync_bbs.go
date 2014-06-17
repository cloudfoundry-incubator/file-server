// This file was generated by counterfeiter
package fake_bbs

import (
	. "github.com/cloudfoundry-incubator/runtime-schema/bbs"

	"github.com/cloudfoundry-incubator/runtime-schema/models"

	"sync"
)

type FakeNsyncBBS struct {
	DesireLRPStub        func(models.DesiredLRP) error
	desireLRPMutex       sync.RWMutex
	desireLRPArgsForCall []struct {
		arg1 models.DesiredLRP
	}
	desireLRPReturns struct {
		result1 error
	}
	RemoveDesiredLRPByProcessGuidStub        func(guid string) error
	removeDesiredLRPByProcessGuidMutex       sync.RWMutex
	removeDesiredLRPByProcessGuidArgsForCall []struct {
		arg1 string
	}
	removeDesiredLRPByProcessGuidReturns struct {
		result1 error
	}
	GetAllDesiredLRPsStub        func() ([]models.DesiredLRP, error)
	getAllDesiredLRPsMutex       sync.RWMutex
	getAllDesiredLRPsArgsForCall []struct{}
	getAllDesiredLRPsReturns struct {
		result1 []models.DesiredLRP
		result2 error
	}
	ChangeDesiredLRPStub        func(change models.DesiredLRPChange) error
	changeDesiredLRPMutex       sync.RWMutex
	changeDesiredLRPArgsForCall []struct {
		arg1 models.DesiredLRPChange
	}
	changeDesiredLRPReturns struct {
		result1 error
	}
}

func (fake *FakeNsyncBBS) DesireLRP(arg1 models.DesiredLRP) error {
	fake.desireLRPMutex.Lock()
	defer fake.desireLRPMutex.Unlock()
	fake.desireLRPArgsForCall = append(fake.desireLRPArgsForCall, struct {
		arg1 models.DesiredLRP
	}{arg1})
	if fake.DesireLRPStub != nil {
		return fake.DesireLRPStub(arg1)
	} else {
		return fake.desireLRPReturns.result1
	}
}

func (fake *FakeNsyncBBS) DesireLRPCallCount() int {
	fake.desireLRPMutex.RLock()
	defer fake.desireLRPMutex.RUnlock()
	return len(fake.desireLRPArgsForCall)
}

func (fake *FakeNsyncBBS) DesireLRPArgsForCall(i int) models.DesiredLRP {
	fake.desireLRPMutex.RLock()
	defer fake.desireLRPMutex.RUnlock()
	return fake.desireLRPArgsForCall[i].arg1
}

func (fake *FakeNsyncBBS) DesireLRPReturns(result1 error) {
	fake.desireLRPReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeNsyncBBS) RemoveDesiredLRPByProcessGuid(arg1 string) error {
	fake.removeDesiredLRPByProcessGuidMutex.Lock()
	defer fake.removeDesiredLRPByProcessGuidMutex.Unlock()
	fake.removeDesiredLRPByProcessGuidArgsForCall = append(fake.removeDesiredLRPByProcessGuidArgsForCall, struct {
		arg1 string
	}{arg1})
	if fake.RemoveDesiredLRPByProcessGuidStub != nil {
		return fake.RemoveDesiredLRPByProcessGuidStub(arg1)
	} else {
		return fake.removeDesiredLRPByProcessGuidReturns.result1
	}
}

func (fake *FakeNsyncBBS) RemoveDesiredLRPByProcessGuidCallCount() int {
	fake.removeDesiredLRPByProcessGuidMutex.RLock()
	defer fake.removeDesiredLRPByProcessGuidMutex.RUnlock()
	return len(fake.removeDesiredLRPByProcessGuidArgsForCall)
}

func (fake *FakeNsyncBBS) RemoveDesiredLRPByProcessGuidArgsForCall(i int) string {
	fake.removeDesiredLRPByProcessGuidMutex.RLock()
	defer fake.removeDesiredLRPByProcessGuidMutex.RUnlock()
	return fake.removeDesiredLRPByProcessGuidArgsForCall[i].arg1
}

func (fake *FakeNsyncBBS) RemoveDesiredLRPByProcessGuidReturns(result1 error) {
	fake.removeDesiredLRPByProcessGuidReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeNsyncBBS) GetAllDesiredLRPs() ([]models.DesiredLRP, error) {
	fake.getAllDesiredLRPsMutex.Lock()
	defer fake.getAllDesiredLRPsMutex.Unlock()
	fake.getAllDesiredLRPsArgsForCall = append(fake.getAllDesiredLRPsArgsForCall, struct{}{})
	if fake.GetAllDesiredLRPsStub != nil {
		return fake.GetAllDesiredLRPsStub()
	} else {
		return fake.getAllDesiredLRPsReturns.result1, fake.getAllDesiredLRPsReturns.result2
	}
}

func (fake *FakeNsyncBBS) GetAllDesiredLRPsCallCount() int {
	fake.getAllDesiredLRPsMutex.RLock()
	defer fake.getAllDesiredLRPsMutex.RUnlock()
	return len(fake.getAllDesiredLRPsArgsForCall)
}

func (fake *FakeNsyncBBS) GetAllDesiredLRPsReturns(result1 []models.DesiredLRP, result2 error) {
	fake.getAllDesiredLRPsReturns = struct {
		result1 []models.DesiredLRP
		result2 error
	}{result1, result2}
}

func (fake *FakeNsyncBBS) ChangeDesiredLRP(arg1 models.DesiredLRPChange) error {
	fake.changeDesiredLRPMutex.Lock()
	defer fake.changeDesiredLRPMutex.Unlock()
	fake.changeDesiredLRPArgsForCall = append(fake.changeDesiredLRPArgsForCall, struct {
		arg1 models.DesiredLRPChange
	}{arg1})
	if fake.ChangeDesiredLRPStub != nil {
		return fake.ChangeDesiredLRPStub(arg1)
	} else {
		return fake.changeDesiredLRPReturns.result1
	}
}

func (fake *FakeNsyncBBS) ChangeDesiredLRPCallCount() int {
	fake.changeDesiredLRPMutex.RLock()
	defer fake.changeDesiredLRPMutex.RUnlock()
	return len(fake.changeDesiredLRPArgsForCall)
}

func (fake *FakeNsyncBBS) ChangeDesiredLRPArgsForCall(i int) models.DesiredLRPChange {
	fake.changeDesiredLRPMutex.RLock()
	defer fake.changeDesiredLRPMutex.RUnlock()
	return fake.changeDesiredLRPArgsForCall[i].arg1
}

func (fake *FakeNsyncBBS) ChangeDesiredLRPReturns(result1 error) {
	fake.changeDesiredLRPReturns = struct {
		result1 error
	}{result1}
}

var _ NsyncBBS = new(FakeNsyncBBS)
