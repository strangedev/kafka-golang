LBIN=./bin
BIN=~/bin

ALL: clean newschema schemaexplorer

clean:
	rm -f $(LBIN)/*

newschema: cmd/schema/newschema.go
	go build -o $(LBIN)/$@ $^

schemaexplorer: cmd/schema/explorer.go
	go build -o $(LBIN)/$@ $^