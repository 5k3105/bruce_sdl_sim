package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"os"
	"time"
	"unsafe"
	//"github.com/tfriedel6/canvas"
	//"github.com/tfriedel6/canvas/sdlcanvas"
)

const (
	winTitle            string = "VIC II"
	VisibleScreenWidth  int32  = 403
	VisibleScreenHeight int32  = 284
	cols                       = VisibleScreenWidth
	rows                       = VisibleScreenHeight
	sizeof_uint32              = 4
	m                          = 3
)

var (
	window        *sdl.Window
	renderer      *sdl.Renderer
	texture       *sdl.Texture
	frame         [rows * cols]uint32
	color_palette [16]uint32
	format        *sdl.PixelFormat
)

func run() int {
	frame = [rows * cols]uint32{}
	/*
		var err error
		wnd, cv, err = sdlcanvas.CreateWindow(screenw, screenh, "Tile Map")
		if err != nil {
			log.Println(err)
			return
		}
		defer wnd.Destroy()
	*/

	var err error
	if err = sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize SDL: %s\n", err)
		os.Exit(1)
	}
	defer sdl.Quit()

	if window, err = sdl.CreateWindow(winTitle, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, VisibleScreenWidth*m, VisibleScreenHeight*m, sdl.WINDOW_SHOWN); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create window: %s\n", err)
		os.Exit(2)
	}
	defer window.Destroy()

	if renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create renderer: %s\n", err)
		os.Exit(3) // don't use os.Exit(3); otherwise, previous deferred calls will never run
	}
	renderer.Clear()
	defer renderer.Destroy()

	texture, err = renderer.CreateTexture(sdl.PIXELFORMAT_ARGB8888, sdl.TEXTUREACCESS_STREAMING, cols, rows)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create texture: %s\n", err)
		os.Exit(3)
	}

	format, err = sdl.AllocFormat(sdl.PIXELFORMAT_ARGB8888)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create format: %s\n", err)
		os.Exit(3)
	}

	color_palette = init_color_palette(format)
	/*
		n := 16
		for i := 0; i < n; i++ {
			screen_draw_rect(10+i, 10+i, 60, i)
		}
	*/

	/// screen_mem_  uint16 = 0x0400
	/// AddrColorRAM uint16 = 0xd800
	/// char_mem_    uint16 = 0xd000

	mem_ram_ = make(map[uint16]uint8)

	//mapnumber := "1"

	//load_ram("blmap"+mapnumber, screen_mem_)
	//load_ram("blcol"+mapnumber, color_mem_)

	load_level_data("LevelData")
	println("# maps: ", len(levelmaps))
	//mem_ram_[screen_mem_+uint16(inc)] = data
	//mem_ram_[baseaddr+uint16(inc)] = data

	lvl := 6
	//delay := 3
	chrset_number := "1"

	for i := 5; i < lvl; i++ {
		switch i {
		//case < 10 :
		case 12:
			chrset_number = "3"
		case 11, 13:
			chrset_number = "2"
			//default:
			//	chrset_number = "1"
		}
		load_ram("chrset"+chrset_number, char_mem_)
		load_level_mem(i)
		draw_screen2()
		screen_refresh()
		//duration := time.Second
		duration := time.Duration(3) * time.Second
		time.Sleep(duration)
	}

	//draw_screen()
	//draw_raster_char_mode()

	/*
		for i:=0;i < 90;i++{
			draw_mcchar(50,10+i,uint8(i),3)
		}
	*/

	//sdl.Delay(3000)

	/* 	 /// screen shot*/
	pixels := unsafe.Pointer(&frame)
	pitch := int(cols) * sizeof_uint32

	surface, err := sdl.CreateRGBSurfaceWithFormatFrom(pixels, VisibleScreenWidth, VisibleScreenHeight, 1, int32(pitch), sdl.PIXELFORMAT_ARGB8888)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create surface: %s\n", err)
		os.Exit(3)
	}

	//r := &sdl.Rect{10,10,10*16+60,10*16+60}
	//y := surface.SetClipRect(r)
	//if y {
	surface.SaveBMP("xyz.bmp")
	//}

	return 0
}

func screen_draw_rect_(x, y, n, color int) {
	for i := 0; i <= n; i++ {
		screen_update_pixel(x+i, y, color)
		screen_update_pixel(x, y+i, color)
		screen_update_pixel(x+i, y+n, color)
		screen_update_pixel(x+n, y+i, color)
	}
}

func screen_update_pixel(x, y, color int) {
	frame[y*int(cols)+x] = color_palette[color&0xf]
}

func screen_refresh() {
	texture.Update(nil, (*[sizeof_uint32]byte)(unsafe.Pointer(&frame))[:], int(cols)*sizeof_uint32)
	renderer.Clear()
	renderer.Copy(texture, nil, nil)
	renderer.Present()

	/* process SDL events once every frame -- keyboard */
	//process_events();
	/* perform vertical refresh sync */
	//vsync();
}

func init_color_palette(format_ *sdl.PixelFormat) [16]uint32 {
	color_palette := [16]uint32{}
	color_palette[0] = sdl.MapRGB(format_, 0x00, 0x00, 0x00)  /// black
	color_palette[1] = sdl.MapRGB(format_, 0xff, 0xff, 0xff)  /// white
	color_palette[2] = sdl.MapRGB(format_, 0xab, 0x31, 0x26)  /// red
	color_palette[3] = sdl.MapRGB(format_, 0x66, 0xda, 0xff)  /// cyan
	color_palette[4] = sdl.MapRGB(format_, 0xbb, 0x3f, 0xb8)  /// violet/purple
	color_palette[5] = sdl.MapRGB(format_, 0x55, 0xce, 0x58)  /// green
	color_palette[6] = sdl.MapRGB(format_, 0x1d, 0x0e, 0x97)  /// blue
	color_palette[7] = sdl.MapRGB(format_, 0xea, 0xf5, 0x7c)  /// yellow
	color_palette[8] = sdl.MapRGB(format_, 0xb9, 0x74, 0x18)  /// orange
	color_palette[9] = sdl.MapRGB(format_, 0x78, 0x53, 0x00)  /// brown
	color_palette[10] = sdl.MapRGB(format_, 0xdd, 0x93, 0x87) /// lt red
	color_palette[11] = sdl.MapRGB(format_, 0x5b, 0x5b, 0x5b) /// dk grey
	color_palette[12] = sdl.MapRGB(format_, 0x8b, 0x8b, 0x8b) /// grey 2
	color_palette[13] = sdl.MapRGB(format_, 0xb0, 0xf4, 0xac) /// lt green
	color_palette[14] = sdl.MapRGB(format_, 0xaa, 0x9d, 0xef) /// lt blue
	color_palette[15] = sdl.MapRGB(format_, 0xb8, 0xb8, 0xb8) /// ly grey
	return color_palette
}

func main() {
	// os.Exit(..) must run AFTER sdl.Main(..) below; so keep track of exit
	// status manually outside the closure passed into sdl.Main(..) below
	var exitcode int
	sdl.Main(func() {
		exitcode = run()
	})
	// os.Exit(..) must run here! If run in sdl.Main(..) above, it will cause
	// premature quitting of sdl.Main(..) function; resource cleaning deferred
	// calls/closing of channels may never run
	os.Exit(exitcode)
}
