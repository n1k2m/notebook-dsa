package bst

import "notebook-dsa/event"

type Node struct {
	Data        event.Event
	Left, Right *Node
}

type BST struct {
	root *Node
}

func New() *BST {
	return &BST{}
}

// --- Методы объекта ---

func (t *BST) Add(e event.Event) {
	t.root = insert(t.root, &e)
}

func (t *BST) FindByDatetime(y, mo, d, h, mi int) SearchResult {
	key := event.Event{Year: y, Month: mo, Day: d, Hour: h, Minute: mi}
	res := &SearchResult{}
	search(t.root, &key, res)
	return *res
}

func (t *BST) DeleteByDatetime(y, mo, d, h, mi int) {
	key := event.Event{Year: y, Month: mo, Day: d, Hour: h, Minute: mi}
	t.root = delete(t.root, &key)
}

func (t *BST) Inorder(fn func(*event.Event)) {
	inorder(t.root, fn)
}

func (t *BST) FilterPlace(substr string, fn func(*event.Event)) {
	filterPlace(t.root, substr, fn)
}

func (t *BST) Count() int {
	return count(t.root)
}

func (t *BST) Height() int {
	return height(t.root)
}

func (t *BST) Collect() []event.Event {
	var result []event.Event
	inorder(t.root, func(e *event.Event) {
		result = append(result, *e)
	})
	return result
}

// --- Внутренние рекурсивные функции ---

func insert(root *Node, e *event.Event) *Node {
	if root == nil {
		return &Node{Data: *e}
	}
	c := event.Cmp(e, &root.Data)
	if c < 0 {
		root.Left = insert(root.Left, e)
	} else {
		root.Right = insert(root.Right, e)
	}
	return root
}

func minNode(n *Node) *Node {
	for n.Left != nil {
		n = n.Left
	}
	return n
}

func delete(root *Node, e *event.Event) *Node {
	if root == nil {
		return nil
	}
	c := event.Cmp(e, &root.Data)
	if c < 0 {
		root.Left = delete(root.Left, e)
	} else if c > 0 {
		root.Right = delete(root.Right, e)
	} else {
		if root.Left == nil {
			return root.Right
		}
		if root.Right == nil {
			return root.Left
		}
		succ := minNode(root.Right)
		root.Data = succ.Data
		root.Right = delete(root.Right, &succ.Data)
	}
	return root
}

type SearchResult struct {
	Events []*event.Event
	Ops    int64
}

func search(root *Node, key *event.Event, res *SearchResult) {
	if root == nil {
		return
	}
	res.Ops++
	c := event.Cmp(key, &root.Data)
	if c == 0 {
		res.Events = append(res.Events, &root.Data)
		search(root.Left, key, res)
		search(root.Right, key, res)
	} else if c < 0 {
		search(root.Left, key, res)
	} else {
		search(root.Right, key, res)
	}
}

func inorder(root *Node, fn func(*event.Event)) {
	if root == nil {
		return
	}
	inorder(root.Left, fn)
	fn(&root.Data)
	inorder(root.Right, fn)
}

func filterPlace(root *Node, substr string, fn func(*event.Event)) {
	if root == nil {
		return
	}
	filterPlace(root.Left, substr, fn)
	if event.PlaceContains(&root.Data, substr) {
		fn(&root.Data)
	}
	filterPlace(root.Right, substr, fn)
}

func count(root *Node) int {
	if root == nil {
		return 0
	}
	return 1 + count(root.Left) + count(root.Right)
}

func height(root *Node) int {
	if root == nil {
		return 0
	}
	lh := height(root.Left)
	rh := height(root.Right)
	if lh > rh {
		return lh + 1
	}
	return rh + 1
}
