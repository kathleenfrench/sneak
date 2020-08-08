package store

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/kathleenfrench/common/fs"
	kval "github.com/kval-access-language/kval-boltdb"
	"github.com/olekukonko/tablewriter"
)

type manager struct {
	kb         kval.Kvalboltdb
	bucket     string
	currentLoc string
	lastLoc    string
	rootBucket bool
	viewer     formatter
}

type formatter interface {
	DumpBuckets([]bucket)
	DumpBucketItems(string, []item)
}

type item struct {
	Key   string
	Value string
}

type bucket struct {
	Name string
}

// Audit gives the ability to see what's happening in the DB
func Audit(dbFilepath string) error {
	color.HiBlue("db filepath: %s", dbFilepath)
	var m manager

	m = manager{
		viewer: &dbDisplay{},
	}

	if !fs.FileExists(dbFilepath) {
		return errors.New("no database exists")
	}

	m.connect(dbFilepath)
	defer kval.Disconnect(m.kb)

	m.readInput()
	return nil
}

func (m *manager) readInput() {
	m.bucketlist()
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		bucket := s.Text()
		fmt.Fprintln(os.Stdout, "")
		switch bucket {
		case "\x18":
			return
		case "\x02":
			if !strings.Contains(m.currentLoc, "") || !strings.Contains(m.currentLoc, ">>") {
				fmt.Fprintf(os.Stdout, "%s\n", "> going back...")
				m.currentLoc = ""
				m.bucketlist()
			} else {
				m.bucketItems(bucket, true)
			}
		case "":
			m.bucketlist()
		default:
			m.bucketItems(bucket, false)
		}

		bucket = ""
	}
}

func (m *manager) updateLoc(bucket string, goBack bool) string {
	// we've probably an invalid value and want to display
	// ourselves again...
	if bucket == m.lastLoc {
		m.currentLoc = bucket
		return m.currentLoc
	}

	// handle goback
	if goBack {
		s := strings.Split(m.currentLoc, ">>")
		m.currentLoc = strings.Join(s[:len(s)-1], ">>")
		m.bucket = strings.Trim(s[len(s)-2], " ")
		return m.currentLoc
	}

	// handle location on merit...
	if m.currentLoc == "" {
		m.currentLoc = bucket
		m.bucket = bucket
	} else {
		m.currentLoc = m.currentLoc + " >> " + bucket
		m.bucket = bucket
	}
	return m.currentLoc
}

func (m *manager) bucketlist() {
	m.rootBucket = true
	m.currentLoc = ""

	buckets := []bucket{}

	res, err := kval.Query(m.kb, "GET _") // KVAL: "GET _" will return ROOT
	if err != nil {
		panic(err)
	}
	for k := range res.Result {
		buckets = append(buckets, bucket{Name: string(k) + "*"})
	}

	fmt.Fprint(os.Stdout, "DB Layout:\n\n")
	m.viewer.DumpBuckets(buckets)
	fmt.Fprintf(os.Stdout, "\n%s\n\n", "> Enter bucket to explore (CTRL-X to quit, CTRL-B to go back, ENTER to go back to ROOT Bucket):")
}

func (m *manager) bucketItems(bucketName string, goBack bool) {
	items := []item{}
	getQuery := m.updateLoc(bucketName, goBack)
	if getQuery != "" {
		fmt.Fprintf(os.Stdout, "Query: "+getQuery+"\n\n")
		res, err := kval.Query(m.kb, "GET "+getQuery)
		if err != nil {
			if err.Error() != "Cannot GOTO bucket, bucket not found" {
				log.Fatal(err)
			} else {
				fmt.Fprintf(os.Stdout, "> Bucket not found\n")
				if m.rootBucket == true {
					m.bucketlist()
					return
				}
				m.bucketItems(m.currentLoc, true)
			}
		}
		if len(res.Result) == 0 {
			fmt.Fprintf(os.Stdout, "Invalid request.\n\n")
			m.bucketItems(m.lastLoc, false)
			return
		}

		for k, v := range res.Result {
			if v == kval.Nestedbucket {
				k = k + "*"
				v = ""
			}
			items = append(items, item{Key: string(k), Value: string(v)})
		}
		fmt.Fprintf(os.Stdout, "Bucket: %s\n", bucketName)
		m.viewer.DumpBucketItems(m.bucket, items)
		m.rootBucket = false // success this far means we're not at ROOT
		m.lastLoc = getQuery // so we can also set the query cache for paging
		fmt.Fprintf(os.Stdout, "\n%s\n\n", "> Enter bucket to explore (CTRL-X to quit, CTRL-B to go back, ENTER to go back to ROOT Bucket):")
	}
}

// connect establishes a connection with the bolt DB file
func (m *manager) connect(file string) {
	var err error

	m.kb, err = kval.Connect(file)
	if err != nil {
		log.Fatal(err)
	}
}

type dbDisplay struct{}

func (d dbDisplay) DumpBuckets(bs []bucket) {
	t := tablewriter.NewWriter(os.Stdout)
	t.SetHeader([]string{"buckets"})
	for _, bucket := range bs {
		row := []string{bucket.Name}
		t.Append(row)
	}
	t.Render()
}

func (d dbDisplay) DumpBucketItems(bucket string, items []item) {
	t := tablewriter.NewWriter(os.Stdout)
	t.SetHeader([]string{"Key", "Value"})
	for _, i := range items {
		row := []string{i.Key, i.Value}
		t.Append(row)
	}

	t.Render()
}
