[![Go Report Card](https://goreportcard.com/badge/maxzender/jv)](https://goreportcard.com/report/maxzender/jv)
[![Build Status](https://travis-ci.org/maxzender/jv.svg?branch=master)](https://travis-ci.org/maxzender/jv)

# jv
jv (for jsonviewer) helps you view your JSON.

<p align="center">
    <img src="http://huangnauh.github.io/examples/jv/example.svg">
    <img src="http://huangnauh.github.io/examples/jv/jsonline.svg">
</p>

## Installation
```
git clone git@github.com:huangnauh/jv.git $GOPATH/src/github.com/maxzender/jv
cd $GOPATH/src/github.com/maxzender/jv && go get ./...
```

## usage

```
Usage: ./jv [-o] <query> [file]
  -h	print usage
  -help
    	print usage
  -o	pretty output
```

## Example usage
```
cat file.json | ./jv aaa
```


