package main

/*
#cgo pkg-config: Qt5Core Qt5Widgets Qt5Gui
#include <QtWidgets/QApplication>
#include <QtWidgets/QMainWindow>
#include <QtWidgets/QVBoxLayout>
#include <QtWidgets/QHBoxLayout>
#include <QtWidgets/QPushButton>
#include <QtWidgets/QLabel>
#include <QtWidgets/QLineEdit>
#include <QtWidgets/QTextEdit>
#include <QtWidgets/QListView>
#include <QtWidgets/QListWidget>
#include <QtWidgets/QStackedWidget>
#include <QtWidgets/QFrame>
#include <QtWidgets/QScrollArea>
#include <QtWidgets/QGroupBox>
#include <QtWidgets/QComboBox>
#include <QtWidgets/QMessageBox>
#include <QtWidgets/QToolBar>
#include <QtWidgets/QStatusBar>
#include <QtWidgets/QMenu>
#include <QtWidgets/QMenuBar>
#include <QtWidgets/QSplitter>
#include <QtWidgets/QTreeWidget>
#include <QtWidgets/QTreeWidgetItem>
#include <QtWidgets/QHeaderView>
#include <QtWidgets/QTableWidget>
#include <QtWidgets/QTableWidgetItem>
#include <QtWidgets/QSortFilterProxyModel>
*/
import "C"

