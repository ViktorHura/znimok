package main

import (
 "image/png"
 "image"
 "os"
 "os/exec"
 "fmt"
 "github.com/disintegration/imaging"
 "strings"
 "strconv"
 "path/filepath"
 "github.com/faiface/pixel"
 "github.com/faiface/pixel/pixelgl"
 "github.com/faiface/pixel/imdraw"
 "golang.org/x/image/colornames"
 "github.com/kbinani/screenshot"
)



var selectionstate bool = false
var min pixel.Vec
var max pixel.Vec
var j int = -4

func run() {

  var(
    winwidth = 0
    winheight = 0
  )


  n := screenshot.NumActiveDisplays()

  for i := 0; i < n; i++ {

    bounds := screenshot.GetDisplayBounds(i)

    winwidth += bounds.Dx()
    if (bounds.Dy() > winheight){
      winheight = bounds.Dy()
    }

  }

  monitors := pixelgl.Monitors()

  var xpos float64 = 0

  for j := 0; j < len(monitors); j++ {

    x, _ := monitors[j].Position()

    if ( x < xpos ){
      xpos = x
    }

  }

  //create long window that spans all screens
  cfg := pixelgl.WindowConfig{
		Title:  "Znimok",
		Bounds: pixel.R(0, 0, float64(winwidth), float64(winheight)),
    Resizable: false,
    Undecorated: true,
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
  win.SetPos(pixel.V(xpos,0))
	win.Clear(colornames.Skyblue)



  var screens []string

  root := "shots/"
    walkerr := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
        screens = append(screens, path)
        return nil
    })
    if walkerr != nil {
        panic(walkerr)
    }

    //display all screens
    for k := 1; k < len(screens); k++ {

        screen := screens[k]

        name := strings.TrimPrefix(screen, `shots\`)
        name = strings.TrimSuffix(name, `.png`)
        fmt.Println(name)

        pic, err := loadPicture(screen)
        if err != nil {
             panic(err)
        }

        screenxpos, err := strconv.Atoi(name)
        if err != nil {
            panic(err)
        }



        sprite := pixel.NewSprite(pic, pic.Bounds())
        sprite.Draw(win, pixel.IM.Moved(pixel.V( (0 - xpos) + float64(screenxpos) + pic.Bounds().Size().X / 2, win.Bounds().Max.Y - pic.Bounds().Size().Y/2 )) )

    }

    pixels := win.Canvas().Pixels()
    for l:=0; l < len(pixels); l++ { //darken image

      pixels[l] = uint8(float64(pixels[l]) / float64(1.2))

    }
    win.Canvas().SetPixels(pixels)

    imd := imdraw.New(nil)

	for !win.Closed() {
    win.Clear(colornames.Aliceblue)
    imd.Clear()
    win.Canvas().SetPixels(pixels)
    var actualmin pixel.Vec = min
    var actualmax pixel.Vec = max


    if (min.X > max.X){
      actualmin.X = max.X
    }
    if (min.Y > max.Y){
      actualmin.Y = max.Y
    }

    if (max.X < min.X){
      actualmax.X = min.X
    }
    if (max.Y < min.Y){
      actualmax.Y = min.Y
    }


    if win.Pressed(pixelgl.MouseButtonLeft) {

      if (selectionstate == false){
        min = win.MousePosition()
        max = min
      }

      max = win.MousePosition()
      selectionstate = true
		}else {

      if (selectionstate == true){
        max = win.MousePosition()

      }

      selectionstate = false
    }

    if win.JustPressed(pixelgl.KeySpace) {

      winpixels := pixels


      for l:=0; l < len(pixels); l++ { //undo darken effect

        winpixels[l] = uint8(float64(winpixels[l]) * float64(1.2))

      }

      imgf := image.NewRGBA(image.Rect(0, 0, winwidth, winheight))
      imgf.Pix = winpixels

      fmt.Println("", actualmin, actualmax)

      img := imaging.Crop(imgf, image.Rect(int(actualmin.X), int(actualmin.Y), int(actualmax.X), int(actualmax.Y)))
      imgc := imaging.FlipV(img)

      savpath := "screenshot.png"

      f, _ := os.OpenFile(savpath, os.O_WRONLY|os.O_CREATE, 0600)
      png.Encode(f, imgc)
      f.Close()
      fmt.Println(savpath)

      output, err := exec.Command("image-clipboard", "screenshot.png").CombinedOutput()
      if err != nil {
        os.Stderr.WriteString(err.Error())
      }

      fmt.Println(string(output))

      win.Destroy()

		}



    imd.Color = colornames.White
	  imd.EndShape = imdraw.RoundEndShape
  	imd.Push(pixel.V(actualmin.X, actualmin.Y), pixel.V(actualmin.X, actualmax.Y))
  	imd.Push(pixel.V(actualmin.X, actualmin.Y), pixel.V(actualmax.X, actualmin.Y))

    imd.Push(pixel.V(actualmax.X, actualmax.Y), pixel.V(actualmax.X, actualmin.Y))
    imd.Push(pixel.V(actualmax.X, actualmax.Y), pixel.V(actualmin.X, actualmax.Y))
  	imd.Line(2)




    imd.Draw(win)
		win.Update()
	}

}

func main() {

  RemoveContents("shots/")

  n := screenshot.NumActiveDisplays()

  for i := 0; i < n; i++ {
    bounds := screenshot.GetDisplayBounds(i)

    img, err := screenshot.CaptureRect(bounds)
    if err != nil {
      panic(err)
    }

    screenpath := "shots/"

    pathErr := os.MkdirAll(screenpath ,0777)

  	//check if you need to panic, fallback or report
  	if pathErr == nil {
  		fmt.Println(pathErr)
  	}

    fileName := fmt.Sprintf("%d.png", bounds.Min.X)
    file, _ := os.Create(screenpath + fileName)
    png.Encode(file, img)
    file.Close()

    fmt.Printf("#%d : %v \"%s\"\n", i, bounds, fileName)

  }

	pixelgl.Run(run)
}

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

func RemoveContents(dir string) {
    d, _ := os.Open(dir)

    defer d.Close()
    names, _ := d.Readdirnames(-1)

    for _, name := range names {
        err := os.RemoveAll(filepath.Join(dir, name))
        if err != nil {
            fmt.Println(err)
        }
    }

}
