package balancer

const (
	//  This is the version of MDP/Client we implement
	MDPC_CLIENT = "MDPC01"

	//  This is the version of MDP/Worker we implement
	MDPW_WORKER = "MDPW01"

	//  MDP/Server commands, as strings
	SIGNAL_READY      = "\001"
	SIGNAL_REQUEST    = "\002"
	SIGNAL_REPLY      = "\003"
	SIGNAL_HEARTBEAT  = "\004"
	SIGNAL_DISCONNECT = "\005"
)

var Commands = []string{"", "READY", "REQUEST", "REPLY", "HEARTBEAT", "DISCONNECT"}
