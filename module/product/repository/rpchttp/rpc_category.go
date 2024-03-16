package rpchttp

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"io"
	"io/ioutil"
	"log"
	"my-app/module/product/query"
	"net/http"
	"strings"
	"time"
)

type rpcGetCategoriesByIds struct {
	url string
}

func NewRpcGetCategoriesByIds(url string) rpcGetCategoriesByIds {
	return rpcGetCategoriesByIds{url: url}
}

func (rpc rpcGetCategoriesByIds) FindWithIds(ctx context.Context, ids []uuid.UUID) ([]query.CategoryDTO, error) {
	url := rpc.url
	method := "GET"

	data := struct {
		Ids []uuid.UUID
	}{
		Ids: ids,
	}

	dataByte, _ := json.Marshal(data)

	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest(method, url, bytes.NewReader(dataByte))

	if err != nil {
		log.Println(err)
		return nil, errors.New("cannot get categories")
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return nil, errors.New("cannot get categories")
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return nil, errors.New("cannot get categories")
	}

	var resData struct {
		Data []query.CategoryDTO `json:"data"`
	}
	if err := json.Unmarshal(body, &resData); err != nil {
		log.Println(err)
		return nil, errors.New("cannot get categories")
	}

	return resData.Data, nil
}

func a() {
	url := "http://localhost:3000/v1/rpc/query-categories-ids"
	method := "GET"

	payload := strings.NewReader(`{
	"ids": [
        "018e26fd-577e-775d-ab81-5daacaed3d53",
        "018e2701-b565-7cf8-a962-8ce3270be718"
    ]
}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}
