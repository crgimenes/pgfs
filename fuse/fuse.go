package fuse

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	_ "bazil.org/fuse/fs/fstestutil"
	"bazil.org/fuse/fuseutil"
	"github.com/crgimenes/pgfs/adapters/postgres"
	"github.com/nuveo/log"
)

type FS struct {
	Nodes map[string]*Node
}

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
	fmt.Println("Root", len(f.Nodes))
	return &Node{fs: f}, nil
}

func (n *Node) Lookup(ctx context.Context, name string) (fs.Node, error) {
	fmt.Println("Lookup", n.Name, n.Inode)
	node, ok := n.fs.Nodes[name]
	if ok {
		return node, nil
	}
	return nil, fuse.ENOENT
}

func (n *Node) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	fmt.Println("ReadDirAll", n.Name, n.Inode)
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

func (n *Node) Attr(ctx context.Context, a *fuse.Attr) error {
	fmt.Println("Attr", n.Name, n.Inode)
	a.Inode = 1
	a.Mode = os.ModeDir | 0555
	if n.Name != "" {
		ext := filepath.Ext(n.Name)
		tableName := strings.TrimSuffix(n.Name, ext)

		t, err := postgres.LoadTable(tableName)
		if err != nil {
			return err
		}

		switch ext {
		case ".json":
			n.Content, err = json.MarshalIndent(t, "", "\t")
			if err != nil {
				return err
			}
		case ".csv":
			n.Content = []byte("not implemented")
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
	err := c.Close()
	if err != nil {
		log.Errorln(err)
	}
}

func (n *Node) Open(ctx context.Context, req *fuse.OpenRequest, resp *fuse.OpenResponse) (fs.Handle, error) {
	fmt.Println("Open", n.Name, n.Inode)
	if !req.Flags.IsReadOnly() {
		return nil, fuse.Errno(syscall.EACCES)
	}
	resp.Flags |= fuse.OpenKeepCache
	return n, nil
}

func (n *Node) Read(ctx context.Context, req *fuse.ReadRequest, resp *fuse.ReadResponse) error {
	fmt.Printf("Reading file %q from %v to %v, inode %v\n", n.Name, req.Offset, req.Size, n.Inode)
	fuseutil.HandleRead(req, resp, n.Content)
	return nil
}

func Run(mountpoint string) (err error) {
	c, err := fuse.Mount(
		mountpoint,
		fuse.FSName("pgfs"),
		fuse.Subtype("pgfs"),
		fuse.LocalVolume(),
		fuse.VolumeName("Postgresql filesystem"),
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
						Content: []byte("test file 4\n"),
					},
					t.Name + ".csv": &Node{
						Name:    t.Name + ".csv",
						fuse:    srv,
						Inode:   inode + 2,
						Type:    fuse.DT_File,
						Content: []byte("test file 4\n"),
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
	if err != nil {
		return
	}

	// Check if the mount process has an error to report.
	<-c.Ready
	err = c.MountError
	return
}
