package paxos

import (
	"coms4113/hw5/pkg/base"
)

const (
	Propose = "propose"
	Accept  = "accept"
	Decide  = "decide"
)

type Proposer struct {
	N             int
	Phase         string
	N_a_max       int
	V             interface{}
	SuccessCount  int
	ResponseCount int
	// To indicate if response from peer is received, should be initialized as []bool of len(server.peers)
	Responses []bool
	// Use this field to check if a message is latest.
	SessionId int

	// in case node will propose again - restore initial value
	InitialValue interface{}
}

type ServerAttribute struct {
	peers []base.Address
	me    int

	// Paxos parameter
	n_p int
	n_a int
	v_a interface{}

	// final result
	agreedValue interface{}

	// Propose parameter
	proposer Proposer

	// retry
	timeout *TimeoutTimer
}

type Server struct {
	base.CoreNode
	ServerAttribute
}

func NewServer(peers []base.Address, me int, proposedValue interface{}) *Server {
	response := make([]bool, len(peers))
	return &Server{
		CoreNode: base.CoreNode{},
		ServerAttribute: ServerAttribute{
			peers: peers,
			me:    me,
			proposer: Proposer{
				InitialValue: proposedValue,
				Responses:    response,
                N_a_max: 0,
			},
			timeout: &TimeoutTimer{},
		},
	}
}

func (server *Server) ResetCount() {
    server.proposer.SuccessCount = 0
    server.proposer.ResponseCount = 0
}

func (newServer *Server) HandlePropReq(message base.Message) []base.Node {
    from := message.From()
    to := message.To()

    if message.(*ProposeRequest).N <= newServer.n_p {
      result := &ProposeResponse {
        CoreMessage: base.MakeCoreMessage(to, from),
        Ok:          false,
        N_p:         newServer.n_p,
        N_a:         newServer.n_a,
        V_a:         newServer.v_a, //nil for 1,2; v1 for 0
        SessionId:   message.(*ProposeRequest).SessionId,
      }
      newServer.Response = []base.Message{
        result,
      }
      return []base.Node{newServer}
    }

    v_a := newServer.v_a
    n_a := newServer.n_a
    newServer.n_p = message.(*ProposeRequest).N


    result := &ProposeResponse {
      CoreMessage: base.MakeCoreMessage(to, from),
      Ok:          true,
      N_p:         newServer.n_p,
      N_a:         n_a,
      V_a:         v_a, //nil for 1,2; v1 for 0
      SessionId:   message.(*ProposeRequest).SessionId,
    }

    newServer.Response = []base.Message{
      result,
    }

    newServer.proposer.Responses = []bool{false, false, false}
    return []base.Node{newServer}
}

func (newServer *Server) HandlePropRes(message base.Message) []base.Node {
    from := message.From()

    // examine SessionId 
    newServer.proposer.ResponseCount++
    if message.(*ProposeResponse).SessionId != newServer.proposer.SessionId {
      return []base.Node{}
    }

    // examine Ok
    if !message.(*ProposeResponse).Ok {
      // Resp N_p bigger;
      //newServer.proposer.N = message.(*ProposeResponse).N_p

      if message.(*ProposeResponse).N_a > newServer.proposer.N_a_max {
        newServer.proposer.N_a_max = message.(*ProposeResponse).N_a
        newServer.proposer.V = message.(*ProposeResponse).V_a
      }
      return []base.Node{newServer}
    }
    newServer.proposer.SuccessCount++


    if len(newServer.proposer.Responses)==0 {
      newServer.proposer.Responses = []bool{false, false, false}
    }


    v := newServer.proposer.V
    // examine N_a, V_a 
    if message.(*ProposeResponse).N_a > newServer.proposer.N_a_max {
      newServer.proposer.N_a_max = message.(*ProposeResponse).N_a
      v = message.(*ProposeResponse).V_a
    }

    if !base.IsNil(newServer.agreedValue) {
      v = newServer.agreedValue
    }

    //if message.(*ProposeResponse).N_a > newServer.n_a {
    //  v = message.(*ProposeResponse).V_a
    //}

    index := newServer.getIndex(from)
    newServer.proposer.Responses[index] = true
    newServer.proposer.V = v

    if newServer.proposer.SuccessCount< 2 {
        return []base.Node{newServer}
    }

    n_p := newServer.proposer.N
    anotherServer := newServer.copy()
    anotherServer.Response = []base.Message{
      &AcceptRequest {
        CoreMessage: base.MakeCoreMessage(newServer.peers[newServer.me], newServer.peers[0]),
        N:           n_p,
        V:           v,
        SessionId:   newServer.proposer.SessionId,
      },
      &AcceptRequest {
        CoreMessage: base.MakeCoreMessage(newServer.peers[newServer.me], newServer.peers[1]),
        N:           n_p,
        V:           v,
        SessionId:   newServer.proposer.SessionId,
      },
      &AcceptRequest {
        CoreMessage: base.MakeCoreMessage(newServer.peers[newServer.me], newServer.peers[2]),
        N:           n_p,
        V:           v,
        SessionId:   newServer.proposer.SessionId,
      },
    }

    anotherServer.proposer.V = v
    anotherServer.proposer.Phase = Accept
    anotherServer.proposer.ResponseCount = 0
    anotherServer.proposer.SuccessCount = 0
    anotherServer.proposer.Responses = []bool{false, false, false}

    response := []base.Node{anotherServer, newServer}

    return response
}

