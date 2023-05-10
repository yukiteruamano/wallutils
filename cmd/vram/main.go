package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/urfave/cli"
	"github.com/xyproto/wallutils"
)

func getVRAMAction(c *cli.Context) error {
	if c.IsSet("list") {
		// Retrieve a slice of GPU structs, or exit with an error, including integrated graphic cards ("VGA" in lspci output)
		gpus, err := wallutils.GPUs(true)
		if err != nil {
			return err
		}
		for _, gpu := range gpus {
			fmt.Printf("[%s] %s, %d MiB\n", gpu.Bus, gpu.Name, gpu.VRAM)
		}
		return nil
	}

	includeIntegrated := c.IsSet("integrated")

	// Retrieve a slice of GPU structs, or exit with an error, by default excluding integrated graphic cards ("VGA" in lspci output)
	gpus, err := wallutils.GPUs(includeIntegrated)
	if err != nil {
		return err
	}

	if len(gpus) == 0 {
		if includeIntegrated {
			fmt.Fprintln(os.Stderr, "error: could not find any available GPUs")
			return errors.New("could not find any available GPUs")
		}
		allGpus, err := wallutils.GPUs(true)
		if err != nil {
			return err
		}
		nonIntegratedGpus, err := wallutils.GPUs(false)
		if err != nil {
			return err
		}
		switch len(allGpus) {
		case 0:
			if len(nonIntegratedGpus) == 0 {
				fmt.Fprintln(os.Stderr, "error: could not find any GPU")
			} else {
				fmt.Fprintf(os.Stderr, "error: found no GPUs while at the same time finding %d non-integrated GPUs\n", len(nonIntegratedGpus))
			}
		case 1:
			fmt.Fprintf(os.Stderr, "error: found one GPUs, and %d of them are non-integrated\n", len(nonIntegratedGpus))
		default:
			fmt.Fprintf(os.Stderr, "error: found %d GPUs, where %d of them are non-integrated\n", len(allGpus), len(nonIntegratedGpus))
		}
	}

	// Output the minimum amount of VRAM in MiB
	minimum := uint(0)
	for _, gpu := range gpus {
		if minimum == 0 || gpu.VRAM < minimum {
			minimum = gpu.VRAM
		}
	}

	// Output the minimum about of VRAM for non-integrated GPUs (if possible), in MiB
	fmt.Printf("%d MiB\n", minimum)
	return nil
}

func main() {
	app := cli.NewApp()

	app.Name = "vram"
	app.Usage = "get the minimum amount of VRAM for all non-integrated GPUs.\n          If only integrated GPUs are available, the minimum amount of VRAM for these are returned instead."
	app.UsageText = "vram [options]"

	app.Version = wallutils.VersionString
	app.HideHelp = true

	cli.VersionFlag = cli.BoolFlag{
		Name:  "version, V",
		Usage: "output version information",
	}

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "list, l",
			Usage: "list bus ID, description and the amount of VRAM for each GPU",
		},
		cli.BoolFlag{
			Name:  "integrated, i",
			Usage: "include integrated GPUs when finding the minimum amount of available VRAM",
		},
	}

	app.Action = getVRAMAction
	if err := app.Run(os.Args); err != nil {
		wallutils.Quit(err)
	}
}
