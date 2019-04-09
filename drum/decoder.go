package drum

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// DecodeFile decodes the drum machine file found at the provided path
// and returns a pointer to a parsed pattern which is the entry point to the
// rest of the data.
// TODO: implement
func DecodeFile(path string) (*Pattern, error) {

	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	b := bytes.NewReader(buf)

	//fmt.Printf("Contents: \n%v\n", hex.Dump(buf))
	p := &Pattern{
		Name: filepath.Base(path),
	}

	if err := binary.Read(b, binary.LittleEndian, &p.Header); err != nil {
		fmt.Println("binary.Read failed:", err)
	}

	for {
		m := Measure{}
		err := binary.Read(b, binary.LittleEndian, &m.Id)
		if err != nil {
			break
		}
		m.Name = decodeName(b)
		err = binary.Read(b, binary.LittleEndian, &m.Steps)
		if err != nil {
			break
		}
		p.Measures = append(p.Measures, m)
	}
	return p, err
}

func decodeName(r io.Reader) string {
	var buffer bytes.Buffer
	var b byte


	binary.Read(r, binary.LittleEndian, &b)
	var len = int(b)
	for i := 0; i < len; i++ {
		err := binary.Read(r, binary.LittleEndian, &b)
		if err != nil || b == 0 {
			break
		}
		buffer.WriteByte(b)
	}
	return buffer.String()
}

// Pattern is the high level representation of the
// drum pattern contained in a .splice file.
// TODO: implement
type Pattern struct {
	Name		string
	Header 		PatternHeader
	Measures	[]Measure
}

func (p Pattern) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%v\n",p.Header))
	for _, m := range p.Measures {
		sb.WriteString(fmt.Sprintf("%v\n", m))
	}
	return sb.String()
}
type PatternHeader struct {
	FileType 	[6]byte
	_			uint64
	Version 	[32]byte
	Tempo	 	float32
}

func (h PatternHeader) String() string {
	return fmt.Sprintf("Saved with HW Version: %s\nTempo: %g", string(bytes.Trim(h.Version[:], "\x00")), h.Tempo)
}

func (p PatternHeader) hexDump() string {
	return fmt.Sprintf("Version: %#x\nTempo: %v", p.Version, p.Tempo)
}

type Measure struct {
	Id 			uint32
	Name 		string
	Steps		[16]byte
}

func (m Measure) String() string {
	return fmt.Sprintf("(%d) %s\t%s", m.Id, m.Name, m.stepsToString(m.Steps))
}

func (m Measure) stepsToString(s [16]byte) string {
	var sb strings.Builder

	for i :=0; i < 16; i++ {
		if i % 4 == 0 {
			sb.WriteString("|")
		}

		if s[i] == 0 {
			sb.WriteString("-")
		} else {
			sb.WriteString("x")
		}
	}
	sb.WriteString("|")
	return sb.String()
}