package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/getlantern/systray"
	"os"
	"reflect"
)

func main() {
	// Should be called at the very beginning of main().
	systray.Run(onReady)
}

type Item struct {
	Title   string `json:"title"`
	Tooltip string `json:"tooltip"`
	Enabled bool   `json:"enabled"`
	Checked bool   `json:"checked"`
}
type Menu struct {
	Icon    string `json:"icon"`
	Title   string `json:"title"`
	Tooltip string `json:"tooltip"`
	Items   []Item `json:"items"`
}

type Action struct {
	Type  string `json:"type"`
	Item  Item   `json:"item"`
	Menu  Menu   `json:"menu"`
	SeqId int    `json:"seq_id"`
}

func readLine(reader *bufio.Reader) string {
	input, err := reader.ReadBytes('\n')
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	return string(input[0 : len(input)-1])
}
func onReady() {
	// We can manipulate the systray in other goroutines
	go func() {
		items := make([]*systray.MenuItem, 0)
		// fmt.Println(items)
		fmt.Println(`{"type": "ready"}`)
		reader := bufio.NewReader(os.Stdin)
		// println(readLine(reader))
		var menu Menu
		json.Unmarshal([]byte(readLine(reader)), &menu)
		// fmt.Println("menu", menu)
		icon, err := base64.StdEncoding.DecodeString(menu.Icon)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		systray.SetIcon(icon)
		systray.SetTitle(menu.Title)
		systray.SetTooltip(menu.Tooltip)
		for i := 0; i < len(menu.Items); i++ {
			item := menu.Items[i]
			menuItem := systray.AddMenuItem(item.Title, item.Tooltip)
			if item.Checked {
				menuItem.Check()
			} else {
				menuItem.Uncheck()
			}
			if item.Enabled {
				menuItem.Enable()
			} else {
				menuItem.Disable()
			}
			items = append(items, menuItem)
		}
		// {"type": "update-item", "item": {"Title":"aa3","Tooltip":"bb","Enabled":true,"Checked":true}, "seqId": 0}
		for {
			cases := make([]reflect.SelectCase, len(items))
			for i, ch := range items {
				cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch.ClickedCh)}
			}

			remaining := len(cases)
			for remaining > 0 {
				chosen, _, ok := reflect.Select(cases)
				if !ok {
					// The chosen channel has been closed, so zero out the channel to disable the case
					cases[chosen].Chan = reflect.ValueOf(nil)
					remaining -= 1
					continue
				}
				menuItem := items[chosen]
				data, err := json.Marshal(Action{
					Type:  "clicked",
					Item:  menu.Items[chosen],
					SeqId: chosen,
				})
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
				}
				fmt.Println(string(data))
				var replyAction Action
				json.Unmarshal([]byte(readLine(reader)), &replyAction)
				// fmt.Println(replyAction)
				updateItem := func(action Action) {
					item := action.Item
					menu.Items[chosen] = item
					if item.Checked {
						menuItem.Check()
					} else {
						menuItem.Uncheck()
					}
					if item.Enabled {
						menuItem.Enable()
					} else {
						menuItem.Disable()
					}
					menuItem.SetTitle(item.Title)
					menuItem.SetTooltip(item.Tooltip)
					// fmt.Println("Done")
					// fmt.Printf("Read from channel %#v and received %s\n", items[chosen], value.String())
				}
				updateMenu := func(action Action) {
					m := action.Menu
					if menu.Title != m.Title {
						menu.Title = m.Title
						systray.SetTitle(menu.Title)
					}
					if menu.Icon != m.Icon {
						menu.Icon = m.Icon
						icon, err := base64.StdEncoding.DecodeString(menu.Icon)
						if err != nil {
							fmt.Fprintln(os.Stderr, err)
						}
						systray.SetIcon(icon)
					}
					if menu.Tooltip != m.Tooltip {
						menu.Tooltip = m.Tooltip
						systray.SetTooltip(menu.Tooltip)
					}
				}
				switch replyAction.Type {
				case "update-item":
					updateItem(replyAction)
				case "update-menu":
					updateMenu(replyAction)
				case "update-item-and-menu":
					updateItem(replyAction)
					updateMenu(replyAction)
				}
			}
		}
	}()
}
