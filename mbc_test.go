package gbMBC

import "testing"

const (
	DrMario = "roms/mbc0/dr_mario.gb"
	Tetris  = "roms/mbc0/tetris.gb"
)

func TestMBC_Title(t *testing.T) {
	var mbc = NewMBC(DrMario)
	if mbc.Title() != "DR.MARIO" {
		t.Errorf("Expected DR.MARIO, got %s", mbc.Title())
	}
	mbc = NewMBC(Tetris)
	if mbc.Title() != "TETRIS" {
		t.Errorf("Expected TETRIS, got %s", mbc.Title())
	}
}
