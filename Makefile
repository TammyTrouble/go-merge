go: sync.go
	go run sync.go /home/tim/Dropbox/dev/sync/A /home/tim/Dropbox/dev/sync/B

app.out: sync.go
	go build -o app.out sync.go

run: app.out
	./app.out /home/tim/Dropbox/dev/sync/A /home/tim/Dropbox/dev/sync/B

fail: app.out
	./app.out /home/tim/Dropbox/dev/sync/A/dne /home/tim/Dropbox/dev/sync/B
	./app.out /home/tim/Dropbox/dev/sync/A/Whitfield.jpg /home/tim/Dropbox/sync/B

clean:
	rm *.out
