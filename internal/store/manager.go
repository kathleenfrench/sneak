package store

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/fatih/color"
	"github.com/kathleenfrench/common/fs"
	"github.com/kathleenfrench/sneak/pkg/htb"
	kval "github.com/kval-access-language/kval-boltdb"
	"github.com/olekukonko/tablewriter"
)

// h/t: https://github.com/hasit/bolter

var helpText = `<CONTROLS>
[:q/CTRL+C to exit] [:b to go back]
[:help for query help] 
[ENTER to return to root bucket]
`

const kvalHelpText = `
<FUNCTIONS>

INS   Insert
GET   Get values
LIS   Check existence
DEL   Delete
REN   Rename

<OPERATORS>

>>    Bucket:Bucket relationship
>>>>  Bucket:Key relationship
::    Key::Value releationship
=>    Name assignment
_     Wildcard

<CAPABILITIES>
{PAT} Given a regex {PAT} for Key XOR Value, find match.

<RESTRICTONS>
Must be >= 1 Buckets for data. 
{PAT} is not a valid option for an INS query.
`

func printRunQuery() {
	fmt.Fprintf(os.Stdout, "\n%s\n\n", color.YellowString("> Enter bucket name"))
}

func printHelpText() {
	fmt.Fprintf(os.Stdout, "\n%s\n", helpText)
}

func printKvalHelpText() {
	color.HiBlue("KVAL (Key Value Access Language) - see full specs at: https://github.com/kval-access-language/kval-language-specification")
	fmt.Fprintf(os.Stdout, "%s\n", kvalHelpText)
	printRunQuery()
}

type manager struct {
	kb         kval.Kvalboltdb
	bucket     string
	currentLoc string
	lastLoc    string
	rootBucket bool
	viewer     formatter
}

type formatter interface {
	DumpBuckets(io.Writer, []bucket)
	DumpBucketItems(io.Writer, string, []item)
}

type item struct {
	Key    string
	Value  string
	Nested bool
}

type box struct {
	Key   string
	Value htb.Box
}

type bucket struct {
	Name string
}

// Audit gives the ability to see what's happening in the DB
func Audit(dbFilepath string) error {
	m := manager{
		viewer: &dbDisplay{},
	}

	if !fs.FileExists(dbFilepath) {
		return errors.New("no database exists")
	}

	m.connect(dbFilepath)
	defer kval.Disconnect(m.kb)

	printHelpText()
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
		case ":q", "\x18":
			color.HiCyan("exiting...")
			return
		case ":b":
			if !strings.Contains(m.currentLoc, "") || !strings.Contains(m.currentLoc, ">>") {
				fmt.Fprintf(os.Stdout, "%s\n", "> going back...")
				m.currentLoc = ""
				m.bucketlist()
			} else {
				m.bucketItems(bucket, true)
			}
		case ":help":
			printKvalHelpText()
		case "":
			m.bucketlist()
		default:
			m.bucketItems(bucket, false)
		}

		bucket = ""
	}
}

func parseBucket(query string) string {
	split := strings.Split(query, " ")
	if len(split) >= 2 {
		return split[1]
	}

	return ""
}

func (m *manager) updateLoc(bucket string, goBack bool) string {
	if bucket == m.lastLoc {
		m.currentLoc = bucket
		return bucket
	}

	if goBack {
		s := strings.Split(m.currentLoc, ">>")
		m.currentLoc = strings.Join(s[:len(s)-1], ">>")
		m.bucket = strings.Trim(s[len(s)-2], " ")
		return m.currentLoc
	}

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
	color.Yellow("ROOT BUCKET")
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

	m.viewer.DumpBuckets(os.Stdout, buckets)
	printRunQuery()
}

func (m *manager) bucketItems(bucket string, goBack bool) {
	items := []item{}
	dbQuery := m.updateLoc(bucket, goBack)
	if dbQuery != "" {
		dbQuery := fmt.Sprintf("GET %s", bucket)
		color.Green("\n[RUNNING]: %s\n", dbQuery)

		res, err := kval.Query(m.kb, dbQuery)
		if err != nil {
			if err.Error() != "Cannot GOTO bucket, bucket not found" {
				log.Fatal(color.RedString(fmt.Sprintf("%v", err)))
			} else {
				fmt.Fprintf(os.Stdout, color.RedString("> Bucket not found\n"))
				if m.rootBucket == true {
					m.bucketlist()
					return
				}
				m.bucketItems(m.currentLoc, true)
			}
		}

		color.HiBlue("# OF RESULTS FOUND: %d", len(res.Result))

		if len(res.Result) == 0 {
			fmt.Fprintf(os.Stdout, color.RedString("No results found\n\n"))
			m.bucketItems(m.lastLoc, false)
			return
		}

		fmt.Fprintf(os.Stdout, fmt.Sprintf("\n%s\nKEYS IN BUCKET:%d\nB+ TREE DEPTH: %d\nINLINE BUCKETS: %d\n\n", color.HiBlueString("STATS"), res.Stats.KeyN, res.Stats.Depth, res.Stats.InlineBucketN))

		for k, v := range res.Result {
			item := item{}
			if v == kval.Nestedbucket {
				color.Red("IS NESTED BUCKET")
				item.Key = strings.TrimSpace(string(k)) + "*"
				item.Value = v
				item.Nested = true
			} else {
				item.Key = strings.TrimSpace(string(k))
				item.Value = strings.TrimSpace(string(v))
			}

			items = append(items, item)
		}

		m.viewer.DumpBucketItems(os.Stdout, m.bucket, items)
		m.rootBucket = false // success this far means we're not at ROOT
		m.lastLoc = dbQuery  // so we can also set the query cache for paging
		printHelpText()
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

func (d dbDisplay) DumpBuckets(w io.Writer, bs []bucket) {
	t := tablewriter.NewWriter(w)
	t.SetHeader([]string{"Buckets"})
	for _, bucket := range bs {
		row := []string{bucket.Name}
		t.Append(row)
	}
	t.SetAutoWrapText(true)
	t.Render()
}

func (d dbDisplay) DumpBucketItems(w io.Writer, bucket string, items []item) {
	color.Yellow("[BUCKET]: %s", bucket)
	t := tablewriter.NewWriter(w)
	t.SetHeader([]string{"Key", "Value"})
	for _, i := range items {
		row := []string{}
		if i.Nested {
			row = append(row, i.Key, "")
		} else {
			reader := bytes.NewReader([]byte(i.Value))
			decoder := gob.NewDecoder(reader)
			var bx htb.Box
			decoder.Decode(&bx)
			row = append(row, i.Key, spew.Sdump(bx))
		}

		t.Append(row)
	}

	t.SetAutoWrapText(true)
	t.Render()
}
