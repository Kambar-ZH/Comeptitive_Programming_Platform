package datastruct

type Standings []*StandingsRow

type StandingsRow struct {
	*Participant
	ProblemResults []*ProblemResult
}
