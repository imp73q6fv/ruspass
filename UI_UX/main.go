package main

import (
	"os"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

func main() {
	// Initialize Qt application
	app := widgets.NewQApplication(len(os.Args), os.Args)

	// Create main window
	mainWindow := widgets.NewQMainWindow(nil, 0)
	mainWindow.SetWindowTitle("RusPass - Password Manager")
	mainWindow.SetMinimumSize2(900, 600)

	// Apply dark theme stylesheet (Enpass-like)
	stylesheet := `
		QMainWindow {
			background-color: #1e1e2e;
		}
		QWidget {
			background-color: #1e1e2e;
			color: #cdd6f4;
			font-family: 'Segoe UI', Arial, sans-serif;
			font-size: 13px;
		}
		QFrame#sidebar {
			background-color: #181825;
			border-right: 1px solid #313244;
		}
		QListWidget {
			background-color: #1e1e2e;
			color: #cdd6f4;
			border: none;
			outline: none;
		}
		QListWidget::item {
			padding: 12px 15px;
			border-bottom: 1px solid #313244;
		}
		QListWidget::item:selected {
			background-color: #313244;
			color: #89b4fa;
		}
		QListWidget::item:hover {
			background-color: #313244;
		}
		QLineEdit, QTextEdit {
			background-color: #313244;
			color: #cdd6f4;
			border: 1px solid #45475a;
			border-radius: 6px;
			padding: 8px 12px;
		}
		QLineEdit:focus, QTextEdit:focus {
			border: 1px solid #89b4fa;
		}
		QPushButton {
			background-color: #89b4fa;
			color: #1e1e2e;
			border: none;
			border-radius: 6px;
			padding: 10px 20px;
			font-weight: bold;
		}
		QPushButton:hover {
			background-color: #b4befe;
		}
		QPushButton#secondary {
			background-color: #45475a;
			color: #cdd6f4;
		}
		QPushButton#danger {
			background-color: #f38ba8;
			color: #1e1e2e;
		}
		QLabel#title {
			font-size: 24px;
			font-weight: bold;
			color: #cdd6f4;
		}
		QLabel#subtitle {
			font-size: 13px;
			color: #a6adc8;
		}
	`
	mainWindow.SetStyleSheet(stylesheet)

	// Create central widget
	centralWidget := widgets.NewQWidget(nil, 0)
	hLayout := widgets.NewQHBoxLayout(centralWidget)
	hLayout.SetContentsMargins(0, 0, 0, 0)
	hLayout.SetSpacing(0)

	// Sidebar
	sidebar := widgets.NewQWidget(nil, 0)
	sidebar.SetObjectName("sidebar")
	sidebar.SetFixedWidth(250)
	vLayout := widgets.NewQVBoxLayout(sidebar)
	vLayout.SetContentsMargins(20, 30, 20, 20)
	vLayout.SetSpacing(15)

	// Title
	titleLabel := widgets.NewQLabel(nil, 0)
	titleLabel.SetText("🔐 RusPass")
	titleLabel.SetObjectName("title")
	titleLabel.SetAlignment(core.Qt__AlignCenter)
	vLayout.AddWidget(titleLabel)

	// Search
	searchField := widgets.NewQLineEdit(nil)
	searchField.SetPlaceholderText("🔍 Search passwords...")
	vLayout.AddWidget(searchField)

	// Categories
	catsLabel := widgets.NewQLabel(nil, 0)
	catsLabel.SetText("CATEGORIES")
	catsLabel.SetObjectName("subtitle")
	vLayout.AddWidget(catsLabel)

	categories := []string{"📋 All", "💼 Work", "👤 Personal", "💰 Finance", "💬 Social", "🛒 Shopping"}
	for _, cat := range categories {
		btn := widgets.NewQPushButton(nil)
		btn.SetText(cat)
		btn.SetStyleSheet(`
			QPushButton {
				background-color: transparent;
				color: #a6adc8;
				text-align: left;
				padding: 8px 12px;
				border-radius: 6px;
				border: none;
			}
			QPushButton:hover {
				background-color: #313244;
				color: #cdd6f4;
			}
		`)
		vLayout.AddWidget(btn)
	}

	vLayout.AddStretch(1)

	// Add button
	addBtn := widgets.NewQPushButton(nil)
	addBtn.SetText("+ Add New Entry")
	vLayout.AddWidget(addBtn)

	hLayout.AddWidget(sidebar)

	// Content area
	content := widgets.NewQWidget(nil, 0)
	contentVLayout := widgets.NewQVBoxLayout(content)
	contentVLayout.SetContentsMargins(30, 30, 30, 30)
	contentVLayout.SetSpacing(20)

	// Header
	headerLabel := widgets.NewQLabel(nil, 0)
	headerLabel.SetText("All Passwords")
	headerLabel.SetObjectName("title")
	contentVLayout.AddWidget(headerLabel)

	// Splitter
	splitter := widgets.NewQSplitter(core.Qt__Horizontal, nil)

	// Entry list
	listWidget := widgets.NewQListWidget(nil)
	for i := 0; i < 5; i++ {
		item := widgets.NewQListWidgetItem(nil, 0)
		item.SetText("Sample Entry")
		listWidget.AddItem(item)
	}
	listWidget.SetMinimumWidth(300)
	splitter.AddWidget(listWidget)

	// Details panel
	detailsPanel := widgets.NewQWidget(nil, 0)
	formLayout := widgets.NewQFormLayout(detailsPanel)
	formLayout.SetSpacing(12)

	fields := []string{"Title", "Username", "Password", "URL"}
	for _, field := range fields {
		label := widgets.NewQLabel(nil, 0)
		label.SetText(field)
		lineEdit := widgets.NewQLineEdit(nil)
		if field == "Password" {
			lineEdit.SetEchoMode(widgets.QLineEdit__Password)
		}
		formLayout.AddRow2(label, lineEdit)
	}

	notesEdit := widgets.NewQTextEdit(nil)
	notesEdit.SetPlaceholderText("Notes...")
	formLayout.AddRow2(widgets.NewQLabel(nil, 0), notesEdit)

	// Buttons
	btnLayout := widgets.NewQHBoxLayout()
	saveBtn := widgets.NewQPushButton(nil)
	saveBtn.SetText("Save")
	copyBtn := widgets.NewQPushButton(nil)
	copyBtn.SetText("Copy Password")
	copyBtn.SetObjectName("secondary")
	deleteBtn := widgets.NewQPushButton(nil)
	deleteBtn.SetText("Delete")
	deleteBtn.SetObjectName("danger")
	btnLayout.AddWidget(saveBtn)
	btnLayout.AddWidget(copyBtn)
	btnLayout.AddWidget(deleteBtn)
	btnLayout.AddStretch(1)

	detailsVLayout := widgets.NewQVBoxLayout(detailsPanel)
	detailsVLayout.AddLayout(formLayout)
	detailsVLayout.AddLayout(btnLayout)

	splitter.AddWidget(detailsPanel)
	splitter.SetStretchFactor(0, 0)
	splitter.SetStretchFactor(1, 1)

	contentVLayout.AddWidget(splitter, 1)
	hLayout.AddWidget(content, 1)

	mainWindow.SetCentralWidget(centralWidget)

	// Menu bar
	menubar := mainWindow.MenuBar()
	fileMenu := menubar.AddMenu("&File")
	fileMenu.AddAction("E&xit").ConnectTriggered(func(bool) {
		mainWindow.Close()
	})

	helpMenu := menubar.AddMenu("&Help")
	helpMenu.AddAction("&About").ConnectTriggered(func(bool) {
		widgets.QMessageBox_About(mainWindow, "About RusPass",
			"RusPass Password Manager\nVersion 1.0\n\nBuilt with Rust + Go Qt\nEncryption: ChaCha20-Poly1305 + Argon2id")
	})

	mainWindow.Show()
	app.Exec()
}
