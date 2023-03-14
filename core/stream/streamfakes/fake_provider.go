// Code generated by counterfeiter. DO NOT EDIT.
package streamfakes

import (
	"sync"

	"github.com/ff14wed/aetherometer/core/stream"
	xivnet "github.com/ff14wed/xivnet/v3"
)

type FakeProvider struct {
	SendRequestStub        func([]byte) ([]byte, error)
	sendRequestMutex       sync.RWMutex
	sendRequestArgsForCall []struct {
		arg1 []byte
	}
	sendRequestReturns struct {
		result1 []byte
		result2 error
	}
	sendRequestReturnsOnCall map[int]struct {
		result1 []byte
		result2 error
	}
	StreamIDStub        func() int
	streamIDMutex       sync.RWMutex
	streamIDArgsForCall []struct {
	}
	streamIDReturns struct {
		result1 int
	}
	streamIDReturnsOnCall map[int]struct {
		result1 int
	}
	SubscribeEgressStub        func() <-chan *xivnet.Block
	subscribeEgressMutex       sync.RWMutex
	subscribeEgressArgsForCall []struct {
	}
	subscribeEgressReturns struct {
		result1 <-chan *xivnet.Block
	}
	subscribeEgressReturnsOnCall map[int]struct {
		result1 <-chan *xivnet.Block
	}
	SubscribeIngressStub        func() <-chan *xivnet.Block
	subscribeIngressMutex       sync.RWMutex
	subscribeIngressArgsForCall []struct {
	}
	subscribeIngressReturns struct {
		result1 <-chan *xivnet.Block
	}
	subscribeIngressReturnsOnCall map[int]struct {
		result1 <-chan *xivnet.Block
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeProvider) SendRequest(arg1 []byte) ([]byte, error) {
	var arg1Copy []byte
	if arg1 != nil {
		arg1Copy = make([]byte, len(arg1))
		copy(arg1Copy, arg1)
	}
	fake.sendRequestMutex.Lock()
	ret, specificReturn := fake.sendRequestReturnsOnCall[len(fake.sendRequestArgsForCall)]
	fake.sendRequestArgsForCall = append(fake.sendRequestArgsForCall, struct {
		arg1 []byte
	}{arg1Copy})
	fake.recordInvocation("SendRequest", []interface{}{arg1Copy})
	fake.sendRequestMutex.Unlock()
	if fake.SendRequestStub != nil {
		return fake.SendRequestStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.sendRequestReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeProvider) SendRequestCallCount() int {
	fake.sendRequestMutex.RLock()
	defer fake.sendRequestMutex.RUnlock()
	return len(fake.sendRequestArgsForCall)
}

func (fake *FakeProvider) SendRequestCalls(stub func([]byte) ([]byte, error)) {
	fake.sendRequestMutex.Lock()
	defer fake.sendRequestMutex.Unlock()
	fake.SendRequestStub = stub
}

func (fake *FakeProvider) SendRequestArgsForCall(i int) []byte {
	fake.sendRequestMutex.RLock()
	defer fake.sendRequestMutex.RUnlock()
	argsForCall := fake.sendRequestArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeProvider) SendRequestReturns(result1 []byte, result2 error) {
	fake.sendRequestMutex.Lock()
	defer fake.sendRequestMutex.Unlock()
	fake.SendRequestStub = nil
	fake.sendRequestReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *FakeProvider) SendRequestReturnsOnCall(i int, result1 []byte, result2 error) {
	fake.sendRequestMutex.Lock()
	defer fake.sendRequestMutex.Unlock()
	fake.SendRequestStub = nil
	if fake.sendRequestReturnsOnCall == nil {
		fake.sendRequestReturnsOnCall = make(map[int]struct {
			result1 []byte
			result2 error
		})
	}
	fake.sendRequestReturnsOnCall[i] = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *FakeProvider) StreamID() int {
	fake.streamIDMutex.Lock()
	ret, specificReturn := fake.streamIDReturnsOnCall[len(fake.streamIDArgsForCall)]
	fake.streamIDArgsForCall = append(fake.streamIDArgsForCall, struct {
	}{})
	fake.recordInvocation("StreamID", []interface{}{})
	fake.streamIDMutex.Unlock()
	if fake.StreamIDStub != nil {
		return fake.StreamIDStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.streamIDReturns
	return fakeReturns.result1
}

func (fake *FakeProvider) StreamIDCallCount() int {
	fake.streamIDMutex.RLock()
	defer fake.streamIDMutex.RUnlock()
	return len(fake.streamIDArgsForCall)
}

func (fake *FakeProvider) StreamIDCalls(stub func() int) {
	fake.streamIDMutex.Lock()
	defer fake.streamIDMutex.Unlock()
	fake.StreamIDStub = stub
}

func (fake *FakeProvider) StreamIDReturns(result1 int) {
	fake.streamIDMutex.Lock()
	defer fake.streamIDMutex.Unlock()
	fake.StreamIDStub = nil
	fake.streamIDReturns = struct {
		result1 int
	}{result1}
}

func (fake *FakeProvider) StreamIDReturnsOnCall(i int, result1 int) {
	fake.streamIDMutex.Lock()
	defer fake.streamIDMutex.Unlock()
	fake.StreamIDStub = nil
	if fake.streamIDReturnsOnCall == nil {
		fake.streamIDReturnsOnCall = make(map[int]struct {
			result1 int
		})
	}
	fake.streamIDReturnsOnCall[i] = struct {
		result1 int
	}{result1}
}

func (fake *FakeProvider) SubscribeEgress() <-chan *xivnet.Block {
	fake.subscribeEgressMutex.Lock()
	ret, specificReturn := fake.subscribeEgressReturnsOnCall[len(fake.subscribeEgressArgsForCall)]
	fake.subscribeEgressArgsForCall = append(fake.subscribeEgressArgsForCall, struct {
	}{})
	fake.recordInvocation("SubscribeEgress", []interface{}{})
	fake.subscribeEgressMutex.Unlock()
	if fake.SubscribeEgressStub != nil {
		return fake.SubscribeEgressStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.subscribeEgressReturns
	return fakeReturns.result1
}

func (fake *FakeProvider) SubscribeEgressCallCount() int {
	fake.subscribeEgressMutex.RLock()
	defer fake.subscribeEgressMutex.RUnlock()
	return len(fake.subscribeEgressArgsForCall)
}

func (fake *FakeProvider) SubscribeEgressCalls(stub func() <-chan *xivnet.Block) {
	fake.subscribeEgressMutex.Lock()
	defer fake.subscribeEgressMutex.Unlock()
	fake.SubscribeEgressStub = stub
}

func (fake *FakeProvider) SubscribeEgressReturns(result1 <-chan *xivnet.Block) {
	fake.subscribeEgressMutex.Lock()
	defer fake.subscribeEgressMutex.Unlock()
	fake.SubscribeEgressStub = nil
	fake.subscribeEgressReturns = struct {
		result1 <-chan *xivnet.Block
	}{result1}
}

func (fake *FakeProvider) SubscribeEgressReturnsOnCall(i int, result1 <-chan *xivnet.Block) {
	fake.subscribeEgressMutex.Lock()
	defer fake.subscribeEgressMutex.Unlock()
	fake.SubscribeEgressStub = nil
	if fake.subscribeEgressReturnsOnCall == nil {
		fake.subscribeEgressReturnsOnCall = make(map[int]struct {
			result1 <-chan *xivnet.Block
		})
	}
	fake.subscribeEgressReturnsOnCall[i] = struct {
		result1 <-chan *xivnet.Block
	}{result1}
}

func (fake *FakeProvider) SubscribeIngress() <-chan *xivnet.Block {
	fake.subscribeIngressMutex.Lock()
	ret, specificReturn := fake.subscribeIngressReturnsOnCall[len(fake.subscribeIngressArgsForCall)]
	fake.subscribeIngressArgsForCall = append(fake.subscribeIngressArgsForCall, struct {
	}{})
	fake.recordInvocation("SubscribeIngress", []interface{}{})
	fake.subscribeIngressMutex.Unlock()
	if fake.SubscribeIngressStub != nil {
		return fake.SubscribeIngressStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.subscribeIngressReturns
	return fakeReturns.result1
}

func (fake *FakeProvider) SubscribeIngressCallCount() int {
	fake.subscribeIngressMutex.RLock()
	defer fake.subscribeIngressMutex.RUnlock()
	return len(fake.subscribeIngressArgsForCall)
}

func (fake *FakeProvider) SubscribeIngressCalls(stub func() <-chan *xivnet.Block) {
	fake.subscribeIngressMutex.Lock()
	defer fake.subscribeIngressMutex.Unlock()
	fake.SubscribeIngressStub = stub
}

func (fake *FakeProvider) SubscribeIngressReturns(result1 <-chan *xivnet.Block) {
	fake.subscribeIngressMutex.Lock()
	defer fake.subscribeIngressMutex.Unlock()
	fake.SubscribeIngressStub = nil
	fake.subscribeIngressReturns = struct {
		result1 <-chan *xivnet.Block
	}{result1}
}

func (fake *FakeProvider) SubscribeIngressReturnsOnCall(i int, result1 <-chan *xivnet.Block) {
	fake.subscribeIngressMutex.Lock()
	defer fake.subscribeIngressMutex.Unlock()
	fake.SubscribeIngressStub = nil
	if fake.subscribeIngressReturnsOnCall == nil {
		fake.subscribeIngressReturnsOnCall = make(map[int]struct {
			result1 <-chan *xivnet.Block
		})
	}
	fake.subscribeIngressReturnsOnCall[i] = struct {
		result1 <-chan *xivnet.Block
	}{result1}
}

func (fake *FakeProvider) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.sendRequestMutex.RLock()
	defer fake.sendRequestMutex.RUnlock()
	fake.streamIDMutex.RLock()
	defer fake.streamIDMutex.RUnlock()
	fake.subscribeEgressMutex.RLock()
	defer fake.subscribeEgressMutex.RUnlock()
	fake.subscribeIngressMutex.RLock()
	defer fake.subscribeIngressMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeProvider) recordInvocation(key string, args []interface{}) {
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

var _ stream.Provider = new(FakeProvider)
