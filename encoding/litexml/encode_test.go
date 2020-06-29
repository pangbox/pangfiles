package litexml

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeSimple(t *testing.T) {
	type Simple struct {
		A string `attr:"a"`
		B string `attr:"b"`
		C string `tag:"nested" attr:"c"`
		D string `content:"inner"`
	}
	type Root struct {
		Info   DocumentInfo
		Simple Simple `tag:"simple"`
	}

	buf := bytes.Buffer{}
	e := NewEncoder(&buf)
	err := e.Encode(Root{
		DocumentInfo{Version: "1.0", Encoding: "utf-8", Standalone: "yes"},
		Simple{
			A: "1234",
			B: "5678",
			C: "C",
			D: "& here is some content",
		}})
	assert.NoError(t, err)

	// TODO: whitespace is currently inaccurate (does not preserve collapse)
	expected := `<?xml version="1.0" encoding="utf-8" standalone="yes" ?>
<simple a="1234" b="5678">
        <nested c="C" />
        &amp; here is some content
</simple>
`
	result := buf.String()
	assert.Equal(t, expected, result)
}

func TestEncodeUpdateList(t *testing.T) {
	type FileInfo struct {
		Filename   string `attr:"fname"`
		Directory  string `attr:"fdir"`
		Size       int    `attr:"fsize"`
		Crc        int    `attr:"fcrc"`
		Date       string `attr:"fdate"`
		Time       string `attr:"ftime"`
		PackedName string `attr:"pname"`
		PackedSize int    `attr:"psize"`
	}

	type UpdateFiles struct {
		Count int        `attr:"count"`
		Files []FileInfo `tag:"fileinfo"`
	}

	type UpdateList struct {
		Info          DocumentInfo
		PatchVer      string      `tag:"patchVer" attr:"value"`
		PatchNum      int         `tag:"patchNum" attr:"value"`
		UpdateListVer string      `tag:"updatelistVer" attr:"value"`
		UpdateFiles   UpdateFiles `tag:"updatefiles"`
	}

	buf := bytes.Buffer{}
	e := NewEncoder(&buf)
	err := e.Encode(UpdateList{
		Info:          DocumentInfo{Version: "1.0", Encoding: "euc-kr", Standalone: "yes"},
		PatchVer:      "KR.Q4.548.00",
		UpdateListVer: "20090331",
		PatchNum:      1,
		UpdateFiles: UpdateFiles{
			Count: 3,
			Files: []FileInfo{
				{
					Filename:   "test01.txt",
					Directory:  "",
					Size:       45,
					Crc:        -90216330,
					Date:       "2020-06-28",
					Time:       "06:01:35",
					PackedName: "test01.txt.zip",
					PackedSize: 154,
				},
				{
					Filename:   "test02.txt",
					Directory:  "",
					Size:       45,
					Crc:        -109573984,
					Date:       "2020-06-28",
					Time:       "06:01:43",
					PackedName: "test02.txt.zip",
					PackedSize: 155,
				},
				{
					Filename:   "test03.txt",
					Directory:  "",
					Size:       45,
					Crc:        -61144858,
					Date:       "2020-06-28",
					Time:       "06:14:13",
					PackedName: "test03.txt.zip",
					PackedSize: 158,
				},
			},
		},
	})
	assert.NoError(t, err)

	expected := `<?xml version="1.0" encoding="euc-kr" standalone="yes" ?>
<patchVer value="KR.Q4.548.00" />
<patchNum value="1" />
<updatelistVer value="20090331" />
<updatefiles count="3">
        <fileinfo fname="test01.txt" fdir="" fsize="45" fcrc="-90216330" fdate="2020-06-28" ftime="06:01:35" pname="test01.txt.zip" psize="154" />
        <fileinfo fname="test02.txt" fdir="" fsize="45" fcrc="-109573984" fdate="2020-06-28" ftime="06:01:43" pname="test02.txt.zip" psize="155" />
        <fileinfo fname="test03.txt" fdir="" fsize="45" fcrc="-61144858" fdate="2020-06-28" ftime="06:14:13" pname="test03.txt.zip" psize="158" />
</updatefiles>
`
	result := buf.String()

	assert.Equal(t, expected, result)
}
