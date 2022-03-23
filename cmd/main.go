package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"gonetsim/internal/dns"
	"gonetsim/internal/web"
	"log"
	"strconv"
	"time"
)

type DnsRequest struct {
	Domain   string
	Received time.Time
}

type HttpRequest struct {
	URI      string
	Received time.Time
	Raw      string
}

type AppState struct {
	dnsRequests  []DnsRequest
	httpRequests []HttpRequest
	dnsList      *widget.List
	httpList     *widget.List
	dnsServer    dns.Server
}

func (state *AppState) AddDnsRequest(domain string) {
	newItem := DnsRequest{Domain: domain, Received: time.Now()}
	state.dnsRequests = append([]DnsRequest{newItem}, state.dnsRequests...)
	state.dnsList.Refresh()
}

func (state *AppState) AddHttpRequest(uri string, raw string) {
	newItem := HttpRequest{URI: uri, Raw: raw, Received: time.Now()}
	state.httpRequests = append([]HttpRequest{newItem}, state.httpRequests...)
	state.httpList.Refresh()
}

func generateDnsRequestLayout(state *AppState) *fyne.Container {
	title := widget.NewLabel("DNS Requests")

	list := widget.NewList(
		func() int {
			return len(state.dnsRequests)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(fmt.Sprintf("%s: %s",
				state.dnsRequests[i].Received.Format(time.RFC3339),
				state.dnsRequests[i].Domain))
		})

	state.dnsList = list
	return container.NewBorder(title, nil, nil, nil, list)
}

func generateHttpRequestLayout(state *AppState, window fyne.Window) *fyne.Container {
	title := widget.NewLabel("HTTP(s) Requests")

	list := widget.NewList(
		func() int {
			return len(state.httpRequests)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(fmt.Sprintf("%s: %s",
				state.httpRequests[i].Received.Format(time.RFC3339),
				state.httpRequests[i].URI))
		})
	list.OnSelected = func(id widget.ListItemID) {
		var modal *widget.PopUp
		modal = widget.NewModalPopUp(
			container.NewVBox(
				widget.NewLabel(state.httpRequests[id].Raw),
				widget.NewButton("Copy", func() {
					window.Clipboard().SetContent(state.httpRequests[id].Raw)
				}),
				widget.NewButton("Close", func() {
					modal.Hide()
					list.Unselect(id)
				}),
			),
			window.Canvas(),
		)
		modal.Show()
	}
	state.httpList = list
	return container.NewBorder(title, nil, nil, nil, list)
}

func main() {
	appState := AppState{
		dnsRequests:  []DnsRequest{},
		httpRequests: []HttpRequest{},
	}
	dnsServer := dns.Server{
		Port: 53,
	}

	webServer := web.Server{
		HttpPort:  80,
		HttpsPort: 443,
	}

	myApp := app.New()
	myWindow := myApp.NewWindow("Network Simulator")

	httpResponse := widget.NewMultiLineEntry()
	httpResponse.SetPlaceHolder("Put your response body here")
	httpResponse.SetText("@_@")

	statusCodeSelect := widget.NewSelect([]string{"200", "400", "404", "500"}, func(value string) {
		log.Println("Status code set to", value)
	})
	statusCodeSelect.SetSelectedIndex(0)

	contentTypeEntry := widget.NewEntry()
	contentTypeEntry.Wrapping = fyne.TextWrapOff
	contentTypeEntry.SetText("text/html")

	dnsPortEntry := widget.NewEntry()
	dnsPortEntry.Wrapping = fyne.TextWrapOff
	dnsPortEntry.SetText("5354")

	httpPortEntry := widget.NewEntry()
	httpPortEntry.Wrapping = fyne.TextWrapOff
	httpPortEntry.SetText("8080")

	httpsPortEntry := widget.NewEntry()
	httpsPortEntry.Wrapping = fyne.TextWrapOff
	httpsPortEntry.SetText("8443")

	content := container.New(
		layout.NewGridLayoutWithRows(2),
		container.New(
			layout.NewGridLayout(2),
			container.New(
				layout.NewVBoxLayout(),
				widget.NewLabel("Setup"),
				container.New(
					layout.NewFormLayout(),
					widget.NewLabel("DNS Server"),
					dnsPortEntry,
					widget.NewLabel("HTTP Server"),
					httpPortEntry,
					widget.NewLabel("HTTPs Server"),
					httpsPortEntry,
				),
				container.New(
					layout.NewHBoxLayout(),
					widget.NewButton("Start", func() {
						dnsPort, err := strconv.Atoi(dnsPortEntry.Text)
						if err == nil {
							dnsServer.Port = dnsPort
						}
						httpPort, err := strconv.Atoi(httpPortEntry.Text)
						if err == nil {
							webServer.HttpPort = httpPort
						}
						httpsPort, err := strconv.Atoi(httpsPortEntry.Text)
						if err == nil {
							webServer.HttpsPort = httpsPort
						}
						dnsServer.Start(appState.AddDnsRequest)
						webServer.Start(appState.AddHttpRequest)
					}),
					widget.NewButton("Stop", func() {
						dnsServer.Stop()
						webServer.Stop()
					}),
					layout.NewSpacer(),
				),
			),
			container.New(
				layout.NewVBoxLayout(),
				widget.NewLabel("Default Response"),
				container.New(
					layout.NewFormLayout(),
					widget.NewLabel("Status Code"),
					statusCodeSelect,
					widget.NewLabel("Content-Type"),
					contentTypeEntry,
				),
				widget.NewLabel("Body"),
				httpResponse,
			),
		),
		container.New(
			layout.NewGridLayout(2),
			generateDnsRequestLayout(&appState),
			generateHttpRequestLayout(&appState, myWindow),
		))

	myWindow.SetContent(content)
	//myWindow.Resize(fyne.NewSize(1024, 500))

	myWindow.ShowAndRun()
}
