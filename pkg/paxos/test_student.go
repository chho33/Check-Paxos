package paxos

import (
	"coms4113/hw5/pkg/base"
)

// Fill in the function to lead the program to a state where A2 rejects the Accept Request of P1
func ToA2RejectP1() []func(s *base.State) bool {

    p1PreparePhase := func(s *base.State) bool {
        s1 := s.Nodes()["s1"].(*Server)
        return s1.proposer.Phase == Propose
    }

    p1PrepareConsensus := func(s *base.State) bool {
        s1 := s.Nodes()["s1"].(*Server)
        return s1.proposer.SuccessCount == 1
    }

    p3PreparePhase := func(s *base.State) bool {
        s3 := s.Nodes()["s3"].(*Server)
        return s3.proposer.Phase == Propose
    }

    p3PrepareConsensus := func(s *base.State) bool {
        s3 := s.Nodes()["s3"].(*Server)
        return s3.proposer.SuccessCount == 1
    }

    p1AcceptPhase := func(s *base.State) bool {
        s1 := s.Nodes()["s1"].(*Server)
        return s1.proposer.Phase == Accept
    }

    p3AcceptPhase := func(s *base.State) bool {
        s3 := s.Nodes()["s3"].(*Server)
        return s3.proposer.Phase == Accept
    }

    p2v3 := func(s *base.State) bool {
        s3 := s.Nodes()["s3"].(*Server)
        s2 := s.Nodes()["s2"].(*Server)
        s1 := s.Nodes()["s1"].(*Server)
        return s3.proposer.Phase == Accept && s2.v_a == "v3" && s1.proposer.N < s2.n_p
    }

    p3v3 := func(s *base.State) bool {
        s3 := s.Nodes()["s3"].(*Server)
        s2 := s.Nodes()["s2"].(*Server)
        s1 := s.Nodes()["s1"].(*Server)
        return s3.proposer.Phase == Accept && s3.v_a == "v3" &&  s2.v_a == "v3" && s1.proposer.N < s2.n_p
    }

    return []func(s *base.State) bool{
        p1PreparePhase,
        p3PreparePhase,
        p1PrepareConsensus,
        p3PrepareConsensus,
        p1AcceptPhase,
        p3AcceptPhase,
        p2v3,
        p3v3,
    }

}

// Fill in the function to lead the program to a state where a consensus is reached in Server 3.
func ToConsensusCase5() []func(s *base.State) bool {

    p3Np2 := func(s *base.State) bool {
        s3 := s.Nodes()["s3"].(*Server)
        s2 := s.Nodes()["s2"].(*Server)
        return s3.proposer.Phase == Accept && s3.proposer.N >= s2.n_p
    }

    p3DecidePhase := func(s *base.State) bool {
        s3 := s.Nodes()["s3"].(*Server)
        return s3.proposer.Phase == Decide
    }

    s3KnowConsensus := func(s *base.State) bool {
        s3 := s.Nodes()["s3"].(*Server)
        return s3.agreedValue == "v3"
    }

    return []func(s *base.State) bool{
        p3Np2,
        p3DecidePhase,
        s3KnowConsensus,
    }

}

// Fill in the function to lead the program to a state where all the Accept Requests of P1 are rejected
func NotTerminate1() []func(s *base.State) bool {
    p1PreparePhase := func(s *base.State) bool {
        s1 := s.Nodes()["s1"].(*Server)
        return s1.proposer.Phase == Propose
    }

    p1AcceptPhase := func(s *base.State) bool {
        s1 := s.Nodes()["s1"].(*Server)
        return s1.proposer.Phase == Accept
    }

    p1AcceptReject1 := func(s *base.State) bool {
        s1 := s.Nodes()["s1"].(*Server)
        return s1.proposer.Phase == Accept && s1.proposer.ResponseCount == 1 && s1.proposer.SuccessCount == 0
    }

    p1AcceptReject2 := func(s *base.State) bool {
        s1 := s.Nodes()["s1"].(*Server)
        return s1.proposer.Phase == Accept && s1.proposer.ResponseCount == 2 && s1.proposer.SuccessCount == 0
    }

    p1AcceptReject := func(s *base.State) bool {
        s1 := s.Nodes()["s1"].(*Server)
        return s1.proposer.Phase == Accept && s1.proposer.ResponseCount == 3 && s1.proposer.SuccessCount == 0
    }

    return []func(s *base.State) bool{
        p1PreparePhase,
        p1AcceptPhase,
        p1AcceptReject1,
        p1AcceptReject2,
        p1AcceptReject,
    }
}

