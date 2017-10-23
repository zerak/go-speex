package speex_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/zerak/go-speex/speex"
)

func TestWeChatDec(t *testing.T) {
	f, err := os.OpenFile("../4.speex", os.O_RDONLY, 0666)
	if err != nil {
		t.Error("open weChat original speex file err:", err)
		return
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		t.Error("weChat speex file info err", err)
		return
	}
	frame := make([]byte, info.Size())
	ret, err := f.Read(frame)
	if err != nil || ret != int(info.Size()) {
		t.Error("read weChat speex err:", err, " file size:", info.Size())
		return
	}

	decoder := speex.NewWeChatDecoder()
	pcm, err := decoder.Decode(frame)
	if err != nil {
		t.Error("decode err:", err)
		return
	}
	ioutil.WriteFile("wechat.wav", pcm, 0666)
}
