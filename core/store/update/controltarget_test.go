package update_test

import (
	"github.com/ff14wed/aetherometer/core/models"
	"github.com/ff14wed/aetherometer/core/store"
	"github.com/ff14wed/aetherometer/core/store/update"
	"github.com/ff14wed/xivnet/v3"
	"github.com/ff14wed/xivnet/v3/datatypes"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gopkg.in/dealancer/validate.v2"
)

var _ = Describe("ControlTarget Update", func() {
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

			notify4Data := &datatypes.ControlTarget{
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

			Expect(validate.Validate(entityEvents)).To(Succeed())
			Expect(validate.Validate(streams)).To(Succeed())
		})

		entityValidationTests(testEnv, false)
	})
})
