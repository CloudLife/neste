package neste

import (
	. "launchpad.net/gocheck"
	"testing"
	"bytes"
	"os"
	"template"
	"path"
)

// Hook up gocheck into the gotest runner.
func Test(t *testing.T) { TestingT(t) }

type S struct{}

func init() {
	Suite(&S{})
}

const (
	// Template base dir
	baseDir = "templates"

	// Template file names
	indexName   = "index.html"
	headName    = "head.html"
	brandName   = "brand.html"
	contentName = "content.html"
	listName    = "list.html"
	footerName  = "footer.html"
)

func (s *S) TestAdd(c *C) {
	tm := New(baseDir, nil)

	t, err := tm.AddFile(indexName)
	c.Assert(err, IsNil)
	c.Assert(len(tm.tFiles), Equals, 1)
	indexPath := path.Join(baseDir, indexName)
	c.Check(tm.tFiles[indexPath].cache, NotNil)
	c.Check(tm.GetFile(indexName), Equals, t)
}

func (s *S) TestRemoveFile(c *C) {
	tm := New(baseDir, nil)
	tm.AddFile(indexName)
	tm.RemoveFile(indexName)

	c.Assert(len(tm.tFiles), Equals, 0)
}

func (s *S) TestClear(c *C) {
	tm := New(baseDir, nil)

	// First try with no templates.
	deleted := tm.Clear()
	c.Check(deleted, Equals, false)

	// Now add something.
	tm.AddFile(indexName)
	tm.AddFile(contentName)
	tm.AddFile(footerName)

	deleted = tm.Clear()
	c.Check(deleted, Equals, true)
	c.Assert(len(tm.tStrings), Equals, 0)
	c.Assert(len(tm.tFiles), Equals, 0)
}

func (s *S) TestExecute(c *C) {
	var data = map[string]string{
		"head":    "<title>Execute</title>",
		"content": "Execute",
		"footer":  "neste template engine"}

	expected :=
`<!DOCTYPE HTML>
<html>
<head><title>Execute</title></head>
<body>
Execute<hr/>neste template engine
</body>
</html>
`

	tm := New(baseDir, nil)
	t, _ := tm.AddFilenl(indexName)

	buf := new(bytes.Buffer)
	err := t.Execute(buf, data)
	c.Assert(err, IsNil)
	c.Assert(string(buf.Bytes()), Equals, expected)
}

func (s *S) TestRender(c *C) {
	var data = map[string]string{
		"head":    "<title>Render</title>",
		"content": "Render",
		"footer":  "neste template engine"}

	expected :=
`<!DOCTYPE HTML>
<html>
<head><title>Render</title></head>
<body>
Render<hr/>neste template engine
</body>
</html>
`

	tm := New(baseDir, nil)
	t, _ := tm.AddFilenl(indexName)

	output, err := t.Render(data)
	c.Assert(err, IsNil)
	c.Assert(output, Equals, expected)
}

func (s *S) TestFormatters(c *C) {
	tstr :=
`
{unesc1 unesc2 unesc3|html}
{unslashed|addSlashes}
{uncapped|capFirst}
{uncapped2|capFirst}
`
	var data = map[string]string{
		"unesc1":    "<hack>",
		"unesc2":    "\\&hack\\",
		"unesc3":    "</hack>",
		"unslashed": `"I'm using neste"`,
		"uncapped":  "neste",
		"uncapped2": "ǿxy"}

	expected :=
`
&lt;hack&gt;\&amp;hack\&lt;/hack&gt;
\"I'm using neste\"
Neste
Ǿxy
`

	tm := New(baseDir, template.FormatterMap{})
	t := tm.MustAdd(tstr, "testFormatters")

	output, err := t.Render(data)
	c.Assert(err, IsNil)
	c.Assert(output, Equals, expected)
}

func (s *S) TestNesting(c *C) {
	expected :=
`<!DOCTYPE HTML>
<html>
<head><title>Page Title</title></head>
<body>
<div id="brand">neste template engine</div><div id="content">
<h1>Page Title</h1>
<p>Example page to demonstrate nested templates.</p>
<ul>
<li>Example</li>
<li>Listing</li>
<li>Area</li>
</ul>
</div><hr/><div id="footer">
Posted : 25th July 2010 12:15
</div>
</body>
</html>
`
	var err os.Error
	var indexData = map[string]string{}
	var headData = map[string]string{"title": "Page Title"}
	var brandData = map[string]string{}
	var footerData = map[string]string{"posted": "25th July 2010 12:15"}
	var listData = map[string]interface{}{"items": &[3]string{"Example", "Listing", "Area"}}
	var contentData = map[string]string{
		"title":   "Page Title",
		"opening": "Example page to demonstrate nested templates."}

	tm := New(baseDir, nil)

	tIndex := tm.MustAddFilenl(indexName)
	tHead := tm.MustAddFile(headName)
	tBrand := tm.MustAddFile(brandName)
	tContent := tm.MustAddFile(contentName)
	tList := tm.MustAddFile(listName)
	tFooter := tm.MustAddFile(footerName)

	contentData["list"], err = tList.Render(listData)
	c.Assert(err, IsNil)

	indexData["head"], err = tHead.Render(headData)
	c.Assert(err, IsNil)

	indexData["brand"], err = tBrand.Render(brandData)
	c.Assert(err, IsNil)

	indexData["content"], err = tContent.Render(contentData)
	c.Assert(err, IsNil)

	indexData["footer"], err = tFooter.Render(footerData)
	c.Assert(err, IsNil)

	output, err := tIndex.Render(indexData)
	c.Assert(err, IsNil)
	c.Assert(output, Equals, expected)
}

