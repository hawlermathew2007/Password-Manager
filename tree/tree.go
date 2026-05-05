package tree

import (
	// "fmt"
	// "tools/account"
	"tools/data"
	"tools/debug"
	"github.com/google/uuid"
	"github.com/rivo/tview"
	"golang.org/x/exp/maps"
)

// Responsible for the Tree UI & Provide details about Accounts
type LogFuncTree func(logType debug.LogType, context debug.LogContext)

type Tree struct {
	Root 								*tview.TreeNode
	NodesList 					map[string]*tview.TreeNode
	ChildsList 					map[*tview.TreeNode]uuid.UUID // Domain ID
	ChildNodeList 			[]data.DomainDetails
	LogFunc							LogFuncTree
	CreateChildContent	func(uuid.UUID) *tview.Frame
	SelectChild					func(string)
	// ProvideChildDetails func(...args) // accDetails := ui.NewAccDetails(&accountHelper) // Notes: take from UI module
}

func (tree *Tree) ProvideCurrentNodeNames() []string {
	return maps.Keys(tree.NodesList)
}

func (tree *Tree) ProvideCurrentNodes() []*tview.TreeNode {
	return maps.Values(tree.NodesList)
}

func (tree *Tree) ProvideChildID(node *tview.TreeNode) uuid.UUID{
	return tree.ChildsList[node]
}

func (tree *Tree) HasDomain(nodeName string) bool {
	// Check if Domain exist
	_, exists := tree.NodesList[nodeName]
	return exists
}

func (tree *Tree) LoadTree() {
	for _, data := range tree.ChildNodeList {
		nodeName := data.DomainName
		childs := data.Usernames
		for _, child := range childs {
			if !tree.HasDomain(nodeName) {
				tree.AddNodeNChild(nodeName, child.Username, child.ID)
			} else{
				tree.AddChild(nodeName, child.Username, child.ID)
			}
		}
	}
}

func (tree *Tree) ExpandNode(node *tview.TreeNode) {
	node.Expand()
}

func (tree *Tree) UnexpandNode(node *tview.TreeNode) {
	node.Collapse()
}

func (tree *Tree) ExpandNodeAll() {
	tree.Root.ExpandAll()
}

func (tree *Tree) UnexpandNodeAll() {
	tree.Root.CollapseAll()
}

func (tree *Tree) AddNodeNChild(nodeName string, childName string, childID uuid.UUID) {

	node := tview.NewTreeNode(nodeName)
	child := tview.NewTreeNode(childName)

	node.SetSelectedFunc(func() {
	  // tree.LogFunc(fmt.Sprintf("Node: %s", node.GetText()))
		tree.LogFunc(debug.TreeNodeSelected, debug.LogContext{
			Domain: nodeName,
		})
	})
	child.SetSelectedFunc(func(){
		// tree.LogFunc(fmt.Sprintf("Child: %s", child.GetText()))
		tree.LogFunc(debug.TreeChildSelected, debug.LogContext{
			Username: childName,
			Domain: nodeName,	
		})
		tree.CreateChildContent(tree.ProvideChildID(child))
		tree.SelectChild("accDet")
	})
	node.AddChild(child)
	
	tree.ChildsList[child] = childID
	tree.NodesList[nodeName] = node
	tree.Root.AddChild(node)
}

func (tree *Tree) AddNode(nodeName string) {
	node := tview.NewTreeNode(nodeName)
	node.SetSelectedFunc(func() {
	  // tree.LogFunc(fmt.Sprintf("Node: %s", node.GetText()))
		tree.LogFunc(debug.TreeNodeSelected, debug.LogContext{
			Domain: nodeName,
		})
	})
	tree.NodesList[nodeName] = node
	tree.Root.AddChild(node)
}

func (tree *Tree) DeleteNode() {

}

func (tree *Tree) MoveNode() {

}

func (tree *Tree) UpdateNode() {

}

func (tree *Tree) AddChild(nodeName string, childName string, childID uuid.UUID) {
	// Should only be called if the Domain for Acc exists
	var _node *tview.TreeNode
	tree.Root.Walk(func(node, parent *tview.TreeNode) bool {
		if node.GetText() == nodeName {
				_node = node
				return false
		}
		return true
	})
	if _node == nil {
		// tree.LogFunc(fmt.Sprintf("Error creating \"%s\" account.", childName))
		// Use Error Handler instead
		return
	}
	child := tview.NewTreeNode(childName)
	child.SetSelectedFunc(func(){
		// tree.LogFunc(fmt.Sprintf("Child: %s", child.GetText()))
		tree.LogFunc(debug.TreeChildSelected, debug.LogContext{
			Username: childName,
			Domain: nodeName,	
		})
		tree.CreateChildContent(tree.ProvideChildID(child))
		tree.SelectChild("accDet")
	})
	tree.ChildsList[child] = childID
	_node.AddChild(child)
}

func (tree *Tree) DeleteChild() {

}

func (tree *Tree) MoveChild() {

}

func (tree *Tree) UpdateChild() {

}
