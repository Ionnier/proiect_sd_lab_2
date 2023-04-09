package main

import (
	"fmt"
	"sync"
)

func mapFunction(intrare <-chan string, iesire chan<- map[string]int) {
	contorizare := make(map[string]int, 0)
	for cuvant := range intrare {
		if checkVowels(cuvant) {
			contorizare["total"] = contorizare["total"] + 1
		}
	}
	contorizare["contor_vector"] = 1
	iesire <- contorizare
	close(iesire)
}

func reduceFunction(intrare <-chan int, iesire chan<- float32) {
	contor := 0
	for n := range intrare {
		contor += n
	}
	iesire <- float32(contor)
	close(iesire)
}

func citireInput(iesire [3]chan<- string) {

	intrare := [][]string{
		{"ana", "parc", "impare", "era", "copil"},
		{"cer", "program", "leu", "alee", "golang","info"},
		{"inima", "impar", "apa", "eleve"},
	}

	for i := range iesire {
		go func(caracter chan<- string, cuvinte []string) {
			for _, c := range cuvinte {
				caracter <- c
			}
			close(caracter)
		}(iesire[i], intrare[i])
	}
}

func genereazaAmestecare(intrare []<-chan map[string]int, iesire [2]chan<- int) {
	
	var sincronizare sync.WaitGroup
	sincronizare.Add(len(intrare))
	for _, caracter := range intrare {
		go func(c <- chan map[string]int){
			defer sincronizare.Done()
			for i:= range c {
				contor_vector, ok := i["contor_vector"]
				if ok {
					iesire[0] <- contor_vector
				}

				total, ok := i["total"]
				if ok {
					iesire[1] <- total
				}
			}
		}(caracter)
	}
	go func() {
		sincronizare.Wait()
		close(iesire[0])
		close(iesire[1])
	}()

}

func scrieMedie(intrare [] <- chan float32) {
	var sincronizare sync.WaitGroup
	sincronizare.Add(len(intrare))
	contor := 0.0
	total := 0.0
	for i:=0; i<len(intrare); i++ {
		go func(pozitieComponenta int, caracter <-chan float32){
			defer sincronizare.Done()
			for medie := range caracter{
				if (pozitieComponenta == 0 ) {
					contor = float64(medie)
				} else {
					total = float64(medie)
				}
			}
		}(i, intrare[i])
	}
	sincronizare.Wait()
	fmt.Printf("%v\n", total/contor)
}

func main() {
	dimensiune := 12
	componentaText1 := make(chan string, dimensiune)
	componentaText2 := make(chan string, dimensiune)
	componentaText3 := make(chan string, dimensiune)

	componentaMap1 := make(chan map[string] int)
	componentaMap2 := make(chan map[string] int)
	componentaMap3 := make(chan map[string] int)

	componentaReduce1 := make(chan int, dimensiune)
	componentaReduce2 := make(chan int, dimensiune)

	componentaMedie1 := make(chan float32, dimensiune)
	componentaMedie2 := make(chan float32, dimensiune)

	go citireInput([3]chan <- string{componentaText1, componentaText2, componentaText3})
	go mapFunction(componentaText1, componentaMap1)
	go mapFunction(componentaText2, componentaMap2)
	go mapFunction(componentaText3, componentaMap3)

	go genereazaAmestecare([]<-chan map[string]int{componentaMap1, componentaMap2, componentaMap3}, [2]chan <- int{componentaReduce1, componentaReduce2})

	go reduceFunction(componentaReduce1, componentaMedie1)
	go reduceFunction(componentaReduce2, componentaMedie2)

	scrieMedie([] <- chan float32{componentaMedie1, componentaMedie2})
}


func checkVowels(cuvant string) bool {
	vowels := []string{"a", "e", "i", "o", "u"}
	if isElementExist(vowels, cuvant[0]) && isElementExist(vowels, cuvant[len(cuvant)-1]) {
		return true
	}
	return false
}


func isElementExist(s []string, str byte) bool {
	for _, v := range s {
	  if v == string(str) {
		return true
	  }
	}
	return false
  }