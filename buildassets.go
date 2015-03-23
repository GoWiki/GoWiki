// +build ignore

package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
)

func main() {
	themesdir, _ := os.Open("themes")
	themes, _ := themesdir.Readdirnames(0)

	asset, _ := os.Create("assets.go")
	fmt.Fprint(asset, `package main

import (
	"reflect"
	"unsafe"
)

func nocopy_bytes(data string) ([]byte, error) {
	var empty [0]byte
	sx := (*reflect.StringHeader)(unsafe.Pointer(&data))
	b := empty[:]
	bx := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	bx.Data = sx.Data
	bx.Len = len(data)
	bx.Cap = bx.Len
	return b, nil
}		

`)

	for _, theme := range themes {
		buf := &bytes.Buffer{}
		gbuf, _ := gzip.NewWriterLevel(buf, gzip.BestCompression)
		tgbuf := tar.NewWriter(gbuf)
		recurseadd(tgbuf, "themes/"+theme+"/static", "static")
		recurseadd(tgbuf, "themes/"+theme+"/templates", "templates")
		tgbuf.Close()
		gbuf.Close()
		fmt.Fprintf(asset, "const _theme_%s_bytes = %q\n", theme, buf.Bytes())
		fmt.Fprintf(asset, `func _theme_%s() ([]byte, error) {
	return nocopy_bytes(_theme_%s_bytes)
}
`, theme, theme)
	}

	fmt.Fprint(asset, `var ThemeTars map[string]func()([]byte, error) = map[string]func()([]byte, error){
`)
	for _, theme := range themes {
		fmt.Fprintf(asset, `	%q: _theme_%s,
`, theme, theme)
	}
	fmt.Fprint(asset, "}")
}

func recurseadd(tarfile *tar.Writer, dir string, tardir string) {
	curdir, _ := os.Open(dir)
	curdiritems, _ := curdir.Readdir(0)
	for _, item := range curdiritems {
		if item.IsDir() {
			recurseadd(tarfile, dir+"/"+item.Name(), tardir+"/"+item.Name())
		} else {
			header, _ := tar.FileInfoHeader(item, "")
			header.Name = tardir + "/" + item.Name()
			tarfile.WriteHeader(header)
			file, _ := os.Open(dir + "/" + item.Name())
			io.Copy(tarfile, file)
		}
	}
}
