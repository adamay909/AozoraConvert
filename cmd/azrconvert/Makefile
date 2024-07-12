default: 
	go build

linux:
	GOOS=linux GOARCH=amd64 go build

macintel:
	GOOS=darwin GOARCH=amd64 go build -o azrconvert_INTEL

applesilicon:
	GOOS=darwin GOARCH=arm64 go build -o azrconvert_APPLESILICON

windows:
	GOOS=windows GOARCH=amd64 go build -o azrconvert.exe

all: linux macintel applesilicon windows
