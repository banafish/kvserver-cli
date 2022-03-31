package client

import (
	"fmt"
	"github.com/banafish/kvserver-cli/util"
	"math/rand"
	"net/rpc"
	"sync"
)

const (
	TCP = "tcp"
)

type Clerk struct {
	servers    map[string]*rpc.Client
	serverIDs  []string
	mu         sync.Mutex
	leaderID   string
	clientID   string
	seq        int
	retryCount int
}

func MakeClerk(serverIDs []string) *Clerk {
	ck := new(Clerk)
	ck.servers = make(map[string]*rpc.Client)
	ck.serverIDs = serverIDs
	ck.retryCount = 10
	ck.clientID = util.GenerateClientID()
	for _, v := range serverIDs {
		ck.servers[v] = nil
	}
	ck.leaderID = serverIDs[0]
	return ck
}

func (ck *Clerk) Get(key string) (string, error) {
	ck.mu.Lock()
	defer ck.mu.Unlock()
	id := ck.leaderID
	args := GetArgs{
		ClientID: ck.clientID,
		Key:      key,
	}

	for i := 0; i < ck.retryCount; i++ {
		var reply GetReply
		if err := ck.sendRPCRequest(id, "KVServerAPI.Get", &args, &reply); err != nil {
			id = ck.getServerIDRandomly()
			//log.Println(err)
			continue
		}

		switch reply.Err {
		case OK:
			ck.leaderID = id
			return reply.Value, nil
		case ErrWrongLeader:
			if reply.LeaderID == "" {
				id = ck.getServerIDRandomly()
			} else {
				id = reply.LeaderID
			}
		default:
			id = ck.getServerIDRandomly()
		}
	}
	return "", fmt.Errorf(ErrRetryCountReached)
}

func (ck *Clerk) PutAppend(key string, value string, op string) error {
	ck.mu.Lock()
	defer ck.mu.Unlock()
	id := ck.leaderID
	ck.seq++
	args := PutAppendArgs{
		ClientID: ck.clientID,
		Seq:      ck.seq,
		Op:       OpType(op),
		Key:      key,
		Value:    value,
	}

	for i := 0; i < ck.retryCount; i++ {
		var reply PutAppendReply
		if err := ck.sendRPCRequest(id, "KVServerAPI.PutAppend", &args, &reply); err != nil {
			id = ck.getServerIDRandomly()
			//log.Println(err)
			continue
		}

		switch reply.Err {
		case OK:
			ck.leaderID = id
			return nil
		case ErrWrongLeader:
			if reply.LeaderID == "" {
				id = ck.getServerIDRandomly()
			} else {
				id = reply.LeaderID
			}
		default:
			id = ck.getServerIDRandomly()
		}
	}
	return fmt.Errorf(ErrRetryCountReached)
}

func (ck *Clerk) Put(key string, value string) error {
	return ck.PutAppend(key, value, "Put")
}
func (ck *Clerk) Append(key string, value string) error {
	return ck.PutAppend(key, value, "Append")
}

func (ck *Clerk) GetRaftStat(isWithLog bool, addr string) string {
	args := GetRaftStatArgs{IsPrintLog: isWithLog}
	var reply GetRaftStatReply
	if err := ck.sendRPCRequest(addr, "RaftAPI.GetRaftStat", &args, &reply); err != nil {
		return err.Error()
	}
	return reply.Stat
}

func (ck *Clerk) GetServerStat(addr string) string {
	var args GetServerStatArgs
	var reply GetServerStatReply
	if err := ck.sendRPCRequest(addr, "KVServerAPI.GetServerStat", &args, &reply); err != nil {
		return err.Error()
	}
	return reply.Stat
}

func (ck *Clerk) sendRPCRequest(serverID string, serviceMethod string, args interface{}, reply interface{}) error {
	s, err := ck.getServer(serverID)
	if err != nil {
		return err
	}
	if err = s.Call(serviceMethod, args, reply); err != nil {
		// 清空连接，重新连接
		ck.servers[serverID] = nil
	}
	return err
}

func (ck *Clerk) getServer(serverID string) (*rpc.Client, error) {
	if ck.servers[serverID] == nil {
		c, err := rpc.DialHTTP(TCP, serverID)
		if err != nil {
			return nil, err
		}
		ck.servers[serverID] = c
		ck.serverIDs = append(ck.serverIDs, serverID)
	}
	return ck.servers[serverID], nil
}

func (ck *Clerk) getServerIDRandomly() string {
	n := rand.Int()
	return ck.serverIDs[n%len(ck.serverIDs)]
}
