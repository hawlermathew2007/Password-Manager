package tree

import (
	"fmt"
	"tools/data"
  "golang.org/x/exp/maps"
	"github.com/rivo/tview"
)

type Tree struct {
	Root 					*tview.TreeNode
	NodesList 		map[string]*tview.TreeNode
	ChildNodeList []data.Storage
	LogFunc				func(string)
}

func (tree *Tree) ProvideCurrentNodes() []string {
	return maps.Keys(tree.NodesList)
}

func (tree *Tree) ProvideChildsInNode(nodeName string) {

}

func (tree *Tree) HasDomain(nodeName string) bool {
	// Check if Domain exist
	_, exists := tree.NodesList[nodeName]
	return exists
}

func (tree *Tree) LoadTree() {
	for _, data := range tree.ChildNodeList {
		nodeName := data.Domain
		childName := data.Username
		if !tree.HasDomain(nodeName) {
			tree.AddNodeNChild(nodeName, childName)
		} else{
			tree.AddChild(nodeName, childName)
		}
	}
}

func (tree *Tree) ExpandNode() {
	
}

func (tree *Tree) UnexpandNode() {

}

func (tree *Tree) AddNodeNChild(nodeName string, childName string) {
	node := tview.NewTreeNode(nodeName)
	node.AddChild(tview.NewTreeNode(childName))
	tree.NodesList[nodeName] = node
	tree.Root.AddChild(node)
}

func (tree *Tree) AddNode(nodeName string) {
	node := tview.NewTreeNode(nodeName)
	tree.NodesList[nodeName] = node
	tree.Root.AddChild(node)
}

func (tree *Tree) DeleteNode() {

}

func (tree *Tree) MoveNode() {

}

func (tree *Tree) UpdateNode() {

}

func (tree *Tree) AddChild(nodeName string, childName string) {
	var _node *tview.TreeNode
	tree.Root.Walk(func(node, parent *tview.TreeNode) bool {
		if node.GetText() == nodeName {
				_node = node
				return false // Stop traversal
		}
		return true
	})
	if _node == nil {
		tree.LogFunc(fmt.Sprintf("Error creating \"%s\" account.", childName))
		return
	}
	_node.AddChild(tview.NewTreeNode(childName))
}

func (tree *Tree) DeleteChild() {

}

func (tree *Tree) MoveChild() {

}

func (tree *Tree) UpdateChild() {

}