func (newServer *Server) HandleAcptReq(message base.Message) []base.Node {
    from := message.From()
    to := message.To()

    if message.(*AcceptRequest).N < newServer.n_p {
      result := &AcceptResponse {
        CoreMessage: base.MakeCoreMessage(to, from),
        Ok:          false,
        SessionId:   message.(*AcceptRequest).SessionId,
      }
      newServer.Response = []base.Message{
        result,
      }
      return []base.Node{newServer}
    }

    newServer.n_p = message.(*AcceptRequest).N
    newServer.n_a = message.(*AcceptRequest).N
    newServer.v_a = message.(*AcceptRequest).V

    result := &AcceptResponse {
      CoreMessage: base.MakeCoreMessage(to, from),
      Ok:          true,
      N_p:         newServer.n_p,
      SessionId:   message.(*AcceptRequest).SessionId,
    }

    newServer.Response = []base.Message{
      result,
    }
    return []base.Node{newServer}
}

func (newServer *Server) HandleAcptRes(message base.Message) []base.Node {
    from := message.From()

    newServer.proposer.ResponseCount++
    if message.(*AcceptResponse).SessionId != newServer.proposer.SessionId {
      return []base.Node{}
    }

    if !message.(*AcceptResponse).Ok {
      return []base.Node{newServer}
    }
    newServer.proposer.SuccessCount++

    if len(newServer.proposer.Responses)==0 {
      newServer.proposer.Responses = []bool{false, false, false}
    }
    index := newServer.getIndex(from)
    newServer.proposer.Responses[index] = true

    response := []base.Node{newServer}
    if newServer.proposer.SuccessCount < 2 {
        return response
    }

    v := newServer.proposer.V

    anotherServer := newServer.copy()
    anotherServer.Response = []base.Message{
      &DecideRequest {
        CoreMessage: base.MakeCoreMessage(newServer.peers[newServer.me], newServer.peers[0]),
        V:           v,
        SessionId:   message.(*AcceptResponse).SessionId,
      },
      &DecideRequest {
        CoreMessage: base.MakeCoreMessage(newServer.peers[newServer.me], newServer.peers[1]),
        V:           v,
        SessionId:   message.(*AcceptResponse).SessionId,
      },
      &DecideRequest{
        CoreMessage: base.MakeCoreMessage(newServer.peers[newServer.me], newServer.peers[2]),
        V:           v,
        SessionId:   message.(*AcceptResponse).SessionId,
      },
    }

    anotherServer.agreedValue = v

    anotherServer.proposer = Proposer {
      N:            message.(*AcceptResponse).N_p,
      Phase:        Decide,
      V:            v,
      SessionId:    message.(*AcceptResponse).SessionId,
      InitialValue: newServer.proposer.InitialValue,
      Responses:    []bool{false, false, false},
      N_a_max: newServer.proposer.N_a_max,
      ResponseCount: 0,
      SuccessCount: 0,
    }

    return []base.Node{anotherServer, response[0]}
}

func (server *Server) MessageHandler(message base.Message) []base.Node {
	//TODO: implement it
	//panic("implement me")
    //While implementing the acceptor states in MessageHandler, I found that there is no phase INFO in ProposeRequest. Then how do we know which state it is now? Also, there is no Va in ProposeRequest, then how do we update the Va for the acceptor?
    //Andrew answered the first question. For your second question, a propose is expected to get the v_a at ProposeResponse.

    newServer := server.copy()
    //newServer.UpdateNV()

    from := message.From()
    to := message.To()

    if server.peers[server.me] != to || !base.IsNil(server.agreedValue) {
      //return []base.Node{newServer}
      return []base.Node{}
    }

    // This happens in NextState 

    _, propReq := message.(*ProposeRequest)
    if propReq {
      return newServer.HandlePropReq(message)
    }

    _, propRes := message.(*ProposeResponse)
    if propRes {
      return newServer.HandlePropRes(message)
    }

    _, acptReq := message.(*AcceptRequest)
    if acptReq {
      return newServer.HandleAcptReq(message)
    }

    _, acptRes := message.(*AcceptResponse)
    if acptRes {
      return newServer.HandleAcptRes(message)
    }

    _, decReq := message.(*DecideRequest)
    if decReq {
      newServer.agreedValue = message.(*DecideRequest).V
      newServer.Response = []base.Message {
        &DecideResponse {
          CoreMessage: base.MakeCoreMessage(to, from),
          Ok:          true,
          SessionId  : message.(*DecideRequest).SessionId,
        },
      }
    }

    return []base.Node{newServer}
}

