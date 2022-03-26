

htmleditor-linux-amd64:
	go build $(ARGS) -o htmleditor-linux-amd64 ./htmleditor

htmleditor-linux-386:
	go build $(ARGS) -o htmleditor-linux-386 ./htmleditor

htmleditor-darwin-amd64:
	go build $(ARGS) -o htmleditor-darwin-amd64 ./htmleditor

htmleditor-windows-amd64:
	go build $(ARGS) -o htmleditor-windows-amd64.exe ./htmleditor

htmleditor-windows-386:
	go build $(ARGS) -o htmleditor-windows-386.exe ./htmleditor

all: htmleditor-linux-amd64 htmleditor-linux-386 htmleditor-darwin-amd64 htmleditor-windows-amd64 htmleditor-windows-386