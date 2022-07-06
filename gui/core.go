package gui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/ozansz/gls/internal"
)

func GetApp() *tview.Application {
	app := tview.NewApplication()
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape || event.Key() == tcell.KeyCtrlC || event.Rune() == 'q' || event.Rune() == 'Q' {
			app.Stop()
		}
		return event
	})
	loadingPage := createLoadingPage(app)
	return app.SetRoot(loadingPage, true).SetFocus(loadingPage)
}

func LoadTreeView(app *tview.Application, node *internal.Node, f internal.SizeFormatter, path string) {
	root := constructTViewTreeFromNodeWithFormatter(node, f)
	root.SetExpanded(true)
	treeView := tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root)
	treeView.SetSelectedFunc(func(node *tview.TreeNode) {
		// Collapse if visible, expand if collapsed.
		node.SetExpanded(!node.IsExpanded())
	})
	app.SetRoot(treeView, true).SetFocus(treeView).Draw()
}

func createLoadingPage(app *tview.Application) tview.Primitive {
	loadingPage := tview.NewTextView().
		SetText("Loading...").
		SetTextAlign(tview.AlignCenter)
	return loadingPage
}

func constructTViewTreeFromNodeWithFormatter(node *internal.Node, f internal.SizeFormatter) *tview.TreeNode {
	treeNode := tview.NewTreeNode(node.InfoWithSizeFormatter(f)).
		SetReference(node).
		SetSelectable(node.IsDir)
	if node.IsDir {
		treeNode.SetColor(DirectoryColor).
			SetExpanded(false)
	}
	for _, child := range node.Children {
		treeNode.AddChild(constructTViewTreeFromNodeWithFormatter(child, f))
	}
	return treeNode
}
