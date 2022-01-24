package client

import (
	"github.com/banafish/kvserver-cli/util"
	"log"
	"math/rand"
	"net/rpc"
	"sync"
)

const (
	TCP = "tcp"
)

type Clerk struct {
	servers   map[string]*rpc.Client
	serverIDs []string
	mu        sync.Mutex
	leaderID  string
	clientID  string
	seq       int
}

func MakeClerk(serverIDs []string) *Clerk {
	ck := new(Clerk)
	ck.servers = make(map[string]*rpc.Client)
	ck.serverIDs = serverIDs
	ck.clientID = util.GenerateClientID()
	for _, v := range serverIDs {
		ck.servers[v] = nil
	}
	ck.leaderID = serverIDs[0]
	return ck
}

func (ck *Clerk) Get(key string) string {
	ck.mu.Lock()
	defer ck.mu.Unlock()
	id := ck.leaderID
	ck.seq++
	args := GetArgs{
		ClientID: ck.clientID,
		Seq:      ck.seq,
		Key:      key,
	}

	for {
		var reply GetReply
		if err := ck.sendRPCRequest(id, "KVServerAPI.Get", &args, &reply); err != nil {
			id = ck.getServerIDRandomly()
			log.Println(err)
			continue
		}

		switch reply.Err {
		case OK:
			ck.leaderID = id
			return reply.Value
		case ErrWrongLeader:
			id = reply.LeaderID
		default:
			id = ck.getServerIDRandomly()
		}
	}
}

func (ck *Clerk) PutAppend(key string, value string, op string) {
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

	for {
		var reply PutAppendReply
		if err := ck.sendRPCRequest(id, "KVServerAPI.PutAppend", &args, &reply); err != nil {
			id = ck.getServerIDRandomly()
			log.Println(err)
			continue
		}

		switch reply.Err {
		case OK:
			ck.leaderID = id
			return
		case ErrWrongLeader:
			id = reply.LeaderID
		default:
			id = ck.getServerIDRandomly()
		}
	}
}

func (ck *Clerk) Put(key string, value string) {
	ck.PutAppend(key, value, "Put")
}
func (ck *Clerk) Append(key string, value string) {
	ck.PutAppend(key, value, "Append")
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
