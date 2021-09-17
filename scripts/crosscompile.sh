version="0.2"
git_root_path=`git rev-parse --show-toplevel`
cd ${git_root_path}/scripts
for os in linux freebsd netbsd openbsd;
do
  for arch in "amd64" "386" "arm" "arm64" 
  do
    mkdir -p ../downloads/${version}/${os}/${arch}
    echo "GOOS=${os} GOARCH=${arch} $opt go build -o ../downloads/${version}/${os}/${arch}/chicha ../chicha.go"
    GOOS=${os} GOARCH=${arch} $opt go build -o ../downloads/${version}/${os}/${arch}/chicha ../chicha.go
    cat ../.env.DEFAULT > ../downloads/${version}/${os}/${arch}/.env.DEFAULT
  done
done

for os in plan9;
do
  for arch in "amd64" "386" "arm"
  do
    mkdir -p ../downloads/${version}/${os}/${arch}
    echo "GOOS=${os} GOARCH=${arch} $opt go build -o ../downloads/${version}/${os}/${arch}/chicha ../chicha.go"
    GOOS=${os} GOARCH=${arch} $opt go build -o ../downloads/${version}/${os}/${arch}/chicha ../chicha.go
    cat ../.env.DEFAULT > ../downloads/${version}/${os}/${arch}/.env.DEFAULT
  done
done

#mac
for os in darwin;
do
  for arch in "amd64" 
  do
    mkdir -p ../downloads/mac/${arch}/${version}
    mkdir -p ../downloads/mac/${arch}/${version}
    echo "GOOS=${os} GOARCH=${arch} CGO_ENABLED=0 go build -o ../downloads/mac/${arch}/${version}/chicha ../chicha.go"
    GOOS=${os} GOARCH=${arch} CGO_ENABLED=0 go build -o ../downloads/mac/${arch}/${version}/chicha ../chicha.go
    cat ../.env.DEFAULT > ../downloads/mac/${arch}/${version}/.env.DEFAULT
  done
done

#dragonfly
for os in dragonfly;
do
  for arch in "amd64"
  do
    mkdir -p ../downloads/${version}/${os}/${arch}
    echo "GOOS=${os} GOARCH=${arch} go build -o ../downloads/${version}/${os}/${arch}/chicha ../chicha.go"
    GOOS=${os} GOARCH=${arch} go build -o ../downloads/${version}/${os}/${arch}/chicha ../chicha.go
    cat ../.env.DEFAULT > ../downloads/${version}/${os}/${arch}/.env.DEFAULT
  done
done

#windows
for os in windows;
do
  for arch in "amd64" "386"
  do
    mkdir -p ../downloads/${version}/${os}/${arch}
    echo "GOOS=${os} GOARCH=${arch} CGO_ENABLED=0 go build -o ../downloads/${version}/${os}/${arch}/chicha.exe ../chicha.go"
    GOOS=${os} GOARCH=${arch} CGO_ENABLED=0 go build -o ../downloads/${version}/${os}/${arch}/chicha.exe ../chicha.go
    cat ../.env.DEFAULT > ../downloads/${version}/${os}/${arch}/.env.DEFAULT
  done
done

