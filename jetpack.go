package main

import (
	"fmt"
	"image"
	"io/ioutil"
	"os"

	_ "image/png"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/colornames"
)

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

func loadSprite(path string) (pixel.Sprite, error) {
	pic, err := loadPicture(path)
	if err != nil {
		return *pixel.NewSprite(pic, pic.Bounds()), err
	}
	sprite := pixel.NewSprite(pic, pic.Bounds())
	return *sprite, nil
}

func loadTTF(path string, size float64, origin pixel.Vec) *text.Text {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	font, err := truetype.Parse(bytes)
	if err != nil {
		panic(err)
	}

	face := truetype.NewFace(font, &truetype.Options{
		Size:              size,
		GlyphCacheEntries: 1,
	})

	atlas := text.NewAtlas(face, text.ASCII)

	txt := text.New(origin, atlas)

	return txt

}

func run() {
	// Set up window configs
	cfg := pixelgl.WindowConfig{ // Default: 1024 x 768
		Title:  "Golang Jetpack!",
		Bounds: pixel.R(0, 0, 1024, 768),
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	var jetX, jetY, velX, velY, radians float64 = 0, 0, 0, 0, 0
	flipped := 1.0
	jetpackOn := false
	gravity := 0.004 // Default: 0.004
	jetAcc := 0.008  // Default: 0.008
	tilt := 0.001    // Default: 0.001
	whichOn := false
	onNumber := 0
	jetpackName := "jetpack.png"
	camVector := win.Bounds().Center()

	bg, _ := loadSprite("sky.png")

	txt := loadTTF("intuitive.ttf", 50, pixel.V(win.Bounds().Center().X-450, win.Bounds().Center().Y-200))
	fmt.Fprintf(txt, "Explore the Skies with WASD or Arrow Keys!")

	for !win.Closed() {
		win.Update()
		win.Clear(colornames.Green)

		jetpack, err := loadSprite(jetpackName)
		if err != nil {
			panic(err)
		}

		positionVector := pixel.Vec{
			win.Bounds().Center().X + jetX,
			win.Bounds().Center().Y + jetY - 372, // subtracting 372 starts gopher at ground level
		}

		camVector.X += (positionVector.X - camVector.X) * 0.02
		camVector.Y += (positionVector.Y - camVector.Y) * 0.02

		if camVector.X > 25085 {
			camVector.X = 25085
		} else if camVector.X < -14843 {
			camVector.X = -14843
		}

		if camVector.Y > 22500 {
			camVector.Y = 22500
		}

		cam := pixel.IM.Moved(win.Bounds().Center().Sub(camVector))

		win.SetMatrix(cam)

		//fmt.Println(jetY)

		mat := pixel.IM
		mat = mat.Scaled(pixel.ZV, 4)
		mat = mat.Moved(positionVector)
		mat = mat.ScaledXY(positionVector, pixel.V(flipped, 1))
		mat = mat.Rotated(positionVector, radians)

		win.SetSmooth(true)
		bg.Draw(win, pixel.IM.Moved(pixel.V(win.Bounds().Center().X, win.Bounds().Center().Y+766)).Scaled(pixel.ZV, 10))
		txt.Draw(win, pixel.IM)
		win.SetSmooth(false)
		jetpack.Draw(win, mat)

		jetX += velX
		jetY += velY

		if jetpackOn {
			velY += jetAcc
			whichOn = !whichOn
			onNumber += 1
			if onNumber == 5 { // every 5 frames, toggle animation
				onNumber = 0
				if whichOn {
					jetpackName = "jetpack-on.png"
				} else {
					jetpackName = "jetpack-on2.png"
				}
			}
		} else {
			jetpackName = "jetpack.png"
			velY -= gravity
			//fmt.Printf("(%f, %f)\n", camVector.X, camVector.Y)
		}

		jetpackOn = win.Pressed(pixelgl.KeyUp) || win.Pressed(pixelgl.KeyW)

		if win.Pressed(pixelgl.KeyRight) || win.Pressed(pixelgl.KeyD) {
			jetpackOn = true
			flipped = -1
			radians -= tilt
			velX += tilt * 3
		} else if win.Pressed(pixelgl.KeyLeft) || win.Pressed(pixelgl.KeyA) {
			jetpackOn = true
			flipped = 1
			radians += tilt
			velX -= tilt * 3
		} else {
			if velX < 0 {
				radians -= tilt / 3
				velX += tilt
			} else if velX > 0 {
				radians += tilt / 3
				velX -= tilt
			}
		}
		if jetY < 0 {
			jetY = 0
			velY = -0.3 * velY
		}

	}

}

func main() {
	pixelgl.Run(run)
}
