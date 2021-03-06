// Code generated by counterfeiter. DO NOT EDIT.
package hookfakes

import (
	"net"
	"sync"
	"time"

	"github.com/ff14wed/aetherometer/core/adapter/hook"
)

type FakeRemoteProcessProvider struct {
	DialPipeStub        func(string, *time.Duration) (net.Conn, error)
	dialPipeMutex       sync.RWMutex
	dialPipeArgsForCall []struct {
		arg1 string
		arg2 *time.Duration
	}
	dialPipeReturns struct {
		result1 net.Conn
		result2 error
	}
	dialPipeReturnsOnCall map[int]struct {
		result1 net.Conn
		result2 error
	}
	InjectDLLStub        func(uint32, string) error
	injectDLLMutex       sync.RWMutex
	injectDLLArgsForCall []struct {
		arg1 uint32
		arg2 string
	}
	injectDLLReturns struct {
		result1 error
	}
	injectDLLReturnsOnCall map[int]struct {
		result1 error
	}
	IsPipeClosedStub        func(error) bool
	isPipeClosedMutex       sync.RWMutex
	isPipeClosedArgsForCall []struct {
		arg1 error
	}
	isPipeClosedReturns struct {
		result1 bool
	}
	isPipeClosedReturnsOnCall map[int]struct {
		result1 bool
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeRemoteProcessProvider) DialPipe(arg1 string, arg2 *time.Duration) (net.Conn, error) {
	fake.dialPipeMutex.Lock()
	ret, specificReturn := fake.dialPipeReturnsOnCall[len(fake.dialPipeArgsForCall)]
	fake.dialPipeArgsForCall = append(fake.dialPipeArgsForCall, struct {
		arg1 string
		arg2 *time.Duration
	}{arg1, arg2})
	fake.recordInvocation("DialPipe", []interface{}{arg1, arg2})
	fake.dialPipeMutex.Unlock()
	if fake.DialPipeStub != nil {
		return fake.DialPipeStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.dialPipeReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeRemoteProcessProvider) DialPipeCallCount() int {
	fake.dialPipeMutex.RLock()
	defer fake.dialPipeMutex.RUnlock()
	return len(fake.dialPipeArgsForCall)
}

func (fake *FakeRemoteProcessProvider) DialPipeCalls(stub func(string, *time.Duration) (net.Conn, error)) {
	fake.dialPipeMutex.Lock()
	defer fake.dialPipeMutex.Unlock()
	fake.DialPipeStub = stub
}

func (fake *FakeRemoteProcessProvider) DialPipeArgsForCall(i int) (string, *time.Duration) {
	fake.dialPipeMutex.RLock()
	defer fake.dialPipeMutex.RUnlock()
	argsForCall := fake.dialPipeArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeRemoteProcessProvider) DialPipeReturns(result1 net.Conn, result2 error) {
	fake.dialPipeMutex.Lock()
	defer fake.dialPipeMutex.Unlock()
	fake.DialPipeStub = nil
	fake.dialPipeReturns = struct {
		result1 net.Conn
		result2 error
	}{result1, result2}
}

func (fake *FakeRemoteProcessProvider) DialPipeReturnsOnCall(i int, result1 net.Conn, result2 error) {
	fake.dialPipeMutex.Lock()
	defer fake.dialPipeMutex.Unlock()
	fake.DialPipeStub = nil
	if fake.dialPipeReturnsOnCall == nil {
		fake.dialPipeReturnsOnCall = make(map[int]struct {
			result1 net.Conn
			result2 error
		})
	}
	fake.dialPipeReturnsOnCall[i] = struct {
		result1 net.Conn
		result2 error
	}{result1, result2}
}

func (fake *FakeRemoteProcessProvider) InjectDLL(arg1 uint32, arg2 string) error {
	fake.injectDLLMutex.Lock()
	ret, specificReturn := fake.injectDLLReturnsOnCall[len(fake.injectDLLArgsForCall)]
	fake.injectDLLArgsForCall = append(fake.injectDLLArgsForCall, struct {
		arg1 uint32
		arg2 string
	}{arg1, arg2})
	fake.recordInvocation("InjectDLL", []interface{}{arg1, arg2})
	fake.injectDLLMutex.Unlock()
	if fake.InjectDLLStub != nil {
		return fake.InjectDLLStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.injectDLLReturns
	return fakeReturns.result1
}

func (fake *FakeRemoteProcessProvider) InjectDLLCallCount() int {
	fake.injectDLLMutex.RLock()
	defer fake.injectDLLMutex.RUnlock()
	return len(fake.injectDLLArgsForCall)
}

func (fake *FakeRemoteProcessProvider) InjectDLLCalls(stub func(uint32, string) error) {
	fake.injectDLLMutex.Lock()
	defer fake.injectDLLMutex.Unlock()
	fake.InjectDLLStub = stub
}

func (fake *FakeRemoteProcessProvider) InjectDLLArgsForCall(i int) (uint32, string) {
	fake.injectDLLMutex.RLock()
	defer fake.injectDLLMutex.RUnlock()
	argsForCall := fake.injectDLLArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeRemoteProcessProvider) InjectDLLReturns(result1 error) {
	fake.injectDLLMutex.Lock()
	defer fake.injectDLLMutex.Unlock()
	fake.InjectDLLStub = nil
	fake.injectDLLReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeRemoteProcessProvider) InjectDLLReturnsOnCall(i int, result1 error) {
	fake.injectDLLMutex.Lock()
	defer fake.injectDLLMutex.Unlock()
	fake.InjectDLLStub = nil
	if fake.injectDLLReturnsOnCall == nil {
		fake.injectDLLReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.injectDLLReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeRemoteProcessProvider) IsPipeClosed(arg1 error) bool {
	fake.isPipeClosedMutex.Lock()
	ret, specificReturn := fake.isPipeClosedReturnsOnCall[len(fake.isPipeClosedArgsForCall)]
	fake.isPipeClosedArgsForCall = append(fake.isPipeClosedArgsForCall, struct {
		arg1 error
	}{arg1})
	fake.recordInvocation("IsPipeClosed", []interface{}{arg1})
	fake.isPipeClosedMutex.Unlock()
	if fake.IsPipeClosedStub != nil {
		return fake.IsPipeClosedStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.isPipeClosedReturns
	return fakeReturns.result1
}

func (fake *FakeRemoteProcessProvider) IsPipeClosedCallCount() int {
	fake.isPipeClosedMutex.RLock()
	defer fake.isPipeClosedMutex.RUnlock()
	return len(fake.isPipeClosedArgsForCall)
}

func (fake *FakeRemoteProcessProvider) IsPipeClosedCalls(stub func(error) bool) {
	fake.isPipeClosedMutex.Lock()
	defer fake.isPipeClosedMutex.Unlock()
	fake.IsPipeClosedStub = stub
}

func (fake *FakeRemoteProcessProvider) IsPipeClosedArgsForCall(i int) error {
	fake.isPipeClosedMutex.RLock()
	defer fake.isPipeClosedMutex.RUnlock()
	argsForCall := fake.isPipeClosedArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeRemoteProcessProvider) IsPipeClosedReturns(result1 bool) {
	fake.isPipeClosedMutex.Lock()
	defer fake.isPipeClosedMutex.Unlock()
	fake.IsPipeClosedStub = nil
	fake.isPipeClosedReturns = struct {
		result1 bool
	}{result1}
}

func (fake *FakeRemoteProcessProvider) IsPipeClosedReturnsOnCall(i int, result1 bool) {
	fake.isPipeClosedMutex.Lock()
	defer fake.isPipeClosedMutex.Unlock()
	fake.IsPipeClosedStub = nil
	if fake.isPipeClosedReturnsOnCall == nil {
		fake.isPipeClosedReturnsOnCall = make(map[int]struct {
			result1 bool
		})
	}
	fake.isPipeClosedReturnsOnCall[i] = struct {
		result1 bool
	}{result1}
}

func (fake *FakeRemoteProcessProvider) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.dialPipeMutex.RLock()
	defer fake.dialPipeMutex.RUnlock()
	fake.injectDLLMutex.RLock()
	defer fake.injectDLLMutex.RUnlock()
	fake.isPipeClosedMutex.RLock()
	defer fake.isPipeClosedMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeRemoteProcessProvider) recordInvocation(key string, args []interface{}) {
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

var _ hook.RemoteProcessProvider = new(FakeRemoteProcessProvider)
