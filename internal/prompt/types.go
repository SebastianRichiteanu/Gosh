package prompt

const (
	runeCtrlL     = 12  // Ctrl+L (clear screen)
	runeEnter     = 13  // Enter (Carriage Return)
	runeBackspace = 127 // Backspace/Delete
	runeTab       = 9   // Tab
	runeEscape    = 27  // Escape (starts CSI sequence)
	runeBracket   = 91  // '[' following ESC in CSI sequences

	runeArrowUp    = 65 // 'A' after ESC [
	runeArrowDown  = 66 // 'B' after ESC [
	runeArrowRight = 67 // 'C' after ESC [
	runeArrowLeft  = 68 // 'D' after ESC [

	myRuneArrowUp    = -1000 // Custom value for up arrow
	myRuneArrowDown  = -1001 // Custom value for down arrow
	myRuneArrowRight = -1002 // Custom value for right arrow
	myRuneArrowLeft  = -1003 // Custom value for left arrow
)
