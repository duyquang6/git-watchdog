package customtypes

type ResultStatus uint8

const (
	QUEUING ResultStatus = iota
	QUEUED
	IN_PROGRESS
	SUCCESS
	FAILURE
)

func (s ResultStatus) String() string {
	switch s {
	case QUEUING:
		return "QUEUING"
	case QUEUED:
		return "QUEUED"
	case IN_PROGRESS:
		return "IN_PROGRESS"
	case SUCCESS:
		return "SUCCESS"
	case FAILURE:
		return "FAILURE"
	}
	return "UNKNOWN"
}
