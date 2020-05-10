// Copyright 2017 Hajime Hoshi
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//

package main

import (
	"encoding/binary"
	"io"
	"log"
	"os"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"
)

func run() error {
	// f, err := os.Open("classic_mpeg2.mp3")
	f, err := os.Open("Test_mpeg2.mp3")
	if err != nil {
		return err
	}
	defer f.Close()

	decoder, err := mp3.NewDecoder(f)
	if err != nil {
		return err
	}

	play(decoder)
	// writeOut(decoder)

	return nil
}

func play(decoder *mp3.Decoder) error {

	c, err := oto.NewContext(decoder.SampleRate(), 2, 2, 8192)
	if err != nil {
		return err
	}
	defer c.Close()

	p := c.NewPlayer()
	defer p.Close()

	// fmt.Printf("Length: %d[bytes]\n", decoder.Length())

	if _, err := io.Copy(p, decoder); err != nil {
		return err
	}

	return nil
}

func writeOut(decoder *mp3.Decoder) {
	out, err := os.Create("output.wav")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	e := wav.NewEncoder(out, decoder.SampleRate(), 16, 2, 1)

	audioBuf, err := newAudioIntBuffer(decoder)
	if err != nil {
		log.Fatal(err)
	}

	if err := e.Write(audioBuf); err != nil {
		log.Fatal(err)
	}
	if err := e.Close(); err != nil {
		log.Fatal(err)
	}
}

func newAudioIntBuffer(decoder *mp3.Decoder) (*audio.IntBuffer, error) {
	buf := audio.IntBuffer{
		Format: &audio.Format{
			NumChannels: 1,
			SampleRate:  decoder.SampleRate(),
		},
	}
	for {
		var sample int16
		err := binary.Read(decoder, binary.LittleEndian, &sample)
		switch {
		case err == io.EOF:
			return &buf, nil
		case err != nil:
			return nil, err
		}
		buf.Data = append(buf.Data, int(sample))
	}
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
