package main

import (
	"fmt"

	"github.com/aapldev00/synthetic_gen_gaas/internal/generator"
	_ "github.com/aapldev00/synthetic_gen_gaas/internal/generator/providers"
)

func main() {
	fmt.Println("--- Synthetic Data Generator (GaaS) ---")

	// Listamos los proveedores que se han registrado solos
	providers := generator.ListProviders()

	fmt.Printf("Registered providers: %v\n", providers)

	if len(providers) > 0 {
		fmt.Println("System ready to generate data.")
	} else {
		fmt.Println("Warning: No providers registered.")
	}
}
