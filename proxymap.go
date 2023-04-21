package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"path"
	"regexp"
	"strings"

	"proxymap/pkg/config"
	"proxymap/pkg/db"
	"proxymap/pkg/filter"

	"github.com/lqqyt2423/go-mitmproxy/proxy"
	log "github.com/sirupsen/logrus"
)

var titleRegexp = regexp.MustCompile("<title>(.+?)</title>")
var sdb db.DB

type SaveResponse struct {
	proxy.BaseAddon
}

func (c *SaveResponse) Response(f *proxy.Flow) {
	var storeinfo db.StoreInfo
	storeinfo.Url = f.Request.URL.String()
	storeinfo.Host = f.Request.URL.Hostname()
	storeinfo.Port = f.Request.URL.Port()
	storeinfo.Method = f.Request.Method
	storeinfo.ReqBody = string(f.Request.Body)
	storeinfo.ReqHeaders = f.Request.Header
	storeinfo.Scheme = f.Request.URL.Scheme

	storeinfo.StatusCode = f.Response.StatusCode
	DecodedBody, err := f.Response.DecodedBody()
	if err != nil {
		fmt.Println(err)
	}
	storeinfo.RespBody = string(DecodedBody)
	match := titleRegexp.FindStringSubmatch(storeinfo.RespBody)
	if len(match) > 1 {
		storeinfo.Title = match[1]
	}
	storeinfo.Length = len(storeinfo.RespBody)
	storeinfo.RespHeaders = f.Response.Header
	storeinfo.ContentType = strings.Split(f.Response.Header.Get("Content-Type"), ";")[0]
	storeinfo.Extension = path.Ext(f.Request.URL.Path)
	if filter.BList.IsPass(&storeinfo) {
		docBytes, err := json.Marshal(storeinfo)
		if err != nil {
			fmt.Println(err)
		}
		hash := md5.Sum([]byte(storeinfo.Method + storeinfo.Url))
		docid := hex.EncodeToString(hash[:])
		sdb.Push(docid, docBytes)
	}
}

func main() {
	cfg := config.LoadConfigFromCli()
	sdb = db.NewMqProducer("webmap", cfg.DB)
	opts := &proxy.Options{
		Addr:              cfg.Addr,
		StreamLargeBodies: 1024 * 1024 * 5,
		Upstream:          cfg.Upstream,
		SslInsecure:       cfg.SslInsecure,
	}

	p, err := proxy.NewProxy(opts)
	if err != nil {
		log.Fatal(err)
	}

	p.AddAddon(&SaveResponse{})

	log.Fatal(p.Start())
}
