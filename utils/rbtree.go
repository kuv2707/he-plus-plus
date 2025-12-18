package utils

/**
- The root of the tree is always black.
- Red nodes cannot have red children.
- Every path from a node to its descendant NULL nodes has the same number of black nodes.
- All leaves are black.

These ensure that longest path <= 2*shortest path

ReferencE: https://github.com/zarif98sjs/RedBlackTree-An-Intuitive-Approach
*/

type RBTreeMap[K comparable, V any] struct {
	root       *RBNode[K, V]
	comparator func(a, b K) int // -1, 0 1
	keyNode    map[K]*RBNode[K, V]
	size       int
	nilNode    *RBNode[K, V]
}

func (rbt *RBTreeMap[K, V]) Size() int {
	return rbt.size
}

func MakeRBTreeMap[K comparable, V any]() RBTreeMap[K, V] {
	nilN := RBNode[K, V]{}
	nilN.children = []*RBNode[K, V]{&nilN, &nilN}
	return RBTreeMap[K, V]{
		root:       &nilN,
		size:       0,
		comparator: func(a, b K) int { return 0 },
		keyNode:    make(map[K]*RBNode[K, V]),
		nilNode:    &nilN,
	}
}

// a>b => return 1
func (rbt *RBTreeMap[K, V]) setComparator(f func(a, b K) int) {
	rbt.comparator = f
}

type RBColor = bool

var RED = true
var BLACK = false

var LEFT = 0
var RIGHT = 1

type RBNode[K comparable, V any] struct {
	color    RBColor
	key      K
	value    V
	children []*RBNode[K, V]
}

func (rbt *RBTreeMap[K, V]) isRed(n *RBNode[K, V]) bool {
	if n == rbt.nilNode {
		return false
	}
	return n.color == RED
}

func (rbt *RBTreeMap[K, V]) flipColor(n *RBNode[K, V]) {
	n.color = !n.color
	if n.children[LEFT] != rbt.nilNode {
		n.children[LEFT].color = !n.children[LEFT].color
	}
	if n.children[RIGHT] != rbt.nilNode {
		n.children[RIGHT].color = !n.children[RIGHT].color
	}

}

func (n *RBNode[K, V]) rotate(dir int) *RBNode[K, V] {
	t := n.children[1-dir]
	n.children[1-dir] = t.children[dir]
	t.children[dir] = n

	t.color = n.color
	n.color = RED
	return t
}

func (n *RBNode[K, V]) doubleRotate(dir int) *RBNode[K, V] {
	// 	   B(G)
	//    /
	// R(P)
	//    \
	//     R(X)
	n.children[1-dir] = n.children[1-dir].rotate(1 - dir)
	// 	    B(G)
	//        /
	//      R(X)
	//      /
	//    R(P)
	return n.rotate(dir)
}

func (rbt *RBTreeMap[K, V]) Insert(key K, value V) {
	if n, exists := rbt.keyNode[key]; exists {
		// essentially an update
		n.value = value
		return
	}
	rbt.root = rbt.nodeInsert(rbt.root, key, value)
	rbt.root.color = BLACK
	rbt.size++
}

func (rbt *RBTreeMap[K, V]) nodeInsert(n *RBNode[K, V], key K, value V) *RBNode[K, V] {
	if n == rbt.nilNode {
		newn := &RBNode[K, V]{key: key, value: value, color: RED, children: []*RBNode[K, V]{rbt.nilNode, rbt.nilNode}}
		rbt.keyNode[key] = newn
		return newn
	}
	dir := rbt.comparator(key, n.key) // key > n.key => true === RIGHT
	if dir == -1 {
		dir = 0
	}
	n.children[dir] = rbt.nodeInsert(n.children[dir], key, value)
	return rbt.insertFixUp(n, dir)
}

