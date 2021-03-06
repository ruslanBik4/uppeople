// Code generated by qtc from "td_img.qtpl". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

// All the text outside function templates is treated as comments,
// i.e. it is just ignored by quicktemplate compiler (`qtc`). It is for humans.
//
// .

//line views/td/td_img.qtpl:5
package td

//line views/td/td_img.qtpl:5
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line views/td/td_img.qtpl:5
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line views/td/td_img.qtpl:5
func StreamTdImgStyle(qw422016 *qt422016.Writer, tableName, colName string, id int) {
//line views/td/td_img.qtpl:5
	qw422016.N().S(`
<style>
.hiddenInput`)
//line views/td/td_img.qtpl:7
	qw422016.E().S(colName)
//line views/td/td_img.qtpl:7
	qw422016.N().D(id)
//line views/td/td_img.qtpl:7
	qw422016.N().S(` > input[type=file] {
    width: 100%;
    height: 100%;
    opacity: 0;
    cursor: pointer;
}
form > input[type=hidden] + span + img{
    left:-10000px;
}
form > input[type=hidden] + span:hover + img{
  position: fixed;
    top: 1%;
    left: 1%;
    right: 1%;
    max-width: 25%;
    max-height: 75%;
}

.hiddenInput`)
//line views/td/td_img.qtpl:25
	qw422016.E().S(colName)
//line views/td/td_img.qtpl:25
	qw422016.N().D(id)
//line views/td/td_img.qtpl:25
	qw422016.N().S(` {
    border: 1px solid #ccc;
    width: 100%;
    height: 100%;
    max-height: 100px;
    display: inline-block;
    overflow: hidden;
    cursor: pointer;
    background: center center no-repeat scroll;
    background-size: contain;
}
.hiddenInputImg`)
//line views/td/td_img.qtpl:36
	qw422016.E().S(colName)
//line views/td/td_img.qtpl:36
	qw422016.N().D(id)
//line views/td/td_img.qtpl:36
	qw422016.N().S(` {
    background-image: url("/api/v1/blob/`)
//line views/td/td_img.qtpl:37
	qw422016.E().S(tableName)
//line views/td/td_img.qtpl:37
	qw422016.N().S(`?id=`)
//line views/td/td_img.qtpl:37
	qw422016.N().D(id)
//line views/td/td_img.qtpl:37
	qw422016.N().S(`&name=`)
//line views/td/td_img.qtpl:37
	qw422016.E().S(colName)
//line views/td/td_img.qtpl:37
	qw422016.N().S(`");
}
.hiddenInputImgNew`)
//line views/td/td_img.qtpl:39
	qw422016.E().S(colName)
//line views/td/td_img.qtpl:39
	qw422016.N().D(id)
//line views/td/td_img.qtpl:39
	qw422016.N().S(` {
    background-image: url("/api/v1/blob/`)
//line views/td/td_img.qtpl:40
	qw422016.E().S(tableName)
//line views/td/td_img.qtpl:40
	qw422016.N().S(`?id=`)
//line views/td/td_img.qtpl:40
	qw422016.N().D(id)
//line views/td/td_img.qtpl:40
	qw422016.N().S(`&name=`)
//line views/td/td_img.qtpl:40
	qw422016.E().S(colName)
//line views/td/td_img.qtpl:40
	qw422016.N().S(`&v=true");
}
</style>
`)
//line views/td/td_img.qtpl:43
}

//line views/td/td_img.qtpl:43
func WriteTdImgStyle(qq422016 qtio422016.Writer, tableName, colName string, id int) {
//line views/td/td_img.qtpl:43
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/td/td_img.qtpl:43
	StreamTdImgStyle(qw422016, tableName, colName, id)
//line views/td/td_img.qtpl:43
	qt422016.ReleaseWriter(qw422016)
//line views/td/td_img.qtpl:43
}

//line views/td/td_img.qtpl:43
func TdImgStyle(tableName, colName string, id int) string {
//line views/td/td_img.qtpl:43
	qb422016 := qt422016.AcquireByteBuffer()
//line views/td/td_img.qtpl:43
	WriteTdImgStyle(qb422016, tableName, colName, id)
//line views/td/td_img.qtpl:43
	qs422016 := string(qb422016.B)
//line views/td/td_img.qtpl:43
	qt422016.ReleaseByteBuffer(qb422016)
//line views/td/td_img.qtpl:43
	return qs422016
//line views/td/td_img.qtpl:43
}