func (server *Server) getIndex(target base.Address) int {
    indices := map[base.Address] int{}
    for i, address := range(server.peers) {
      indices[address] = i
    }
    return indices[target]
}

// To start a new round of Paxos.
func (server *Server) StartPropose() {
	//TODO: implement it

    if base.IsNil(server.proposer.InitialValue) {
        return
    }

    n := server.proposer.N + 1
    sessionid := server.proposer.SessionId + 1

    server.Response = []base.Message{
        &ProposeRequest{
            CoreMessage: base.MakeCoreMessage(server.peers[server.me], server.peers[0]),
            N:           n,
            SessionId:   sessionid,
        },
        &ProposeRequest{
            CoreMessage: base.MakeCoreMessage(server.peers[server.me], server.peers[1]),
            N:           n,
            SessionId:   sessionid,
        },
        &ProposeRequest{
            CoreMessage: base.MakeCoreMessage(server.peers[server.me], server.peers[2]),
            N:           n,
            SessionId:   sessionid,
        },
    }

	response := make([]bool, len(server.peers))
    server.proposer =  Proposer {
		InitialValue: server.proposer.InitialValue,
		Responses:    response,
        N_a_max: 0,
	}
    server.proposer.N = n
    server.proposer.SessionId = sessionid
    //server.proposer.V = server.peers[server.me]
    if base.IsNil(server.proposer.V) {
        server.proposer.V = server.proposer.InitialValue
    }
    if !base.IsNil(server.agreedValue) {
        server.proposer.V = server.agreedValue
    }
    server.proposer.Phase = Propose
	//server.proposer.SuccessCount = 0
	//server.proposer.ResponseCount = 0
    //server.proposer.Responses = []bool{false, false, false}
}

// Returns a deep copy of server node
func (server *Server) copy() *Server {
	response := make([]bool, len(server.peers))
	for i, flag := range server.proposer.Responses {
		response[i] = flag
	}

	var copyServer Server
	copyServer.me = server.me
	// shallow copy is enough, assuming it won't change
	copyServer.peers = server.peers
	copyServer.n_a = server.n_a
	copyServer.n_p = server.n_p
	copyServer.v_a = server.v_a
	copyServer.agreedValue = server.agreedValue
	copyServer.proposer = Proposer{
		N:             server.proposer.N,
		Phase:         server.proposer.Phase,
		N_a_max:       server.proposer.N_a_max,
		V:             server.proposer.V,
		SuccessCount:  server.proposer.SuccessCount,
		ResponseCount: server.proposer.ResponseCount,
		Responses:     response,
		InitialValue:  server.proposer.InitialValue,
		SessionId:     server.proposer.SessionId,
	}

	// doesn't matter, timeout timer is state-less
	copyServer.timeout = server.timeout

	return &copyServer
}

func (server *Server) NextTimer() base.Timer {
	return server.timeout
}

// A TimeoutTimer tick simulates the situation where a proposal procedure times out.
// It will close the current Paxos round and start a new one if no consensus reached so far,
// i.e. the server after timer tick will reset and restart from the first phase if Paxos not decided.
// The timer will not be activated if an agreed value is set.
func (server *Server) TriggerTimer() []base.Node {
	if server.timeout == nil {
		return nil
	}

	subNode := server.copy()
	subNode.StartPropose()

	return []base.Node{subNode}
}

func (server *Server) Attribute() interface{} {
	return server.ServerAttribute
}

func (server *Server) Copy() base.Node {
	return server.copy()
}

func (server *Server) Hash() uint64 {
	return base.Hash("paxos", server.ServerAttribute)
}

func (server *Server) Equals(other base.Node) bool {
	otherServer, ok := other.(*Server)

	if !ok || server.me != otherServer.me ||
		server.n_p != otherServer.n_p || server.n_a != otherServer.n_a || server.v_a != otherServer.v_a ||
		(server.timeout == nil) != (otherServer.timeout == nil) {
		return false
	}

	if server.proposer.N != otherServer.proposer.N || server.proposer.V != otherServer.proposer.V ||
		server.proposer.N_a_max != otherServer.proposer.N_a_max || server.proposer.Phase != otherServer.proposer.Phase ||
		server.proposer.InitialValue != otherServer.proposer.InitialValue ||
		server.proposer.SuccessCount != otherServer.proposer.SuccessCount ||
		server.proposer.ResponseCount != otherServer.proposer.ResponseCount {
		return false
	}

	for i, response := range server.proposer.Responses {
		if response != otherServer.proposer.Responses[i] {
			return false
		}
	}

	return true
}

func (server *Server) Address() base.Address {
	return server.peers[server.me]
}

func (server *Server) GetProposer() Proposer {
	return server.proposer
}
