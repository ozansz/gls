package gui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/ozansz/gls/internal"
	"github.com/ozansz/gls/internal/info"
	"github.com/ozansz/gls/log"
)

var (
	currGrid          *tview.Grid            = nil
	currTreeView      *tview.TreeView        = nil
	currFileInfoTab   *tview.Table           = nil
	currPath          string                 = ""
	currSizeFormatter internal.SizeFormatter = nil
)

func GetApp(path string, f internal.SizeFormatter) *tview.Application {
	currPath = path
	currSizeFormatter = f
	app := tview.NewApplication()
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape || event.Key() == tcell.KeyCtrlC || event.Rune() == 'q' || event.Rune() == 'Q' {
			app.Stop()
		}
		if event.Rune() == 'c' || event.Rune() == 'C' {
			currTreeView.GetRoot().CollapseAll()
		}
		if event.Rune() == 'e' || event.Rune() == 'E' {
			currTreeView.GetRoot().ExpandAll()
		}
		if event.Key() == tcell.KeyBackspace || event.Key() == tcell.KeyDEL {
			if currTreeView == nil {
				log.Warning("Tree view is nil")
				return event
			}
			cNode := currTreeView.GetCurrentNode()
			if cNode == currTreeView.GetRoot() {
				showCannotRemoveRootWarning(app, cNode)
				return event
			}
			if len(cNode.GetChildren()) > 0 {
				showCannotRemoveFolderWarning(app, cNode)
				return event
			}
			askRemoveFile(app, cNode)
		}
		return event
	})
	loadingPage := createLoadingPage(app)
	return app.SetRoot(loadingPage, true).SetFocus(loadingPage)
}

func LoadTreeView(app *tview.Application, node *internal.Node, path string) {
	root := constructTViewTreeFromNodeWithFormatter(node)
	root.SetExpanded(true)
	treeView := tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root).
		SetSelectedFunc(func(node *tview.TreeNode) {
			// Collapse if visible, expand if collapsed.
			node.SetExpanded(!node.IsExpanded())
		})
	treeView.SetBorder(true).
		SetTitle(fmt.Sprintf("[ %s ]", path)).
		SetTitleAlign(tview.AlignCenter).
		SetTitleColor(TreeViewTitleColor)
	treeView.SetChangedFunc(func(node *tview.TreeNode) {
		updateFileInfoTab(app, node.GetReference().(*internal.Node))
	})
	treeView.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyTAB {
			// app.SetFocus(currFileInfoTab)
			cNode := treeView.GetCurrentNode()
			if cNode.IsExpanded() {
				cNode.CollapseAll()
			} else {
				cNode.ExpandAll()
			}
		}
	})

	grid := tview.NewGrid().SetRows(-1, -1, -1, -1, -1).SetColumns(-1, -1, -1, -1, -1)
	grid.SetBorder(true).
		SetTitle(fmt.Sprintf("[ %s ]", info.ProjectNameWithVersion())).
		SetTitleColor(GridTitleColor)

	fileInfoTab := createFileInfoTable(app)

	grid.AddItem(treeView, 0, 0, 10, 5, 0, 0, true)
	grid.AddItem(fileInfoTab, 10, 0, 2, 5, 0, 0, true)

	currTreeView = treeView
	currFileInfoTab = fileInfoTab
	currGrid = grid

	updateFileInfoTab(app, node)

	// app.SetRoot(treeView, true).SetFocus(treeView).Draw()
	app.SetRoot(grid, true).SetFocus(grid).Draw()
}

func createFileInfoTable(app *tview.Application) *tview.Table {
	table := tview.NewTable()
	attrHeader := tview.NewTableCell("Attribute").
		SetTextColor(FileInfoAttrColor).
		SetSelectable(false)
	valueHeader := tview.NewTableCell("Value").
		SetTextColor(FileInfoValueColor).
		SetSelectable(false)
	table.SetCell(0, 0, attrHeader).
		SetCell(0, 1, valueHeader).
		SetFixed(1, 0)
	table.SetSelectable(true, false).
		SetSeparator('|').
		SetBordersColor(BorderColor).
		SetTitle("Loading...").
		SetTitleColor(FileInfoTitleColor).
		SetBorder(true)
	// table.SetDoneFunc(func(key tcell.Key) {
	// 	if key == tcell.KeyTAB {
	// 		app.SetFocus(currTreeView)
	// 	}
	// })
	return table
}

