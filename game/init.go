package game

import (
	"embed"
	"fmt"
	"image"
	"image/color"
	"math"

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
	g.world.Clear()
	shadowImage.Fill(color.Black)

	
	// potential good link for writing tiles
	// https://github.com/hajimehoshi/ebiten/blob/main/examples/camera/main.go#L75

	// Draw background
	screen.Fill(color.Black)
	screen.DrawImage(bgImage, &ebiten.DrawImageOptions{})

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

	if g.showRays {
		// Draw rays
		for _, r := range rays {
			vector.StrokeLine(g.world, float32(r.x1), float32(r.y1), float32(r.x2), float32(r.y2), 1, color.RGBA{255, 255, 0, 150}, true)
		}
	}

	// Draw shadow
	ox := &ebiten.DrawImageOptions{}
	ox.ColorScale.ScaleAlpha(0.7)
	g.world.DrawImage(shadowImage, ox)

	// Draw walls
	for _, obj := range g.objects {
		for _, w := range obj.walls {
			vector.StrokeLine(g.world, float32(w.x1), float32(w.y1), float32(w.x2), float32(w.y2), 1, color.RGBA{255, 0, 0, 255}, true)
		}
	}

	

	// Draw player as a rect
	vector.DrawFilledRect(g.world, float32(g.px)-2, float32(g.py)-2, 6, 6, color.Black, true)
	vector.DrawFilledRect(g.world, float32(g.px)-1, float32(g.py)-1, 4, 4, color.RGBA{255, 100, 100, 255}, true)

	g.camera.Render(g.world, screen)

	if g.showRays {
		ebitenutil.DebugPrintAt(g.world, "R: hide rays", meta.Padding, 0)
	} else {
		ebitenutil.DebugPrintAt(g.world, "R: show rays", meta.Padding, 0)
	}

	worldX, worldY := g.camera.ScreenToWorld(int(g.px), int(g.py))

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
		camera: Camera{
			ViewPort: f64.Vec2{meta.ScreenWidth, meta.ScreenHeight},
			Position: f64.Vec2{
				meta.ScreenWidth / 2,
				meta.ScreenHeight / 2,
			},
		},
	}

	alphas := image.Pt(meta.FOV*2, meta.FOV*2)
	a := image.NewAlpha(image.Rectangle{image.Pt(0, 0), alphas})
	for j := 0; j < alphas.Y; j++ {
		for i := 0; i < alphas.X; i++ {
			// d is the distance between (i, j) and the (circle) center.
			d := math.Sqrt(float64((i-meta.FOV)*(i-meta.FOV) + (j-meta.FOV)*(j-meta.FOV)))
			// Alphas around the center are 0 and values outside of the circle are 0xff.
			b := uint8(max(0, min(0xff, int(3*d*0xff/meta.FOV)-2*0xff)))
			a.SetAlpha(i, j, color.Alpha{b})
		}
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
