package types

import (
    "sync"
    "container/list"
)

// 节点容器
type Nodes struct {
    // 读写锁
    locker sync.RWMutex
    // 数据存储
    member *list.List
    // 编号
    number uint64
}

// 对比时间片
func (n *Nodes) wCheck(Term uint64) bool {
    // 如果当前时间片小于容器时间片，则认为是无效数据
    if Term < n.number {
        return false
    }
    // 如果当前时间片大于容器时间片，则更新时间片并重置队列
    if Term > n.number {
        if n.member == nil {
            n.member = list.New()
        } else {
            n.member.Init()
        }
        n.number = Term
    }
    if n.member == nil {
        n.member = &list.List{}
    }
    
    return true
}

// 将无素从队尾入队列
func (n *Nodes) PushBack(Term, UniqueId uint64) {
    n.locker.Lock()
    defer n.locker.Unlock()
    if n.wCheck(Term) {
        n.member.PushBack(UniqueId)
    }
}

// 将队列从队尾入队列
func (n *Nodes) PushBackList(Term uint64, List *list.List) {
    n.locker.Lock()
    defer n.locker.Unlock()
    if n.wCheck(Term) {
        n.member.PushBackList(List)
    }
}

// 将无素从队首入队列
func (n *Nodes) PushFront(Term, UniqueId uint64) {
    n.locker.Lock()
    defer n.locker.Unlock()
    if n.wCheck(Term) {
        n.member.PushFront(UniqueId)
    }
}

// 将队列从队首入队列
func (n *Nodes) PushFrontList(Term uint64, List *list.List) {
    n.locker.Lock()
    defer n.locker.Unlock()
    if n.wCheck(Term) {
        n.member.PushBackList(List)
    }
}

// 将元素从队列中移除
func (n *Nodes) Remove(Term, UniqueId uint64) {
    n.locker.Lock()
    defer n.locker.Unlock()
    if n.wCheck(Term) {
        for Next := n.Front(Term); Next != nil; Next = Next.Next() {
            if v, ok := Next.Value.(uint64); ok && v == UniqueId {
                n.member.Remove(Next)
            }
        }
    }
}

// 计算队列长度
func (n *Nodes) Len(Term uint64) int {
    n.locker.Lock()
    defer n.locker.Unlock()
    if n.member == nil {
        return 0
    }
    return n.member.Len()
}

// 首元素
func (n *Nodes) Front(Term uint64) *list.Element {
    if n.wCheck(Term) {
        return n.member.Front()
    }
    return nil
}

// 尾元素
func (n *Nodes) Back(Term uint64) *list.Element {
    if n.wCheck(Term) {
        return n.member.Back()
    }
    return nil
}