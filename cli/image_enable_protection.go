package cli

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newImageEnableProtectionCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "enable-protection [FLAGS] IMAGE PROTECTIONLEVEL [PROTECTIONLEVEL...]",
		Short:                 "Enable resource protection for an image",
		Args:                  cobra.MinimumNArgs(2),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.ensureToken,
		RunE:                  cli.wrap(runImageEnableProtection),
	}
	return cmd
}

func runImageEnableProtection(cli *CLI, cmd *cobra.Command, args []string) error {
	imageID, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.New("invalid image ID")
	}
	image := &hcloud.Image{ID: imageID}

	var unknown []string
	opts := hcloud.ImageChangeProtectionOpts{}
	for _, arg := range args[1:] {
		switch strings.ToLower(arg) {
		case "delete":
			opts.Delete = hcloud.Bool(true)
		default:
			unknown = append(unknown, arg)
		}
	}
	if len(unknown) > 0 {
		return fmt.Errorf("unknown protection level: %s", strings.Join(unknown, ", "))
	}

	action, _, err := cli.Client().Image.ChangeProtection(cli.Context, image, opts)
	if err != nil {
		return err
	}

	if err := cli.ActionProgress(cli.Context, action); err != nil {
		return err
	}

	fmt.Printf("Resource protection enabled for image %d\n", image.ID)
	return nil
}
