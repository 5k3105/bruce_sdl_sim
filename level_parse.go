package main

import (
	//"encoding/hex"
	"fmt"
	"io/ioutil"
)

type level struct {
	compressed []byte
	screen     [1000]byte
	color      [1000]byte
	lanterns   []byte
	code       []byte
}

var levelmaps map[int]level

func load_level_mem(lvl int) {
	for inc, data := range levelmaps[lvl].screen {
		mem_ram_[screen_mem_+uint16(inc)] = data
	}
	//for inc, _ := range levelmaps[lvl].color {
	//	mem_ram_[color_mem_+uint16(inc)] = uint8(inc) //data
	//}

	for inc, data := range levelmaps[lvl].color {
		mem_ram_[color_mem_+uint16(inc)] = data
	}

	/*
		for inc, data := range levelmaps[0].color {
			mem_ram_[screen_mem_+uint16(inc)] = data
		}
		for inc, _ := range levelmaps[0].screen {
			mem_ram_[color_mem_+uint16(inc)] = uint8(inc) //data
		}
	*/
}

func load_level_data(filename string) {
	file, err := ioutil.ReadFile("data/" + filename)
	if err != nil {
		println(err.Error)
	}

	parse_levels(file)

}

func parse_levels(bytes []byte) {
	var (
		r                  bool
		ccol, cols         uint16 = 40, 40
		clevel             int
		soffs, boffs, offs uint16
		x, c, b, s         uint8
	)

	var screen, color = [1000]byte{}, [1000]byte{}
	levelmaps = make(map[int]level)

top:
	if int(boffs) > len(bytes) {
		return
	}

	if boffs == 0x10f {
		ro := fmt.Sprintf("%#x", boffs)
		h := fmt.Sprintf("%#x", b)
		fmt.Printf("10f:  %s  %s\n", h, ro)
	}

	b = bytes[boffs] /// 16 8D
	/*
		fmt.Println()
		o := fmt.Sprintf("%#x", boffs + 0x68F9) /// next level: 6A09
		ro := fmt.Sprintf("%#x", boffs)
		h := fmt.Sprintf("%#x", b)
		fmt.Printf("top: %s   %s  %s\n", o, h, ro)
	*/
	if !ISSET_BIT(b, 7) { /// bpl
		r = true
		goto init
	}
	if b == 0x80 { //00
		goto next_level
	}
	b &= 0x7F

	r = false

init:
	x = b
	/*
		{
			ro := fmt.Sprintf("%#x", boffs)
			h := fmt.Sprintf("%#x", x)
			fmt.Printf(":::  %s  %s\n", ro, h)
		}
	*/
	s = b
	boffs++

data_read:
	b = bytes[boffs] /// 8A
	if !ISSET_BIT(b, 7) {
		goto write_screen
	}
	c = 0xF8
	color[soffs] = c
	color[soffs+cols] = c
	///mem_ram_[color_mem+soffs] = c
	///mem_ram_[color_mem+soffs+cols] = c
	b &= 0x7F

write_screen:
	screen[soffs] = b
	screen[soffs+cols] = b | 0x80
	///mem_ram_[screen_mem+soffs] = b
	///mem_ram_[screen_mem+soffs+cols] = b || 0x80
	if !r {
		boffs++
	}
	soffs++

	ccol--
	if ccol == 0 {

		println(s, boffs, soffs)
		//fmt.Printf("%s", hex.Dump(screen[soffs-uint16(cols):soffs]))

		ccol = cols
		soffs += cols

		println("row 2:")
		//fmt.Printf("%s", hex.Dump(screen[soffs-uint16(cols):soffs]))

	}

	x--
	if x != 0 {
		goto data_read
	}

	//	println(s, boffs, soffs)
	//	fmt.Printf("%s", hex.Dump(screen[soffs-uint16(s):soffs]))

	if r {
		boffs++
	}

	goto top

next_level:
	println("--- next level: ", clevel)
	lvl := level{
		compressed: bytes[offs:boffs],
		screen:     screen,
		color:      color,
	}
	println(offs, boffs)
	levelmaps[clevel] = lvl
	soffs = 0
	offs = boffs
	screen, color = [1000]byte{}, [1000]byte{}

	clevel++
	if clevel == 14 {
		return
	}
	boffs++
	goto top

}

/* level offsets
68F9
6A09
6B3E
6C9D
6E1B
6F82
70DD
7227
7311
7491
75D8
774C
7843
79C6
7A46
7BB4
7C5C
7D8D
7E90
7FE0
*/
