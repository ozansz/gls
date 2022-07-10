package gui

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/ozansz/gls/internal"
	"github.com/ozansz/gls/internal/info"
	"github.com/ozansz/gls/internal/types"
	"github.com/ozansz/gls/log"
)

var (
	currGrid          *tview.Grid                  = nil
	currTreeView      *tview.TreeView              = nil
	currFileInfoTab   *tview.Table                 = nil
	currLastLogView   *tview.TextView              = nil
	currPath          string                       = ""
	currSizeFormatter types.SizeFormatter          = nil
	originalRootNode  *types.Node                  = nil
	isFormInputActive bool                         = false
	markedFiles       map[*tview.TreeNode]struct{} = make(map[*tview.TreeNode]struct{})
)

func GetApp(path string, f types.SizeFormatter) *tview.Application {
	currPath = path
	currSizeFormatter = f
	app := tview.NewApplication()
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlC {
			app.Stop()
		}
		if !isFormInputActive {
			if event.Rune() == 'q' || event.Rune() == 'Q' || event.Key() == tcell.KeyEscape {
				app.Stop()
			}
			if event.Rune() == 'c' || event.Rune() == 'C' {
				currTreeView.GetRoot().CollapseAll()
			}
			if event.Rune() == 'e' || event.Rune() == 'E' {
				currTreeView.GetRoot().ExpandAll()
			}
			if event.Rune() == 's' || event.Rune() == 'S' {
				showSearchNameForm(app, false)
			}
			if event.Rune() == 'r' || event.Rune() == 'R' {
				showSearchNameForm(app, true)
			}
			if event.Rune() == 'x' || event.Rune() == 'X' {
				restoreOriginalRoot(app)
			}
			if event.Rune() == 'm' || event.Rune() == 'M' {
				markUnmarkFile(app)
			}
			if event.Rune() == 'u' || event.Rune() == 'U' {
				unmarkAll(app)
			}
			// Commands below here are about the current hovered file.
			if currTreeView == nil {
				log.Warning("Tree view is nil")
				return event
			}
			cNode := currTreeView.GetCurrentNode()
			if event.Rune() == 'o' || event.Rune() == 'O' {
				relPath := cNode.GetReference().(*types.Node).RelativePath(currPath)
				if err := internal.OpenFile(relPath); err != nil {
					log.Errorf("Could not open file %q: %v", relPath, err)
					showMessage(app, fmt.Sprintf("Could not open file %q: %v", relPath, err), nil)
					return event
				}
			}
			if event.Rune() == 'p' || event.Rune() == 'P' {
				relPath := cNode.GetReference().(*types.Node).RelativePath(currPath)
				askOpenFileWithProgram(app, relPath)
			}
			if event.Key() == tcell.KeyBackspace || event.Key() == tcell.KeyDEL {
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
		}
		return event
	})
	loadingPage := createLoadingPage(app)
	return app.SetRoot(loadingPage, true).SetFocus(loadingPage)
}

func LoadTreeView(app *tview.Application, node *types.Node, path string) {
	lastLogTextView := tview.NewTextView().
		SetText("OK.").
		SetTextColor(tcell.ColorWhite).
		SetWrap(true)
	currLastLogView = lastLogTextView

	originalRootNode = node
	root := constructNativeTree(node)
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
		updateFileInfoTab(app, node.GetReference().(*types.Node))
	})
	treeView.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyTAB {
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

	helpSideBar := createHelpSideBar(app)

	grid.AddItem(treeView, 0, 0, 19, 4, 0, 0, true)
	grid.AddItem(fileInfoTab, 19, 0, 4, 4, 0, 0, false)
	grid.AddItem(lastLogTextView, 23, 0, 2, 5, 0, 0, false)
	grid.AddItem(helpSideBar, 0, 4, 23, 1, 0, 0, false)

	currTreeView = treeView
	currFileInfoTab = fileInfoTab
	currGrid = grid

	updateFileInfoTab(app, node)

	app.SetRoot(grid, true).SetFocus(grid).Draw()
}

func createHelpSideBar(app *tview.Application) *tview.Table {
	table := tview.NewTable()
	keyHeader := tview.NewTableCell("Key").
		SetTextColor(FileInfoAttrColor).
		SetSelectable(false)
	commandHeader := tview.NewTableCell("Command").
		SetTextColor(FileInfoValueColor).
		SetSelectable(false)
	table.SetCell(0, 0, keyHeader).
		SetCell(0, 1, commandHeader).
		SetFixed(1, 0)
	table.SetSelectable(false, false).
		SetSeparator('|').
		SetBordersColor(BorderColor).
		SetTitle("[ Shortcuts ]").
		SetTitleColor(FileInfoTitleColor).
		SetBorder(true).
		SetBorderColor(BorderColor)
	for i, s := range keyboardShortcuts {
		table.SetCell(i+1, 0, tview.NewTableCell(s.Key).SetTextColor(FileInfoAttrColor)).
			SetCell(i+1, 1, tview.NewTableCell(s.Command).SetTextColor(FileInfoValueColor))
	}
	return table
}

