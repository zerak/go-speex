Golang binding for speex(https://github.com/winlinvip/speex)

## Usage

First, get the source code:

```
go get -d github.com/zerak/go-speex
```

Then, compile the speex:

```
cd $GOPATH/src/github.com/zerak/go-speex &&
git clone https://github.com/winlinvip/speex.git speex-lib &&
cd speex-lib/ && bash autogen.sh && ./configure --prefix=`pwd`/objs --enable-static && make && make install &&
cd ..
```

Done, import and use the package:

* [speex decoder](dec/example_test.go), decode the speex frame to PCM samples.

To run all examples:

```
cd $GOPATH/src/github.com/zerak/go-speex && go test ./...
```

the file 4.speex is WeChat format the TestWeChatDec show the code for decode.
