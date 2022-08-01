function test() {
    echo $1 
    go test -run $1 
}

test TestBasic
test TestBasic2
test TestBfs1
test TestBfs2
test TestBfs3
test TestInvariant
test TestPartition1
test TestPartition2

# student
test TestCase5Failures
test TestNotTerminate
test TestConcurrentProposer
test TestFailChecks
