package sub

import (
	"encoding/base64"
	"reflect"
	"testing"

	"github.com/mhsanaei/3x-ui/v3/internal/database/model"
)

func TestMihomoECHOpts(t *testing.T) {
	config := base64.StdEncoding.EncodeToString([]byte("test-ECHConfigList"))
	for _, tc := range []struct {
		name  string
		input string
		want  map[string]any
	}{
		{"fixed", config, map[string]any{"enable": true, "config": config}},
		{"dynamic host", "cloudflare-ech.com+https://dns.alidns.com/dns-query", map[string]any{"enable": true, "query-server-name": "cloudflare-ech.com"}},
		{"dynamic DoH only", "https://dns.alidns.com/dns-query", map[string]any{"enable": true}},
		{"empty", "", nil},
		{"invalid", "not-an-ech-config", nil},
		{"bad dynamic", "cloudflare-ech.com+https://", nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := mihomoECHOpts(tc.input)
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("mihomoECHOpts(%q) = %#v, want %#v", tc.input, got, tc.want)
			}
		})
	}
}

func TestClashECHInheritedAndPerHostOverride(t *testing.T) {
	fixed := base64.StdEncoding.EncodeToString([]byte("server-ECHConfigList"))
	streamJSON := `{"network":"ws","security":"tls","wsSettings":{"path":"/ws"},"tlsSettings":{"serverName":"direct.example.test","alpn":["http/1.1"],"settings":{"fingerprint":"chrome","echConfigList":"` + fixed + `"}}}`
	svc := &SubClashService{SubService: &SubService{}}
	base := svc.streamData(streamJSON)
	inbound := &model.Inbound{
		Listen:   "direct.example.test",
		Port:     2083,
		Protocol: model.VLESS,
		Settings: `{"encryption":"none"}`,
	}
	client := model.Client{ID: "11111111-2222-4333-8444-555555555555"}
	hostFixed := base64.StdEncoding.EncodeToString([]byte("host-ECHConfigList"))

	for _, tc := range []struct {
		name string
		ep   map[string]any
		want map[string]any
	}{
		{"inherited direct ECH", nil, map[string]any{"enable": true, "config": fixed}},
		{"per-host Cloudflare ECH", map[string]any{"echConfigList": "cloudflare-ech.com+https://dns.alidns.com/dns-query"}, map[string]any{"enable": true, "query-server-name": "cloudflare-ech.com"}},
		{"per-host fixed ECH", map[string]any{"echConfigList": hostFixed}, map[string]any{"enable": true, "config": hostFixed}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			stream := cloneStreamForExternalProxy(base)
			if tc.ep != nil {
				applyExternalProxyTLSToStream(tc.ep, stream, "tls")
			}
			got := svc.buildProxy(svc.SubService, inbound, client, stream, tc.ep)
			if got == nil {
				t.Fatal("VLESS WS proxy was filtered")
			}
			if !reflect.DeepEqual(got["ech-opts"], tc.want) {
				t.Errorf("ech-opts = %#v, want %#v", got["ech-opts"], tc.want)
			}
			if got["network"] != "ws" || got["type"] != "vless" {
				t.Errorf("unexpected proxy: %#v", got)
			}
		})
	}
}

func TestClashECHNotEmittedWithoutConfiguration(t *testing.T) {
	svc := &SubClashService{}
	proxy := map[string]any{"type": "vless"}
	if !svc.applySecurity(proxy, "tls", map[string]any{"tlsSettings": map[string]any{"serverName": "plain.example.test"}}) {
		t.Fatal("TLS proxy was rejected")
	}
	if _, ok := proxy["ech-opts"]; ok {
		t.Fatalf("unexpected ECH: %#v", proxy)
	}
}
