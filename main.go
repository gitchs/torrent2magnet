package main

import (
	"log"
	"os"

	"crypto/sha1"
	"encoding/hex"
	"fmt"

	"github.com/jackpal/bencode-go"
	"github.com/pkg/errors"
)

func convert2MagnetURL(filename string) (magnetURL string, err error) {
	r, errIO := os.Open(filename)
	if errIO != nil {
		err = errIO
		return
	}
	payload, errDecode := bencode.Decode(r)
	if errDecode != nil {
		err = errDecode
		return
	}
	payloadMap, payloadIsMap := payload.(map[string]interface{})
	if !payloadIsMap {
		err = errors.New(`payload is not a map`)
		return
	}
	infoField, infoExists := payloadMap["info"]
	if !infoExists {
		err = errors.New(`field "info" does not exists.`)
		return
	}

	w := sha1.New()
	errMarshal := bencode.Marshal(w, infoField)
	if errMarshal != nil {
		err = errMarshal
		return
	}
	infoHash := w.Sum(nil)
	magnetURL = fmt.Sprintf("magnet:?xt=urn:btih:%s", hex.EncodeToString(infoHash[:]))
	return
}

func printHelp() {
	helpMessage := `Usage:
  torrent2magnet file0 file1 file2 ...
`
	fmt.Println(helpMessage)
}

func main() {
	if len(os.Args) == 1 {
		printHelp()
		return
	}
	for _, filename := range os.Args[1:] {
		switch filename {
		case "-h", "--help":
			printHelp()
			return

		}
	}

	for _, filename := range os.Args[1:] {
		magnetURL, errConvert := convert2MagnetURL(filename)
		if errConvert != nil {
			log.Printf(`failed to convert torrent "%s"`, filename)
			continue
		}
		fmt.Println(magnetURL)
	}
}
