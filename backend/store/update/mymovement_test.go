package update_test

import (
	"math"
	"time"

	"github.com/ff14wed/sibyl/backend/models"
	"github.com/ff14wed/sibyl/backend/store"
	"github.com/ff14wed/sibyl/backend/store/update"
	"github.com/ff14wed/xivnet/v2"
	"github.com/ff14wed/xivnet/v2/datatypes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gstruct"
	"github.com/onsi/gomega/types"
)

var _ = Describe("MyMovement Update", func() {
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

		matchExpectedLocation = gstruct.MatchAllFields(gstruct.Fields{
			"X":           BeNumerically("~", 100.1, 1e-4),
			"Y":           BeNumerically("~", 200.2, 1e-4),
			"Z":           BeNumerically("~", 300.3, 1e-4),
			"Orientation": BeNumerically("~", math.Pi),
			"LastUpdated": Equal(time.Unix(12, 0)),
		})

		movementData := &datatypes.MyMovement{
			Direction: 0,
			U1:        1,
			U2:        2,
			X:         100.1,
			Y:         200.2,
			Z:         300.3,
		}
		b.Data = movementData
	})

	It("generates an update that sets the entity's location", func() {
		u := generator.Generate(streamID, true, b)
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
	})

	entityValidationTests(testEnv, true)
})
