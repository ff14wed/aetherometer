package hub_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/ff14wed/aetherometer/core/hub"
)

var _ = Describe("EventBatcher", func() {
	It("emits one event if every notification was within the batching interval", func() {
		b := hub.NewEventBatcher(20 * time.Millisecond)
		for i := 0; i < 10; i++ {
			b.Notify()
		}
		// The batcher should emit only one event
		Eventually(b.BatchedEvents()).Should(Receive())
		Consistently(b.BatchedEvents()).ShouldNot(Receive())
	})

	It("emits more than one event if notifications are within different batching intervals", func() {
		b := hub.NewEventBatcher(20 * time.Millisecond)
		for i := 0; i < 5; i++ {
			b.Notify()
		}
		time.Sleep(50 * time.Millisecond)
		for i := 0; i < 5; i++ {
			b.Notify()
		}

		// The batcher should emit only two events
		Eventually(b.BatchedEvents()).Should(Receive())
		Eventually(b.BatchedEvents()).Should(Receive())
		Consistently(b.BatchedEvents()).ShouldNot(Receive())
	})
})
