package handlers

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"khoa-example.com/product-api/data"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		p.getProducts(rw, r)
		return
	}

	if r.Method == http.MethodPost {
		p.addProducts(rw, r)
		return
	}

	if r.Method == http.MethodPut {
		//expect the ID in the URI
		ro := regexp.MustCompile("/([0-9]+)")
		g := ro.FindAllStringSubmatch(r.URL.Path, -1)
		p.l.Println(g)

		if len(g[0]) != 2 {
			p.l.Println("Invalid URL, more than 1 capture group")
			http.Error(rw, "Invalid URL", http.StatusInternalServerError)
			return
		}

		idString := g[0][1]
		id, err := strconv.Atoi(idString)
		if err != nil {
			p.l.Println("Invalid URL, unable to convert to num ", idString)
			http.Error(rw, "Invalid URL", http.StatusInternalServerError)
			return
		}

		p.updateProducts(id, rw, r)
		return
	}

	// catch all
	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func (p *Products) getProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle Get Product")
	lp := data.GetProducts()
	// d, err := json.Marshal(lp)
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshall json", http.StatusInternalServerError)
	}
	// rw.Write(d)
}

func (p *Products) addProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle Post Product")

	prod := &data.Product{}
	err := prod.FromJSON(r.Body)
	if err != nil {
		fmt.Println(err)
		http.Error(rw, "Unable to unmarshall json", http.StatusInternalServerError)
	}

	p.l.Printf("Prod: %#v", prod)
	data.AddProduct(prod)
}

func (p *Products) updateProducts(id int, rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle Put Product")
	prod := &data.Product{}
	err := prod.FromJSON(r.Body)
	if err != nil {
		fmt.Println(err)
		http.Error(rw, "Unable to unmarshall json", http.StatusInternalServerError)
	}

	err = data.UpdateProduct(id, prod)
	if err == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Product not found", http.StatusInternalServerError)
		return
	}

}
