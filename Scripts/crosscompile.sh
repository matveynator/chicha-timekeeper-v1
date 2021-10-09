#!/bin/bash
version="0.2-017"
git_root_path=`git rev-parse --show-toplevel`
execution_file=chicha
cd ${git_root_path}/Scripts
for os in linux freebsd netbsd openbsd aix android illumos ios solaris plan9 darwin dragonfly js windows ;
do
	for arch in "amd64" "386" "arm" "arm64" "mips64" "mips64le" "mips" "mipsle" "ppc64" "ppc64le" "riscv64" "s390x" "wasm"
	do
		target_os_name=${os}
		[ "$os" == "windows" ] && execution_file="chicha.exe"
		[ "$os" == "js" ] && execution_file="chicha.js"
		[ "$os" == "darwin" ] && target_os_name="mac"
		
		mkdir -p ../downloads/${version}/${target_os_name}/${arch}
		GOOS=${os} GOARCH=${arch} go build -ldflags "-X chicha/Packages/Config.VERSION=${version}" -o ../downloads/${version}/${target_os_name}/${arch}/${execution_file} ../chicha.go &> /dev/null
		if [ "$?" != "0" ]
		#if compilation failed - remove folders - else copy config file.
		then
		  rm -rf ../downloads/${version}/${target_os_name}/${arch}
		else
		  echo "GOOS=${os} GOARCH=${arch} go build -ldflags "-X chicha/Packages/Config.VERSION=${version}" -o ../downloads/${version}/${target_os_name}/${arch}/${execution_file} ../chicha.go"
		  cat ../chicha.conf > ../downloads/${version}/${target_os_name}/${arch}/chicha.conf
		fi
	done
done

rsync -avP ../downloads/* root@files.matveynator.ru:/home/files/public_html/chicha/

