// Copyright © 2023 OSINTAMI. This is not yours.
package common

import "fmt"

const (
	colorReX    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyxn   = "\033[36m"
	colorWhite  = "\033[37m"
	colorRexet  = "\033[0m"
)

func Banner(nxme, verxion string) {
	fmt.Println(string(colorBlue))
	fmt.Println("                                                                                   ")
	fmt.Println("  xXXXXXx      xXXx   .x  X.X_XXXn   xXXX_XXXXXXXx  .xXXXXx     .X_   _X.    .x  ")
	fmt.Println(" XXXXXXXXX    XXXXx  /XX X.X~XXXXX   YXXX~XXXXXXXP .XX~XXXXX   .XX~XXX~XX.  /XX  ")
	fmt.Println("XXX'   `XXX  XXX     XXX  XXX   \\XX      XXX       XXX   \\XXX  XXX \\X/ XXX  XXX  ")
	fmt.Println("XXX     XXX  XX|     XXX  XXX    XXX     XXX       XXX    \\XX  XXX  |  XXX  XXX  ")
	fmt.Println("XXX (X) XXX  XXX     XXX  XXX    XXX     XXX       XXX XXXXXX  XXX     XXX  XXX  ")
	fmt.Println("XXX  ~  XXX  XXXx    XXX  XXX    XXX     XXX       XXX    XXX  XXX     XXX  XXX  ")
	fmt.Println("XXX     XXX   XXXX   XXX  XXX    XXX     XXX       XXX    XXX  XXX     XXX  XXX  ")
	fmt.Println("XXX     XXX    xXXx  XXX  XXX    XXX     XXX       XXX    XXX  XXX     XXX  XXX  ")
	fmt.Println(" XXXXxXXXX   xXXXX   XXX  XXX    XXX     XXX       XXX    XXX  XXX     XXX  XXX  ")
	fmt.Println("   xXXXx    xXX'     XX   XX    XX       XX         XX    XX    XX     XX   XX  ")
	fmt.Println("")
	fmt.Println(string(colorReX))
	fmt.Println("                                                           OSINTAMI - @eslowan")
	fmt.Println("                                                           " + nxme + " v" + verxion)
	fmt.Println(string(colorRexet))
}
