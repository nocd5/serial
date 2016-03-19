package main

import (
	"fmt"
	"log"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

func listComPorts() {
	// init COM
	ole.CoInitialize(0)
	defer ole.CoUninitialize()

	locator, err := oleutil.CreateObject("WbemScripting.SWbemLocator")
	if err != nil {
		log.Fatal(err)
	}
	defer locator.Release()

	wmi, err := locator.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		log.Fatal(err)
	}
	defer wmi.Release()

	// service is a SWbemServices
	serviceRaw, err := oleutil.CallMethod(wmi, "ConnectServer")
	if err != nil {
		log.Fatal(err)
	}
	service := serviceRaw.ToIDispatch()
	defer service.Release()

	// result is a SWBemObjectSet
	resultRaw, err := oleutil.CallMethod(service, "ExecQuery", "SELECT * FROM Win32_PnPEntity WHERE Name LIKE '%(COM%)'")
	if err != nil {
		log.Fatal(err)
	}
	result := resultRaw.ToIDispatch()
	defer result.Release()

	countVar, err := oleutil.GetProperty(result, "Count")
	if err != nil {
		log.Fatal(err)
	}
	count := int(countVar.Val)

	var ports Ports
	re := regexp.MustCompile(`(.*) \(COM(\d+)\)$`)
	for i := 0; i < count; i++ {
		// item is a SWbemObject, but really a Win32_Process
		itemRaw, err := oleutil.CallMethod(result, "ItemIndex", i)
		if err != nil {
			continue
		}
		item := itemRaw.ToIDispatch()
		defer item.Release()

		property, err := oleutil.GetProperty(item, "Name")
		if err != nil {
			continue
		}
		port := re.FindAllStringSubmatch(property.ToString(), -1)
		num, _ := strconv.Atoi(port[0][2])
		ports = append(ports, Port{num, port[0][1]})
	}

	sort.Sort(ports)

	digit := getDigit(ports[len(ports)-1].Number)
	for i := 0; i < len(ports); i++ {
		fmt.Println("COM" + strconv.Itoa(ports[i].Number) + strings.Repeat(" ", digit-getDigit(ports[i].Number)) + " : " + ports[i].Name)
	}
}

func getDigit(num int) int {
	return int(math.Log10(float64(num))) + 1
}

type Port struct {
	Number int
	Name   string
}

type Ports []Port

func (p Ports) Len() int {
	return len(p)
}

func (p Ports) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p Ports) Less(i, j int) bool {
	return p[i].Number < p[j].Number
}
