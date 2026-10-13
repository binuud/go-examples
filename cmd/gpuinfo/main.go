package main

import (
	"fmt"
	"time"

	"github.com/andersrennermalm/gpuinfo"
)

func main() {
	// gpuinfo reads GPU stats via IOKit — no root/sudo needed
	ch, err := gpuinfo.New(gpuinfo.WithInterval(1 * time.Second))
	if err != nil {
		panic(err)
	}

	for g := range ch {
		fmt.Printf("GPU: %s\n", g.Name)
		fmt.Printf("  Usage:   %.1f%%\n", g.UsedPercent)
		fmt.Printf("  VRAM:    %s / %s\n", g.UsedMemory, g.TotalMemory)
		fmt.Printf("  Encoder: %.1f%%  Decoder: %.1f%%\n", g.UsedEncoder, g.UsedDecoder)
		fmt.Println()
	}
}
