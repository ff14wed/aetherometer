package update_test

import (
	"math"

	"github.com/ff14wed/aetherometer/core/models"
	"github.com/ff14wed/aetherometer/core/store"
	"github.com/ff14wed/aetherometer/core/store/update"
	"github.com/ff14wed/xivnet/v3"
	"github.com/ff14wed/xivnet/v3/datatypes"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gstruct"
	"github.com/onsi/gomega/types"
	"gopkg.in/dealancer/validate.v2"
)

var _ = Describe("Movement Update", func() {
	var (
		testEnv = new(testVars)

		b         *xivnet.Block
		streams   *store.Streams
		streamID  int
		subjectID uint64
		entity    *models.Entity
		generator update.Generator

		matchExpectedLocation types.GomegaMatcher
	)

	BeforeEach(func() {
		*testEnv = genericSetup()
		b = testEnv.b
		streams = testEnv.streams
		streamID = testEnv.streamID
		subjectID = testEnv.subjectID
		entity = testEnv.entity
		generator = testEnv.generator

		matchExpectedLocation = gstruct.PointTo(gstruct.MatchAllFields(gstruct.Fields{
			"X":           BeNumerically("~", 100.1, 0.02),
			"Y":           BeNumerically("~", 200.2, 0.02),
			"Z":           BeNumerically("~", 300.3, 0.02),
			"Orientation": BeNumerically("~", math.Pi),
			"LastUpdated": Equal(b.Time),
		}))

		movementData := &datatypes.Movement{
			Direction: 128,
		}
		movementData.Position.X.SetFloat(100.1)
		movementData.Position.Y.SetFloat(200.2)
		movementData.Position.Z.SetFloat(300.3)
		b.Data = movementData
	})

	It("generates an update that sets the entity's location", func() {
		u := generator.Generate(streamID, false, b)
		Expect(u).ToNot(BeNil())
		streamEvents, entityEvents, err := u.ModifyStore(streams)
		Expect(err).ToNot(HaveOccurred())
		Expect(streamEvents).To(BeEmpty())

		Expect(entityEvents).To(HaveLen(1))
		Expect(entityEvents[0].StreamID).To(Equal(streamID))
		Expect(entityEvents[0].EntityID).To(Equal(subjectID))
		eventType, assignable := entityEvents[0].Type.(models.UpdateLocation)
		Expect(assignable).To(BeTrue())
		Expect(eventType.Location).To(matchExpectedLocation)

		Expect(entity.Location).To(matchExpectedLocation)

		Expect(validate.Validate(entityEvents)).To(Succeed())
		Expect(validate.Validate(streams)).To(Succeed())
	})

	entityValidationTests(testEnv, false)
})
