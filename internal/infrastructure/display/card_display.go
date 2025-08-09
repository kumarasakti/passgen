package display

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/kumarasakti/passgen/internal/domain/entities"
)

const (
	// Card styling constants - ensure perfect symmetry with comfortable padding
	totalCardWidth   = 55  // Total width in runes (visual width) - increased for better spacing
	contentWidth     = 53  // Content width (totalCardWidth - 2 for left/right borders)  
	cornerTopLeft    = "┌"
	cornerTopRight   = "┐"
	cornerBottomLeft = "└"
	cornerBottomRight = "┘"
	horizontal       = "─"
	vertical         = "│"
	space            = " "
)

// visualWidth calculates the actual visual width of a string in the terminal
// This handles emojis and wide characters that may take 2 columns
func visualWidth(s string) int {
	width := 0
	for _, r := range s {
		switch {
		case r < 32: // Control characters
			// Don't count control characters
		case r < 127: // Basic ASCII
			width++
		case r >= 0x1F600 && r <= 0x1F64F: // Emoticons
			width += 2
		case r >= 0x1F300 && r <= 0x1F5FF: // Misc Symbols and Pictographs
			width += 2
		case r >= 0x1F680 && r <= 0x1F6FF: // Transport and Map
			width += 2
		case r >= 0x2600 && r <= 0x26FF: // Misc symbols
			width += 2
		case r >= 0x2700 && r <= 0x27BF: // Dingbats
			width += 2
		case r >= 0xFE00 && r <= 0xFE0F: // Variation selectors
			// Don't count variation selectors
		default:
			// For other characters, assume width 1 but this could be enhanced
			width++
		}
	}
	return width
}

// CardDisplayer handles the enhanced card-style display for password metadata
type CardDisplayer struct{}

// NewCardDisplayer creates a new card displayer
func NewCardDisplayer() *CardDisplayer {
	return &CardDisplayer{}
}

// DisplayPasswordCard renders password metadata in enhanced card style
func (d *CardDisplayer) DisplayPasswordCard(metadata *entities.PasswordMetadata) {
	// Create card header with service name - ensure perfect symmetry
	header := fmt.Sprintf("─ %s ", metadata.Service)
	// Calculate remaining space using Unicode-aware counting
	headerRunes := utf8.RuneCountInString(header)
	remainingRunes := contentWidth - headerRunes
	headerPadding := strings.Repeat("─", remainingRunes)
	
	fmt.Printf("%s%s%s%s\n", cornerTopLeft, header, headerPadding, cornerTopRight)
	
	// Display fields with proper spacing
	d.displayField("👤", metadata.Username)
	d.displayField("🌐", metadata.URL)
	d.displayField("📝", metadata.Notes)
	
	// Add separator line if we have content above
	if d.hasBasicContent(metadata) {
		d.displayEmptyLine()
	}
	
	// Display dates and strength on one line
	dateStrength := d.formatDateAndStrength(metadata)
	d.displayContentLine(dateStrength)
	
	// Display auto-rotation if enabled
	if metadata.AutoRotation != nil && metadata.AutoRotation.Enabled {
		rotationInfo := d.formatRotationInfo(metadata.AutoRotation)
		d.displayContentLine(rotationInfo)
	}
	
	// Close card with perfect symmetry
	fmt.Printf("%s%s%s\n", cornerBottomLeft, strings.Repeat(horizontal, contentWidth), cornerBottomRight)
	
	// Display access options
	fmt.Printf("\n🔐 passgen store get %s --copy | --show\n", metadata.Service)
}

// displayField shows a field only if it has content
func (d *CardDisplayer) displayField(icon, content string) {
	if content != "" {
		line := fmt.Sprintf("%s %s", icon, content)
		d.displayContentLine(line)
	}
}

// displayContentLine displays a line of content within the card
func (d *CardDisplayer) displayContentLine(content string) {
	// Use visual width for accurate emoji handling
	contentVisualWidth := visualWidth(content)
	maxVisualWidth := contentWidth - 4 // Account for padding (2 left + 2 right)
	
	if contentVisualWidth > maxVisualWidth {
		// Truncate to fit with ellipsis
		maxWidth := maxVisualWidth - 3 // Account for "..."
		truncated := ""
		currentWidth := 0
		
		for _, r := range content {
			runeWidth := visualWidth(string(r))
			if currentWidth + runeWidth > maxWidth {
				break
			}
			truncated += string(r)
			currentWidth += runeWidth
		}
		content = truncated + "..."
		contentVisualWidth = visualWidth(content)
	}
	
	// Calculate padding needed for right alignment
	paddingNeeded := maxVisualWidth - contentVisualWidth
	if paddingNeeded < 0 {
		paddingNeeded = 0
	}
	padding := strings.Repeat(space, paddingNeeded)
	
	fmt.Printf("%s  %s%s  %s\n", vertical, content, padding, vertical)
}

// displayEmptyLine displays an empty line within the card
func (d *CardDisplayer) displayEmptyLine() {
	padding := strings.Repeat(space, contentWidth)
	fmt.Printf("%s%s%s\n", vertical, padding, vertical)
}

// hasBasicContent checks if metadata has username, URL, or notes
func (d *CardDisplayer) hasBasicContent(metadata *entities.PasswordMetadata) bool {
	return metadata.Username != "" || metadata.URL != "" || metadata.Notes != ""
}

