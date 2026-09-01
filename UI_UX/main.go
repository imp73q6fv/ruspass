package main

import (
"os"
"strings"

"github.com/therecipe/qt/core"
"github.com/therecipe/qt/gui"
"github.com/therecipe/qt/widgets"
)

// Mock data for demonstration
type Entry struct {
ID       int
Title    string
Username string
Password string
URL      string
Category string
Icon     string
}

var mockEntries = []Entry{
{ID: 1, Title: "GitHub", Username: "user@example.com", Password: "gh_p_123456", URL: "github.com", Category: "Work", Icon: "github"},
{ID: 2, Title: "Gmail", Username: "john.doe@gmail.com", Password: "secure_pass_123", URL: "gmail.com", Category: "Personal", Icon: "google"},
{ID: 3, Title: "Chase Bank", Username: "johndoe", Password: "bank_secure_99", URL: "chase.com", Category: "Finance", Icon: "bank"},
{ID: 4, Title: "Netflix", Username: "family@netflix.com", Password: "movie_time!", URL: "netflix.com", Category: "Social", Icon: "video"},
{ID: 5, Title: "Amazon", Username: "shopper@amazon.com", Password: "prime_user_2024", URL: "amazon.com", Category: "Shopping", Icon: "cart"},
}

func main() {
// Enable High DPI scaling
core.QCoreApplication_SetAttribute(core.Qt__AA_EnableHighDpiScaling, true)

app := widgets.NewQApplication(len(os.Args), os.Args)
app.SetApplicationName("RusPass")
app.SetOrganizationName("RusPassTeam")

// --- Stylesheet (Enpass-like Dark Theme) ---
darkStyle := `
QMainWindow, QDialog {
background-color: #1e1e2e;
color: #cdd6f4;
font-family: "Segoe UI", "Roboto", sans-serif;
font-size: 14px;
}
QWidget#centralWidget {
background-color: #1e1e2e;
}
QWidget#sidebar {
background-color: #181825;
border-right: 1px solid #313244;
}
QWidget#header {
background-color: #1e1e2e;
border-bottom: 1px solid #313244;
}
QWidget#contentArea {
background-color: #1e1e2e;
}
QListWidget {
background-color: #1e1e2e;
color: #cdd6f4;
border: none;
outline: none;
}
QListWidget::item {
height: 60px;
border-radius: 8px;
margin: 4px 8px;
padding-left: 10px;
}
QListWidget::item:hover {
background-color: #313244;
}
QListWidget::item:selected {
background-color: #89b4fa;
color: #11111b;
font-weight: bold;
}
QLineEdit {
background-color: #313244;
border: 1px solid #45475a;
border-radius: 6px;
padding: 8px 12px;
color: #cdd6f4;
selection-background-color: #89b4fa;
}
QLineEdit:focus {
border: 1px solid #89b4fa;
}
QPushButton {
background-color: #89b4fa;
color: #11111b;
border: none;
border-radius: 6px;
padding: 8px 16px;
font-weight: bold;
}
QPushButton:hover {
background-color: #b4befe;
}
QPushButton#secondaryBtn {
background-color: #313244;
color: #cdd6f4;
}
QPushButton#secondaryBtn:hover {
background-color: #45475a;
}
QLabel#titleLabel {
font-size: 18px;
font-weight: bold;
color: #cdd6f4;
}
QLabel#subtitleLabel {
font-size: 12px;
color: #a6adc8;
}
QScrollArea {
border: none;
background: transparent;
}
QScrollBar:vertical {
background: #1e1e2e;
width: 10px;
border-radius: 5px;
}
QScrollBar::handle:vertical {
background: #45475a;
border-radius: 5px;
min-height: 20px;
}
QScrollBar::add-line:vertical, QScrollBar::sub-line:vertical {
height: 0px;
}
`
app.SetStyleSheet(darkStyle)

// --- Main Window ---
mainWindow := widgets.NewQMainWindow(nil, 0)
mainWindow.SetWindowTitle("RusPass - Vault")
mainWindow.SetMinimumSize2(1000, 700)
mainWindow.SetStyleSheet("background-color: #1e1e2e;")

// Central Widget
centralWidget := widgets.NewQWidget(nil, 0)
centralWidget.SetObjectName("centralWidget")
mainWindow.SetCentralWidget(centralWidget)

mainLayout := widgets.NewQHBoxLayout(centralWidget)
mainLayout.SetContentsMargins(0, 0, 0, 0)
mainLayout.SetSpacing(0)

// --- Sidebar (Left) ---
sidebar := widgets.NewQWidget(nil, 0)
sidebar.SetObjectName("sidebar")
sidebar.SetFixedWidth(260)
sidebarLayout := widgets.NewQVBoxLayout(sidebar)
sidebarLayout.SetContentsMargins(15, 20, 15, 20)
sidebarLayout.SetSpacing(15)

// Logo / App Name
logoLabel := widgets.NewQLabel(nil, 0)
logoLabel.SetText("🔒 RusPass")
logoLabel.SetStyleSheet("font-size: 24px; font-weight: bold; color: #89b4fa; padding: 10px;")
sidebarLayout.AddWidget(logoLabel, 0, core.Qt__AlignHCenter)

// Search Bar
searchInput := widgets.NewQLineEdit(nil, 0)
searchInput.SetPlaceholderText("Search vault...")
sidebarLayout.AddWidget(searchInput, 0, core.Qt__AlignTop)

// Categories List
categories := []string{"All Items", "Favorites", "Work", "Personal", "Finance", "Social", "Shopping"}
listWidget := widgets.NewQListWidget(nil, 0)
listWidget.SetVerticalScrollMode(widgets.QAbstractItemView__ScrollPerPixel)

for _, cat := range categories {
item := widgets.NewQListWidgetItem(cat, nil, 0)
item.SetSizeHint(core.NewQSize2(200, 40))
listWidget.AddItem(item)
}
listWidget.Item(0).SetSelected(true) // Select "All Items" by default

sidebarLayout.AddWidget(listWidget, 1, 0) // Stretch factor 1

// Add Button at bottom
addBtn := widgets.NewQPushButton(nil, 0)
addBtn.SetText("+ New Item")
addBtn.SetFixedHeight(45)
sidebarLayout.AddWidget(addBtn, 0, core.Qt__AlignBottom)

mainLayout.AddWidget(sidebar, 0, 0)

// --- Content Area (Right) ---
contentArea := widgets.NewQWidget(nil, 0)
contentArea.SetObjectName("contentArea")
contentLayout := widgets.NewQVBoxLayout(contentArea)
contentLayout.SetContentsMargins(30, 30, 30, 30)
contentLayout.SetSpacing(20)

// Header
headerWidget := widgets.NewQWidget(nil, 0)
headerWidget.SetObjectName("header")
headerLayout := widgets.NewQHBoxLayout(headerWidget)
headerLayout.SetContentsMargins(0, 0, 0, 15)

titleLabel := widgets.NewQLabel(nil, 0)
titleLabel.SetObjectName("titleLabel")
titleLabel.SetText("All Items")

headerLayout.AddWidget(titleLabel, 1, 0)
headerLayout.AddStretch(1)

contentLayout.AddWidget(headerWidget, 0, 0)

// Entry List View (Main List)
entryListWidget := widgets.NewQListWidget(nil, 0)
entryListWidget.SetVerticalScrollMode(widgets.QAbstractItemView__ScrollPerPixel)
entryListWidget.SetStyleSheet("background: transparent;")

// Populate Mock Data
updateEntryList := func(filter string) {
entryListWidget.Clear()
for _, entry := range mockEntries {
if filter != "" && !strings.Contains(strings.ToLower(entry.Title), strings.ToLower(filter)) {
continue
}

itemWidget := widgets.NewQWidget(nil, 0)
itemWidget.SetStyleSheet("background-color: #313244; border-radius: 8px; margin: 2px;")
itemLayout := widgets.NewQHBoxLayout(itemWidget)
itemLayout.SetContentsMargins(15, 10, 15, 10)

// Icon Placeholder
iconLabel := widgets.NewQLabel(nil, 0)
iconLabel.SetText("🔑")
iconLabel.SetFixedSize2(40, 40)
iconLabel.SetStyleSheet("font-size: 24px; background: #45475a; border-radius: 20px; qproperty-alignment: AlignCenter;")

// Text Info
infoWidget := widgets.NewQWidget(nil, 0)
infoLayout := widgets.NewQVBoxLayout(infoWidget)
infoLayout.SetContentsMargins(10, 0, 0, 0)
infoLayout.SetSpacing(4)

nameLabel := widgets.NewQLabel(nil, 0)
nameLabel.SetText(entry.Title)
nameLabel.SetStyleSheet("font-weight: bold; font-size: 15px; color: #cdd6f4;")

userLabel := widgets.NewQLabel(nil, 0)
userLabel.SetText(entry.Username)
userLabel.SetStyleSheet("color: #a6adc8; font-size: 13px;")

infoLayout.AddWidget(nameLabel, 0, 0)
infoLayout.AddWidget(userLabel, 0, 0)

itemLayout.AddWidget(iconLabel, 0, 0)
itemLayout.AddWidget(infoWidget, 1, 0)

listItem := widgets.NewQListWidgetItem(nil, 0)
listItem.SetSizeHint(core.NewQSize2(300, 70))
entryListWidget.AddItem(listItem)
entryListWidget.setItemWidget(listItem, itemWidget)
}
}

updateEntryList("")
contentLayout.AddWidget(entryListWidget, 1, 0)

// Detail Panel (Bottom/Right Preview)
detailWidget := widgets.NewQWidget(nil, 0)
detailWidget.SetFixedHeight(180)
detailWidget.SetStyleSheet("background-color: #181825; border-radius: 12px; border: 1px solid #313244;")
detailLayout := widgets.NewQHBoxLayout(detailWidget)
detailLayout.SetContentsMargins(20, 20, 20, 20)

// Left side of detail: Info
detailInfoLayout := widgets.NewQVBoxLayout()
detailInfoLayout.SetSpacing(10)

dTitle := widgets.NewQLabel(nil, 0)
dTitle.SetText("GitHub")
dTitle.SetStyleSheet("font-size: 20px; font-weight: bold; color: #89b4fa;")

dUser := widgets.NewQLabel(nil, 0)
dUser.SetText("user@example.com")
dUser.SetStyleSheet("color: #cdd6f4;")

dURL := widgets.NewQLabel(nil, 0)
dURL.SetText("https://github.com")
dURL.SetStyleSheet("color: #89b4fa; text-decoration: underline;")

detailInfoLayout.AddWidget(dTitle, 0, 0)
detailInfoLayout.AddWidget(dUser, 0, 0)
detailInfoLayout.AddWidget(dURL, 0, 0)
detailInfoLayout.AddStretch(1)

detailLayout.AddLayout(detailInfoLayout, 1)

// Right side of detail: Actions
actionLayout := widgets.NewQVBoxLayout()
actionLayout.SetSpacing(10)
actionLayout.SetAlignment(core.Qt__AlignRight)

copyPassBtn := widgets.NewQPushButton(nil, 0)
copyPassBtn.SetText("Copy Password")
copyPassBtn.SetFixedWidth(140)

editBtn := widgets.NewQPushButton(nil, 0)
editBtn.SetText("Edit")
editBtn.SetFixedWidth(140)
editBtn.SetObjectName("secondaryBtn")

actionLayout.AddStretch(1)
actionLayout.AddWidget(copyPassBtn, 0, 0)
actionLayout.AddWidget(editBtn, 0, 0)

detailLayout.AddLayout(actionLayout, 0)

contentLayout.AddWidget(detailWidget, 0, 0)

mainLayout.AddWidget(contentArea, 1, 0)

// --- Interactions ---

// Search Filter
searchInput.TextChangedConnect(func(text string) {
updateEntryList(text)
})

// Category Selection
listWidget.ItemClickedConnect(func(item *widgets.QListWidgetItem) {
titleLabel.SetText(item.Text())
// In a real app, filter by category here
updateEntryList("") 
})

// Entry Selection
entryListWidget.ItemClickedConnect(func(item *widgets.QListWidgetItem) {
widget := entryListWidget.itemWidget(item)
if widget != nil {
// Extract data from mock based on index for demo simplicity
idx := entryListWidget.Row(item)
if idx >= 0 && idx < len(mockEntries) {
e := mockEntries[idx]
dTitle.SetText(e.Title)
dUser.SetText(e.Username)
dURL.SetText("https://" + e.URL)
}
}
})

// Copy Password Action
copyPassBtn.ClickedConnect(func() {
clipboard := gui.QGuiApplication_clipboard()
clipboard.SetText("gh_p_123456")

msg := widgets.NewQMessageBox(nil, 0)
msg.SetText("Password Copied!")
msg.SetInformativeText("The password has been copied to your clipboard.")
msg.SetStyleSheet(darkStyle)
msg.Exec()
})

mainWindow.Show()

app.Exec()
}
