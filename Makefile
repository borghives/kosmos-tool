build:
	go build -o kmongo ./kmongo
	go build -o ksecret ./ksecret
	go build -o kpage ./kpage

install:
	go install ./kmongo
	go install ./ksecret
	go install ./kpage

clean:
	rm -f kmongo/kmongo
	rm -f kmongo/kmongo.exe
	rm -f ksecret/ksecret
	rm -f ksecret/ksecret.exe
	rm -f kpage/kpage
	rm -f kpage/kpage.exe