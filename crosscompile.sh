version="0.1"

for os in linux freebsd netbsd openbsd plan9;
do
  for arch in amd64 "386" arm arm64 
  do
    mkdir -p ../../download/$os/$arch/$version
    echo "GOOS=$os GOARCH=$arch $opt go build -o ../../download/$os/$arch/$version/chicha chicha.go"
    GOOS=$os GOARCH=$arch $opt go build -o ../../download/$os/$arch/$version/chicha chicha.go
    cat .env.DEFAULT > ../../download/$os/$arch/$version/.env.DEFAULT
  done
done

#mac
for os in darwin;
do
  for arch in amd64 "386" arm64
  do
    mkdir -p ../../download/mac/$arch/$version
    echo "GOOS=$os GOARCH=$arch CGO_ENABLED=0 go build -o ../../download/mac/$arch/$version/chicha chicha.go"
    GOOS=$os GOARCH=$arch CGO_ENABLED=0 go build -o ../../download/mac/$arch/$version/chicha chicha.go
    cat .env.DEFAULT > ../../download/$os/$arch/$version/.env.DEFAULT
  done
done

#dragonfly
for os in dragonfly;
do
  for arch in amd64
  do
    mkdir -p ../../download/$os/$arch/$version
    echo "GOOS=$os GOARCH=$arch go build -o ../../download/$os/$arch/$version/chicha chicha.go"
    GOOS=$os GOARCH=$arch go build -o ../../download/$os/$arch/$version/chicha chicha.go
    cat .env.DEFAULT > ../../download/$os/$arch/$version/.env.DEFAULT
  done
done

#windows
for os in windows;
do
  for arch in amd64 "386"
  do
    mkdir -p ../../download/$os/$arch/$version
    echo "GOOS=$os GOARCH=$arch CGO_ENABLED=0 go build -o ../../download/$os/$arch/$version/chicha.exe chicha.go"
    GOOS=$os GOARCH=$arch CGO_ENABLED=0 go build -o ../../download/$os/$arch/$version/chicha.exe chicha.go
    cat .env.DEFAULT > ../../download/$os/$arch/$version/.env.DEFAULT
  done
done

#rsync -avP --delete ../../download root@matveynator.ru:/home/chicha/public_html/
