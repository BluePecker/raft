package types

import (
    "testing"
)

func TestNodes(t *testing.T) {
    var Term uint64 = 1
    
    TestData := []uint64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
    
    nodes := &Nodes{}
    
    for _, UniqueId := range TestData {
        nodes.PushBack(Term, UniqueId)
    }
    
    if nodes.Len(Term) != len(TestData) {
        t.Fatal("error node.PushBack")
    }
    
    var Container []uint64
    
    for Next := nodes.Front(Term); Next != nil; Next = Next.Next() {
        if v, ok := Next.Value.(uint64); ok {
            Container = append(Container, v)
        }
    }
    
    if len(Container) != len(TestData) {
        t.Fatal("error nodes.front")
    }
}
