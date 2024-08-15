package game

func (g *Game) clear() {
	g.objects = g.objects[:0]
}