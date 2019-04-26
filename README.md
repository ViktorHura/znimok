# znimok
Znimok - yet another screenshot application, written in go

Using https://github.com/bamontelucas/image-clipboard to copy images due to limitations in go

Usage: Make a selection with mouse and press Space to copy to clipboard

Works on multiple displays as long as the tops of the displays are aligned


Building: 

      go build -ldflags -H=windowsgui main.go
