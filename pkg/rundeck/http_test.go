package rundeck

import (
	"testing"

	httpclient "github.com/ahuret/go-rundeck/pkg/httpclient"
	"github.com/ahuret/go-rundeck/pkg/rundeck/responses"
	"github.com/stretchr/testify/require"
)

func TestHTTP404(t *testing.T) {
	client, server, _ := newTestRundeckClient([]byte("hello"), "application/json", 404)
	defer server.Close()
	funcs := map[string]func(string, ...httpclient.RequestOption) ([]byte, error){
		"get":  client.httpGet,
		"put":  client.httpPut,
		"post": client.httpPost,
	}
	for n, f := range funcs {
		res, err := f("/", requestExpects(200))
		require.Nil(t, res, n+" body should be nil")
		require.Error(t, err, n+" should return an error")
		require.IsType(t, ErrMissingResource, err, n+" should return ErrMissingResource")
	}
	_, err := client.httpDelete("/f", requestExpects(204))
	require.Error(t, err, "delete should return an error")
	require.IsType(t, ErrMissingResource, err, "delete should return ErrMissingResource")
}

func TestRDErrorResponse(t *testing.T) {
	jsonfile, err := responses.GetTestData(responses.ErrorResponseTestFile)
	if err != nil {
		t.Fatalf(err.Error())
	}

	client, server, _ := newTestRundeckClient(jsonfile, "application/json", 500)
	defer server.Close()
	funcs := map[string]func(string, ...httpclient.RequestOption) ([]byte, error){
		"get":  client.httpGet,
		"put":  client.httpPut,
		"post": client.httpPost,
	}
	for n, f := range funcs {
		res, reserr := f("/", requestExpects(200))
		require.Nil(t, res, n+" body should be nil")
		require.Error(t, reserr, n+" should return an error")
		require.Equal(t, "something blew up", reserr.Error())
	}
	_, reserr := client.httpDelete("/f", requestExpects(204))
	require.Error(t, reserr, "delete should return an error")
	require.Equal(t, "something blew up", reserr.Error())
}
