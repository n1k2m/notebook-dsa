package avl

import "notebook-dsa/event"

type Node struct {
	Data        event.Event
	Left, Right *Node
	Height      int
}

type AVL struct {
	root *Node
}

func New() *AVL {
	return &AVL{}
}

// --- Методы объекта ---

func (t *AVL) Add(e event.Event) {
	t.root = insert(t.root, &e)
}

func (t *AVL) FindByDatetime(y, mo, d, h, mi int) SearchResult {
	key := event.Event{Year: y, Month: mo, Day: d, Hour: h, Minute: mi}
	res := &SearchResult{}
	search(t.root, &key, res)
	return *res
}

func (t *AVL) DeleteByDatetime(y, mo, d, h, mi int) {
	key := event.Event{Year: y, Month: mo, Day: d, Hour: h, Minute: mi}
	t.root = delete(t.root, &key)
}

func (t *AVL) Inorder(fn func(*event.Event)) {
	inorder(t.root, fn)
}

func (t *AVL) FilterPlace(substr string, fn func(*event.Event)) {
	filterPlace(t.root, substr, fn)
}

func (t *AVL) Count() int {
	return count(t.root)
}

func (t *AVL) Height() int {
	return nodeHeight(t.root)
}

func (t *AVL) Collect() []event.Event {
	var result []event.Event
	inorder(t.root, func(e *event.Event) {
		result = append(result, *e)
	})
	return result
}

// --- Внутренние рекурсивные функции ---

func nodeHeight(n *Node) int {
	if n == nil {
		return 0
	}
	return n.Height
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func updateHeight(n *Node) {
	n.Height = maxInt(nodeHeight(n.Left), nodeHeight(n.Right)) + 1
}

func balanceFactor(n *Node) int {
	if n == nil {
		return 0
	}
	return nodeHeight(n.Left) - nodeHeight(n.Right)
}

func rotateRight(y *Node) *Node {
	x := y.Left
	t2 := x.Right
	x.Right = y
	y.Left = t2
	updateHeight(y)
	updateHeight(x)
	return x
}

func rotateLeft(x *Node) *Node {
	y := x.Right
	t2 := y.Left
	y.Left = x
	x.Right = t2
	updateHeight(x)
	updateHeight(y)
	return y
}

func balance(n *Node) *Node {
	if n == nil {
		return nil
	}
	updateHeight(n)
	bf := balanceFactor(n)
	if bf > 1 {
		if balanceFactor(n.Left) < 0 {
			n.Left = rotateLeft(n.Left)
		}
		return rotateRight(n)
	}
	if bf < -1 {
		if balanceFactor(n.Right) > 0 {
			n.Right = rotateRight(n.Right)
		}
		return rotateLeft(n)
	}
	return n
}

func insert(root *Node, e *event.Event) *Node {
	if root == nil {
		return &Node{Data: *e, Height: 1}
	}
	c := event.Cmp(e, &root.Data)
	if c < 0 {
		root.Left = insert(root.Left, e)
	} else if c > 0 {
		root.Right = insert(root.Right, e)
	} else {
		return root
	}
	return balance(root)
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
	return balance(root)
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
