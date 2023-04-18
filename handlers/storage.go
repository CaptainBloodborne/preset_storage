package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type PresetHolder struct {
	Token    string
	Preset   []Nm
	Total    int
	Subjects []string
	Brands   []string
}

type Nm struct {
	NmId        int     `json:"nm_id"`
	SubjectId   int     `json:"subject_id"`
	BrandId     int     `json:"brand_id"`
	StockExists bool    `json:"stock_exists"`
	Score       float64 `json:"score"`
}

type Data struct {
	Data []PresetHolder `json:"data"`
}

var store = map[string]PresetHolder{}
var lock sync.RWMutex

func HandleGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	key := r.URL.Query().Get("key")
	subjectsFilter := r.URL.Query().Get("subject")
	brandsFilter := r.URL.Query().Get("brand")

	page, err := strconv.ParseInt(r.URL.Query().Get("page"), 10, 64)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Parsing page number has failed!")
		page = 1
	}
	offset, err := strconv.ParseInt(r.URL.Query().Get("offset"), 10, 64)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Parsing offset has failed!")
	}
	lock.RLock()
	var filteredNms []Nm
	if subjectsFilter != "" {
		subjects := strings.Split(subjectsFilter, "|")
		var intSubjects []int64
		intSubjects = convertStringSliceToIntSlice(&subjects)

		for _, value := range store[key].Preset {
			if filterBySubject(value, intSubjects) {
				filteredNms = append(filteredNms, value)
			}
		}

	}

	var filteredBrandNms []Nm

	if brandsFilter != "" {
		brands := strings.Split(brandsFilter, "|")
		var intBrands []int64
		intBrands = convertStringSliceToIntSlice(&brands)

		if len(filteredNms) > 0 {
			for _, value := range filteredNms {
				if filterByBrand(value, intBrands) {
					filteredBrandNms = append(filteredBrandNms, value)
				}
			}
		} else {
			for _, value := range store[key].Preset {
				if filterByBrand(value, intBrands) {
					filteredBrandNms = append(filteredBrandNms, value)
				}
			}
		}

	}

	var value PresetHolder

	if len(filteredNms) == 0 && len(filteredBrandNms) == 0 {
		if len(store[key].Preset) == 0 {
			value = PresetHolder{Token: key, Preset: []Nm{}}
		} else if len(store[key].Preset) < 100 && len(store[key].Preset) > 0 {
			//value = store[key].Preset[0 : len(store[key].Preset)-1]
			value = PresetHolder{
				Token:    key,
				Preset:   store[key].Preset[0 : len(store[key].Preset)-1],
				Subjects: store[key].Subjects,
				Brands:   store[key].Brands,
			}
		} else {
			value = PresetHolder{
				Token:    key,
				Preset:   store[key].Preset[(page-1)*offset : offset*page],
				Subjects: store[key].Subjects,
				Brands:   store[key].Brands,
			}
		}

	} else if len(filteredNms) > 0 && len(filteredBrandNms) == 0 {
		var outerBound int
		var innerBound int
		if int(offset*page) > len(filteredNms) {
			outerBound = len(filteredNms) - 1
			innerBound = 0
		} else {
			outerBound = int(offset * page)
			innerBound = int((page - 1) * offset)
		}

		//value = filteredNms[innerBound:outerBound]
		value = PresetHolder{
			Token:    key,
			Preset:   store[key].Preset[innerBound:outerBound],
			Subjects: store[key].Subjects,
			Brands:   store[key].Brands,
		}
	} else if len(filteredNms) == 0 && len(filteredBrandNms) > 0 {
		var outerBound int
		var innerBound int
		if int(offset*page) > len(filteredBrandNms) {
			outerBound = len(filteredBrandNms) - 1
			innerBound = 0
		} else {
			outerBound = int(offset * page)
			innerBound = int((page - 1) * offset)
		}

		//value = filteredBrandNms[innerBound:outerBound]
		value = PresetHolder{
			Token:    key,
			Preset:   store[key].Preset[innerBound:outerBound],
			Subjects: store[key].Subjects,
			Brands:   store[key].Brands,
		}
	} else if len(filteredNms) > 0 && len(filteredBrandNms) > 0 {
		var outerBound int
		var innerBound int
		if int(offset*page) > len(filteredBrandNms) {
			outerBound = len(filteredBrandNms) - 1
			innerBound = 0
		} else {
			outerBound = int(offset * page)
			innerBound = int((page - 1) * offset)
		}

		//value = filteredBrandNms[innerBound:outerBound]
		value = PresetHolder{
			Token:    key,
			Preset:   store[key].Preset[innerBound:outerBound],
			Subjects: store[key].Subjects,
			Brands:   store[key].Brands,
		}
	}
	lock.RUnlock()
	json.NewEncoder(w).Encode(value)
	fmt.Println(len(value.Preset))
}

func HandleSet(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method)
	decoder := json.NewDecoder(r.Body)

	var preset PresetHolder
	err := decoder.Decode(&preset)

	if err != nil {
		fmt.Println("Error decoding data")
		return
	}

	lock.Lock()
	store[preset.Token] = preset
	lock.Unlock()
	fmt.Println(preset.Token)
	fmt.Println("Preset len is:", len(preset.Preset))
	//if r.Method == "POST" {
	//	err := r.ParseForm()
	//
	//	if err != nil {
	//		http.Error(w, "Error parsing form data", http.StatusBadRequest)
	//		return
	//	}
	//	fmt.Println(r)
	//	key := r.FormValue("sos")
	//	fmt.Println(key)
	//	//value := r.FormValue("value")
	//	fmt.Fprintf(w, "Received POST request with key %v", key)
	//	return
	//}
	//key := r.URL.Query().Get("key")
	//value := r.URL.Query().Get("value")
	//lock.Lock()
	//store[key] = value
	//lock.Unlock()
	w.Write([]byte("OK"))
}

func filterBySubject(nm Nm, subjects []int64) bool {
	for _, value := range subjects {
		if value == int64(nm.SubjectId) {
			return true
		}

	}
	return false
}

func filterByBrand(nm Nm, brands []int64) bool {
	for _, value := range brands {
		if value == int64(nm.BrandId) {
			return true
		}

	}
	return false
}

func convertStringSliceToIntSlice(stringSlice *[]string) []int64 {
	var newIntSlice []int64

	for _, value := range *stringSlice {
		intValue, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			fmt.Println(err)
		}
		newIntSlice = append(newIntSlice, intValue)
	}
	return newIntSlice
}
