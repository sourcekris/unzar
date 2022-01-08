BINARY=unzar

build:
	go build -o ${BINARY} ${BINARY}.go

run:
	go run ${BINARY}.go

clean:
	go clean ..
	rm bin/${BINARY}-*
	rm release/${BINARY}-*

compile:
	echo "Compiling for every OS and Platform"
	GOOS=linux GOARCH=amd64 go build -o bin/${BINARY}-linux-amd64 ${BINARY}.go
	GOOS=linux GOARCH=386 go build -o bin/${BINARY}-linux-i386 ${BINARY}.go
	GOOS=linux GOARCH=arm go build -o bin/${BINARY}-linux-arm ${BINARY}.go
	GOOS=linux GOARCH=arm64 go build -o bin/${BINARY}-linux-arm64 ${BINARY}.go
	GOOS=freebsd GOARCH=386 go build -o bin/${BINARY}-freebsd-i386 ${BINARY}.go
	GOOS=freebsd GOARCH=amd64 go build -o bin/${BINARY}-freebsd-amd64 ${BINARY}.go
	GOOS=darwin GOARCH=amd64 go build -o bin/${BINARY}-mac-amd64 ${BINARY}.go
	GOOS=darwin GOARCH=arm64 go build -o bin/${BINARY}-mac-arm64 ${BINARY}.go
	GOOS=windows GOARCH=amd64 go build -o bin/${BINARY}-windows-amd64.exe ${BINARY}.go

package:
	echo "Packaging releases into releases/*.zip"
	cd release && cp ../bin/${BINARY}-linux-amd64 ${BINARY} &&  zip ${BINARY}-linux-amd64.zip ${BINARY} && rm ${BINARY}
	cd release && cp ../bin/${BINARY}-linux-i386 ${BINARY} &&  zip ${BINARY}-linux-i386.zip ${BINARY} && rm ${BINARY}
	cd release && cp ../bin/${BINARY}-linux-arm ${BINARY} &&  zip ${BINARY}-linux-arm.zip ${BINARY} && rm ${BINARY}
	cd release && cp ../bin/${BINARY}-linux-arm64 ${BINARY} &&  zip ${BINARY}-linux-arm64.zip ${BINARY} && rm ${BINARY}
	cd release && cp ../bin/${BINARY}-freebsd-amd64 ${BINARY} &&  zip ${BINARY}-freebsd-amd64.zip ${BINARY} && rm ${BINARY}
	cd release && cp ../bin/${BINARY}-freebsd-i386 ${BINARY} &&  zip ${BINARY}-freebsd-i386.zip ${BINARY} && rm ${BINARY}
	cd release && cp ../bin/${BINARY}-mac-amd64 ${BINARY} &&  zip ${BINARY}-mac-amd64.zip ${BINARY} && rm ${BINARY}
	cd release && cp ../bin/${BINARY}-mac-arm64 ${BINARY} &&  zip ${BINARY}-mac-arm64.zip ${BINARY} && rm ${BINARY}
	cd release && cp ../bin/${BINARY}-windows-amd64.exe ${BINARY}.exe &&  zip ${BINARY}-windows-amd64.zip ${BINARY}.exe && rm ${BINARY}.exe

release: compile package
	
all: build