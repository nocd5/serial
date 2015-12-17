package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

// copy & paste from
// http://stackoverflow.com/questions/20365286/query-wmi-from-go
func listComPorts() {
	// init COM
	ole.CoInitialize(0)
	defer ole.CoUninitialize()

	unknown, _ := oleutil.CreateObject("WbemScripting.SWbemLocator")
	defer unknown.Release()

	wmi, _ := unknown.QueryInterface(ole.IID_IDispatch)
	defer wmi.Release()

	// service is a SWbemServices
	serviceRaw, _ := oleutil.CallMethod(wmi, "ConnectServer")
	service := serviceRaw.ToIDispatch()
	defer service.Release()

	// result is a SWBemObjectSet
	resultRaw, _ := oleutil.CallMethod(service, "ExecQuery", "SELECT * FROM WIN32_SerialPort")
	result := resultRaw.ToIDispatch()
	defer result.Release()

	countVar, _ := oleutil.GetProperty(result, "Count")
	count := int(countVar.Val)

	var ports [][]string
	re := regexp.MustCompile(`(.*) \((.*?)\)$`)
	for i := 0; i < count; i++ {
		// item is a SWbemObject, but really a Win32_Process
		itemRaw, _ := oleutil.CallMethod(result, "ItemIndex", i)
		item := itemRaw.ToIDispatch()
		defer item.Release()

		property, _ := oleutil.GetProperty(item, "Name")

		port := re.FindAllStringSubmatch(property.ToString(), -1)
		ports = append(ports, port[0][1:3])
	}

	// " " pading
	maxlen := 0
	for i := 0; i < len(ports); i++ {
		if maxlen < len(ports[i][1]) {
			maxlen = len(ports[i][1])
		}
	}
	for i := 0; i < len(ports); i++ {
		fmt.Println(ports[i][1] + strings.Repeat(" ", maxlen-len(ports[i][1])+1) + ": " + ports[i][0])
	}
}
