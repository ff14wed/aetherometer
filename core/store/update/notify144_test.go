package update_test

import (
	"github.com/ff14wed/sibyl/backend/models"
	"github.com/ff14wed/sibyl/backend/store"
	"github.com/ff14wed/sibyl/backend/store/update"
	"github.com/ff14wed/xivnet/v3"
	"github.com/ff14wed/xivnet/v3/datatypes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Notify144 Update", func() {
	Describe("type 0x32", func() {
		var (
			testEnv = new(testVars)

			b         *xivnet.Block
			streams   *store.Streams
			streamID  int
			subjectID uint64
			entity    *models.Entity
			generator update.Generator

			expectedTarget uint64
		)

		BeforeEach(func() {
			*testEnv = genericSetup()
			b = testEnv.b
			streams = testEnv.streams
			streamID = testEnv.streamID
			subjectID = testEnv.subjectID
			entity = testEnv.entity
			generator = testEnv.generator

			expectedTarget = 0xABCDEF01

			notify4Data := &datatypes.Notify144{
				Type:     0x32,
				TargetID: uint32(expectedTarget),
			}

			b.Data = notify4Data
		})

		It("generates an update that sets the entity's target", func() {
			u := generator.Generate(streamID, false, b)
			Expect(u).ToNot(BeNil())
			streamEvents, entityEvents, err := u.ModifyStore(streams)
			Expect(err).ToNot(HaveOccurred())
			Expect(streamEvents).To(BeEmpty())

			Expect(entityEvents).To(ConsistOf(models.EntityEvent{
				StreamID: streamID,
				EntityID: subjectID,
				Type: models.UpdateTarget{
					TargetID: expectedTarget,
				},
			}))

			Expect(entity.TargetID).To(Equal(expectedTarget))
		})

		entityValidationTests(testEnv, false)
	})
})
