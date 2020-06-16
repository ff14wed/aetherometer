// Code generated by counterfeiter. DO NOT EDIT.
package modelsfakes

import (
	"sync"

	"github.com/ff14wed/aetherometer/core/models"
)

type FakeStreamEventSource struct {
	SubscribeStub        func() (chan *models.StreamEvent, uint64)
	subscribeMutex       sync.RWMutex
	subscribeArgsForCall []struct {
	}
	subscribeReturns struct {
		result1 chan *models.StreamEvent
		result2 uint64
	}
	subscribeReturnsOnCall map[int]struct {
		result1 chan *models.StreamEvent
		result2 uint64
	}
	UnsubscribeStub        func(uint64)
	unsubscribeMutex       sync.RWMutex
	unsubscribeArgsForCall []struct {
		arg1 uint64
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeStreamEventSource) Subscribe() (chan *models.StreamEvent, uint64) {
	fake.subscribeMutex.Lock()
	ret, specificReturn := fake.subscribeReturnsOnCall[len(fake.subscribeArgsForCall)]
	fake.subscribeArgsForCall = append(fake.subscribeArgsForCall, struct {
	}{})
	fake.recordInvocation("Subscribe", []interface{}{})
	fake.subscribeMutex.Unlock()
	if fake.SubscribeStub != nil {
		return fake.SubscribeStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.subscribeReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeStreamEventSource) SubscribeCallCount() int {
	fake.subscribeMutex.RLock()
	defer fake.subscribeMutex.RUnlock()
	return len(fake.subscribeArgsForCall)
}

func (fake *FakeStreamEventSource) SubscribeCalls(stub func() (chan *models.StreamEvent, uint64)) {
	fake.subscribeMutex.Lock()
	defer fake.subscribeMutex.Unlock()
	fake.SubscribeStub = stub
}

func (fake *FakeStreamEventSource) SubscribeReturns(result1 chan *models.StreamEvent, result2 uint64) {
	fake.subscribeMutex.Lock()
	defer fake.subscribeMutex.Unlock()
	fake.SubscribeStub = nil
	fake.subscribeReturns = struct {
		result1 chan *models.StreamEvent
		result2 uint64
	}{result1, result2}
}

func (fake *FakeStreamEventSource) SubscribeReturnsOnCall(i int, result1 chan *models.StreamEvent, result2 uint64) {
	fake.subscribeMutex.Lock()
	defer fake.subscribeMutex.Unlock()
	fake.SubscribeStub = nil
	if fake.subscribeReturnsOnCall == nil {
		fake.subscribeReturnsOnCall = make(map[int]struct {
			result1 chan *models.StreamEvent
			result2 uint64
		})
	}
	fake.subscribeReturnsOnCall[i] = struct {
		result1 chan *models.StreamEvent
		result2 uint64
	}{result1, result2}
}

func (fake *FakeStreamEventSource) Unsubscribe(arg1 uint64) {
	fake.unsubscribeMutex.Lock()
	fake.unsubscribeArgsForCall = append(fake.unsubscribeArgsForCall, struct {
		arg1 uint64
	}{arg1})
	fake.recordInvocation("Unsubscribe", []interface{}{arg1})
	fake.unsubscribeMutex.Unlock()
	if fake.UnsubscribeStub != nil {
		fake.UnsubscribeStub(arg1)
	}
}

func (fake *FakeStreamEventSource) UnsubscribeCallCount() int {
	fake.unsubscribeMutex.RLock()
	defer fake.unsubscribeMutex.RUnlock()
	return len(fake.unsubscribeArgsForCall)
}

func (fake *FakeStreamEventSource) UnsubscribeCalls(stub func(uint64)) {
	fake.unsubscribeMutex.Lock()
	defer fake.unsubscribeMutex.Unlock()
	fake.UnsubscribeStub = stub
}

func (fake *FakeStreamEventSource) UnsubscribeArgsForCall(i int) uint64 {
	fake.unsubscribeMutex.RLock()
	defer fake.unsubscribeMutex.RUnlock()
	argsForCall := fake.unsubscribeArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeStreamEventSource) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.subscribeMutex.RLock()
	defer fake.subscribeMutex.RUnlock()
	fake.unsubscribeMutex.RLock()
	defer fake.unsubscribeMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeStreamEventSource) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ models.StreamEventSource = new(FakeStreamEventSource)
