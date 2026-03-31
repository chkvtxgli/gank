package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gank/internal/config"
	"gank/internal/extractor"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

type model struct {
	step         int
	pdfPath      string
	bank         string
	account      string
	outputPath   string
	configDir    string
	transactions []extractor.Transaction
	cfg          *config.BankConfig
	err          error
	form         *huh.Form
	width        int
	height       int
}

func initialModel(configDir string) model {
	if configDir == "" {
		configDir = "banks"
	}
	return model{
		step:      0,
		configDir: configDir,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	}

	switch m.step {
	case 0:
		return m.updateFileSelect(msg)
	case 1:
		return m.updateBankSelect(msg)
	case 2:
		return m.updateAccountSelect(msg)
	case 3:
		return m.updateOutputSelect(msg)
	case 4:
		return m.updatePreview(msg)
	}

	return m, nil
}

func (m model) updateFileSelect(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "q" {
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	formModel, cmd := m.form.Update(msg)
	if f, ok := formModel.(*huh.Form); ok {
		m.form = f
	}
	if m.form.State == huh.StateCompleted {
		m.pdfPath = m.form.GetString("file")
		m.step = 1
		return m.showBankForm()
	}
	return m, cmd
}

func (m model) updateBankSelect(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "q" {
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	formModel, cmd := m.form.Update(msg)
	if f, ok := formModel.(*huh.Form); ok {
		m.form = f
	}
	if m.form.State == huh.StateCompleted {
		m.bank = m.form.GetString("bank")
		m.step = 2
		return m.showAccountForm()
	}
	return m, cmd
}

func (m model) updateAccountSelect(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "q" {
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	formModel, cmd := m.form.Update(msg)
	if f, ok := formModel.(*huh.Form); ok {
		m.form = f
	}
	if m.form.State == huh.StateCompleted {
		m.account = m.form.GetString("account")
		m.step = 3
		return m.showOutputForm()
	}
	return m, cmd
}

func (m model) updateOutputSelect(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "q" {
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	formModel, cmd := m.form.Update(msg)
	if f, ok := formModel.(*huh.Form); ok {
		m.form = f
	}
	if m.form.State == huh.StateCompleted {
		m.outputPath = m.form.GetString("output")
		m.step = 4
		return m.processAndPreview()
	}
	return m, cmd
}

func (m model) updatePreview(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "y", "enter":
			if m.outputPath != "" {
				err := os.WriteFile(m.outputPath, []byte(extractor.FormatJournal(m.transactions, m.cfg)), 0644)
				if err != nil {
					m.err = err
					return m, nil
				}
			}
			return m, tea.Quit
		case "n", "esc":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) showFileForm() (tea.Model, tea.Cmd) {
	var fileOpts []huh.Option[string]
	filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() && strings.HasSuffix(strings.ToLower(path), ".pdf") {
			fileOpts = append(fileOpts, huh.NewOption(path, path))
		}
		return nil
	})

	if len(fileOpts) == 0 {
		m.err = fmt.Errorf("no PDF files found in current directory")
		return m, nil
	}

	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select a PDF file").
				Options(fileOpts...).
				Key("file"),
		),
	).WithTheme(createTheme())

	m.step = 0
	return m, m.form.Init()
}

func (m model) showBankForm() (tea.Model, tea.Cmd) {
	var bankOpts []huh.Option[string]
	entries, err := os.ReadDir(m.configDir)
	if err != nil {
		m.err = fmt.Errorf("cannot read config directory '%s': %w", m.configDir, err)
		return m, nil
	}
	for _, entry := range entries {
		if entry.IsDir() {
			bankOpts = append(bankOpts, huh.NewOption(entry.Name(), entry.Name()))
		}
	}

	if len(bankOpts) == 0 {
		m.err = fmt.Errorf("no bank configurations found in '%s'", m.configDir)
		return m, nil
	}

	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select bank").
				Options(bankOpts...).
				Key("bank"),
		),
	).WithTheme(createTheme())

	return m, m.form.Init()
}

func (m model) showAccountForm() (tea.Model, tea.Cmd) {
	var accountOpts []huh.Option[string]

	accountDir := filepath.Join(m.configDir, m.bank)
	entries, err := os.ReadDir(accountDir)
	if err != nil {
		m.err = fmt.Errorf("cannot read bank directory '%s': %w", accountDir, err)
		return m, nil
	}
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".yaml") {
			name := strings.TrimSuffix(entry.Name(), ".yaml")
			accountOpts = append(accountOpts, huh.NewOption(name, name))
		}
	}

	accountOpts = append([]huh.Option[string]{huh.NewOption("(default)", "")}, accountOpts...)

	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select account type").
				Options(accountOpts...).
				Key("account"),
		),
	).WithTheme(createTheme())

	return m, m.form.Init()
}

