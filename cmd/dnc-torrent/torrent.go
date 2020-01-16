package main

import (
	"dnc-torrent/internal/service"
	"flag"
	"github.com/spf13/cobra"
	"github.com/marksamman/bencode"
	"os"
)

var (
	torrentPath string
)

func main() {
	service.InitLogger()

	cmd := &cobra.Command{
		Use: "",
		Long: ``,
		Run: run,
	}

	// include standard flags
	cmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)

	flag := cmd.Flags()
	flag.StringVar(&torrentPath, "torrentPath", "", "torrent file path")

	if err := cmd.Execute(); err != nil {
		service.Sugar.Fatal(err)
	}

}

func run(cmd *cobra.Command, args []string) {

	if len(torrentPath) == 0 {
		service.Sugar.Fatal("torrentPath is mandatory")
	}

	file, err := os.Open(torrentPath)
	if err != nil {
		service.Sugar.Error(err)
	}
	defer file.Close()

	dict, err := bencode.Decode(file)
	if err != nil {
		service.Sugar.Error(err)
	}

	service.Sugar.Debugf("Announce: %s", dict["announce"].(string))
}