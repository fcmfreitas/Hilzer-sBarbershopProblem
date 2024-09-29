package main

import (
	"fmt"
	"programasGo/FPPDSemaforo"
)

const (
	maxCustomers  = 20
	dayCostumers = 100
	sofaCapacity  = 4
	barberChairs  = 3
)

var (
	sofa           *FPPDSemaforo.Semaphore
	chair          *FPPDSemaforo.Semaphore
	barber         *FPPDSemaforo.Semaphore
	customer       *FPPDSemaforo.Semaphore
	cash           *FPPDSemaforo.Semaphore
	receipt        *FPPDSemaforo.Semaphore
	totalCustomers *FPPDSemaforo.Semaphore
	fim = make(chan int)
	cont = 0;
)

func customerProcess(id int) {
	totalCustomers.Wait() // Só permite entrar se < 20 clientes
	enterShop(id)

	// Aguarda vaga no sofá
	sofa.Wait()
	sitOnSofa(id)

	// Aguarda vaga na cadeira do barbeiro
	chair.Wait()
	sitInBarberChair(id)
	sofa.Signal() // Libera espaço no sofá

	// Avisar o barbeiro que está pronto
	customer.Signal()
	barber.Wait() // Espera o barbeiro

	// Cortanado cabelo
	getHairCut(id)

	// Cliente vai pagar
	chair.Signal() //libera cadeira 
	cash.Signal()
	receipt.Wait() // Aguarda o recibo
	pay(id)

	// Cliente deixa a barbearia
	exitShop(id)
	totalCustomers.Signal() // Libera a entrada de mais um cliente
}

func barberProcess(id int) {
	for {
		// Espera cliente
		customer.Wait()
		// Barbeiro pronto para cortar o cabelo
		barber.Signal()
		
		cutHair(id)

		cash.Wait()
		acceptPayment(id)

		// Emite recibo
		receipt.Signal()
	}
}

func enterShop(id int) {
	fmt.Printf("Cliente %d entrou na barbearia.\n", id)
}

func sitOnSofa(id int) {
	fmt.Printf("Cliente %d sentou no sofá.\n", id)
}

func sitInBarberChair(id int) {
	fmt.Printf("Cliente %d sentou na cadeira do barbeiro.\n", id)
}

func getHairCut(id int) {
	fmt.Printf("Cliente %d está cortando o cabelo.\n", id)
}

func pay(id int) {
	fmt.Printf("Cliente %d pagou e está saindo.\n", id)
}

func cutHair(id int) {
	fmt.Printf("Barbeiro %d está cortando cabelo.\n", id)
}

func acceptPayment(id int) {
	fmt.Printf("Barbeiro %d recebeu o pagamento.\n", id)
}

func exitShop(id int) {
	fmt.Printf("Cliente %d saiu da barbearia.\n", id)
	cont++;
	if(cont == dayCostumers){
		fim <- 1
	}
}

func main() {

	sofa = FPPDSemaforo.NewSemaphore(sofaCapacity)       // Capacidade do sofá
	chair = FPPDSemaforo.NewSemaphore(barberChairs)      // Cadeiras do barbeiro
	barber = FPPDSemaforo.NewSemaphore(0)                // Sincroniza o barbeiro
	customer = FPPDSemaforo.NewSemaphore(0)              // Sincroniza clientes
	cash = FPPDSemaforo.NewSemaphore(0)                  // Sincroniza pagamento
	receipt = FPPDSemaforo.NewSemaphore(0)               // Sincroniza recibo
	totalCustomers = FPPDSemaforo.NewSemaphore(maxCustomers) // Limita a barbearia em 20 clientes


	// Inicia os barbeiros
	for i := 1; i <= barberChairs; i++ {
		go barberProcess(i)
	}

	// Inicia os clientes conforme o numero de clientes no dia (ex: 100)
	for i := 1; i <= dayCostumers; i++ {
		go customerProcess(i)
	}

	<- fim //após contar que todos os clientes do dia foram embora, finaliza o programa
}
