package hub_test

import (
	"fmt"
	"sync"

	"github.com/ff14wed/sibyl/backend/hub"
	"github.com/ff14wed/sibyl/backend/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("EntityHub", func() {
	It("broadcasts messages to multiple subscribers", func() {
		h := hub.NewEntityHub(5)
		expected := models.EntityEvent{StreamID: 1234}
		receivedChan := make(chan models.EntityEvent, 5)
		defer close(receivedChan)

		wg := sync.WaitGroup{}

		subscriber := func() {
			defer GinkgoRecover()
			sub, id := h.Subscribe()
			defer h.Unsubscribe(id)

			wg.Done()
			val := <-sub
			receivedChan <- val
		}
		wg.Add(5)
		for i := 0; i < 5; i++ {
			go subscriber()
		}
		// Wait for all the subscribers to finish adding their subscriber
		wg.Wait()

		h.Broadcast(models.EntityEvent{StreamID: 1234})
		for i := 0; i < 5; i++ {
			var receivedMsg models.EntityEvent
			Eventually(receivedChan).Should(
				Receive(&receivedMsg),
				fmt.Sprintf("Subscriber %d did not receive an expected message", i),
			)
			Expect(receivedMsg).To(Equal(expected))
		}
	})

	It("no longer broadcasts messages to removed subscribers", func() {
		h := hub.NewEntityHub(5)
		sub, id := h.Subscribe()
		h.Unsubscribe(id)
		h.Broadcast(models.EntityEvent{StreamID: 1234})
		Consistently(sub).ShouldNot(Receive())
	})
})
