# go-speex fork from(https://github.com/winlinvip/speex)

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
cd $GOPATH/src/github.com/winlinvip/go-speex && go test ./...
```

There are an example of SPEEX audio packets in FLV:

* [avatar speex over FLV](doc/speex_data.go), user can use this file to decode to PCM.
* [audio resample](https://github.com/winlinvip/go-aresample).

For more information about SPEEX codec, read:

* [github.com](https://github.com/winlinvip/speex), source code of speex codec.
* [examples](http://www.speex.org/docs/manual/speex-manual/node13.html), encoder and decoder example.

Winlin 2016
