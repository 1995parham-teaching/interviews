package main

import "math"

type ListNode struct {
	Val  int
	Next *ListNode
}

func mergeKLists(lists []*ListNode) *ListNode {
	var head *ListNode
	var it *ListNode

	for {

		minIndex := -1
		min := math.MaxInt

		for index, list := range lists {
			if list != nil && list.Val < min {
				min = list.Val
				minIndex = index
			}
		}

		if minIndex == -1 {
			break
		}

		lists[minIndex] = lists[minIndex].Next

		node := &ListNode{
			Val:  min,
			Next: nil,
		}

		if head == nil {
			head = node
			it = node
		} else {
			it.Next = node
			it = node
		}
	}

	return head
}
