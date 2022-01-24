package client

const (
	OK                    = ""
	ErrNoKey              = "ErrNoKey"
	ErrWrongLeader        = "ErrWrongLeader"
	ErrFailReachAgreement = "ErrFailReachAgreement"

	OpTypeGet    OpType = "Get"
	OpTypePut    OpType = "Put"
	OpTypeAppend OpType = "Append"
)

type Err string

// Put or Append
type PutAppendArgs struct {
	ClientID string
	Seq      int
	Key      string
	Value    string
	Op       OpType // "Put" or "Append"
}

type PutAppendReply struct {
	Err      Err
	LeaderID string
}

type GetArgs struct {
	ClientID string
	Seq      int
	Key      string
}

type GetReply struct {
	Err      Err
	Value    string
	LeaderID string
}

type OpType string

type GetRaftStatArgs struct {
	IsPrintLog bool
}

type GetRaftStatReply struct {
	Stat string
}

type GetServerStatArgs struct {
}

type GetServerStatReply struct {
	Stat string
}
