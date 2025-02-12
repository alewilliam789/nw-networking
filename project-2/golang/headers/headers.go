package headers

import (
  "bytes"
  "encoding/binary"
)

type UDPHeader struct {
  SrcPort uint16
  DestPort uint16
  Length uint16
  Checksum uint16
}

func(udp *UDPHeader) Marshall() []byte {
  
  buf := new(bytes.Buffer)

  binary.Write(buf, binary.BigEndian, udp.SrcPort)
  binary.Write(buf, binary.BigEndian, udp.DestPort)
  binary.Write(buf, binary.BigEndian, udp.Length)
  binary.Write(buf, binary.BigEndian, udp.Checksum)

  out := buf.Bytes()

  pad := 8 - len(out)

  for i := 0; i < pad; i++ {
    out = append(out, 0)
  }

  return out
}

