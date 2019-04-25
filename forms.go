package main

import (
	"html/template"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
)

type St struct {
	Outputs []Content `json:"outputs"`
}

type Content struct{
	Id string `json:"id"`
	CreatedAt string `json:"created_at"`
	Data DataStr `json:"data"'`
	Model ModelStr `json:"model"'`
	Status StatusStr `json:"status"'`
	Input string `json:"input"'`
}

type InputStr struct {
	Data DataInputStr `json:"data"'`
	Id string `json:"id"'`
}

type DataInputStr struct {
	Image ImageStr `json:"image"'`
}

type ImageStr struct {
	Url string `json:"url"'`
}

type ModelStr struct {
	AppId string `json:"app_id"'`
	Id string `json:"id"'`
	Name string `json:"name"'`
	OutputInfo OutputInfoStr `json:"output_info"'`
	DisplayName string `json:"display_name"'`
	CreatedAt string `json:"created_at"`
	ModelVersion ModelVersionStr `json:"model_version"`
}

type ModelVersionStr struct {
	Id string `json:"id"'`
	Status StatusStr `json:"status"'`
	CreatedAt string `json:"created_at"`
}

type StatusStr struct {
	Code int `json:"Code"'`
	Description string `json:"description"'`
}

type OutputInfoStr struct {
	Message string `json:"message"'`
	TypeExt string `json:"type_ext"'`
	Type string `json:"type"'`
}

type DataStr struct{
	Concepts []ConceptStr `json:"concepts"'`
}

type ConceptStr struct {
	AppId string `json:"app_id"'`
	Id string `json:"id"'`
	Name string `json:"name"'`
	Value float64 `json:"value"'`
}

type ImageWeight struct {
	Url string
	Weight float64
}

var m map[string][]ImageWeight

type SearchString struct {
	SearchWord   string
}

func main() {

	m = make(map[string] []ImageWeight)
	for range [1]int{} {
		file, err := os.Open("imagesAll.txt")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		imageNumber := 0
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			fmt.Println(scanner.Text(), imageNumber)
			imageNumber++
			img := `{"inputs": [{"data": {"image": {"url": "` + scanner.Text() + `" }}}]}`
			imgData := []byte(img)

			urlAPI := "https://api.clarifai.com/v2/models/aaa03c23b3724a16a56b629203edc62c/outputs"
			req, err := http.NewRequest("POST", urlAPI, bytes.NewBuffer(imgData))
			req.Header.Set("Authorization", "Key Enter your API key here")
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				panic(err)
			}

			body, err := ioutil.ReadAll(resp.Body)
			bodyString := string(body)
			defer resp.Body.Close()


			var r St
			json.Unmarshal([]byte(bodyString), &r)
			if len(r.Outputs) > 0 {
				for _, e := range r.Outputs[0].Data.Concepts {
					m[strings.ToLower(e.Name)] = append(m[strings.ToLower(e.Name)], ImageWeight{Url:scanner.Text(), Weight:e.Value})
				}
			}
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

	}

	keys := make([]string, len(m))
	for k, _ := range m{
		sort.Slice(m[k], func(i, j int) bool {
			return m[k][i].Weight > m[k][j].Weight
		})
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})


	//Creating Server and Displaying results

	tmpl := template.Must(template.ParseFiles("forms.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			tmpl.Execute(w, nil)
			return
		}

		details := SearchString{
			SearchWord:   r.FormValue("word"),
		}
		fmt.Println("Search word is : ",details.SearchWord)
		word := strings.ToLower(details.SearchWord)
		i := sort.SearchStrings(keys, word)


		if i >= len(keys) {
			i--
		}
		fmt.Println( keys[i])

		dist := make(map[string]bool)
		result := make([]string, 0)

		j, k := 1, 1

		for _, v := range m[keys[i]] {
			if _, ok := dist[v.Url]; !ok {
				result = append(result, v.Url)
				dist[v.Url] = true
				if len(result) == 10 {
					break
				}
			}
		}

		for len(result) < 10 {
			if i + j < len(keys) && ( i - k >= 0 && strings.Compare(keys[i + j], keys[i - k]) < 0 ) {
				for _, v := range m[keys[i + j]] {
					if _, ok := dist[v.Url]; !ok {
						result = append(result, v.Url)
						dist[v.Url] = true
						if len(result) == 10 {
							break
						}
					}
				}
				j++
			} else if i - k >= 0 {
				for _, v := range m[keys[i - k]] {
					if _, ok := dist[v.Url]; !ok {
						result = append(result, v.Url)
						dist[v.Url] = true
						if len(result) == 10 {
							break
						}
					}
				}
				k++
			} else if i + j < len(keys) {
				for _, v := range m[keys[i + j]] {
					if _, ok := dist[v.Url]; !ok {
						result = append(result, v.Url)
						dist[v.Url] = true
						if len(result) == 10 {
							break
						}
					}
				}
				j++
			} else {
				break
			}
		}

		var resArray []string

		//fmt.Println("Result Map is : ", result)
		//fmt.Println("Length of Result Map is : ", len(result))
		for _, k := range result{
			resArray = append(resArray, k)
		}
		//fmt.Println("Res Array is : ", resArray)
		//fmt.Println("Length of Res Array is : ", len(resArray))

		fmt.Println("R1 is : ", resArray[0])

		tmpl.Execute(w, struct{ Success bool
		DisplayWord string
		OutArray []string
		}{true, word,resArray})
	})

	http.ListenAndServe(":8080", nil)
}