func (m model) showOutputForm() (tea.Model, tea.Cmd) {
	defaultOutput := strings.TrimSuffix(filepath.Base(m.pdfPath), filepath.Ext(m.pdfPath)) + ".journal"

	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Output file path (leave empty for stdout)").
				Placeholder(defaultOutput).
				Key("output"),
		),
	).WithTheme(createTheme())

	return m, m.form.Init()
}

func (m model) processAndPreview() (tea.Model, tea.Cmd) {
	var err error
	m.cfg, err = config.LoadBank(m.bank, m.account)
	if err != nil {
		m.err = err
		return m, nil
	}

	pages, err := extractor.ExtractText(m.pdfPath)
	if err != nil {
		m.err = err
		return m, nil
	}

	m.transactions = extractor.ParseTransactions(pages, m.cfg)
	return m, nil
}

func createTheme() *huh.Theme {
	t := huh.ThemeCharm()
	t.Focused.Base = t.Focused.Base.BorderForeground(ColorAccent)
	t.Focused.Title = t.Focused.Title.Foreground(ColorHighlight)
	t.Focused.SelectedOption = t.Focused.SelectedOption.Foreground(ColorHighlight)
	return t
}

func (m model) View() string {
	if m.err != nil {
		return ErrorStyle.Render(fmt.Sprintf("Error: %v\n\nPress q or ctrl+c to exit.", m.err)) + "\n"
	}

	switch m.step {
	case 0, 1, 2, 3:
		if m.form != nil {
			return m.form.View()
		}
		return ""
	case 4:
		return m.renderPreview()
	}

	return ""
}

func (m model) renderPreview() string {
	var b strings.Builder

	b.WriteString(TitleStyle.Render(" Gank - Transaction Preview ") + "\n\n")

	b.WriteString(LabelStyle.Render("Source: ") + ValueStyle.Render(m.pdfPath) + "\n")
	b.WriteString(LabelStyle.Render("Bank: ") + ValueStyle.Render(m.bank))
	if m.account != "" {
		b.WriteString(" (" + m.account + ")")
	}
	b.WriteString("\n")
	b.WriteString(LabelStyle.Render("Transactions: ") + ValueStyle.Render(fmt.Sprintf("%d", len(m.transactions))) + "\n\n")

	if len(m.transactions) > 0 {
		b.WriteString(HeaderStyle.Render("Preview:") + "\n")
		previewCount := min(5, len(m.transactions))
		for i := 0; i < previewCount; i++ {
			t := m.transactions[i]
			dateStr := t.Date.Format("2006-01-02")
			absAmt := t.Amount
			if absAmt < 0 {
				absAmt = -absAmt
			}
			b.WriteString(MutedStyle.Render(fmt.Sprintf("  %s  %-40s  %.2f\n", dateStr, t.Description, absAmt)))
		}
		if len(m.transactions) > 5 {
			b.WriteString(MutedStyle.Render(fmt.Sprintf("  ... and %d more\n", len(m.transactions)-5)))
		}
		b.WriteString("\n")
	}

	if m.outputPath != "" {
		b.WriteString(SuccessStyle.Render(fmt.Sprintf("Output: %s", m.outputPath)) + "\n")
	} else {
		b.WriteString(MutedStyle.Render("Output: stdout") + "\n")
	}

	b.WriteString("\n")
	b.WriteString(BoxStyle.Render("Save to file? [y/n]"))

	return b.String()
}

func Run(configDir string) error {
	m := initialModel(configDir)
	showModel, initCmd := m.showFileForm()
	if showModel.(model).err != nil {
		p := tea.NewProgram(showModel, tea.WithAltScreen())
		_, err := p.Run()
		_ = initCmd
		return err
	}
	p := tea.NewProgram(showModel, tea.WithAltScreen())
	_, err := p.Run()
	if err != nil {
		return err
	}
	_ = initCmd
	return nil
}

func RunWithFile(pdfPath string, configDir string) error {
	m := initialModel(configDir)
	m.pdfPath = pdfPath
	m.step = 1
	showModel, cmd := m.showBankForm()
	if showModel.(model).err != nil {
		p := tea.NewProgram(showModel, tea.WithAltScreen())
		_, err := p.Run()
		_ = cmd
		return err
	}
	_ = cmd
	p := tea.NewProgram(showModel, tea.WithAltScreen())
	_, err := p.Run()
	return err
}