func updateFileInfoTab(app *tview.Application, node *internal.Node) {
	// log.Infof("updateFileInfoTab is called with: %v", node)
	if currFileInfoTab == nil {
		log.Warning("updateFileInfoTab: currFileInfoTab is nil")
		return
	}
	currFileInfoTab.SetTitle(fmt.Sprintf("[ %s ]", node.Name))
	relativePath := node.RelativePath(currPath)
	pathAttrCell := tview.NewTableCell("Path").
		SetMaxWidth(FileInfoTabAttrWidth).
		SetTextColor(FileInfoAttrColor)
	pathValueCell := tview.NewTableCell(relativePath).
		SetTextColor(FileInfoValueColor)
	sizeAttrCell := tview.NewTableCell("Size").
		SetMaxWidth(FileInfoTabAttrWidth).
		SetTextColor(FileInfoAttrColor)
	sizeValueCell := tview.NewTableCell(fmt.Sprintf("%s (%d)", currSizeFormatter(node.Size), node.Size)).
		SetTextColor(FileInfoValueColor)
	typeAttrCell := tview.NewTableCell("Type").
		SetMaxWidth(FileInfoTabAttrWidth).
		SetTextColor(FileInfoAttrColor)
	typeValueCell := tview.NewTableCell(node.Mode.Type().String()).
		SetTextColor(FileInfoValueColor)
	permAttrCell := tview.NewTableCell("Permissions").
		SetMaxWidth(FileInfoTabAttrWidth).
		SetTextColor(FileInfoAttrColor)
	permValueCell := tview.NewTableCell(node.Mode.Perm().String()).
		SetTextColor(FileInfoValueColor)
	modifiedAttrCell := tview.NewTableCell("Modified").
		SetMaxWidth(FileInfoTabAttrWidth).
		SetTextColor(FileInfoAttrColor)
	modifiedValueCell := tview.NewTableCell(node.LastModification.String()).
		SetTextColor(FileInfoValueColor)
	currFileInfoTab.SetCell(1, 0, pathAttrCell).
		SetCell(1, 1, pathValueCell).
		SetCell(2, 0, sizeAttrCell).
		SetCell(2, 1, sizeValueCell).
		SetCell(3, 0, typeAttrCell).
		SetCell(3, 1, typeValueCell).
		SetCell(4, 0, permAttrCell).
		SetCell(4, 1, permValueCell).
		SetCell(5, 0, modifiedAttrCell).
		SetCell(5, 1, modifiedValueCell)
	// app.QueueUpdateDraw(func() {
	// 	// currFileInfoTab.Select(1, 0)
	// 	currFileInfoTab.ScrollToBeginning()
	// })
}

func createLoadingPage(app *tview.Application) tview.Primitive {
	loadingPage := tview.NewTextView().
		SetText("Loading...").
		SetTextAlign(tview.AlignCenter)
	return loadingPage
}

func constructTViewTreeFromNodeWithFormatter(node *internal.Node) *tview.TreeNode {
	treeNode := tview.NewTreeNode(node.InfoWithSizeFormatter(currSizeFormatter)).
		SetReference(node).
		SetSelectable(true)
	if node.IsDir {
		treeNode.SetColor(DirectoryColor).
			SetExpanded(false)
	}
	for _, child := range node.Children {
		treeNode.AddChild(constructTViewTreeFromNodeWithFormatter(child))
	}
	return treeNode
}

func showCannotRemoveFolderWarning(app *tview.Application, tnode *tview.TreeNode) {
	modal := tview.NewModal().
		SetText(fmt.Sprintf("Cannot remove folder %q", tnode.GetReference().(*internal.Node).Name)).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "OK" {
				app.SetRoot(currGrid, true).SetFocus(currGrid)
			}
		})
	app.SetRoot(modal, true).SetFocus(modal)
}

func showCannotRemoveRootWarning(app *tview.Application, tnode *tview.TreeNode) {
	modal := tview.NewModal().
		SetText(fmt.Sprintf("Cannot remove root folder %q", tnode.GetReference().(*internal.Node).Name)).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "OK" {
				app.SetRoot(currGrid, true).SetFocus(currGrid)
			}
		})
	app.SetRoot(modal, true).SetFocus(modal)
}

func askRemoveFile(app *tview.Application, tnode *tview.TreeNode) {
	node := tnode.GetReference().(*internal.Node)
	modal := tview.NewModal().
		SetText(fmt.Sprintf("Are you sure to remove %q?", node.RelativePath(currPath))).
		AddButtons([]string{"Cancel", "Yes"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Yes" {
				if err := tnode.GetReference().(*internal.Node).Remove(currPath); err != nil {
					log.Errorf("Could not remove file %q: %v", node.Name, err)
					showCannotRemoveFileError(app, node.Name, err.Error())
					return
				}
				newRoot := constructTViewTreeFromNodeWithFormatter(currTreeView.GetRoot().GetReference().(*internal.Node))
				newRoot.SetExpanded(true)
				currTreeView.SetRoot(newRoot).
					SetCurrentNode(newRoot)
			}
			app.SetRoot(currGrid, true).SetFocus(currGrid)
		})
	app.SetRoot(modal, true).SetFocus(modal)
}

func showCannotRemoveFileError(app *tview.Application, name string, err string) {
	modal := tview.NewModal().
		SetText(fmt.Sprintf("Cannot remove file %q: %s", name, err)).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "OK" {
				app.SetRoot(currGrid, true).SetFocus(currGrid)
			}
		})
	app.SetRoot(modal, true).SetFocus(modal)
}
