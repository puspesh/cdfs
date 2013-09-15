package main

type Counters struct {
	totalSize uint32
	totalTime uint32
}

type GoogleDrive struct {
	Store      string
	authString string
	Counters
}

type FileApi interface {
	Upload(string)
	CheckSize() (error, bool)
}
