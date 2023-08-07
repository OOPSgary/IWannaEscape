package app

/*
func (g *Game) _syncCharacter() {
	defer g.mainWorld.MainCharacter.Obj.Update()
	dx, dy := g.mainWorld.MainCharacter.SpeedX, g.mainWorld.MainCharacter.SpeedY
	for _, sl := range strikeList {
		if sl.Online {
			if contactSet := g.mainWorld.MainCharacter.Obj.Shape.Intersection(dx, dy, sl.Shape); contactSet != nil {
				Dead = true
			}
		}
	}
	if col := g.mainWorld.MainCharacter.Obj.Check(dx, dy, "Stopper"); col != nil {
		if co := col.ContactWithObject(col.Objects[0]); co != nil {
			if g.mainWorld.MainCharacter.SpeedY != 0 {
				log.Println(g.mainWorld.MainCharacter.SpeedY)
				if (col.Objects[0].Y <= g.mainWorld.MainCharacter.Obj.Y+g.mainWorld.MainCharacter.Obj.H+g.mainWorld.MainCharacter.SpeedY) ||
					(col.Objects[0].Y+col.Objects[0].H >= g.mainWorld.MainCharacter.Obj.Y+g.mainWorld.MainCharacter.SpeedY) {
					// g.mainWorld.MainCharacter.Obj.X += co.X()
					// g.mainWorld.MainCharacter.SpeedX = 0
					// log.Println(co.Y())
					g.mainWorld.MainCharacter.Obj.Y += co.Y()
					g.mainWorld.MainCharacter.Obj.X += dx
					g.mainWorld.MainCharacter.SpeedY = 0
					g.mainWorld.MainCharacter.Jump.Reset()
					return
				}
				// g.mainWorld.MainCharacter.SpeedX = 0
				g.mainWorld.MainCharacter.Obj.X += co.X()
				g.mainWorld.MainCharacter.Obj.Y += g.mainWorld.MainCharacter.SpeedY
				// g.mainWorld.MainCharacter.Obj.Y += co.Y()
				// g.mainWorld.MainCharacter.SpeedY = 0
				// g.mainWorld.MainCharacter.SpeedX = 0
				return
			}
			g.mainWorld.MainCharacter.SpeedX += co.X()
			dx = g.mainWorld.MainCharacter.SpeedX

		} else if slide := col.SlideAgainstCell(col.Cells[0], "Stopper"); slide != nil && math.Abs(slide.X()) <= 8 {
			// If we are able to slide here, we do so. No contact was made, and vertical speed (dy) is maintained upwards.
			g.mainWorld.MainCharacter.Obj.X += slide.X()
			// g.mainWorld.MainCharacter.Obj.Y += slide.Y()
			log.Println("Slide Against")
			return
		}

	}
	g.mainWorld.MainCharacter.Obj.X += dx
	g.mainWorld.MainCharacter.Obj.Y += dy
	if g.mainWorld.MainCharacter.Obj.Check(dx, dy+1, "Stopper") != nil {
		return
	}
	if g.mainWorld.MainCharacter.SpeedY <= 5 && g.mainWorld.MainCharacter.SpeedY >= 0 {
		g.mainWorld.MainCharacter.SpeedY += 0.25

	} else if g.mainWorld.MainCharacter.SpeedY <= 0 {
		g.mainWorld.MainCharacter.SpeedY += 0.2
	}
	// y = g.mainWorld.MainCharacter.SpeedY

	// }

}
*/