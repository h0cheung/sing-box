package sniff_test

import (
	"context"
	"strings"
	"testing"

	"github.com/sagernet/sing-box/common/sniff"

	"github.com/stretchr/testify/require"
)

func TestSniffHTTP1(t *testing.T) {
	t.Parallel()
	pkt := "GET / HTTP/1.1\r\nHost: www.google.com\r\nAccept: */*\r\n\r\n"
	sniffdata := make(chan sniff.SniffData, 1)
	var data sniff.SniffData
	sniff.HTTPHost(context.Background(), strings.NewReader(pkt), sniffdata)
	data, ok := <-sniffdata
	if ok {
		metadata := data.GetMetadata()
		err := data.GetErr()
		require.NoError(t, err)
		require.Equal(t, metadata.Domain, "www.google.com")
	}
}

func TestSniffHTTP1WithPort(t *testing.T) {
	t.Parallel()
	pkt := "GET / HTTP/1.1\r\nHost: www.gov.cn:8080\r\nAccept: */*\r\n\r\n"
	sniffdata := make(chan sniff.SniffData, 1)
	var data sniff.SniffData
	sniff.HTTPHost(context.Background(), strings.NewReader(pkt), sniffdata)
	data, ok := <-sniffdata
	if ok {
		metadata := data.GetMetadata()
		err := data.GetErr()
		require.NoError(t, err)
		require.Equal(t, metadata.Domain, "www.gov.cn")
	}
}
