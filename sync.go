package main

import (          // You can import many packages at once by enclosing all
	"fmt"     // packages inside of parenthesis
	"os"
	"os/exec"
	"bufio"
	"strings"
	"time"
	"log"
)

type root struct {  // There are no classes or inheritance Go
	            // Structs provide variable aggregation
	root string
	dirs, dupes, notin, diff *directory
	md5 map[string] string
	numFiles, numDirs int
		    // Go is a procedural OOP hybrid language but
}                   // but OOP is done Googly with duck typing.

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
	PRINT = true
)

func main() {  // All functions are preceded by the keyword func and 
	       // the main function takes no arguments and returns no values

	treeA, treeB := root{}, root{}  // You may initialize two variable using commas

	if len( os.Args ) == 1 { // The os package allows use of command line arguments

		fmt.Println("Enter path A:")
		fmt.Scan(&treeA.root)
		fmt.Println("Enter path B:")
		fmt.Scan(&treeB.root)
	} else {
		treeA.root = os.Args[1]
		treeB.root = os.Args[2]
	}

	// This if statement runs the BuildTree function and assigns the returned
	// error to err. The 'if' part comes after the ;
	if err := BuildTree(&treeA); err != nil {  

		log.Println(err)
	}

	if err := BuildTree(&treeB); err != nil { // It is bad form to throw out the errors

		log.Println(err)
	}

	if err := CrossCheck(&treeA, &treeB); err != nil {

		log.Println(err)
	}
}

//   Name         Input Variables      Return value
func CrossCheck ( A *root, B *root ) ( err error ) {  // Named returns require parenthesis
	
	A.notin = NewDir(B.root)
	B.notin = NewDir(A.root)
	A.diff = NewDir(A.root)
	B.diff = NewDir(B.root)

	// For loops are versatile and can iterate over values and keys in a map

	for k, v := range B.md5 { // Test for files in B that are not in A
		
		if _, test := A.md5[k]; !test {

			temp := File{Name: B.root + v, Hash: k}
			A.notin.files[A.notin.fileCtr] = &temp
			A.notin.fileCtr++

			if A.notin.fileCtr == len(A.notin.files) { // If Files full resize

				ResizeFile( A.notin )
			}
		}
	}

	for k, v := range A.md5 { // Test for files in A that are not in B

		if _, test := B.md5[k]; !test {

			temp := File{Name: A.root + v, Hash: k}
			B.notin.files[B.notin.fileCtr] = &temp
			B.notin.fileCtr++

			if B.notin.fileCtr == len(B.notin.files) { // If Files full resize
				ResizeFile( B.notin )
			}
		}
	}

	if PRINT {

        fmt.Println("Files not in:", A.root)
        for i := 0; i < A.notin.fileCtr; i++ {

                fmt.Println(A.notin.files[i].Name)
        }

        fmt.Println("Files not in:", B.root)
        for i := 0; i < B.notin.fileCtr; i++ {

                fmt.Println(B.notin.files[i].Name)
        }
	}

	var temp [] string
	var counter int = 0

	if A.numFiles > B.numFiles {

		temp = make ( [] string, A.numFiles )
	} else {
		temp = make ( [] string, B.numFiles )
	}

	for k, _ := range A.md5 {

		if _, test := B.md5[k]; test {

			temp[counter] = k
			counter++
		}
	}

	counter = 0
	for _, i := range temp {

		At := A.md5[i]
		Bt := B.md5[i]

		if ! SamePath(At, A.root, Bt, B.root) {

			ta := File{ Name: A.root + At, Hash: i }
			tb := File{ Name: B.root + Bt, Hash: i }

			A.diff.files[A.diff.fileCtr] = &ta
			B.diff.files[B.diff.fileCtr] = &tb
			counter++
			A.diff.fileCtr++
			B.diff.fileCtr++

			if A.diff.fileCtr == len(A.diff.files) { // If Files full resize

				ResizeFile( A.diff )
			}

			if B.diff.fileCtr == len(B.diff.files) { // If Files full resize

				ResizeFile( B.diff )
			}
		}
	}

	if PRINT {

		fmt.Println("Files in diff locations")
        for i := 0; i < A.diff.fileCtr; i++ {

            fmt.Println(A.diff.files[i].Name)
            fmt.Println(B.diff.files[i].Name)
        }
	}
	return err
}

// Since this function returns an unnamed boolean it requires no parenthesis
func SamePath ( fileA string, rootA string, fileB string, rootB string ) bool {

	A := strings.TrimPrefix ( fileA, rootA )
	B := strings.TrimPrefix ( fileB, rootB )

	if A == B {
		return true
	} else {
		return false
	}
}

func BuildTree ( Root *root ) ( err error ) {

	Root.root = strings.TrimSuffix( Root.root, "/" )

	Root.dirs = NewDir(Root.root)
	Root.dupes = NewDir(Root.root)
	
	if test := CheckDir( Root.root ); test == true {
		
		if err := ExploreTree( Root.dirs ); err != nil {
			
			fmt.Println(err)
		}
	} else {
		os.Exit(1)
	}

	Root.numFiles, Root.numDirs = Count ( Root.dirs )

	Root.md5 = make ( map[string]string, Root.numFiles )

	Md5Map( Root.root, Root.dirs, Root.md5, Root.dupes )

	if PRINT {
		
		fmt.Println("Dupes:")
		for i := 0; i < Root.dupes.fileCtr; i++ {

			fmt.Println(Root.dupes.files[i].Name, Root.dupes.files[i].Hash)
		}
	}
	return err
}