func (rbt *RBTreeMap[K, V]) insertFixUp(n *RBNode[K, V], dir int) *RBNode[K, V] {
	// X = inserted node (red)
	if rbt.isRed(n.children[dir]) { // red parent (wrt n ka grandchild)
		if rbt.isRed(n.children[1-dir]) { // red uncle
			if rbt.isRed(n.children[dir].children[dir]) || rbt.isRed(n.children[dir].children[1-dir]) { // maybe unnecessary
				// CASE 1: grandchild itself red, so flip cols of grandpa father uncle
				//   B(G)
				// 	/   \
				// R(P)   R(U)
				// 	|
				// R(X)
				rbt.flipColor(n)
				// now its:
				//   R(G)
				// 	/   \
				// B(P)   B(U)
				// 	|
				// R(X)
			}
		} else { // black uncle
			if rbt.isRed(n.children[dir].children[dir]) {
				//   B(G)
				// 	/
				// 	R(P)
				// 	/
				// R(X)
				n = n.rotate(1 - dir)
				//  B(P)
				// /   \
				// R(X)   R(G)
			} else if rbt.isRed(n.children[dir].children[1-dir]) {
				// 	   B(G)
				//    /
				// R(P)
				//    \
				//     R(X)
				n = n.doubleRotate(1 - dir)
				// 	B(X)
				//    /   \
				// R(P)   R(G)
			}
		}
	}
	return n
}

func (rbt *RBTreeMap[K, V]) Delete(key K) (V, bool) {
	if _, ex := rbt.keyNode[key]; !ex {
		var zero V
		return zero, false
	}
	ok := false
	nr, val := rbt.nodeDelete(rbt.root, key, &ok)
	rbt.root = nr
	if nr != rbt.nilNode {
		nr.color = BLACK
	}
	delete(rbt.keyNode, key)
	rbt.size--;
	return val, true
}

func (rbt *RBTreeMap[K, V]) nodeDelete(n *RBNode[K, V], key K, ok *bool) (*RBNode[K, V], V) {
	var zero V
	if n == rbt.nilNode {
		*ok = true
		return rbt.nilNode, zero
	}
	if n.key == key {
		if n.children[LEFT] == rbt.nilNode || n.children[RIGHT] == rbt.nilNode {
			// one or less children
			ch := n.children[LEFT]
			if ch == rbt.nilNode {
				ch = n.children[RIGHT]
			}
			if rbt.isRed(n) {
				// can simply delete a red node
				*ok = true
			} else if rbt.isRed(ch) {
				ch.color = BLACK
				*ok = true
			}
			// delete this node
			return ch, n.value
		} else {
			// both children exist
			// replace with inorder predecessor
			pred := getInorderPredecessor(n.children[LEFT], rbt.nilNode)
			n.key = pred.key
			n.value = pred.value
			key = pred.key // new wanted node for deletion
			rbt.keyNode[n.key] = n

		}
	}
	dir := rbt.comparator(key, n.key)
	if dir == -1 {
		dir = 0
	}
	del, val := rbt.nodeDelete(n.children[dir], key, ok)
	n.children[dir] = del
	if *ok {
		return n, val
	} else {
		return rbt.deleteFixUp(n, dir, ok), val
	}
}

func (rbt *RBTreeMap[K, V]) deleteFixUp(n *RBNode[K, V], dir int, ok *bool) *RBNode[K, V] {
	parent := n
	sibling := n.children[1-dir]

	if rbt.isRed(sibling) {
		// convert to black sibling
		n = n.rotate(dir)
		sibling = parent.children[1-dir]
	}
	if sibling != rbt.nilNode {
		// 1. black sibling with black children
		if !rbt.isRed(sibling.children[LEFT]) && !rbt.isRed(sibling.children[RIGHT]) {
			if rbt.isRed(parent) {
				*ok = true
			}
			parent.color = BLACK
			sibling.color = RED
		} else {
			// black sibling with non-black children
			initParentColor := parent.color
			isRedSiblingBlackConv := !(n == parent)
			if rbt.isRed(sibling.children[1-dir]) {
				parent = parent.rotate(dir)
			} else {
				parent = parent.doubleRotate(dir)
			}
			parent.color = initParentColor
			parent.children[LEFT].color = BLACK
			parent.children[RIGHT].color = BLACK

			if isRedSiblingBlackConv {
				n.children[dir] = parent
			} else {
				n = parent
			}
			*ok = true
		}
	}
	return n
}

func getInorderPredecessor[K comparable, V any](n *RBNode[K, V], nilNode *RBNode[K, V]) *RBNode[K, V] {
	var ret *RBNode[K, V] = nilNode
	for n != nilNode {
		ret = n
		n = n.children[RIGHT]
	}
	return ret
}
