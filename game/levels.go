package game

type level struct {
	ipx, ipy float64 // initial player (x,y)
}

func (g *Game) clear() {
	g.objects = g.objects[:0]
}

func (g *Game) loadLevels() {

}