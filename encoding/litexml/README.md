# LiteXML
This library implements a variant of XML that is more lenient and less formal than typical XML, allowing multiple root elements and only supporting the basic functionality of XML. It does not require escaping in places where XML would; for example, bare ampersands are allowed in string literals.

This format is based on the one used by update lists.

## Usage
LiteXML structs are normal Go structs with some tags. Unlike many other Go encoding libraries, this library is not designed to work with arbitrary structures, but only ones that are specifically meant for LiteXML usage. These are called LiteXML structs.

LiteXML structures support the tags `tag`, `attr`, and `content`. Here is the structure of the PangYa UpdateList as LiteXML structs:

```go
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
```

...which will output an XML document like this:

```xml
<?xml version="1.0" encoding="utf-8" standalone="yes" ?>
<patchVer value="RG.R7.920.00" />
<patchNum value="1" />
<updatelistVer value="20090331" />
<updatefiles count="1">
        <fileinfo fname="Game.exe" fdir="" fsize="946321" fcrc="-83219464" fdate="2020-06-28" ftime="06:01:35" pname="Game.exe.zip" psize="638612" />
</updatefiles>
```

## Shortcomings
For now, the library is using UTF-8 encoding only. The actual game uses EUC-KR. This is OK for UpdateList because it is ASCII-safe, but needs to be reconciled sooner or later.

LiteXML structures are very rigid, but the parsing routine is designed to be somewhat lenient. Therefore, the decoding process is highly lossy. Notably, decoding does not support ordering at all. If you have multiple distinct sections of content in a block, it will be impossible to decode. (Worse, right now, it will silently decode very incorrectly.) If you don't care about content or only care about very simple content, this should be OK.

This library significantly lacks testing and safeguards. This is partly by nature due to its limited use, but it would be nice to fix this over time anyways.