// Fill in the function to lead the program to a state where all the Accept Requests of P3 are rejected
func NotTerminate2() []func(s *base.State) bool {
    p3PreparePhase := func(s *base.State) bool {
        s3 := s.Nodes()["s3"].(*Server)
        return s3.proposer.Phase == Propose
    }

    p3AcceptPhase := func(s *base.State) bool {
        s3 := s.Nodes()["s3"].(*Server)
        return s3.proposer.Phase == Accept
    }

    p3AcceptReject1 := func(s *base.State) bool {
        s3 := s.Nodes()["s3"].(*Server)
        return s3.proposer.Phase == Accept && s3.proposer.ResponseCount == 1 && s3.proposer.SuccessCount == 0
    }

    p3AcceptReject2 := func(s *base.State) bool {
        s3 := s.Nodes()["s3"].(*Server)
        return s3.proposer.Phase == Accept && s3.proposer.ResponseCount == 2 && s3.proposer.SuccessCount == 0
    }

    p3AcceptReject := func(s *base.State) bool {
        s3 := s.Nodes()["s3"].(*Server)
        return s3.proposer.Phase == Accept && s3.proposer.ResponseCount == 3 && s3.proposer.SuccessCount == 0
    }

    return []func(s *base.State) bool{
        p3PreparePhase,
        p3AcceptPhase,
        p3AcceptReject1,
        p3AcceptReject2,
        p3AcceptReject,
    }
}

// Fill in the function to lead the program to a state where all the Accept Requests of P1 are rejected again.
func NotTerminate3() []func(s *base.State) bool {
    return NotTerminate1()
}

// Fill in the function to lead the program to make P1 propose first, then P3 proposes, but P1 get rejects in
// Accept phase
func concurrentProposer1() []func(s *base.State) bool {
    p1PreparePhase := func(s *base.State) bool {
        s1 := s.Nodes()["s1"].(*Server)
        return s1.proposer.Phase == Propose
    }

    p1PrepareConsensus := func(s *base.State) bool {
        s1 := s.Nodes()["s1"].(*Server)
        return s1.proposer.SuccessCount == 1
    }

    p3PreparePhase := func(s *base.State) bool {
        s3 := s.Nodes()["s3"].(*Server)
        return s3.proposer.Phase == Propose
    }

    p3PrepareConsensus := func(s *base.State) bool {
        s3 := s.Nodes()["s3"].(*Server)
        return s3.proposer.SuccessCount == 1
    }

    p1AcceptPhase := func(s *base.State) bool {
        s1 := s.Nodes()["s1"].(*Server)
        return s1.proposer.Phase == Accept
    }

    p3AcceptPhase := func(s *base.State) bool {
        s3 := s.Nodes()["s3"].(*Server)
        return s3.proposer.Phase == Accept
    }

    p1AcceptReject1 := func(s *base.State) bool {
        s1 := s.Nodes()["s1"].(*Server)
        return s1.proposer.Phase == Accept && s1.proposer.ResponseCount == 1 && s1.proposer.SuccessCount == 0
    }

    p1AcceptReject2 := func(s *base.State) bool {
        s1 := s.Nodes()["s1"].(*Server)
        return s1.proposer.Phase == Accept && s1.proposer.ResponseCount == 2 && s1.proposer.SuccessCount == 0
    }

    p1AcceptReject := func(s *base.State) bool {
        s1 := s.Nodes()["s1"].(*Server)
        return s1.proposer.Phase == Accept && s1.proposer.ResponseCount == 3 && s1.proposer.SuccessCount == 0
    }

    return []func(s *base.State) bool{
        p1PreparePhase,
        p1PrepareConsensus,
        p3PreparePhase,
        p3PrepareConsensus,
        p1AcceptPhase,
        p3AcceptPhase,
        p1AcceptReject1,
        p1AcceptReject2,
        p1AcceptReject,
    }
}

// Fill in the function to lead the program continue  P3's proposal  and reaches consensus at the value of "v3".
func concurrentProposer2() []func(s *base.State) bool {
    p3PreparePhase := func(s *base.State) bool {
        s3 := s.Nodes()["s3"].(*Server)
        return s3.proposer.Phase == Propose
    }

    p3PrepareConsensus := func(s *base.State) bool {
        s3 := s.Nodes()["s3"].(*Server)
        return s3.proposer.SuccessCount == 1
    }

    p3AcceptPhase := func(s *base.State) bool {
        s3 := s.Nodes()["s3"].(*Server)
        return s3.proposer.Phase == Accept
    }

    p3AcceptConsensus := func(s *base.State) bool {
        s3 := s.Nodes()["s3"].(*Server)
        return s3.proposer.Phase == Accept && s3.proposer.SuccessCount == 1
    }

    p3DecidePhase := func(s *base.State) bool {
        s3 := s.Nodes()["s3"].(*Server)
        return s3.proposer.Phase == Decide
    }

    s3KnowConsensus := func(s *base.State) bool {
        s3 := s.Nodes()["s3"].(*Server)
        return s3.agreedValue == "v3"
    }

    return []func(s *base.State) bool{
        p3PreparePhase,
        p3PrepareConsensus,
        p3AcceptPhase,
        p3AcceptConsensus,
        p3DecidePhase,
        s3KnowConsensus,
    }
}
