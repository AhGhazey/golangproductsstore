package rest

import (
	"cmd/ims.server/pkg/adding"
	"cmd/ims.server/pkg/listing"
	"cmd/ims.server/pkg/updating"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/jszwec/csvutil"
	"github.com/julienschmidt/httprouter"
)

func Handler(a adding.Service, l listing.Service, u updating.Service) http.Handler {
	router := httprouter.New()
	router.GET("/", home())
	router.POST("/addProduct", addProduct(a))
	router.POST("/updatebulck", updateBulkRecords(u))
	router.GET("/getproductbysku/:sku", getProductBySku(l))
	router.POST("/consumeproduct", consumeProduct(u))

	return router
}

func home() func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		fmt.Fprint(w, "Hello Jumia")
	}
}

func getProductBySku(l listing.Service) func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		sku := ps.ByName("sku")
		products, err := l.GetProductBySku(sku)

		if err != nil {
			log.Println("error getting product", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		if len(*products) == 0 {
			http.Error(w, err.Error(), http.StatusNotFound)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(products)
	}
}

func consumeProduct(u updating.Service) func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		decoder := json.NewDecoder(r.Body)

		var product updating.Product
		err := decoder.Decode(&product)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		isvalid, err := u.ConsumeProduct(product)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		if isvalid {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode("completed successfuly")
		} else if err == nil && !isvalid {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode("stock is not enough")
		}

	}
}

func addProduct(s adding.Service) func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		decoder := json.NewDecoder(r.Body)

		var product adding.Product
		err := decoder.Decode(&product)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = s.AddProduct(product)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode("New Product added.")
	}
}

func updateBulkRecords(s updating.Service) func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		//max 10 mega byte size
		err := r.ParseMultipartForm(10 << 20)

		if err != nil {
			http.Error(w, fmt.Errorf("file size greater than 10 Mb:%w", err).Error(), http.StatusBadRequest)
			return
		}

		file, _, err := r.FormFile("myFile")
		if err != nil {
			fmt.Println("Error Retrieving the File")
			fmt.Println(err)
			return
		}

		defer file.Close()

		csvReader := csv.NewReader(file)

		productHeader, err := csvReader.Read()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		dec, err := csvutil.NewDecoder(csvReader, productHeader...)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var products []updating.ProductCSV
		for {
			var p updating.ProductCSV

			if err := dec.Decode(&p); err == io.EOF {
				break
			} else if err != nil {
				log.Fatal(err)
			}
			products = append(products, p)
		}
		err = s.UpdateBulkRecords(products)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode("products updated.")
	}
}
