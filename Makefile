build:
	go build -o kmongo ./kmongo
	go build -o ksecret ./ksecret
	go build -o kpage ./kpage
	go build -o kyake ./kyake

install:
	go install ./kmongo
	go install ./ksecret
	go install ./kpage
	go install ./kyake

clean:
	rm -f kmongo/kmongo
	rm -f kmongo/kmongo.exe
	rm -f ksecret/ksecret
	rm -f ksecret/ksecret.exe
	rm -f kpage/kpage
	rm -f kpage/kpage.exe