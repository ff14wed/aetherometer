package hub_test

import (
	"sync"

	"github.com/ff14wed/aetherometer/core/hub"
	"go.uber.org/atomic"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NotifyHub", func() {
	It("broadcasts messages to multiple subscribers", func() {
		h := hub.NewNotifyHub(5)

		var countReceived atomic.Uint32
		wg := sync.WaitGroup{}

		subscriber := func() {
			defer GinkgoRecover()
			sub, id := h.Subscribe()
			defer h.Unsubscribe(id)

			wg.Done()
			<-sub
			countReceived.Add(1)
		}
		wg.Add(5)
		for i := 0; i < 5; i++ {
			go subscriber()
		}
		// Wait for all the subscribers to finish adding their subscriber
		wg.Wait()

		h.Broadcast()
		Eventually(func() int {
			return int(countReceived.Load())
		}).Should(Equal(5), "Not all subscribers received messages")
	})

	It("no longer broadcasts messages to removed subscribers", func() {
		h := hub.NewNotifyHub(5)
		sub, id := h.Subscribe()
		h.Unsubscribe(id)
		h.Broadcast()
		Consistently(sub).ShouldNot(Receive())
	})
})
