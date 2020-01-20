package internal

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"github.com/jackpal/bencode-go"
)

type bencodeInfo struct {
	Pieces      string `bencode:"pieces"`
	// piece size (bytes)
	PieceLength int    `bencode:"piece length"`
	// file size (bytes)
	Length      int    `bencode:"length"`
	// file name
	Name        string `bencode:"name"`
}

type BencodeTorrent struct {
	Announce string      `bencode:"announce"`
	Info     bencodeInfo `bencode:"info"`
}

type BencodeTrackerResp struct {
	Interval int `bencode:"interval"`
	Peers string `bencode:"peers"`
}

func (b *BencodeTorrent) ToTorrentFile() *TorrentFile {

	tf := &TorrentFile{
		Announce:   b.Announce,
		Filename:   b.Info.Name,
		Filesize:   b.Info.Length,
		Blocksize:  b.Info.PieceLength,
		InfoHash:   b.Info.hash(),
		PiecesHash: b.Info.splitPieceHashes(),
	}

	return tf
}

func (i *bencodeInfo) hash() ([20]byte) {
	var buf bytes.Buffer
	err := bencode.Marshal(&buf, *i)
	if err != nil {
		return [20]byte{}
	}
	h := sha1.Sum(buf.Bytes())
	return h
}

func (i *bencodeInfo) splitPieceHashes() ([][20]byte) {
	hashLen := 20 // Length of SHA-1 hash
	buf := []byte(i.Pieces)
	if len(buf)%hashLen != 0 {
		fmt.Errorf("Received malformed pieces of length %d", len(buf))
		return nil
	}
	numHashes := len(buf) / hashLen
	hashes := make([][20]byte, numHashes)

	for i := 0; i < numHashes; i++ {
		copy(hashes[i][:], buf[i*hashLen:(i+1)*hashLen])
	}
	return hashes
}

func (tr *BencodeTrackerResp) SplitPeers() ([][6]byte) {
	hashLen := 6
	buf := []byte(tr.Peers)
	if len(buf)%hashLen != 0 {
		fmt.Errorf("Received malformed pieces of length %d", len(buf))
		return nil
	}
	numHashes := len(buf) / hashLen
	hashes := make([][6]byte, numHashes)

	for i := 0; i < numHashes; i++ {
		copy(hashes[i][:], buf[i*hashLen:(i+1)*hashLen])
	}
	return hashes
}