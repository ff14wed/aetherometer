package update_test

import (
	"github.com/ff14wed/aetherometer/core/datasheet"
	"github.com/ff14wed/aetherometer/core/models"
	"github.com/ff14wed/aetherometer/core/store"
	"github.com/ff14wed/aetherometer/core/store/update"
	"github.com/ff14wed/aetherometer/core/testassets"
	"github.com/ff14wed/xivnet/v3"
	"github.com/ff14wed/xivnet/v3/datatypes"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gopkg.in/dealancer/validate.v2"
)

var _ = Describe("EquipChange Update", func() {
	var (
		testEnv = new(testVars)

		b         *xivnet.Block
		streams   *store.Streams
		d         *datasheet.Collection
		streamID  int
		subjectID uint64
		entity    *models.Entity
		generator update.Generator

		expectedClass *models.ClassJob
		expectedLevel int
	)

	BeforeEach(func() {
		*testEnv = genericSetup()
		b = testEnv.b
		streams = testEnv.streams
		d = testEnv.d
		streamID = testEnv.streamID
		subjectID = testEnv.subjectID
		entity = testEnv.entity
		generator = testEnv.generator

		d.ClassJobData = testassets.ExpectedClassJobData

		expectedClass = &models.ClassJob{
			ID: 0x12,
		}

		expectedLevel = 0x34

		equipChangeData := &datatypes.EquipChange{
			ClassJob: 0x12,
			Level:    0x34,
		}
		b.Data = equipChangeData
	})

	It("generates an update that sets the entity's ClassJob and Level", func() {
		u := generator.Generate(streamID, false, b)
		Expect(u).ToNot(BeNil())
		streamEvents, entityEvents, err := u.ModifyStore(streams)
		Expect(err).ToNot(HaveOccurred())
		Expect(streamEvents).To(BeEmpty())

		Expect(entityEvents).To(ConsistOf(models.EntityEvent{
			StreamID: streamID,
			EntityID: subjectID,
			Type: models.UpdateClass{
				ClassJob: expectedClass,
				Level:    expectedLevel,
			},
		}))

		Expect(entity.ClassJob).To(Equal(expectedClass))
		Expect(entity.Level).To(Equal(expectedLevel))

		Expect(validate.Validate(entityEvents)).To(Succeed())
		Expect(validate.Validate(streams)).To(Succeed())
	})

	Context("when the ClassJob can be found in the datasheet", func() {
		BeforeEach(func() {
			expectedClass = &models.ClassJob{
				ID:           1,
				Name:         "gladiator",
				Abbreviation: "GLA",
			}

			equipChangeData := &datatypes.EquipChange{
				ClassJob: 1,
			}
			b.Data = equipChangeData
		})

		It("populates the name and abbreviation field too", func() {
			u := generator.Generate(streamID, false, b)
			Expect(u).ToNot(BeNil())
			streamEvents, entityEvents, err := u.ModifyStore(streams)
			Expect(err).ToNot(HaveOccurred())
			Expect(streamEvents).To(BeEmpty())

			Expect(entityEvents).To(ConsistOf(models.EntityEvent{
				StreamID: streamID,
				EntityID: subjectID,
				Type: models.UpdateClass{
					ClassJob: expectedClass,
				},
			}))

			Expect(entity.ClassJob).To(Equal(expectedClass))

			Expect(validate.Validate(entityEvents)).To(Succeed())
			Expect(validate.Validate(streams)).To(Succeed())
		})
	})

	entityValidationTests(testEnv, false)
})
