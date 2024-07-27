package main

import (
	"encoding/json"
	"io"
	"net/http"
)

type Cep struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Unidade     string `json:"unidade"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

func main() {
	http.HandleFunc("/", HandlerCepShearch)
	http.ListenAndServe(":8080", nil)
}

func HandlerCepShearch(writer http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	cepParam := req.URL.Query().Get("cep")
	if cepParam == "" {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	cep, err := CepSearch(cepParam)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	json.NewEncoder(writer).Encode(cep)
}

func CepSearch(cep string) (*Cep, error) {
	res, err := http.Get("https://viacep.com.br/ws/" + cep + "/json/")
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var viaCep Cep
	err = json.Unmarshal(body, &viaCep)
	if err != nil {
		return nil, err
	}
	return &viaCep, nil
}
