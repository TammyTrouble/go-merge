package main

import (
	"fmt"
	"os"
	"os/exec"
	"bufio"
	"strings"
	"time"
)

type directory struct {

	root string
	dirs []*directory
	files []*File
	fileCtr, dirCtr, file_inc, dir_inc int

}

type File struct {

	Name, Hash string
	size int64
	modification_time time.Time	

}

const (
	DIR = 10 
	FILE = 50
	layout = "Jan 2, 2006 3:04pm"
)

func main() {

	treeA := directory { file_inc: 1, dir_inc: 1 }
	treeA.dirs = make( []*directory, DIR )
	treeA.files = make( []*File, FILE )

	treeB := directory { file_inc: 1, dir_inc: 1 }
	treeB.dirs = make( []*directory, DIR )
	treeB.files = make( []*File, FILE )

	if len( os.Args ) == 1 {

		fmt.Println("Enter path A:")
		fmt.Scan(&treeA.root)
		fmt.Println("Enter path B:")
		fmt.Scan(&treeB.root)
	} else {
		treeA.root = os.Args[1]
		treeB.root = os.Args[2]
	}

	if test := CheckDir( treeA.root ); test == true {
		
		if err := ExploreTree( &treeA ); err != nil {
			
			fmt.Println(err)
		}
	} else {
		os.Exit(1)
	}

	if test := CheckDir( treeB.root ); test == true {
		
		if err := ExploreTree( &treeB ); err != nil {

			fmt.Println(err)
		}
	} else {
		os.Exit(1)
	}

	filesA, dirsA := Count (&treeA)
	filesB, dirsB := Count (&treeB)

	fmt.Println(dirsA, dirsB)

	dupesA := NewDir(treeA.root)
	dupesB := NewDir(treeB.root)
	var md5A = make(map[string]string, filesA)
	Md5Map( treeA.root, &treeA, md5A, dupesA )

	var md5B = make(map[string]string, filesB)
	Md5Map( treeB.root, &treeB, md5B, dupesB )

	fmt.Println("A")
	for k, v := range md5A {
		fmt.Println("k:", k, "v:", v)
	}
	fmt.Println("B")
	for k, v := range md5B {
		fmt.Println("k:", k, "v:", v)
	}

	fmt.Println("A")
	for i := 0; i < dupesA.fileCtr; i++ {
		fmt.Println("Dupes:", dupesA.files[i].Name, dupesA.files[i].Hash)
	}
	fmt.Println("B")
	for i := 0; i < dupesB.fileCtr; i++ {
		fmt.Println("Dupes:", dupesB.files[i].Name, dupesB.files[i].Hash)
	}

	// Create maps.

	// Make recommendations


}

func Md5Map ( prefix string, tree *directory, md5 map[string] string, dupes *directory )  {

	pre := strings.TrimPrefix( tree.root, prefix )

	for i := 0; i < tree.dirCtr; i++ {
		Md5Map ( prefix, tree.dirs[i], md5, dupes )
	}

	for i := 0; i < tree.fileCtr; i++ {

		temp := tree.files[i].Hash

		if _, test := md5[temp]; test { // Test if duplicate hash exists

			temp := File{Name: prefix + "/" + tree.files[i].Name, Hash: temp}
			dupes.files[dupes.fileCtr] = &temp
			dupes.fileCtr++

		} else { // Else insert hash into map
			md5[ temp ] = pre + "/" + tree.files[i].Name
		}
	}

}

func FillDupes ( tree *directory, md5 map[string] string, dupes *directory ) {
	
	for i := 0; i < dupes.fileCtr; i++ {

		temp := md5[dupes.files[i].Hash]
		tempFile := File{Name: tree.root + "/" + temp, Hash: dupes.files[i].Hash}
		dupes.files[dupes.fileCtr] = &tempFile
		dupes.fileCtr++
		fmt.Println(tree.root + "/" + temp)
	}
}

func CheckDir ( x string ) bool {

	info, _ := os.Stat( x )

	if info.IsDir() {
		return true
	} else {
		fmt.Println(x, "is not a directory.")
		return false
	}
	return false
}

func ExploreTree( tree *directory ) error {

	file, err := os.Open(tree.root)

	defer file.Close()

	out, err := file.Readdirnames(0)

	for _, i := range out {

		fInfo, _ := os.Stat( tree.root + "/" + i ) // Get current file info

		if fInfo.IsDir() { // If IsDir create and explore new directory

			tree.dirs[ tree.dirCtr ] = NewDir( tree.root + "/" + i )
			tree.dirCtr++

			ExploreTree( tree.dirs[ tree.dirCtr-1 ] )

			if tree.dirCtr == len(tree.dirs) { // If directories full resize
				ResizeDir( tree )
			}
		} else {

			tree.files[ tree.fileCtr ], _ = NewFile( tree.root + "/" + i )
			tree.fileCtr++

			if tree.fileCtr == len(tree.files) { // If Files full resize
				ResizeFile( tree )
			}
		} // End if IsDir
	} // End for loop

	return err;
}

func PrintTree( tree *directory ) {

	fmt.Println( tree.root )


	for i := 0; i < tree.dirCtr; i++ {
		PrintTree( tree.dirs[i] )
	}

	for i := 0; i < tree.fileCtr; i++ {
		fmt.Println( tree.files[i] )
	}

}

func ( file *File ) String() string {

	return fmt.Sprintf("  %s\tsize: %d\tmd5:%s\tmod_time: %s", file.Name, file.size, file.Hash, file.modification_time.Format(layout) )
}

func ( dir *directory ) String() string {

	return fmt.Sprint( dir.root )
}

func Count( dir *directory ) ( files int, dirs int ) {

	files, dirs = dir.fileCtr, dir.dirCtr

	for i := 0; i < dir.dirCtr; i++ {

		ftemp, dtemp := Count( dir.dirs[i] )

		files += ftemp
		dirs += dtemp
	}

	return files, dirs
}

func NewFile( file string ) ( *File, error ) {
	
	temp := File{}

	fInfo, fail := os.Stat(file)

		temp.size = fInfo.Size()
		temp.Name = fInfo.Name()
		temp.modification_time = fInfo.ModTime()

		cmd := exec.Command("md5sum", file)
		stdout, fail := cmd.StdoutPipe()
		cmd.Start()
		r := bufio.NewReader(stdout)
		line, fail := r.ReadString(' ')

		temp.Hash = line

	return &temp, fail
}

func NewDir( direct string ) ( *directory ) {

	temp := directory{ root: direct, file_inc: 1, dir_inc: 1 }

	temp.dirs = make( []*directory, DIR )
	temp.files = make( []*File, FILE )

	return &temp
}

func ResizeDir( tree *directory ) {

	tree.dir_inc++
	
	a := make([]*directory, DIR * tree.dir_inc)
	for i := range tree.dirs {
		a[i] = tree.dirs[i]
	}

	tree.dirs = a
}

func ResizeFile( tree *directory ) {

	tree.file_inc++

	a := make([]*File, FILE * tree.file_inc)
	for i := range tree.files {
		a[i] = tree.files[i]
	}

	tree.files = a
}

/*
type FileInfo interface {
        Name() string       // base name of the file
        Size() int64        // length in bytes for regular files; system-dependent for others
        Mode() FileMode     // file mode bits
        ModTime() time.Time // modification time
        IsDir() bool        // abbreviation for Mode().IsDir()
        Sys() interface{}   // underlying data source (can return nil)
}
*/