func Md5Map ( prefix string, tree *directory, md5 map[string] string, dupes *directory ) {

	pre := strings.TrimPrefix( tree.root, prefix )

	for i := 0; i < tree.dirCtr; i++ {

		Md5Map ( prefix, tree.dirs[i], md5, dupes )
	}

	for i := 0; i < tree.fileCtr; i++ {

		temp := tree.files[i].Hash

		if _, test := md5[temp]; test { // Test if duplicate hash exists

			temp := File{ Name: prefix + "/" + tree.files[i].Name, Hash: temp }
			dupes.files[dupes.fileCtr] = &temp
			dupes.fileCtr++

			if dupes.fileCtr == len(dupes.files) { // If Files full resize

				ResizeFile( dupes )
			}

			temp2 := md5[tree.files[i].Hash]
			tempFile := File{ Name: prefix + temp2, Hash: tree.files[i].Hash }
			dupes.files[dupes.fileCtr] = &tempFile
			dupes.fileCtr++

			if dupes.fileCtr == len(dupes.files) { // If Files full resize

				ResizeFile( dupes )
			}

		} else { // Else insert hash into map
			md5[ temp ] = pre + "/" + tree.files[i].Name
		}
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

func ExploreTree( tree *directory ) ( err error ) {

	file, err := os.Open ( tree.root ) // err shadow variable
	
	if err != nil {
		return err
	}

	defer file.Close()

	out, err := file.Readdirnames(0) // err shadow variable
	
	if err != nil {
		return err
	}

	for _, i := range out {

		fInfo, _ := os.Stat( tree.root + "/" + i ) // Get current file info

		if fInfo.IsDir() { // If IsDir create and explore new directory

			tree.dirs[ tree.dirCtr ] = NewDir( tree.root + "/" + i )
			tree.dirCtr++

			if err := ExploreTree( tree.dirs[ tree.dirCtr-1 ] ); err != nil {

				return err
			}

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

// Duck typing in action.  These two functions fulfill the Stringer interface
// for both the File and directory structs

func ( file *File ) String() string {

	return fmt.Sprintf("%s\tmd5:%s", file.Name, file.Hash )
}

func ( dir *directory ) String() string {

	return fmt.Sprint( dir.root )
}

// Multiple return values require parenthesis around them.
func Count( dir *directory ) ( files int, dirs int ) {

	files, dirs = dir.fileCtr, dir.dirCtr

	for i := 0; i < dir.dirCtr; i++ {

		ftemp, dtemp := Count( dir.dirs[i] )

		files += ftemp  // += does not support multiple assignments
		dirs += dtemp
	}

	return files, dirs
}

func NewFile( file string ) ( *File, error ) {
	
	temp := File {}

	fInfo, fail := os.Stat(file)

	temp.size = fInfo.Size()
	temp.Name = fInfo.Name()
	temp.modification_time = fInfo.ModTime()

	// I use os.exec to execute md5sum on every file within the directory
	// Go does not block and wait for these operations to complete
	// Once all of the operations have completed the next function
	// will start.
/*
Md5sum
d895987270daac28023d5c1894f6c0c0  /down/Where.the.Buffalo.Roam.1980.720p.HDTVRip.mkv
real    3m48.386s
user    1m2.624s
sys     0m9.588s

997764dcd14b95941ead3e347682adef  /down/This.is.the.End.2013.720p.BluRay.x264-SPARKS.mkv
real    1m58.489s
user    0m33.759s
sys     0m5.511s

Go Imp


*/
	cmd := exec.Command("md5sum", file)
	stdout, _ := cmd.StdoutPipe()
	cmd.Start()
	r := bufio.NewReader(stdout)
	line, _ := r.ReadString(' ')
	temp.Hash = line

	return &temp, fail  // Go's garbage collection is not heavy handed and only
}                           // deletes junk when they are no longer referenced
                            // not when they are no longer in scope.
                            // So go nuts and return pointers to temporary variables


func NewDir( direct string ) *directory {

	temp := directory{ root: direct, file_inc: 1, dir_inc: 1 }

	temp.dirs = make( []*directory, DIR )
	temp.files = make( []*File, FILE )

	return &temp
}

func ResizeDir( tree *directory ) {

	tree.dir_inc++
	
	a := make( []*directory, DIR * tree.dir_inc )
	for i := 0; i < tree.dirCtr; i++ {
		a[i] = tree.dirs[i]
	}

	tree.dirs = a
}

func ResizeFile( tree *directory ) {

	tree.file_inc++

	a := make( []*File, FILE * tree.file_inc )
	for i := 0; i < tree.fileCtr; i++ {
		a[i] = tree.files[i]
	}

	tree.files = a
}