import (
	"os"
	"path/filepath"
	
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

// RusPassApp represents the main application
type RusPassApp struct {
	app         *widgets.QApplication
	window      *widgets.QMainWindow
	backend     *BackendService
	dbPath      string
	masterToken string
	isLoggedIn  bool
	
	// UI Components
	stackedWidget    *widgets.QStackedWidget
	loginPage        *widgets.QWidget
	mainPage         *widgets.QWidget
	entryListWidget  *widgets.QTableWidget
	entryDetailForm  map[string]*widgets.QLineEdit
	categoryCombo    *widgets.QComboBox
	searchInput      *widgets.QLineEdit
	statusBar        *widgets.QStatusBar
}

// NewRusPassApp creates a new RusPass application instance
func NewRusPassApp() *RusPassApp {
	return &RusPassApp{
		entryDetailForm: make(map[string]*widgets.QLineEdit),
	}
}

// Initialize sets up the application
func (r *RusPassApp) Initialize(args []string) int {
	// Set application attributes
	core.QCoreApplication_SetAttribute(core.Qt__AA_EnableHighDpiScaling, true)
	
	// Create application
	r.app = widgets.NewQApplication(len(args), args)
	
	// Set application name and style
	r.app.SetApplicationName("RusPass")
	r.app.SetOrganizationName("RusPass")
	
	// Set dark theme stylesheet (Enpass-like)
	darkStyle := `
		QMainWindow {
			background-color: #1e1e1e;
			color: #ffffff;
		}
		QWidget {
			background-color: #2d2d2d;
			color: #ffffff;
			font-size: 14px;
		}
		QPushButton {
			background-color: #4a9eff;
			color: white;
			border: none;
			padding: 8px 16px;
			border-radius: 4px;
			font-weight: bold;
		}
		QPushButton:hover {
			background-color: #3a8eef;
		}
		QPushButton:pressed {
			background-color: #2a7edf;
		}
		QLineEdit, QTextEdit, QComboBox {
			background-color: #3d3d3d;
			color: white;
			border: 1px solid #555555;
			border-radius: 4px;
			padding: 6px;
		}
		QLineEdit:focus, QTextEdit:focus {
			border: 1px solid #4a9eff;
		}
		QTableWidget {
			background-color: #2d2d2d;
			color: white;
			gridline-color: #3d3d3d;
			border: none;
		}
		QTableWidget::item {
			padding: 8px;
		}
		QTableWidget::item:selected {
			background-color: #4a9eff;
		}
		QHeaderView::section {
			background-color: #3d3d3d;
			color: white;
			padding: 8px;
			border: none;
		}
		QLabel {
			color: #ffffff;
			padding: 4px;
		}
		QGroupBox {
			border: 1px solid #555555;
			border-radius: 4px;
			margin-top: 12px;
			padding-top: 12px;
			font-weight: bold;
		}
		QGroupBox::title {
			subcontrol-origin: margin;
			left: 10px;
			padding: 0 4px;
		}
		QScrollArea {
			border: none;
			background-color: transparent;
		}
		QSplitter::handle {
			background-color: #555555;
			width: 2px;
		}
		QMenu {
			background-color: #2d2d2d;
			color: white;
			border: 1px solid #555555;
		}
		QMenu::item:selected {
			background-color: #4a9eff;
		}
		QMenuBar {
			background-color: #1e1e1e;
			color: white;
		}
		QStatusBar {
			background-color: #1e1e1e;
			color: #888888;
		}
	`
	r.app.SetStyleSheet(darkStyle)
	
	// Determine database path
	homeDir, _ := os.UserHomeDir()
	r.dbPath = filepath.Join(homeDir, ".ruspass", "ruspass.db")
	
	// Ensure directory exists
	os.MkdirAll(filepath.Dir(r.dbPath), 0755)
	
	// Initialize backend service
	r.backend = NewBackendService("") // Will be set when backend is compiled
	
	// Create main window
	r.createMainWindow()
	
	// Show login page initially
	r.showLoginPage()
	
	return r.app.Exec()
}

// createMainWindow sets up the main window structure
func (r *RusPassApp) createMainWindow() {
	r.window = widgets.NewQMainWindow(nil, 0)
	r.window.SetWindowTitle("RusPass - Password Manager")
	r.window.SetMinimumSize2(1000, 700)
	r.window.Resize2(1200, 800)
	
	// Create central widget and layout
	centralWidget := widgets.NewQWidget(nil, 0)
	mainLayout := widgets.NewQVBoxLayout(centralWidget)
	mainLayout.SetContentsMargins(0, 0, 0, 0)
	mainLayout.SetSpacing(0)
	
	// Create menu bar
	r.createMenuBar()
	
	// Create stacked widget for page navigation
	r.stackedWidget = widgets.NewQStackedWidget(nil)
	mainLayout.AddWidget(r.stackedWidget, 1, 0)
	
	// Create pages
	r.createLoginPage()
	r.createMainPage()
	
	// Create status bar
	r.statusBar = widgets.NewQStatusBar(nil)
	r.window.SetStatusBar(r.statusBar)
	r.statusBar.ShowMessage("Welcome to RusPass", 0)
	
	r.window.SetCentralWidget(centralWidget)
}

// createMenuBar creates the application menu bar
func (r *RusPassApp) createMenuBar() {
	menuBar := r.window.MenuBar()
	
	// File menu
	fileMenu := menuBar.AddMenu("&File")
	
	newEntryAction := fileMenu.AddAction("&New Entry")
	newEntryAction.SetShortcut(gui.NewQKeySequence2("Ctrl+N"))
	newEntryAction.ConnectTriggered(func(checked bool) {
		r.showAddEntryForm()
	})
	
	exitAction := fileMenu.AddAction("E&xit")
	exitAction.SetShortcut(gui.NewQKeySequence2("Ctrl+Q"))
	exitAction.ConnectTriggered(func(checked bool) {
		r.app.Quit()
	})
	
	// Edit menu
	editMenu := menuBar.AddMenu("&Edit")
	
	editAction := editMenu.AddAction("&Edit Entry")
	editAction.SetShortcut(gui.NewQKeySequence2("Ctrl+E"))
	editAction.ConnectTriggered(func(checked bool) {
		r.editSelectedEntry()
	})
	
	deleteAction := editMenu.AddAction("&Delete Entry")
	deleteAction.SetShortcut(gui.NewQKeySequence2("Delete"))
	deleteAction.ConnectTriggered(func(checked bool) {
		r.deleteSelectedEntry()
	})
	
	// View menu
	viewMenu := menuBar.AddMenu("&View")
	
	refreshAction := viewMenu.AddAction("&Refresh")
	refreshAction.SetShortcut(gui.NewQKeySequence2("F5"))
	refreshAction.ConnectTriggered(func(checked bool) {
		r.loadEntries()
	})
	
	// Help menu
	helpMenu := menuBar.AddMenu("&Help")
	
	aboutAction := helpMenu.AddAction("&About")
	aboutAction.ConnectTriggered(func(checked bool) {
		widgets.QMessageBox_Information(r.window, "About RusPass", 
			"RusPass Password Manager\nVersion 1.0.0\n\nSecure password management with ChaCha20-Poly1305 encryption and Argon2id key derivation.", 
			widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
	})
}

// createLoginPage creates the login page
func (r *RusPassApp) createLoginPage() {
	r.loginPage = widgets.NewQWidget(nil, 0)
	layout := widgets.NewQVBoxLayout(r.loginPage)
	layout.SetAlignment(core.Qt__AlignCenter)
	
	// Add spacer at top
	layout.AddStretch(1)
	
	// Logo/Title
	titleLabel := widgets.NewQLabel(nil, 0)
	titleLabel.SetText("🔐 RusPass")
	titleLabel.SetStyleSheet("font-size: 48px; font-weight: bold; color: #4a9eff; padding: 20px;")
	titleLabel.SetAlignment2(core.Qt__AlignCenter)
	layout.AddWidget(titleLabel, 0, core.Qt__AlignCenter)
	
	// Subtitle
	subtitleLabel := widgets.NewQLabel(nil, 0)
	subtitleLabel.SetText("Secure Password Manager")
	subtitleLabel.SetStyleSheet("font-size: 18px; color: #888888; padding: 10px;")
	subtitleLabel.SetAlignment2(core.Qt__AlignCenter)
	layout.AddWidget(subtitleLabel, 0, core.Qt__AlignCenter)
	
	// Add spacer
	layout.AddStretch(1)
	
	// Password input container
	passwordContainer := widgets.NewQWidget(nil, 0)
	passwordLayout := widgets.NewQVBoxLayout(passwordContainer)
	passwordLayout.SetMaximumWidth(400)
	
	// Master password label
	passwordLabel := widgets.NewQLabel(nil, 0)
	passwordLabel.SetText("Master Password")
	passwordLabel.SetStyleSheet("font-size: 16px; padding: 10px;")
	passwordLayout.AddWidget(passwordLabel, 0, core.Qt__AlignCenter)
	
	// Master password input
	masterPasswordInput := widgets.NewQLineEdit(nil)
	masterPasswordInput.SetEchoMode(widgets.QLineEdit__Password)
	masterPasswordInput.SetPlaceholderText("Enter your master password")
	masterPasswordInput.SetStyleSheet("font-size: 16px; padding: 12px;")
	masterPasswordInput.SetMinimumHeight(50)
	passwordLayout.AddWidget(masterPasswordInput, 0, 0)
	
	// Login button
	loginButton := widgets.NewQPushButton2("Unlock Vault", nil)
	loginButton.SetStyleSheet("font-size: 18px; padding: 12px; margin-top: 20px;")
	loginButton.SetMinimumHeight(50)
	loginButton.ConnectClicked(func() {
		r.attemptLogin(masterPasswordInput.Text())
	})
	passwordLayout.AddWidget(loginButton, 0, core.Qt__AlignCenter)
	
	layout.AddWidget(passwordContainer, 0, core.Qt__AlignCenter)
	
	// Add spacer at bottom
	layout.AddStretch(2)
	
	r.stackedWidget.AddWidget(r.loginPage)
}

// createMainPage creates the main application page
func (r *RusPassApp) createMainPage() {
	r.mainPage = widgets.NewQWidget(nil, 0)
	mainLayout := widgets.NewQHBoxLayout(r.mainPage)
	mainLayout.SetContentsMargins(10, 10, 10, 10)
	mainLayout.SetSpacing(10)
	
	// Left panel - Entry list
	leftPanel := r.createLeftPanel()
	mainLayout.AddWidget(leftPanel, 1, 0)
	
	// Right panel - Entry details
	rightPanel := r.createRightPanel()
	mainLayout.AddWidget(rightPanel, 1, 0)
	
	r.stackedWidget.AddWidget(r.mainPage)
}

// createLeftPanel creates the left panel with entry list
func (r *RusPassApp) createLeftPanel() *widgets.QWidget {
	panel := widgets.NewQWidget(nil, 0)
	layout := widgets.NewQVBoxLayout(panel)
	layout.SetContentsMargins(0, 0, 0, 0)
	layout.SetSpacing(10)
	
	// Search box
	searchContainer := widgets.NewQWidget(nil, 0)
	searchLayout := widgets.NewQHBoxLayout(searchContainer)
	searchLayout.SetContentsMargins(0, 0, 0, 0)
	
	r.searchInput = widgets.NewQLineEdit(nil)
	r.searchInput.SetPlaceholderText("🔍 Search passwords...")
	r.searchInput.SetMinimumHeight(40)
	r.searchInput.ConnectTextChanged(func(text string) {
		r.filterEntries(text)
	})
	searchLayout.AddWidget(r.searchInput, 1, 0)
	
	layout.AddWidget(searchContainer, 0, 0)
	
	// Entry table
	r.entryListWidget = widgets.NewQTableWidget(nil)
	r.entryListWidget.SetColumnCount(3)
	r.entryListWidget.SetHorizontalHeaderLabels([]string{"Title", "Username", "Category"})
	r.entryListWidget.HorizontalHeader().SetSectionResizeMode(0, widgets.QHeaderView__Stretch)
	r.entryListWidget.HorizontalHeader().SetSectionResizeMode(1, widgets.QHeaderView__Stretch)
	r.entryListWidget.HorizontalHeader().SetSectionResizeMode(2, widgets.QHeaderView__ResizeToContents)
	r.entryListWidget.SetSelectionBehavior(widgets.QAbstractItemView__SelectRows)
	r.entryListWidget.SetSelectionMode(widgets.QAbstractItemView__SingleSelection)
	r.entryListWidget.SetEditTriggers(widgets.QAbstractItemView__NoEditTriggers)
	r.entryListWidget.SetAlternatingRowColors(true)
	r.entryListWidget.SetShowGrid(false)
	r.entryListWidget.ConnectCellClicked(func(row, column int) {
		r.showEntryDetails(row)
	})
	layout.AddWidget(r.entryListWidget, 1, 0)
	
	// Action buttons
	buttonContainer := widgets.NewQWidget(nil, 0)
	buttonLayout := widgets.NewQHBoxLayout(buttonContainer)
	buttonLayout.SetContentsMargins(0, 0, 0, 0)
	
	addButton := widgets.NewQPushButton2("+ Add New", nil)
	addButton.SetMinimumHeight(40)
	addButton.ConnectClicked(func() {
		r.showAddEntryForm()
	})
	buttonLayout.AddWidget(addButton, 0, 0)
	
	buttonLayout.AddStretch(1)
	
	layout.AddWidget(buttonContainer, 0, 0)
	
	return panel
}

// createRightPanel creates the right panel with entry details
func (r *RusPassApp) createRightPanel() *widgets.QWidget {
	panel := widgets.NewQWidget(nil, 0)
	layout := widgets.NewQVBoxLayout(panel)
	layout.SetContentsMargins(0, 0, 0, 0)
	layout.SetSpacing(15)
	
	// Title
	titleLabel := widgets.NewQLabel(nil, 0)
	titleLabel.SetText("Entry Details")
	titleLabel.SetStyleSheet("font-size: 20px; font-weight: bold; padding: 10px;")
	layout.AddWidget(titleLabel, 0, 0)
	
	// Scrollable form area
	scrollArea := widgets.NewQScrollArea(nil)
	scrollArea.SetWidgetResizable(true)
	scrollArea.SetStyleSheet("border: none; background-color: transparent;")
	
	formContainer := widgets.NewQWidget(nil, 0)
	formLayout := widgets.NewQVBoxLayout(formContainer)
	formLayout.SetSpacing(15)
	
	// Title field
	formLayout.AddWidget(widgets.NewQLabel2("Title:"), 0, 0)
	titleInput := widgets.NewQLineEdit(nil)
	titleInput.SetMinimumHeight(40)
	r.entryDetailForm["title"] = titleInput
	formLayout.AddWidget(titleInput, 0, 0)
	
	// Username field
	formLayout.AddWidget(widgets.NewQLabel2("Username:"), 0, 0)
	usernameInput := widgets.NewQLineEdit(nil)
	usernameInput.SetMinimumHeight(40)
	r.entryDetailForm["username"] = usernameInput
	formLayout.AddWidget(usernameInput, 0, 0)
	
	// Password field
	formLayout.AddWidget(widgets.NewQLabel2("Password:"), 0, 0)
	passwordContainer := widgets.NewQWidget(nil, 0)
	passwordLayout := widgets.NewQHBoxLayout(passwordContainer)
	passwordLayout.SetContentsMargins(0, 0, 0, 0)
	
	passwordInput := widgets.NewQLineEdit(nil)
	passwordInput.SetEchoMode(widgets.QLineEdit__Password)
	passwordInput.SetMinimumHeight(40)
	r.entryDetailForm["password"] = passwordInput
	passwordLayout.AddWidget(passwordInput, 1, 0)
	
	showPasswordBtn := widgets.NewQPushButton2("👁", nil)
	showPasswordBtn.SetMaximumWidth(50)
	showPasswordBtn.ConnectClicked(func() {
		if passwordInput.EchoMode() == widgets.QLineEdit__Password {
			passwordInput.SetEchoMode(widgets.QLineEdit__Normal)
			showPasswordBtn.SetText("🙈")
		} else {
			passwordInput.SetEchoMode(widgets.QLineEdit__Password)
			showPasswordBtn.SetText("👁")
		}
	})
	passwordLayout.AddWidget(showPasswordBtn, 0, 0)
	formLayout.AddWidget(passwordContainer, 0, 0)
	
	// URL field
	formLayout.AddWidget(widgets.NewQLabel2("URL:"), 0, 0)
	urlInput := widgets.NewQLineEdit(nil)
	urlInput.SetMinimumHeight(40)
	r.entryDetailForm["url"] = urlInput
	formLayout.AddWidget(urlInput, 0, 0)
	
	// Category field
	formLayout.AddWidget(widgets.NewQLabel2("Category:"), 0, 0)
	r.categoryCombo = widgets.NewQComboBox(nil)
	r.categoryCombo.SetMinimumHeight(40)
	r.categoryCombo.AddItems([]string{"General", "Work", "Personal", "Finance", "Social", "Shopping", "Travel", "Other"})
	formLayout.AddWidget(r.categoryCombo, 0, 0)
	
	// Notes field
	formLayout.AddWidget(widgets.NewQLabel2("Notes:"), 0, 0)
	notesInput := widgets.NewQTextEdit(nil)
	notesInput.SetMinimumHeight(100)
	r.entryDetailForm["notes"] = nil // QTextEdit handled separately
	r.entryDetailForm["notes_textedit"] = notesInput
	formLayout.AddWidget(notesInput, 1, 0)
	
	formLayout.AddStretch(1)
	
	scrollArea.SetWidget(formContainer)
	layout.AddWidget(scrollArea, 1, 0)
	
	// Action buttons
	buttonContainer := widgets.NewQWidget(nil, 0)
	buttonLayout := widgets.NewQHBoxLayout(buttonContainer)
	buttonLayout.SetContentsMargins(0, 0, 0, 0)
	
	saveButton := widgets.NewQPushButton2("💾 Save", nil)
	saveButton.SetMinimumHeight(40)
	saveButton.ConnectClicked(func() {
		r.saveEntry()
	})
	buttonLayout.AddWidget(saveButton, 0, 0)
	
	buttonLayout.AddStretch(1)
	
	deleteButton := widgets.NewQPushButton2("🗑 Delete", nil)
	deleteButton.SetMinimumHeight(40)
	deleteButton.SetStyleSheet("background-color: #dc3545;")
	deleteButton.ConnectClicked(func() {
		r.deleteSelectedEntry()
	})
	buttonLayout.AddWidget(deleteButton, 0, 0)
	
	layout.AddWidget(buttonContainer, 0, 0)
	
	return panel
}

// showLoginPage displays the login page
func (r *RusPassApp) showLoginPage() {
	r.stackedWidget.SetCurrentWidget(r.loginPage)
	r.statusBar.ShowMessage("Please enter your master password", 0)
}

// showMainPage displays the main page
func (r *RusPassApp) showMainPage() {
	r.stackedWidget.SetCurrentWidget(r.mainPage)
	r.loadEntries()
	r.statusBar.ShowMessage("Vault unlocked", 0)
}

// attemptLogin tries to authenticate with the master password
func (r *RusPassApp) attemptLogin(password string) {
	if password == "" {
		widgets.QMessageBox_Warning(r.window, "Warning", "Please enter a master password", 
			widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
		return
	}
	
	// For demo purposes, accept any non-empty password
	// In production, this would call the Rust backend
	r.masterToken = "demo_token"
	r.isLoggedIn = true
	
	r.showMainPage()
}

// loadEntries loads all entries into the table
func (r *RusPassApp) loadEntries() {
	r.entryListWidget.SetRowCount(0)
	
	// Demo data
	demoEntries := []struct {
		title, username, category string
	}{
		{"GitHub", "user@example.com", "Work"},
		{"Google", "user@gmail.com", "Personal"},
		{"Bank of America", "user123", "Finance"},
		{"Facebook", "user@facebook.com", "Social"},
		{"Amazon", "user@amazon.com", "Shopping"},
	}
	
	for i, entry := range demoEntries {
		r.entryListWidget.InsertRow(i)
		
		titleItem := widgets.NewQTableWidgetItem2(entry.title, 0)
		r.entryListWidget.SetItem(i, 0, titleItem)
		
		usernameItem := widgets.NewQTableWidgetItem2(entry.username, 0)
		r.entryListWidget.SetItem(i, 1, usernameItem)
		
		categoryItem := widgets.NewQTableWidgetItem2(entry.category, 0)
		r.entryListWidget.SetItem(i, 2, categoryItem)
	}
	
	r.statusBar.ShowMessage(core.QString_FromStdString("Loaded %d entries"), 0)
}

// filterEntries filters the entry list based on search query
func (r *RusPassApp) filterEntries(query string) {
	// Implementation would filter the table based on query
	// For demo, just reload all entries
	r.loadEntries()
}

// showEntryDetails shows details of selected entry
func (r *RusPassApp) showEntryDetails(row int) {
	if row < 0 {
		return
	}
	
	titleItem := r.entryListWidget.Item(row, 0)
	usernameItem := r.entryListWidget.Item(row, 1)
	categoryItem := r.entryListWidget.Item(row, 2)
	
	if titleItem != nil {
		r.entryDetailForm["title"].SetText(titleItem.Text())
	}
	if usernameItem != nil {
		r.entryDetailForm["username"].SetText(usernameItem.Text())
	}
	if categoryItem != nil {
		index := r.categoryCombo.FindText(categoryItem.Text(), core.Qt__MatchExactly)
		if index >= 0 {
			r.categoryCombo.SetCurrentIndex(index)
		}
	}
	
	// Clear password and other fields for security
	r.entryDetailForm["password"].SetText("••••••••••••")
	r.entryDetailForm["url"].SetText("")
	if notesEdit, ok := r.entryDetailForm["notes_textedit"]; ok && notesEdit != nil {
		notesEdit.SetText("")
	}
}

// showAddEntryForm clears the form for adding a new entry
func (r *RusPassApp) showAddEntryForm() {
	r.entryDetailForm["title"].SetText("")
	r.entryDetailForm["username"].SetText("")
	r.entryDetailForm["password"].SetText("")
	r.entryDetailForm["url"].SetText("")
	r.categoryCombo.SetCurrentIndex(0)
	if notesEdit, ok := r.entryDetailForm["notes_textedit"]; ok && notesEdit != nil {
		notesEdit.SetText("")
	}
	r.statusBar.ShowMessage("Creating new entry", 0)
}

// saveEntry saves the current entry
func (r *RusPassApp) saveEntry() {
	title := r.entryDetailForm["title"].Text()
	if title == "" {
		widgets.QMessageBox_Warning(r.window, "Warning", "Title is required", 
			widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
		return
	}
	
	// In production, this would save to the backend
	widgets.QMessageBox_Information(r.window, "Success", "Entry saved successfully", 
		widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
	
	r.loadEntries()
	r.statusBar.ShowMessage("Entry saved", 0)
}

// editSelectedEntry prepares the form for editing
func (r *RusPassApp) editSelectedEntry() {
	currentRow := r.entryListWidget.CurrentRow()
	if currentRow < 0 {
		widgets.QMessageBox_Warning(r.window, "Warning", "Please select an entry to edit", 
			widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
		return
	}
	r.showEntryDetails(currentRow)
}

// deleteSelectedEntry deletes the selected entry
func (r *RusPassApp) deleteSelectedEntry() {
	currentRow := r.entryListWidget.CurrentRow()
	if currentRow < 0 {
		widgets.QMessageBox_Warning(r.window, "Warning", "Please select an entry to delete", 
			widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
		return
	}
	
	result := widgets.QMessageBox_Question(r.window, "Confirm Delete", 
		"Are you sure you want to delete this entry?", 
		widgets.QMessageBox__Yes|widgets.QMessageBox__No, widgets.QMessageBox__No)
	
	if result == widgets.QMessageBox__Yes {
		r.entryListWidget.RemoveRow(currentRow)
		r.statusBar.ShowMessage("Entry deleted", 0)
	}
}
