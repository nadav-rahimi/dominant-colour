## Dominant Colour
Made in reference to [this](https://pdfs.semanticscholar.org/fc50/a3950d6ce54717b945079329069dcd8ccb7a.pdf) paper

### Description
Finds *n* most dominant colours in an image

### Example Code
```go
package main

import dc "github.com/nadav-rahimi/dominant-colour"

func main() {
    // Get the 6 most dominant colours from the image
    colours := dc.FindDominantColoursBT("path/to/image.jpg", 6)
    // Draw a rectangle of these colours
    dc.DrawRectangle(colours)
    // Recreate the image using these colours
    dc.RecreateImage("path/to/image.jpg", colours)
}
```

### Example Results
#### Input
![example input](images/skin.jpg)

#### Output (8 Colours)
![example input](images/skin_render.jpeg)