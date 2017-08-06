# systray-portable
A portable version of [go systray](https://github.com/getlantern/systray), using stdin/stdout to communicate with other language


## Protocol

Each line is a json string.

tray binary =>  
=> ready  `{"type": "ready"}`  
<= init menu
```json
{
  "icon": "<base64 string of image>",
  "title": "Title",
  "tooltip": "Tooltips",
  "items":[{
    "title": "aa",
    "tooltip":"bb",
    "checked": true,
    "enabled": true
  }, {
    "title": "aa2",
    "tooltip":"bb",
    "checked": false,
    "enabled": true
  }]}
```
=> clicked  
```json
{
  "type":"clicked",
  "item":{"title":"aa","tooltip":"bb","enabled":true,"checked":true},
  "menu":{"icon":"","title":"","tooltip":"","items":null},
  "seq_id":0
}
```
<= update-item / update-menu / update-item-and-menu
```json
{
  "type": "update-item",
  "item": {"title":"aa3","tooltip":"bb","enabled":true,"checked":true},
  "seq_id": 0
}
```

## Binary
main_xxx_release: `go build -ldflags "-s -w" tray.go`  
main_xxx: `go build tray.go`
