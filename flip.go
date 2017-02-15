package main

import(
	"fmt"
	"os"
	"os/exec"
	"bufio"
	"os/signal"
	"syscall"
	"time"
	"path/filepath"
	"strings"
)

func main() {

	// tratamiento de un Ctrl+C
	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		os.Exit(0)
	}()
	
	// nombre+extension
	filename := filepath.Base(os.Args[0])
	
	filename = strings.TrimSuffix(filename, filepath.Ext(filename))
	// ff = nombre sin extension
	ff := filename

	fmt.Println("Archivo " + ff + ".go:")
	
	// elimina el fichero gemelo
	err := os.Remove("./" + opposite(ff) + ".go")
	if err != nil {
		fmt.Println("\tImposible eliminar " + opposite(ff) + ".go")
	} else {
		fmt.Println("\t" + opposite(ff)+".go eliminado, esperando ...")
	}

	time.Sleep(time.Second * 3) // tiempo de espera

	// crear reader
	// de aqui leemos el archivo actual para copiarlo
	fr, err := os.Open("./" + ff + ".go")
	if err != nil {
		panic(err)
	}
	defer fr.Close()

	scanner := bufio.NewScanner(fr)

	// crear writer
	// crea nuevo fichero gemelo 
	fw, err := os.Create("./" + opposite(ff) + ".go")
	if err != nil {
		panic(err)
	}
	defer fw.Close()

	writer := bufio.NewWriter(fw)

	//copiar

	fmt.Println("\tCopiando ...")

	for scanner.Scan() {

		writer.WriteString(scanner.Text())
		writer.WriteString("\n")
		writer.Flush()
	}
	fr.Close()
	fw.Close()

	// Ejecucion y cierre

	fmt.Println("\tPreparando programa gemelo ...")

	cmd := exec.Command("go", "run", opposite(ff)+".go")
	
	fmt.Println("Exit " + ff + ".go\n")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()

}

func opposite(s string) (ret string){
	if s == "flip"{
		ret = "flop"
	}
	if s == "flop"{
		ret = "flip"
	}

	return
}

