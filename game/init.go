package game

import (
	"embed"
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/tristanisham/cope/meta"
	"golang.org/x/image/math/f64"
)

type Game struct {
	showRays bool
	px, py   float64
	objects  []object
	assets   embed.FS
	camera   Camera
	world    *ebiten.Image
	// lc       int // level cursor
	// levels   []level
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return ErrGameExit
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.showRays = !g.showRays
	}

	g.handleMovement()

	// Update the camera position to follow the player
	g.camera.Position[0] = g.px - g.camera.ViewPort[0]/2
	g.camera.Position[1] = g.py - g.camera.ViewPort[1]/2
	
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// potential good link for writing tiles
	// https://github.com/hajimehoshi/ebiten/blob/main/examples/camera/main.go#L75

	shadowImage.Fill(color.Black)
	rays := rayCasting(float64(g.px), float64(g.py), g.objects)

	// Subtract ray triangles from shadow
	opt := &ebiten.DrawTrianglesOptions{}
	opt.Address = ebiten.AddressRepeat
	opt.Blend = ebiten.BlendSourceOut
	for i, line := range rays {
		nextLine := rays[(i+1)%len(rays)]

		// Draw triangle of area between rays
		v := rayVertices(float64(g.px), float64(g.py), nextLine.x2, nextLine.y2, line.x2, line.y2)
		shadowImage.DrawTriangles(v, []uint16{0, 1, 2}, triangleImage, opt)
	}

	// Draw background
	screen.DrawImage(bgImage, nil)

	if g.showRays {
		// Draw rays
		for _, r := range rays {
			vector.StrokeLine(screen, float32(r.x1), float32(r.y1), float32(r.x2), float32(r.y2), 1, color.RGBA{255, 255, 0, 150}, true)
		}
	}

	// Draw shadow
	op := &ebiten.DrawImageOptions{}
	op.ColorScale.ScaleAlpha(0.7)
	screen.DrawImage(shadowImage, op)

	// Draw walls
	for _, obj := range g.objects {
		for _, w := range obj.walls {
			vector.StrokeLine(screen, float32(w.x1), float32(w.y1), float32(w.x2), float32(w.y2), 1, color.RGBA{255, 0, 0, 255}, true)
		}
	}

	// Draw player as a rect
	vector.DrawFilledRect(screen, float32(g.px)-2, float32(g.py)-2, 6, 6, color.Black, true)
	vector.DrawFilledRect(screen, float32(g.px)-1, float32(g.py)-1, 4, 4, color.RGBA{255, 100, 100, 255}, true)

	if g.showRays {
		ebitenutil.DebugPrintAt(screen, "R: hide rays", meta.Padding, 0)
	} else {
		ebitenutil.DebugPrintAt(screen, "R: show rays", meta.Padding, 0)
	}

	g.camera.Render(g.world, screen)

	worldX, worldY := g.camera.ScreenToWorld(ebiten.CursorPosition())

	ebitenutil.DebugPrintAt(screen, "WASD: move", 160, 0)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("TPS: %0.2f", ebiten.ActualTPS()), 51, 51)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Rays: 2*%d", len(rays)/2), meta.Padding, 222)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Pos: (%0.0f,%0.0f)", g.px, g.py), meta.Padding, 233)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Cam: (%0.0f,%0.0f)", worldX, worldY), meta.Padding, 244)


}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return meta.ScreenWidth, meta.ScreenHeight
}

func NewGame(assets embed.FS) *Game {
	g := &Game{
		px:     meta.ScreenWidth / 2,
		py:     meta.ScreenHeight / 2,
		assets: assets,
		camera: Camera{ViewPort: f64.Vec2{meta.ScreenWidth, meta.ScreenHeight}},
	}

	g.loadLevels()

	g.world = ebiten.NewImage(meta.WorldWidth, meta.WorldHeight)

	// Add outer walls
	g.objects = append(g.objects, object{rect(meta.Padding, meta.Padding, meta.ScreenWidth-2*meta.Padding, meta.ScreenHeight-2*meta.Padding)})

	// Angled wall
	g.objects = append(g.objects, object{[]line{{50, 110, 100, 150}}})

	// Rectangles
	g.objects = append(g.objects, object{rect(45, 50, 70, 20)})
	g.objects = append(g.objects, object{rect(150, 50, 30, 60)})

	return g
}

func init() {
	bgImage.Fill(color.RGBA{0, 0, 139, 255})
	triangleImage.Fill(color.White)
}
