package noip

type ApiError string

const (
	NoHostErr   ApiError = "nohost"
	BadAuthErr  ApiError = "badauth"
	BadAgentErr ApiError = "badagent"
	DonatoErr   ApiError = "!donato"
	AbuseErr    ApiError = "abuse"
	FatalErr    ApiError = "911"
)

func (n ApiError) String() string {
	switch n {
	case NoHostErr:
		return "Hostname supplied does not exist under specified account, client exit and require user to enter new login credentials before performing and additional request."
	case BadAuthErr:
		return "Invalid username password combination"
	case BadAgentErr:
		return "Client disabled. Client should exit and not perform any more updates without user intervention."
	case DonatoErr:
		return "An update request was sent including a feature that is not available to that particular user such as offline options."
	case AbuseErr:
		return "Username is blocked due to abuse. Either for not following our update specifications or disabled due to violation of the No-IP terms of service. Our terms of service can be viewed here. Client should stop sending updates."
	case FatalErr:
		return "A fatal error on our side such as a database outage. Retry the update no sooner 30 minutes."
	default:
		return "Unknown error"
	}
}