func setLastLog(text string) {
	if currLastLogView == nil {
		return
	}
	currLastLogView.SetText(fmt.Sprintf("%s %s", time.Now().Format("2006-01-02 15:04:05"), text))
}

func setError(text string) {
	setLastLog(fmt.Sprintf("[ERROR] %s", text))
}

func setInfo(text string) {
	setLastLog(fmt.Sprintf("[INFO] %s", text))
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
	return table
}

func updateFileInfoTab(app *tview.Application, node *types.Node) {
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
	sizeValueCell := tview.NewTableCell(fmt.Sprintf("%s real, %s on disk (%d)", currSizeFormatter(node.Size), currSizeFormatter(node.SizeOnDisk), node.Size)).
		SetTextColor(FileInfoValueColor)
	typeAttrCell := tview.NewTableCell("Type").
		SetMaxWidth(FileInfoTabAttrWidth).
		SetTextColor(FileInfoAttrColor)
	fileType, err := node.GetFileType(currPath)
	if err != nil {
		log.Errorf("Failed to get type for file %q: %v", node.Name, err)
		fileType = fmt.Sprintf("<error: %v>", err)
	}
	typeValueCell := tview.NewTableCell(fileType).
		SetTextColor(FileInfoValueColor)
	permAttrCell := tview.NewTableCell("Permissions").
		SetMaxWidth(FileInfoTabAttrWidth).
		SetTextColor(FileInfoAttrColor)
	permValueCell := tview.NewTableCell(node.Mode.String()).
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
}

func createLoadingPage(app *tview.Application) tview.Primitive {
	loadingPage := tview.NewTextView().
		SetText("Loading...").
		SetTextAlign(tview.AlignCenter)
	return loadingPage
}

func constructNativeTree(node *types.Node) *tview.TreeNode {
	t := constructTViewTreeFromNodeWithFormatter(node, currSizeFormatter)
	setInfo(fmt.Sprintf("Constructed tree with %d files", node.FileCount()))
	return t
}

func constructTViewTreeFromNodeWithFormatter(node *types.Node, f types.SizeFormatter) *tview.TreeNode {
	treeNode := tview.NewTreeNode(node.InfoWithSizeFormatter(f)).
		SetReference(node).
		SetSelectable(true).
		SetColor(UnmarkedFileColor)
	if node.IsDir {
		treeNode.SetColor(DirectoryColor).
			SetExpanded(false)
	}
	for _, child := range node.Children {
		treeNode.AddChild(constructTViewTreeFromNodeWithFormatter(child, f))
	}
	return treeNode
}

func showCannotRemoveFolderWarning(app *tview.Application, tnode *tview.TreeNode) {
	modal := tview.NewModal().
		SetText(fmt.Sprintf("Cannot remove folder %q", tnode.GetReference().(*types.Node).Name)).
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
		SetText(fmt.Sprintf("Cannot remove root folder %q", tnode.GetReference().(*types.Node).Name)).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "OK" {
				app.SetRoot(currGrid, true).SetFocus(currGrid)
			}
		})
	app.SetRoot(modal, true).SetFocus(modal)
}

func askRemoveFile(app *tview.Application, tnode *tview.TreeNode) {
	node := tnode.GetReference().(*types.Node)
	modal := tview.NewModal().
		SetText(fmt.Sprintf("Are you sure to remove %q?", node.RelativePath(currPath))).
		AddButtons([]string{"Cancel", "Yes"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Yes" {
				if err := tnode.GetReference().(*types.Node).Remove(currPath); err != nil {
					log.Errorf("Could not remove file %q: %v", node.Name, err)
					showMessage(app, fmt.Sprintf("Cannot remove file %q: %v", node.Name, err.Error()), nil)
					return
				}
				newRoot := constructNativeTree(currTreeView.GetRoot().GetReference().(*types.Node))
				newRoot.SetExpanded(true)
				currTreeView.SetRoot(newRoot).
					SetCurrentNode(newRoot)
			}
			app.SetRoot(currGrid, true).SetFocus(currGrid)
		})
	app.SetRoot(modal, true).SetFocus(modal)
}

func showMessage(app *tview.Application, message string, callback func()) {
	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "OK" {
				app.SetRoot(currGrid, true).SetFocus(currGrid)
			}
		})
	app.SetRoot(modal, true).SetFocus(modal)
	if callback != nil {
		callback()
	}
}

