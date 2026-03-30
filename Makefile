build:
	go build -o kmongo ./kmongo
	go build -o ksecret ./ksecret

install:
	go install ./kmongo
	go install ./ksecret

clean:
	rm -f kmongo/kmongo
	rm -f kmongo/kmongo.exe
	rm -f ksecret/ksecret
	rm -f ksecret/ksecret.exe