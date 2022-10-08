package gui

import (
	"fmt"
	"github.com/ozansz/gls/internal"
	"github.com/ozansz/gls/internal/types"
	"github.com/ozansz/gls/log"
	"github.com/rivo/tview"
)

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
