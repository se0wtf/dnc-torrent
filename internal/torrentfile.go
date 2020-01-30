package internal

import (
	"net/url"
	"strconv"
)

type TorrentFile struct {
	Announce string
	Filename string
	Filesize int
	Blocksize int
	InfoHash [20]byte
	PiecesHash [][20]byte
}

func (t *TorrentFile) BuildTrackerURL(peerID [20]byte, port uint16) (string, error) {
	base, err := url.Parse(t.Announce)
	if err != nil {
		return "", err
	}
	params := url.Values{
		"info_hash":  []string{string(t.InfoHash[:])},
		"peer_id":    []string{string(peerID[:])},
		"port":       []string{strconv.Itoa(int(port))},
		"uploaded":   []string{"0"},
		"downloaded": []string{"0"},
		"compact":    []string{"1"},
		"left":       []string{strconv.Itoa(t.Filesize)},
	}
	base.RawQuery = params.Encode()
	return base.String(), nil
}