package filesystem

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"

	_ "bazil.org/fuse/fs/fstestutil"
	"bazil.org/fuse/fuseutil"
	"github.com/crgimenes/pgfs/adapters/postgres"
)

// FS base of filesystem
type FS struct {
	Nodes map[string]*Node
}

// Node in the filesystem
type Node struct {
	fuse    *fs.Server
	fs      *FS
	Inode   uint64
	Name    string
	Type    fuse.DirentType
	Content []byte
}

// Root return root directory
func (f *FS) Root() (fs.Node, error) {
	return &Node{fs: f}, nil
}

// Load required resources for the filesystem
func Load() {
	postgres.Load()
}

// Lookup a node and return
func (n *Node) Lookup(ctx context.Context, name string) (fs.Node, error) {
	log.Println("Lookup", name)
	node, ok := n.fs.Nodes[name]
	if ok {
		return node, nil
	}
	return nil, fuse.ENOENT
}

// ReadDirAll read all files and subdirectories in a directorie
func (n *Node) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	log.Println("ReadDirAll", n.Name)
	var dirDirs []fuse.Dirent
	for _, node := range n.fs.Nodes {
		dirent := fuse.Dirent{
			Inode: node.Inode,
			Name:  node.Name,
			Type:  node.Type,
		}
		dirDirs = append(dirDirs, dirent)
	}
	return dirDirs, nil
}

// Attr return the file attribute
func (n *Node) Attr(ctx context.Context, a *fuse.Attr) (err error) {
	fmt.Println("Attr", n.Name)
	a.Inode = 1
	a.Mode = os.ModeDir | 0555
	if n.Name != "" {
		ext := filepath.Ext(n.Name)
		tableName := strings.TrimSuffix(n.Name, ext)

		switch ext {
		case ".json":
			n.Content, err = postgres.LoadTableJSON(tableName)
			if err != nil {
				return err
			}
		case ".csv":
			n.Content, err = postgres.LoadTableCSV(tableName)
			if err != nil {
				return err
			}
		default:
			n.Content = []byte(n.Name)
		}

	}
	if n.Type == fuse.DT_File {
		a.Inode = n.Inode
		a.Mode = 0444
		a.Size = uint64(len(n.Content))
	}
	return nil
}

func close(c io.Closer) {
	log.Println("closing")
	err := c.Close()
	if err != nil {
		log.Println("error closing", err)
	}
}

// Open file
func (n *Node) Open(ctx context.Context, req *fuse.OpenRequest, resp *fuse.OpenResponse) (fs.Handle, error) {
	log.Println("Open", n.Name)
	if !req.Flags.IsReadOnly() {
		return nil, fuse.Errno(syscall.EACCES)
	}
	resp.Flags |= fuse.OpenKeepCache
	return n, nil
}

// Read file content
func (n *Node) Read(ctx context.Context, req *fuse.ReadRequest, resp *fuse.ReadResponse) error {
	log.Println("Read", n.Name)
	log.Printf("reading file %q from %v to %v, inode %v\n", n.Name, req.Offset, req.Size, n.Inode)
	fuseutil.HandleRead(req, resp, n.Content)
	return nil
}

// Mount the file system
func Mount(mountpoint string) (err error) {
	log.Println("Mounting filesystem")
	c, err := fuse.Mount(
		mountpoint,
		fuse.FSName("pgfs"),
		fuse.Subtype("pgfs"),
		//fuse.ReadOnly(),
		fuse.AllowOther(),
	)
	if err != nil {
		return
	}
	defer close(c)

	if p := c.Protocol(); !p.HasInvalidate() {
		return fmt.Errorf("kernel FUSE support is too old to have invalidations: version %v", p)
	}

	tables, err := postgres.ListTables()
	if err != nil {
		return
	}
	nodes := make(map[string]*Node)
	srv := fs.New(c, nil)

	var inode uint64 = 2
	for _, t := range tables {
		node := Node{
			Name:  t.Name,
			fuse:  srv,
			Inode: inode,
			Type:  fuse.DT_Dir,
			fs: &FS{
				Nodes: map[string]*Node{
					t.Name + ".json": &Node{
						Name:    t.Name + ".json",
						fuse:    srv,
						Inode:   inode + 1,
						Type:    fuse.DT_File,
						Content: []byte(""),
					},
					t.Name + ".csv": &Node{
						Name:    t.Name + ".csv",
						fuse:    srv,
						Inode:   inode + 2,
						Type:    fuse.DT_File,
						Content: []byte(""),
					},
				},
			},
		}
		nodes[t.Name] = &node
		inode += 3
	}

	filesys := &FS{
		Nodes: nodes,
	}

	err = srv.Serve(filesys)
	return
}

// WriteRequest write request
func (n *Node) WriteRequest(req *fuse.WriteRequest) {
	log.Println("WriteRequest", n.Name)
}

// WriteResponse write response
func (n *Node) WriteResponse(resp *fuse.WriteResponse) {
	log.Println("WriteResponse", n.Name)
}

// Write file content
func (n *Node) Write(ctx context.Context, req *fuse.WriteRequest, resp *fuse.WriteResponse) error {
	log.Println("Write", n.Name)
	log.Printf("writing file %q from %v, inode %v\n", n.Name, req.Offset, n.Inode)
	return nil
}

// Unmount the file system
func Unmount(mountpoint string) (err error) {
	log.Println("Unmounting filesystem")
	err = fuse.Unmount(mountpoint)
	return
}
