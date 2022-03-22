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
	"time"
)

type DnsRequest struct {
	Domain   string
	Received time.Time
}

type HttpRequest struct {
	Protocol string
	URI      string
	Received time.Time
}

type AppState struct {
	dnsRequests  []DnsRequest
	httpRequests []HttpRequest
	dnsList      *widget.List
	httpList     *widget.List
}

func (state *AppState) AddDnsRequest(domain string) {
	newItem := DnsRequest{Domain: domain, Received: time.Now()}
	state.dnsRequests = append([]DnsRequest{newItem}, state.dnsRequests...)
	state.dnsList.Refresh()
}

func (state *AppState) AddHttpRequest(protocol string, uri string) {
	newItem := HttpRequest{Protocol: protocol, URI: uri, Received: time.Now()}
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

func generateHttpRequestLayout(state *AppState) *fyne.Container {
	title := widget.NewLabel("HTTP(s) Requests")

	list := widget.NewList(
		func() int {
			return len(state.httpRequests)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(fmt.Sprintf("[%s] %s: %s",
				state.httpRequests[i].Protocol,
				state.httpRequests[i].Received.Format(time.RFC3339),
				state.httpRequests[i].URI))
		})

	state.httpList = list
	return container.NewBorder(title, nil, nil, nil, list)
}

func main() {
	appState := AppState{
		dnsRequests:  []DnsRequest{},
		httpRequests: []HttpRequest{},
	}

	myApp := app.New()
	myWindow := myApp.NewWindow("Network Simulator")

	content := container.New(
		layout.NewGridLayout(2),
		generateDnsRequestLayout(&appState),
		generateHttpRequestLayout(&appState))

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(1024, 500))

	dns.StartDNSServer(53, appState.AddDnsRequest)
	web.StartWebServers(80, 443, appState.AddHttpRequest)

	myWindow.ShowAndRun()
}
