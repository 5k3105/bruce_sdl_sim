package main

import (
	"io/ioutil"
)

var (
	mem_ram_    map[uint16]uint8
	screen_mem_ uint16 = 0x0400
	char_mem_   uint16 = 0xd000

	vic_cols   uint16 = 40
	color_mem_ uint16 = 0xd800
)

func draw_screen2() {
	cols := vic_cols
	lines := 25 * 8

	for line := 0; line < lines; line++ {
		y := line
		// screen_draw_rect(0,y,320,1)
		for column := 0; column < int(cols); column++ {
			x := column * 8
			row := line / 8
			char_row := line % 8
			c := get_screen_char(column, row)
			data := get_char_data(c, char_row)
			color := get_char_color(column, row)
			draw_mcchar(x, y, data, (color & 0x7))
		}
	}
}

func ISSET_BIT(v, b uint8) bool {
	return (v & (1 << uint8(b))) != 0
}

func draw_mcchar(x, y int, data, color uint8) {
	/// $FB,$F1,$FE,$FB,$F1,$FE,$FB,$F2,$F1,$FC,$F2,$F1
	/// bgcolor := [3]byte{0xFC, 0xF2, 0xF1}

	bgcolor := [3]byte{0xFB, 0xF2, 0xF1}

	if y < 5*8 {
		bgcolor = [3]byte{0xFB, 0xF1, 0xFE}
	}

	for i := uint8(0); i < 4; i++ {
		var c uint8                              /// color
		var cs uint8 = ((data >> (i * 2)) & 0x3) /// color source
		switch cs {
		case 0:
			c = bgcolor[0]
		case 1:
			c = bgcolor[1]
		case 2:
			c = bgcolor[2]
		case 3:
			c = color
		}
		xoffs := x + 8 - int(i)*2
		screen_update_pixel(xoffs, y, int(c))
		screen_update_pixel(xoffs+1, y, int(c))
	}
}

func draw_char(x, y int, data, color uint8) {
	for i := 0; i < 8; i++ {
		xoffs := x + 8 - i
		/* draw pixel */
		if ISSET_BIT(data, uint8(i)) {
			screen_update_pixel(xoffs, y, int(color))
		}
	}
}

func screen_draw_rect(x, y, n, color int) {
	for i := 0; i < n; i++ {
		screen_update_pixel(x+i, y, color)
	}
}

func draw_screen() {
	cols := vic_cols
	lines := 25 * 8

	for line := 0; line < lines; line++ {
		y := line
		// screen_draw_rect(0,y,320,1)
		for column := 0; column < int(cols); column++ {
			x := column * 8
			row := line / 8
			char_row := line % 8
			c := get_screen_char(column, row)
			data := get_char_data(c, char_row)
			color := get_char_color(column, row)
			draw_mcchar(x, y, data, (color & 0x7))

			/*
				if ISSET_BIT(color, 3) {
					draw_mcchar(x, y, data, (color & 0x7))
				} else {
					draw_char(x, y, data, color)
				}
			*/

		}
	}
}

func get_screen_char(column, row int) uint8 {
	addr := screen_mem_ + (uint16(row) * vic_cols) + uint16(column)
	return vic_read_byte(addr)
}

func get_char_data(chr uint8, line int) uint8 {
	addr := char_mem_ + (uint16(chr) * 8) + uint16(line)
	return vic_read_byte(addr)
}

func vic_read_byte(addr uint16) uint8 {
	return read_byte_no_io(addr) // fix all direct
	//return read_byte_no_io(vic_addr + (addr & 0x3fff))
}

func get_char_color(column, row int) uint8 {
	addr := color_mem_ + (uint16(row) * vic_cols) + uint16(column)
	return (read_byte_no_io(addr) & 0x0f)
}

func read_byte_no_io(addr uint16) uint8 {
	return mem_ram_[addr]
}

func write_byte_no_io(addr uint16, v uint8) {
	mem_ram_[addr] = v
}

func load_ram(filename string, baseaddr uint16) {
	file, err := ioutil.ReadFile("data/" + filename) /// ([]byte, error)
	if err != nil {
		println(err.Error)
	}

	for inc, data := range file {
		mem_ram_[baseaddr+uint16(inc)] = data
	}
}

/** VIC banking
 * @brief retrieves vic base address
 *
 * PRA bits (0..1)
 *
 *  %00, 0: Bank 3: $C000-$FFFF, 49152-65535
 *  %01, 1: Bank 2: $8000-$BFFF, 32768-49151
 *  %10, 2: Bank 1: $4000-$7FFF, 16384-32767
 *  %11, 3: Bank 0: $0000-$3FFF, 0-16383 (standard)
 */

/*
uint16_t Cia2::vic_base_address()
{
  return ((~pra_&0x3) << 14);
}
*/

/*

func horizontal_scroll() int {
	cr2_ := 0
	return cr2_ & 0x7
}


	/// default memory pointers
  screen_mem_ = Memory::kBaseAddrScreen;
  char_mem_   = Memory::kBaseAddrChars;
  bitmap_mem_ = Memory::kBaseAddrBitmap;

    static const uint16_t kBaseAddrScreen = 0x0400;
    static const uint16_t kBaseAddrChars  = 0xd000;
    static const uint16_t kBaseAddrBitmap = 0x0000;
    static const uint16_t kBasecolor_mem_ = 0xd800;

    static const uint16_t kcolor_mem_ = 0xd800;
*/

/// type Memory

/*
func vic_read_byte(addr uint16) uint8 {
	var v uint8                                           //uint8_t v
	//vic_addr := vic_base_address() + (addr & 0x3fff) /// uint16_t

	if (vic_addr >= 0x1000 && vic_addr < 0x2000) || (vic_addr >= 0x9000 && vic_addr < 0xa000) {
		v = mem_rom_[BaseAddrChars+(vic_addr&0xfff)]
	} else {
		v = read_byte_no_io(vic_addr)
	}
	return v
}
*/

/*
func draw_raster_char_mode() {
	kMCCharMode := true
	graphic_mode_ := true

	y := 8
	/// io_->screen_draw_rect(kGFirstCol,y,kGResX,bgcolor_[0]);

	for column := 0; column < int(kGCols); column++ {
		x := column * 8
		line := 0
		row := line / 8
		char_row := line % 8

		c := get_screen_char(column, row)

		data := get_char_data(c, char_row)

		color := get_char_color(column, row)

		println(x, y, data, (color & 0x7))

		if graphic_mode_ == kMCCharMode && ISSET_BIT(color, 3) {
			draw_mcchar(x, y, data, (color & 0x7))
		} else {
			draw_char(x, y, data, color)
		}
	}
}
*/
