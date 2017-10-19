package speex

/*
#cgo CFLAGS: -I${SRCDIR}/../speex-lib/objs/include
#cgo LDFLAGS: ${SRCDIR}/../speex-lib/objs/lib/libspeex.a -lm
#include "speex/speex.h"

#include <sys/malloc.h>
//#include <malloc.h>
#include <memory.h>

#define PCMMAX(a,b) ((a) > (b) ? (a) : (b))
#define PCMMIN(a,b) ((a) > (b) ? (b) : (a))

typedef struct PCMFifoBuffer {
    unsigned char *buffer;
    unsigned char *rptr, *wptr, *end;
} PCMFifoBuffer;

int pcm_fifo_init(PCMFifoBuffer *f, int size);
void pcm_fifo_free(PCMFifoBuffer *f);
int pcm_fifo_size(PCMFifoBuffer *f);
int pcm_fifo_read(PCMFifoBuffer *f, unsigned char *buf, int buf_size);
int pcm_fifo_generic_read(PCMFifoBuffer *f, int buf_size, void (*func)(void*, void*, int), void* dest);
void pcm_fifo_write(PCMFifoBuffer *f, const unsigned char *buf, int size);
void pcm_fifo_realloc(PCMFifoBuffer *f, unsigned int size);
void pcm_fifo_drain(PCMFifoBuffer *f, int size);

static unsigned char pcm_fifo_peek(PCMFifoBuffer *f, int offs) {
    unsigned char *ptr = f->rptr + offs;
    if (ptr >= f->end)
        ptr -= f->end - f->buffer;
    return *ptr;
}

int pcm_fifo_init(PCMFifoBuffer *f, int size) {
    f->wptr = f->rptr =  f->buffer = (unsigned char*)malloc(size);

    f->end = f->buffer + size;
    if (!f->buffer)
        return -1;
    return 0;
}

void pcm_fifo_free(PCMFifoBuffer *f) {
    free(f->buffer);
}

int pcm_fifo_size(PCMFifoBuffer *f) {
    int size = f->wptr - f->rptr;
    if (size < 0)
        size += f->end - f->buffer;
    return size;
}

int pcm_fifo_read(PCMFifoBuffer *f, unsigned char *buf, int buf_size) {
    return pcm_fifo_generic_read(f, buf_size, 0, buf);
}

void pcm_fifo_realloc(PCMFifoBuffer *f, unsigned int new_size) {
    unsigned int old_size= f->end - f->buffer;

    if(old_size < new_size) {
        int len= pcm_fifo_size(f);
        PCMFifoBuffer f2;

        pcm_fifo_init(&f2, new_size);
        pcm_fifo_read(f, f2.buffer, len);
        f2.wptr += len;
        free(f->buffer);
        *f= f2;
    }
}

void pcm_fifo_write(PCMFifoBuffer *f, const unsigned char *buf, int size) {
    do {
        int len = PCMMIN(f->end - f->wptr, size);
        memcpy(f->wptr, buf, len);
        f->wptr += len;
        if (f->wptr >= f->end)
            f->wptr = f->buffer;
        buf += len;
        size -= len;
    } while (size > 0);
}

int pcm_fifo_generic_read(PCMFifoBuffer *f, int buf_size, void (*func)(void*, void*, int), void* dest) {
    int size = pcm_fifo_size(f);

    if (size < buf_size)
        return -1;
    do {
        int len = PCMMIN(f->end - f->rptr, buf_size);
        if(func) func(dest, f->rptr, len);
        else{
            memcpy(dest, f->rptr, len);
            dest = (unsigned char*)dest + len;
        }
        pcm_fifo_drain(f, len);
        buf_size -= len;
    } while (buf_size > 0);
    return 0;
}

void pcm_fifo_drain(PCMFifoBuffer *f, int size) {
    f->rptr += size;
    if (f->rptr >= f->end)
        f->rptr -= f->end - f->buffer;
}


#ifndef NULL
#define NULL ((void *)0)
#endif

#define MAX_FRAME_SIZE 2000
#define MAX_FRAME_BYTES 2000

typedef struct _TRSpeexDecodeContex {
	void *st;
	SpeexBits bits;
	int frame_size;

	PCMFifoBuffer* pFifo;
} TRSpeexDecodeContex;

int TRSpeexDecodeInit(TRSpeexDecodeContex* stDecode) {
	int modeID = -1;
	const SpeexMode *decmode=NULL;
	int nframes;
	int vbr_enabled;
	int chan;
	int rate;
	void *st;
	int quality;
	int dec_frame_size;
	int complexity;
	int nbBytes;
	int ret;
	int enh_enabled;
	int decrate;
	int declookahead;

	if(stDecode == NULL)
		return -1;

	modeID = SPEEX_MODEID_WB;

	speex_bits_init(&(stDecode->bits));

	decmode = speex_lib_get_mode (modeID);

	stDecode->st = speex_decoder_init(decmode);
	if(stDecode->st == NULL)
		return -1;

	enh_enabled = 1;
	decrate = 16000;

	speex_decoder_ctl(stDecode->st, SPEEX_SET_ENH, &enh_enabled);
	speex_decoder_ctl(stDecode->st, SPEEX_SET_SAMPLING_RATE, &decrate);
	speex_decoder_ctl(stDecode->st, SPEEX_GET_FRAME_SIZE, &dec_frame_size);
	speex_decoder_ctl(stDecode->st, SPEEX_GET_LOOKAHEAD, &declookahead);

	stDecode->frame_size = dec_frame_size;
	stDecode->pFifo = (PCMFifoBuffer*)malloc(sizeof(PCMFifoBuffer));

	if(stDecode->pFifo != NULL)	{
		ret = pcm_fifo_init(stDecode->pFifo, 1024*10000);
		if(ret == -1)
			return -1;
	}
	else
		return -1;

	return 1;
}

int TRSpeexDecode(TRSpeexDecodeContex* stDecode,char* pInput, int nInputSize, char* pOutput, int* nOutSize) {
	int nbBytes;
	char aInputBuffer[MAX_FRAME_BYTES];
	int nFrameNo;
	int nDecSize;
	int nTmpSize;
	int ret = 0;

	if(stDecode == NULL)
		return -1;

	if(pInput == NULL)
		return -1;

	if(pOutput == NULL)
		return -1;

	if(nInputSize < 0)
		return -1;

	if(nInputSize > 1024*10000)
		return -1;

	if(stDecode->pFifo != NULL)
		pcm_fifo_write(stDecode->pFifo, (unsigned char*)pInput, nInputSize);
	else
		return -1;

	nFrameNo = 0;
	nDecSize = 0;
	nTmpSize = 0;

	while(pcm_fifo_size(stDecode->pFifo) >= 60) {
		pcm_fifo_read(stDecode->pFifo, (unsigned char*)aInputBuffer, 60);

		speex_bits_read_from(&(stDecode->bits), aInputBuffer, 60);

		ret = speex_decode_int(stDecode->st, &(stDecode->bits), (spx_int16_t*)pOutput+nFrameNo*(stDecode->frame_size));
		if(ret == -1 || ret == -2) {
			nOutSize = 0;
			return -1;
		}

		nTmpSize += stDecode->frame_size*2;

		nFrameNo ++;
	}

	*nOutSize = nTmpSize;
	return 1;
}

int TRSpeexDecodeRelease(TRSpeexDecodeContex* stDecode) {
	if(stDecode == NULL)
		return -1;

	if (stDecode->st != NULL)
		speex_decoder_destroy(stDecode->st);

	speex_bits_destroy(&(stDecode->bits));

	if(stDecode->pFifo != NULL) {
		pcm_fifo_free(stDecode->pFifo);
		free(stDecode->pFifo);
		stDecode->pFifo = NULL;
	}

	return 1;
}

*/
import "C"