// formatDateAndStrength formats date and strength info on one line
func (d *CardDisplayer) formatDateAndStrength(metadata *entities.PasswordMetadata) string {
	dateStr := metadata.UpdatedAt.Format("Jan 2, 2006")
	return fmt.Sprintf("📅 %s • 💪 %s", dateStr, metadata.StrengthInfo)
}

// formatRotationInfo formats auto-rotation information
func (d *CardDisplayer) formatRotationInfo(rotation *entities.AutoRotationInfo) string {
	nextDate := rotation.NextRotation.Format("Jan 2")
	return fmt.Sprintf("🔄 Rotates every %d days (Next: %s)", rotation.IntervalDays, nextDate)
}

// DisplayPasswordList renders a list of passwords in a clean table format
func (d *CardDisplayer) DisplayPasswordList(passwords []entities.PasswordMetadata, storeName string) {
	if len(passwords) == 0 {
		fmt.Printf("📋 No passwords found in store '%s'\n", storeName)
		return
	}

	fmt.Printf("📋 Passwords in store '%s':\n", storeName)
	
	// Table headers and borders
	fmt.Println("┌──────────────────────────────────────────────────────────┐")
	fmt.Println("│ Service      │ Username     │ Updated    │ Auto-Rotation │")
	fmt.Println("├──────────────────────────────────────────────────────────┤")
	
	// Display each password entry
	for _, password := range passwords {
		service := d.truncateString(password.Service, 12)
		username := d.truncateString(password.Username, 12)
		updated := password.UpdatedAt.Format("Jan 2")
		
		var rotation string
		if password.AutoRotation != nil && password.AutoRotation.Enabled {
			rotation = fmt.Sprintf("%d days", password.AutoRotation.IntervalDays)
		} else {
			rotation = "-"
		}
		rotation = d.truncateString(rotation, 13)
		
		fmt.Printf("│ %-12s │ %-12s │ %-10s │ %-13s │\n", 
			service, username, updated, rotation)
	}
	
	fmt.Println("└──────────────────────────────────────────────────────────┘")
	fmt.Println("\n💡 Use 'passgen store get <service>' to view details")
}

// DisplayRotationStatus shows auto-rotation status for passwords
func (d *CardDisplayer) DisplayRotationStatus(statuses []entities.RotationStatus, storeName string) {
	if len(statuses) == 0 {
		fmt.Printf("🔄 No auto-rotation passwords in store '%s'\n", storeName)
		return
	}

	fmt.Printf("🔄 Auto-rotation status for store '%s':\n", storeName)
	
	// Table headers and borders
	fmt.Println("┌──────────────────────────────────────────────────────────────┐")
	fmt.Println("│ Service      │ Status     │ Next Rotation │ Days Until     │")
	fmt.Println("├──────────────────────────────────────────────────────────────┤")
	
	// Display each rotation status
	for _, status := range statuses {
		service := d.truncateString(status.Service, 12)
		statusIcon := d.getStatusIcon(status.Status)
		statusText := d.truncateString(fmt.Sprintf("%s %s", statusIcon, status.Status), 10)
		nextRotation := status.NextRotation.Format("Jan 2")
		daysUntil := fmt.Sprintf("%d days", status.DaysUntilNext)
		
		fmt.Printf("│ %-12s │ %-10s │ %-13s │ %-14s │\n", 
			service, statusText, nextRotation, daysUntil)
	}
	
	fmt.Println("└──────────────────────────────────────────────────────────────┘")
}

// truncateString truncates a string to fit within maxWidth, adding ellipsis if needed
func (d *CardDisplayer) truncateString(s string, maxWidth int) string {
	if utf8.RuneCountInString(s) <= maxWidth {
		return s
	}
	
	runes := []rune(s)
	if len(runes) <= maxWidth-3 {
		return s
	}
	
	return string(runes[:maxWidth-3]) + "..."
}

// getStatusIcon returns appropriate icon for rotation status
func (d *CardDisplayer) getStatusIcon(status string) string {
	switch status {
	case "Due":
		return "🔴"
	case "Soon":
		return "🟡"
	case "Good":
		return "🟢"
	default:
		return "⚪"
	}
}

// DisplayPasswordBox displays the actual password in a secure box format
func (d *CardDisplayer) DisplayPasswordBox(password string) {
	// Create a symmetric box for the password
	passwordWidth := utf8.RuneCountInString(password)
	boxWidth := passwordWidth + 4 // 2 spaces padding on each side
	
	// Ensure minimum width for better appearance
	if boxWidth < 20 {
		boxWidth = 20
	}
	
	contentPadding := boxWidth - 2 // subtract borders
	
	// Top border
	fmt.Printf("┌%s┐\n", strings.Repeat("─", contentPadding))
	
	// Content with password
	leftPadding := (contentPadding - passwordWidth) / 2
	rightPadding := contentPadding - passwordWidth - leftPadding
	fmt.Printf("│%s%s%s│\n", 
		strings.Repeat(" ", leftPadding), 
		password, 
		strings.Repeat(" ", rightPadding))
	
	// Bottom border
	fmt.Printf("└%s┘\n", strings.Repeat("─", contentPadding))
}
