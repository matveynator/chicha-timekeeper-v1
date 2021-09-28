version="0.2-007"
git_root_path=`git rev-parse --show-toplevel`
cd ${git_root_path}/Scripts
for os in linux freebsd netbsd openbsd;
do
	for arch in "amd64" "386" "arm" "arm64" 
	do
		mkdir -p ../downloads/${version}/${os}/${arch}
		echo "GOOS=${os} GOARCH=${arch} go build -o ../downloads/${version}/${os}/${arch}/chicha ../chicha.go"
		GOOS=${os} GOARCH=${arch} go build -o ../downloads/${version}/${os}/${arch}/chicha ../chicha.go
		cat ../chicha.conf > ../downloads/${version}/${os}/${arch}/chicha.conf
	done
done

for os in plan9;
do
	for arch in "amd64" "386" "arm"
	do
		mkdir -p ../downloads/${version}/${os}/${arch}
		echo "GOOS=${os} GOARCH=${arch} $opt go build -o ../downloads/${version}/${os}/${arch}/chicha ../chicha.go"
		GOOS=${os} GOARCH=${arch} $opt go build -o ../downloads/${version}/${os}/${arch}/chicha ../chicha.go
		cat ../chicha.conf > ../downloads/${version}/${os}/${arch}/chicha.conf
	done
done

#mac
for os in darwin;
do
	for arch in "amd64" 
	do
		mkdir -p ../downloads/${version}/mac/${arch}
		echo "GOOS=${os} GOARCH=${arch} CGO_ENABLED=0 go build -o ../downloads/${version}/mac/${arch}/chicha ../chicha.go"
		GOOS=${os} GOARCH=${arch} CGO_ENABLED=0 go build -o ../downloads/${version}/mac/${arch}/chicha ../chicha.go
		cat ../chicha.conf > ../downloads/${version}/mac/${arch}/chicha.conf
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
		cat ../chicha.conf > ../downloads/${version}/${os}/${arch}/chicha.conf
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
		cat ../chicha.conf > ../downloads/${version}/${os}/${arch}/chicha.conf
	done
done

rsync -avP ../downloads/* root@files.matveynator.ru:/home/files/public_html/chicha/
