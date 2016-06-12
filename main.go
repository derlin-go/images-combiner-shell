package main

import (
    "flag"
    "fmt"
    "os"
    "image/color"
    "image"
    "github.com/derlin-go/combiner"
)

func main() {

    var yGap int;
    var opaque bool;
    var yGapStrColor, bgStrColor string;
    var yGapColor color.Color;
    var bgColor color.Color = color.Transparent;

    // --------- program arguments


    // register and parse program arguments
    flag.IntVar(&yGap, "gap", 0, "gap between images");
    flag.BoolVar(&opaque, "opaque", false, "replace transparency by white or bgColor (if defined)");
    flag.StringVar(&yGapStrColor, "gapColor", "", "color of gap between images, as an hex string");
    flag.StringVar(&bgStrColor, "bgColor", "", "replace alpha to (leave empty to keep transparency");
    flag.Parse();

    // handle the bgColor argument: convert str -> color + set the gap color to the bg color by default
    if opaque || bgStrColor != "" {
        if bgStrColor != "" {
            var err error;
            bgColor, err = combiner.ParseColor(bgStrColor)
            if (err != nil) {
                fmt.Println(err)
                os.Exit(0)
            }
        } else {
            bgColor = color.White;
        }
        yGapColor = bgColor;
    }

    // if the gap color is defined, try to convert the string to a color
    if yGapStrColor != "" {
        // default to bgcolor
        var err error;
        yGapColor, err = combiner.ParseColor(yGapStrColor)
        if (err != nil) {
            fmt.Println(err)
            os.Exit(0)
        }
    }

    // --------- load and resize images

    image_paths := os.Args[len(os.Args) - flag.NArg():]
    images := make([]*image.Image, len(image_paths))

    if len(images) == 0 {
        fmt.Println("No images specified.")
        os.Exit(0)
    }

    // open and decode the images
    for idx, path := range (image_paths) {
        img, _, _ := combiner.OpenAndDecode(path)
        images[idx] = &img
    }

    out, err := os.Create("./output.png")
    if err != nil {
        fmt.Println("Could not create outfile")
        panic(err)
    }

    data, err := combiner.Compose(images, bgColor, opaque, yGap, yGapColor)
    if err != nil {
        panic(err)
    }
    _, err = out.Write(data)

}