func showSearchNameForm(app *tview.Application, isRegex bool) {
	var searchCaseInsensitive bool
	var searchInverse bool
	inputLabel := "Name"
	if isRegex {
		inputLabel = "Regular expression"
	}
	form := tview.NewForm().
		AddInputField(inputLabel, "", 32, nil, nil).
		AddButton("Cancel", func() {
			isFormInputActive = false
			app.SetRoot(currGrid, true).SetFocus(currGrid)
		})
	form.AddCheckbox("Case insensitive", false, func(checked bool) {
		searchCaseInsensitive = checked
	})
	form.AddCheckbox("Inverse (exclude the names)", false, func(checked bool) {
		searchInverse = checked
	})
	form.AddButton("Go", func() {
		defer func() {
			isFormInputActive = false
		}()
		query := form.GetFormItem(0).(*tview.InputField).GetText()
		if query == "" {
			showMessage(app, "Please enter a non-empty query", nil)
			return
		}
		log.Infof("Searching for query: %s", query)
		var err error
		var opts *types.TreeFilterOptions
		if isRegex {
			opts, err = types.NewTreeFilterOpts("", query, searchCaseInsensitive, searchInverse)
		} else {
			opts, err = types.NewTreeFilterOpts(query, "", searchCaseInsensitive, searchInverse)
		}
		if err != nil {
			errStr := fmt.Sprintf("Could not create filter options: %v", err)
			showMessage(app, errStr, nil)
			setError(errStr)
			return
		}
		newRootNode, err := originalRootNode.NewFilteredTree(opts)
		log.Infof("New root node: %v", newRootNode)
		if err != nil {
			log.Errorf("Could not run search for %q: %v", query, err)
			showMessage(app, fmt.Sprintf("Could not run search for %q: %v", query, err.Error()), nil)
			return
		}
		newRoot := constructNativeTree(newRootNode)
		newRoot.SetExpanded(true)
		currTreeView.SetRoot(newRoot).
			SetCurrentNode(newRoot)
		app.SetRoot(currGrid, true).SetFocus(currGrid)
	})
	form.SetBorder(true).
		SetTitle("Search by name").
		SetTitleAlign(tview.AlignCenter).
		SetTitleColor(SearchFormTitleColor)
	isFormInputActive = true
	app.SetRoot(form, true).SetFocus(form)
}

func restoreOriginalRoot(app *tview.Application) {
	root := constructNativeTree(originalRootNode)
	root.SetExpanded(true)
	currTreeView.SetRoot(root).
		SetCurrentNode(root)
	app.SetRoot(currGrid, true).SetFocus(currGrid)
}

func askOpenFileWithProgram(app *tview.Application, relPath string) {
	form := tview.NewForm().
		AddInputField("Executable", "", 32, nil, nil)
	form.AddButton("Open", func() {
		defer func() {
			isFormInputActive = false
		}()
		program := form.GetFormItem(0).(*tview.InputField).GetText()
		if program == "" {
			showMessage(app, "Please enter an executable name", nil)
			return
		}
		log.Infof("Opening %q with %q", relPath, program)
		if err := internal.OpenFileWithProgram(relPath, program); err != nil {
			log.Errorf("Could not open file %q with %q: %v", relPath, err)
			setError(fmt.Sprintf("Could not open file %q with %q: %v", relPath, program, err))
			app.SetRoot(currGrid, true).SetFocus(currGrid)
		}
	})
	form.AddButton("Cancel", func() {
		isFormInputActive = false
		app.SetRoot(currGrid, true).SetFocus(currGrid)
	})
	form.SetBorder(true).
		SetTitle("Open file with program").
		SetTitleAlign(tview.AlignCenter).
		SetTitleColor(SearchFormTitleColor)
	isFormInputActive = true
	app.SetRoot(form, true).SetFocus(form)

}

func markUnmarkFile(app *tview.Application) {
	cNode := currTreeView.GetCurrentNode()
	if _, ok := markedFiles[cNode]; ok { // unmark if already marked
		cNode.SetColor(UnmarkedFileColor)
		delete(markedFiles, cNode)
		log.Debugf("Removed node with name %q from marked files", cNode.GetText())
	} else {
		cNode.SetColor(MarkedFileColor)
		markedFiles[cNode] = struct{}{}
		log.Debugf("Added node with name %q to marked files", cNode.GetText())
	}
}

func unmarkAll(app *tview.Application) {
	for n := range markedFiles {
		n.SetColor(UnmarkedFileColor)
		delete(markedFiles, n)
		log.Debugf("Removed node with name %q from marked files (unmarkAll)", n.GetText())
	}
	log.Debugf("markedFiles after unmarkAll: %#v", markedFiles)
}