import (
	"bytes"
	"errors"
	"fmt"
	"unsafe"
)

type SpeexWeChatDecoder struct {
	context C.TRSpeexDecodeContex
}

func NewSpeexWeChatDecoder() *SpeexWeChatDecoder {
	return &SpeexWeChatDecoder{}
}

//func (v *SpeexWeChatDecoder) Init() (err error) {
//	r := C.TRSpeexDecodeInit(&v.context)
//	if int(r) <= 0 {
//		return fmt.Errorf("init decoder failed, err=%v", int(r))
//	}
//
//	return
//}

//func (v *SpeexWeChatDecoder) Close() {
//	C.TRSpeexDecodeRelease(&v.context)
//}

// @return pcm is nil when EOF.
func (v *SpeexWeChatDecoder) Decode(frame []byte) (pcm []byte, err error) {
	var src = bytes.NewBuffer(frame) // input
	var result = bytes.Buffer{}      // output

	bufSize := 6         // read data size
	nOutSize := C.int(0) // decode size
	nTotalLen := 0       // decode total size

	buf := make([]byte, 44)
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
	buf[29] = 0x7D
	buf[32] = 0x02
	buf[34] = 0x10
	buf[36] = 'd'
	buf[37] = 'a'
	buf[38] = 't'
	buf[39] = 'a'

	initRet := C.TRSpeexDecodeInit(&v.context)
	if int(initRet) <= -1 {
		return pcm, errors.New("init err")
	}
	defer C.TRSpeexDecodeRelease(&v.context)

	// write buf to out result
	ret, err := result.Write(buf)
	if err != nil || ret != len(buf) {
		return
	}

	// read bufSize count
	tmpBufIn := make([]byte, bufSize)
	tmpBufOut := make([]byte, bufSize)
	ret, err = src.Read(tmpBufIn)
	if err != nil || ret != bufSize {
		return pcm, errors.New("read input header err")
	}

	for {
		outBuf := (*C.char)(unsafe.Pointer(&tmpBufOut[0]))
		inBuf := (*C.char)(unsafe.Pointer(&tmpBufIn[0]))
		C.TRSpeexDecode(&v.context, inBuf, C.int(bufSize), outBuf, &nOutSize)
		//fmt.Printf("codec:%v\n", nOutSize)

		ret, err = src.Read(tmpBufIn)
		if err != nil || ret != bufSize {
			fmt.Println("src read err:", err, " total read size:", nTotalLen)
			break
		}

		result.Write(tmpBufOut)

		nTotalLen += int(nOutSize)
	}

	// todo
	// write total len
	// part1 buf[40,41,42,43]
	//tmp := bytes.NewBuffer(make([]byte, 4))
	//binary.Read(tmp, binary.BigEndian, &nTotalLen)
	//result.Write(tmp)
	// part2 buf[4,5,6,7]

	pcm = result.Bytes()
	return
}