//line views/td/td_img.qtpl:44
func StreamTdImgBrowse(qw422016 *qt422016.Writer, preRoute, tableName, colName string, id int) {
//line views/td/td_img.qtpl:44
	qw422016.N().S(`
  `)
//line views/td/td_img.qtpl:45
	StreamTdImgStyle(qw422016, tableName, colName, id)
//line views/td/td_img.qtpl:45
	qw422016.N().S(`
    <form id="`)
//line views/td/td_img.qtpl:46
	qw422016.E().S(tableName)
//line views/td/td_img.qtpl:46
	qw422016.E().S(colName)
//line views/td/td_img.qtpl:46
	qw422016.N().D(id)
//line views/td/td_img.qtpl:46
	qw422016.N().S(`" method="POST" action="`)
//line views/td/td_img.qtpl:46
	qw422016.E().S(preRoute)
//line views/td/td_img.qtpl:46
	qw422016.E().S(tableName)
//line views/td/td_img.qtpl:46
	qw422016.N().S(`/update"
                enctype="multipart/form-data" style="height: 100%;">
    <input type="hidden" value="`)
//line views/td/td_img.qtpl:48
	qw422016.N().D(id)
//line views/td/td_img.qtpl:48
	qw422016.N().S(`" name="id" />
    <span class="hiddenInput`)
//line views/td/td_img.qtpl:49
	qw422016.E().S(colName)
//line views/td/td_img.qtpl:49
	qw422016.N().D(id)
//line views/td/td_img.qtpl:49
	qw422016.N().S(` hiddenInputImg`)
//line views/td/td_img.qtpl:49
	qw422016.E().S(colName)
//line views/td/td_img.qtpl:49
	qw422016.N().D(id)
//line views/td/td_img.qtpl:49
	qw422016.N().S(`">
        <input name="`)
//line views/td/td_img.qtpl:50
	qw422016.E().S(colName)
//line views/td/td_img.qtpl:50
	qw422016.N().S(`" type="file" accept="image"
        onchange="return saveForm(this.parentElement.parentElement, function(d, t_form) {
        $('.hiddenInputImg`)
//line views/td/td_img.qtpl:52
	qw422016.E().S(colName)
//line views/td/td_img.qtpl:52
	qw422016.N().D(id)
//line views/td/td_img.qtpl:52
	qw422016.N().S(`').removeClass('hiddenInputImg`)
//line views/td/td_img.qtpl:52
	qw422016.E().S(colName)
//line views/td/td_img.qtpl:52
	qw422016.N().D(id)
//line views/td/td_img.qtpl:52
	qw422016.N().S(`').addClass('hiddenInputImgNew`)
//line views/td/td_img.qtpl:52
	qw422016.E().S(colName)
//line views/td/td_img.qtpl:52
	qw422016.N().D(id)
//line views/td/td_img.qtpl:52
	qw422016.N().S(`');  });" >
    </span>
    <img src="/api/v1/blob/`)
//line views/td/td_img.qtpl:54
	qw422016.E().S(tableName)
//line views/td/td_img.qtpl:54
	qw422016.N().S(`?id=`)
//line views/td/td_img.qtpl:54
	qw422016.N().D(id)
//line views/td/td_img.qtpl:54
	qw422016.N().S(`&name=`)
//line views/td/td_img.qtpl:54
	qw422016.E().S(colName)
//line views/td/td_img.qtpl:54
	qw422016.N().S(`"/>
    </form>
`)
//line views/td/td_img.qtpl:56
}

//line views/td/td_img.qtpl:56
func WriteTdImgBrowse(qq422016 qtio422016.Writer, preRoute, tableName, colName string, id int) {
//line views/td/td_img.qtpl:56
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/td/td_img.qtpl:56
	StreamTdImgBrowse(qw422016, preRoute, tableName, colName, id)
//line views/td/td_img.qtpl:56
	qt422016.ReleaseWriter(qw422016)
//line views/td/td_img.qtpl:56
}

//line views/td/td_img.qtpl:56
func TdImgBrowse(preRoute, tableName, colName string, id int) string {
//line views/td/td_img.qtpl:56
	qb422016 := qt422016.AcquireByteBuffer()
//line views/td/td_img.qtpl:56
	WriteTdImgBrowse(qb422016, preRoute, tableName, colName, id)
//line views/td/td_img.qtpl:56
	qs422016 := string(qb422016.B)
//line views/td/td_img.qtpl:56
	qt422016.ReleaseByteBuffer(qb422016)
//line views/td/td_img.qtpl:56
	return qs422016
//line views/td/td_img.qtpl:56
}

//line views/td/td_img.qtpl:58
func StreamTdImgView(qw422016 *qt422016.Writer, preRoute, tableName, colName string, id int) {
//line views/td/td_img.qtpl:58
	qw422016.N().S(`
  `)
//line views/td/td_img.qtpl:59
	StreamTdImgStyle(qw422016, tableName, colName, id)
//line views/td/td_img.qtpl:59
	qw422016.N().S(`
    <img src="/api/v1/blob/`)
//line views/td/td_img.qtpl:60
	qw422016.E().S(tableName)
//line views/td/td_img.qtpl:60
	qw422016.N().S(`?id=`)
//line views/td/td_img.qtpl:60
	qw422016.N().D(id)
//line views/td/td_img.qtpl:60
	qw422016.N().S(`&name=`)
//line views/td/td_img.qtpl:60
	qw422016.E().S(colName)
//line views/td/td_img.qtpl:60
	qw422016.N().S(`" class="hiddenInput`)
//line views/td/td_img.qtpl:60
	qw422016.E().S(colName)
//line views/td/td_img.qtpl:60
	qw422016.N().D(id)
//line views/td/td_img.qtpl:60
	qw422016.N().S(`"/>
   </span>
`)
//line views/td/td_img.qtpl:62
}

//line views/td/td_img.qtpl:62
func WriteTdImgView(qq422016 qtio422016.Writer, preRoute, tableName, colName string, id int) {
//line views/td/td_img.qtpl:62
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/td/td_img.qtpl:62
	StreamTdImgView(qw422016, preRoute, tableName, colName, id)
//line views/td/td_img.qtpl:62
	qt422016.ReleaseWriter(qw422016)
//line views/td/td_img.qtpl:62
}

//line views/td/td_img.qtpl:62
func TdImgView(preRoute, tableName, colName string, id int) string {
//line views/td/td_img.qtpl:62
	qb422016 := qt422016.AcquireByteBuffer()
//line views/td/td_img.qtpl:62
	WriteTdImgView(qb422016, preRoute, tableName, colName, id)
//line views/td/td_img.qtpl:62
	qs422016 := string(qb422016.B)
//line views/td/td_img.qtpl:62
	qt422016.ReleaseByteBuffer(qb422016)
//line views/td/td_img.qtpl:62
	return qs422016
//line views/td/td_img.qtpl:62
}
