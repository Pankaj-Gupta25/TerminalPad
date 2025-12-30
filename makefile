build:
	@go build -o TerminalPad.exe .

run: build
	@TerminalPad.exe
