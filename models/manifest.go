package models

type Manifest struct {
	Chunks     []Chunks    `json:"Chunks"`
	Name       string      `json:"Name"`
	ChunkCount int         `json:"ChunkCount"`
	ChunkSize  int         `json:"ChunkSize"`
	MaxSize    int         `json:"MaxSize"`
	Size       int         `json:"Size"`
	Hash       string      `json:"Hash"`
	Iv         interface{} `json:"IV"`
}

type Chunks struct {
	ID   string `json:"Id"`
	Ord  int    `json:"Ord"`
	Hash string `json:"Hash"`
}
