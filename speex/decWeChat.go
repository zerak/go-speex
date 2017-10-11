package speex

import "C"

import (
	"bytes"
	"fmt"
	"unsafe"
)

type SpeexWeChatDecoder struct {
	context C.TRSpeexDecodeContex
}

func NewSpeexWeChatDecoder() *SpeexWeChatDecoder {
	return &SpeexWeChatDecoder{}
}

// @remark only support mono speex(channels must be 1).
func (v *SpeexWeChatDecoder) Init() (err error) {
	r := C.TRSpeexDecodeInit(&v.context)
	if int(r) <= 0 {
		return fmt.Errorf("init decoder failed, err=%v", int(r))
	}

	return
}

func (v *SpeexWeChatDecoder) Close() {
	C.TRSpeexDecodeRelease(&v.context)
}

// @return pcm is nil when EOF.
func (v *SpeexWeChatDecoder) Decode(frame []byte) (pcm []byte, err error) {
	p := (*C.char)(unsafe.Pointer(&frame[0]))
	pSize := C.int(len(frame))

	pIsDone := C.int(0)

	var result bytes.Buffer

	var buf[44] byte
	buf[0] = 'R'
	buf[1] = 'I'
	buf[2] = 'F'
	buf[3] = 'F'

	buf[8] = 'W'
	buf[9] = 'A'
	buf[10] = 'V'
	buf[11] = 'E'
	buf[12] = 'f'
	buf[13] = 'm'
	buf[14] = 't'
	buf[15] = 0x20

	buf[16] = 0x10
	buf[20] = 0x01
	buf[22] = 0x01
	buf[24] = 0x80
	buf[25] = 0x3E
	buf[29]= 0x7D
	buf[32] = 0x02
	buf[34] = 0x10
	buf[36] = 'd'
	buf[37] = 'a'
	buf[38] = 't'
	buf[39] = 'a'

	nOutSize := C.int(0)
	nTotalLen := C.int(0)
	for {
		if pIsDone == 1 {
			break
		}
		
		r := C.TRSpeexDecode(&v.context,frame,&nOutSize)
		if r <= 0 {
			return nil, fmt.Errorf("decode failed, err=%v", int(r))
		}

		nTotalLen += nOutSize

		result.Write(pcmTmp)
	}

	pcm = result.Bytes()

	return
}

func (v *SpeexWeChatDecoder) FrameSize() int {
	return int(C.speexdec_frame_size(&v.context.m))
}

func (v *SpeexWeChatDecoder) SampleRate() int {
	return int(C.speexdec_sample_rate(&v.m))
}

func (v *SpeexWeChatDecoder) Channels() int {
	return 1
}
