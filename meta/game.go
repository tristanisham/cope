// meta exists as a consumable package for values and constants used by multiple packages as multiple levels
// its expressed purpose is to avoid pollution.
package meta

const (
	ScreenWidth = 640
	ScreenHeight = 480
	TileSize = 23
	TitleFontSize = FontSize * 1.5
	FontSize = 24
	SmallFontSize = FontSize / 2
	Padding = 20
)