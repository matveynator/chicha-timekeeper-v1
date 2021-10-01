version="0.2-010"
git_root_path=`git rev-parse --show-toplevel`
execution_file=chicha
cd ${git_root_path}/Scripts
for os in linux freebsd netbsd openbsd aix android illumos ios solaris plan9 darwin dragonfly js windows ;
do
	for arch in "amd64" "386" "arm" "arm64" "mips64" "mips64le" "mips" "mipsle" "ppc64" "ppc64le" "riscv64" "s390x" "wasm"  
	do
		mkdir -p ../downloads/${version}/${os}/${arch}
		[ "$os" == "windows" ] && execution_file="chicha.exe"
		[ "$os" == "js" ] && execution_file="chicha.js"
		GOOS=${os} GOARCH=${arch} go build -o ../downloads/${version}/${os}/${arch}/${execution_file} ../chicha.go &> /dev/null
		if [ "$?" != "0" ] 
		#if compilation failed - remove folders - else copy config file.
		then 
		  rm -rf ../downloads/${version}/${os}/${arch}
		else
		  echo "GOOS=${os} GOARCH=${arch} go build -o ../downloads/${version}/${os}/${arch}/${execution_file} ../chicha.go"
		  cat ../chicha.conf > ../downloads/${version}/${os}/${arch}/chicha.conf
		fi
	done
done

# irsync -avP ../downloads/* root@files.matveynator.ru:/home/files/public_html/chicha/

