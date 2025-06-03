package client

import "github.com/winfsp/cgofuse/fuse"

type fuseFS struct {
	fuse.FileSystemBase
}

func (fs *fuseFS) Open(path string, flags int) (errc int, fh uint64) {
	switch path {
	case "/" + "John.txt":
		return 0, 0
	default:
		return fuse.ENOENT, ^uint64(0)
	}
}

func (fs *fuseFS) Getattr(path string, stat *fuse.Stat_t, fh uint64) (errc int) {
	Log(path)
	switch path {
	case "/":
		stat.Mode = fuse.S_IFDIR | 0555
		return 0
	case "/" + "Jimmy.txt":
		fallthrough
	case "/" + "John.txt":
		stat.Mode = fuse.S_IFREG | 0444
		stat.Size = int64(len("Skib"))
		return 0
	default:
		return fuse.ENOENT
	}
}

func (fs *fuseFS) Read(path string, buff []byte, ofst int64, fh uint64) (n int) {
	endofst := ofst + int64(len(buff))
	if endofst > int64(len("Skib")) {
		endofst = int64(len("Skib"))
	}
	if endofst < ofst {
		return 0
	}

	n = copy(buff, "Skib"[ofst:endofst])
	return
}

func (fs *fuseFS) Release(path string, fh uint64) (errc int) {
	// Nothing to do
	if path != "/"+"John.txt" {
		return fuse.ENOENT
	}
	return 0
}

func (fs *fuseFS) Readdir(path string, fill func(name string, stat *fuse.Stat_t, ofst int64) bool, ofst int64, fh uint64) (errc int) {
	fill("John.txt", nil, 0)
	fill("Jimmy.txt", nil, 0)
	return 0
}

func (c *IClient) Mount() {
	c.FuseFS = &fuseFS{}

	host := fuse.NewFileSystemHost(c.FuseFS)
	host.Mount("", make([]string, 0))
}
