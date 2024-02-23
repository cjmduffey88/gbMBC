package gbMBC

import (
	"fmt"
	"os"
	"strings"
)

const (
	MBC0 = uint8(iota)
	MBC1
	MBC2
	MBC3
	MBC5
	MBC6
	MBC7
)

type MBC struct {
	ROM [][]byte // ROM banks
	RAM [][]byte // RAM banks

	ROMBank0 uint8 // Active ROM bank 0, Only used in MBC1, usually 0x00
	ROMBank1 uint8 // Active ROM bank 1

	RAMEnable bool  // If false, RAM writes are ignored
	RAMBank   uint8 // Active RAM bank
}

func NewMBC(pathROM string) *MBC {
	var mbc = new(MBC)
	mbc.ROM = make([][]byte, 0)
	mbc.RAM = make([][]byte, 0)
	mbc.ROMBank1 = 0x01

	// Read ROM file
	rom, err := os.ReadFile(pathROM)
	if err != nil {
		panic(err)
	}

	// Split ROM into 16KB banks
	for i := 0; i < len(rom); i += 0x4000 {
		mbc.ROM = append(mbc.ROM, rom[i:i+0x4000])
	}

	// Create 16 RAM banks of 8kb each
	for i := 0; i < 16; i++ {
		mbc.RAM = append(mbc.RAM, make([]byte, 0x2000))
	}

	return mbc
}

func (mbc *MBC) Read(address uint16) byte {
	switch {
	// ROM Bank 0
	case address < 0x4000:
		return mbc.ROM[mbc.ROMBank0][address]
	// ROM Bank 1
	case address < 0x8000:
		return mbc.ROM[mbc.ROMBank1][address-0x4000]
	// RAM Bank
	case address >= 0xA000 && address <= 0xBFFF:
		return mbc.RAM[mbc.RAMBank][address-0xA000]
	// Invalid Read
	default:
		panic(fmt.Sprintf("MBC: invalid read at 0x%04X", address))
	}
}

func (mbc *MBC) Write(address uint16, value byte) {
	switch mbc.mbcType() {
	case MBC0:
		// Writing to RAM
		if address >= 0xA000 && address <= 0xBFFF {
			// Does this cartridge have RAM?
			switch mbc.ROM[0][0x0147] {
			case 0x08, 0x09:
				mbc.RAM[mbc.RAMBank][address-0xA000] = value
			}
		}

	case MBC1:
		panic(fmt.Sprintf("MBC: MBC1 not implemented"))

	case MBC2:
		panic(fmt.Sprintf("MBC: MBC2 not implemented"))

	case MBC3:
		panic(fmt.Sprintf("MBC: MBC3 not implemented"))

	case MBC5:
		panic(fmt.Sprintf("MBC: MBC5 not implemented"))

	case MBC6:
		panic(fmt.Sprintf("MBC: MBC6 not implemented"))

	case MBC7:
		panic(fmt.Sprintf("MBC: MBC7 not implemented"))

	// Invalid Write
	default:
		panic(fmt.Sprintf("Write (MBC): invalid write at 0x%04X", address))
	}

}

func (mbc *MBC) Title() string {
	// convert title to string and replace null bytes
	return strings.Replace(string(mbc.ROM[0][0x0134:0x0143]), "\x00", "", -1)
}

func (mbc *MBC) hasRAM() bool {
	// Checks RAM Size according to header, 0 is no ram, 1 was never used
	switch mbc.ROM[0][0x0149] {
	case 0x00, 0x01:
		return false
	default:
		return true
	}
}

func (mbc *MBC) hasBattery() bool {
	// Checks cartridge type to determine if battery backed RAM is present
	switch mbc.ROM[0][0x0147] {
	case 0x03, 0x06, 0x09, 0x0D, 0x0F, 0x10, 0x13, 0x1B, 0x1E, 0x22, 0xFF:
		return true
	default:
		return false
	}
}

func (mbc *MBC) hasTimer() bool {
	// Checks cartridge type to determine if timer is present
	switch mbc.ROM[0][0x0147] {
	case 0x0F, 0x10:
		return true
	default:
		return false
	}
}

func (mbc *MBC) hasRumble() bool {
	// Checks cartridge type to determine if rumble is present
	switch mbc.ROM[0][0x0147] {
	case 0x1C, 0x1D, 0x1E, 0x22:
		return true
	default:
		return false
	}
}

func (mbc *MBC) mbcType() uint8 {
	switch mbc.ROM[0][0x0147] {
	case 0x00, 0x08, 0x09:
		return MBC0
	case 0x01, 0x02, 0x03:
		return MBC1
	case 0x05, 0x06:
		return MBC2
	case 0x0F, 0x10, 0x11, 0x12, 0x13:
		return MBC3
	case 0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E:
		return MBC5
	case 0x20, 0x22:
		return MBC6
	case 0x23, 0x25:
		return MBC7
	default:
		panic(fmt.Sprintf("MBC: invalid type 0x%02X", mbc.ROM[0][0x0147]))
	}
}
