package main

import (
	"dnc-torrent/internal"
	"flag"
	bencode "github.com/jackpal/bencode-go"
	"github.com/spf13/cobra"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

var (
	torrentPath string
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

	sPeerId := make([]byte, 20)
	rand.Read(sPeerId[:])

	var peerId [20]byte
	copy(peerId[:], sPeerId)

	url, err  := tf.BuildTrackerURL(peerId, 3889)
	if err != nil {
		internal.Sugar.Fatalf("Error: %v", err)
	}

	c := &http.Client{Timeout: 15 * time.Second}
	resp, err := c.Get(url)
	if err != nil {
		internal.Sugar.Fatalf("Error: %v", err)
	}
	defer resp.Body.Close()

	btr := internal.BencodeTrackerResp{}
	bencode.Unmarshal(resp.Body, &btr)

	internal.Sugar.Debugf("Interval: %d, NumPeers: %d", btr.Interval, len(btr.SplitPeers()))

	peers, err := internal.UnmarshalPeers([]byte(btr.Peers))
	if err != nil {
		internal.Sugar.Fatalf("Error: %v", err)
	}

	internal.Sugar.Debugf("PeerId: %v (%d), Hash: %+v (%d)", peerId, len(peerId), tf.InfoHash, len(tf.InfoHash))

	for _, p := range peers {
		go TestHandshake(p, tf.InfoHash, peerId)
	}

	time.Sleep(1*time.Minute)
}

func TestHandshake(p *internal.Peer, infohash, peerId [20]byte) {

	_, err := internal.New(*p, peerId, infohash)
	if err != nil {
		log.Printf("Could not handshake with %s. Err: %v. Disconnecting\n", p.Address(), err)
	}else {
		log.Printf("Completed handshake with %s\n", p.Address())
	}

}