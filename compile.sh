#!/bin/sh

project_name="daemon"
release_version="0.0.1"

release_dir=./release
rm -rf $release_dir/*
mkdir -p $release_dir

cd  $(dirname $0)

gofmt -w ./

for goos in "linux" "darwin" "freebsd" "windows"
do
    if [ "$goos" == "windows" ]; then
      obj_name=$project_name.exe
    else
      obj_name=$project_name
    fi

    GOOS=$goos GOARCH=amd64 go build
    zip $release_dir/$project_name-$goos-amd64.zip $obj_name
    GOOS=$goos GOARCH=386 go build
    zip $release_dir/$project_name-$goos-386.zip $obj_name
    rm -f $obj_name
done

cd $release_dir
for file in ./*
do
    md5 -r $file >> sha1sum.txt
done
