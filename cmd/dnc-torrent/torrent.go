package main

import (
	"dnc-torrent/internal"
	"encoding/binary"
	"flag"
	"fmt"
	bencode "github.com/jackpal/bencode-go"
	"github.com/spf13/cobra"
	"math/rand"
	"net"
	"net/http"
	"os"
)

var (
	torrentPath string
	peers [][6]byte
)

func main() {
	internal.InitLogger()

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
		internal.Sugar.Fatal(err)
	}

}

func run(cmd *cobra.Command, args []string) {

	if len(torrentPath) == 0 {
		internal.Sugar.Fatal("torrentPath is mandatory")
	}

	file, err := os.Open(torrentPath)
	if err != nil {
		internal.Sugar.Error(err)
	}
	defer file.Close()

	bto := &internal.BencodeTorrent{}
	bencode.Unmarshal(file, bto)

	tf := bto.ToTorrentFile()

	internal.Sugar.Debugf("Announce: %s, FileSize: %dMb, PieceSize: %dKb", tf.Announce, tf.Filesize/1024/2014, tf.Blocksize/1024)

	peerId := make([]byte, 20)
	rand.Read(peerId)

	url, err  := tf.BuildTrackerURL(peerId, 3889)
	if err != nil {
		internal.Sugar.Fatalf("Error: %v", err)
	}
	resp, err := http.Get(url)
	if err != nil {
		internal.Sugar.Fatalf("Error: %v", err)
	}
	defer resp.Body.Close()

	btr := internal.BencodeTrackerResp{}
	bencode.Unmarshal(resp.Body, &btr)

	internal.Sugar.Debugf("Interval: %d, NumPeers: %d", btr.Interval, len(btr.SplitPeers()))

	peers, err := Unmarshal([]byte(btr.Peers))
	if err != nil {
		internal.Sugar.Fatalf("Error: %v", err)
	}
	for _, p := range peers {
		internal.Sugar.Infof("IP: %s:%d", p.IP.String(), p.Port)
	}
}

// Unmarshal parses peer IP addresses and ports from a buffer
func Unmarshal(peersBin []byte) ([]internal.Peer, error) {
	const peerSize = 6 // 4 for IP, 2 for port
	numPeers := len(peersBin) / peerSize
	if len(peersBin)%peerSize != 0 {
		err := fmt.Errorf("Received malformed peers")
		return nil, err
	}
	peers := make([]internal.Peer, numPeers)
	for i := 0; i < numPeers; i++ {
		offset := i * peerSize
		peers[i].IP = net.IP(peersBin[offset : offset+4])
		peers[i].Port = binary.BigEndian.Uint16(peersBin[offset+4 : offset+6])
	}
	return peers, nil
}