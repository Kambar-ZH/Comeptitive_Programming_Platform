package consts

type ParticipantType int

const (
	CONTESTANT ParticipantType = iota
	VIRTUAL
)

func (p ParticipantType) String() string {
	return map[ParticipantType]string{
		CONTESTANT: "CONTESTANT",
		VIRTUAL:    "VIRTUAL",
	}[p]
}
