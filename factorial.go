package main

import (
	"fmt"
	"math/big"
	"runtime"
	"sync"
	"os"
	"time"
)

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())  // Use all available cores

	numCPU := runtime.NumCPU()
	fmt.Println("\n# CPUs: ",numCPU)      // Mac Mini M2 Pro system uses 6 performance and 4 efficiency cores (10 total)

	n := 100000000                    	  // Change this to compute n!

	fmt.Printf("\nComputing %d!\n", n)
 	
	start1 := time.Now()

	result := Factorial(n)

	elapsed1 := time.Since(start1)

    fmt.Printf("Factorial took %s\n", elapsed1.Truncate(time.Second))
	
	start2 := time.Now()

	fmt.Println("Length of result:", len(result.String()))

	elapsed2 := time.Since(start2)

	fmt.Printf("Calculation time for string length %s\n", elapsed2.Truncate(time.Second) )

	// Write factorial to file as a string

	start3 := time.Now()

	err := writeBigIntToFile("/Volumes/NVME2TB/factorial.txt", result)
	if err != nil {
		fmt.Println("\nError writing to file:", err)
	} else {
		elapsed3 := time.Since(start3)
		fmt.Println("\nSuccessfully wrote string to file in:", elapsed3.Truncate(time.Second))
		elapsed := time.Since(start1)
    	fmt.Printf("\nTotal program time: %s", elapsed.Truncate(time.Second))
	}

}

// Factorial computes n! using the parallel prime swing algorithm

func Factorial(n int) *big.Int {
	if n < 2 {
		return big.NewInt(1)
	}
	return recFactorial(n)
}

func recFactorial(n int) *big.Int {
	if n < 2 {
		return big.NewInt(1)
	}
	f := recFactorial(n / 2)
	f.Mul(f, f)
	f.Mul(f, Swing(n))
	return f
}

func writeBigIntToFile(filename string, n *big.Int) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(n.String())
	return err
}

// Swing computes the swing(n) used in Prime Swing.
func Swing(n int) *big.Int {
	if n < 33 {
		return smallOddSwing(n)
	}

	primes := sieve(n)
	var wg sync.WaitGroup
	ch := make(chan *big.Int, len(primes))

	for _, p := range primes {
		if p > n/2 {
			break
		}
		e := 0
		q := n
		for q /= p; q > 0; q /= p {
			if q%2 == 1 {
				e++
			}
		}
		if e > 0 {
			wg.Add(1)
			go func(p, e int) {
				defer wg.Done()
				ch <- big.NewInt(0).Exp(big.NewInt(int64(p)), big.NewInt(int64(e)), nil)
			}(p, e)
		}
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	result := big.NewInt(1)
	for val := range ch {
		result.Mul(result, val)
	}

	return result
}

// smallOddSwing returns precomputed swing values for small n.
func smallOddSwing(n int) *big.Int {
	// Precomputed small odd swings
	swingTable := []int64{
		1, 1, 1, 2, 2, 4, 6, 12, 20, 36,
		60, 112, 192, 336, 576, 960, 1600, 2688,
		4608, 7936, 13440, 23040, 38400, 64512, 107520,
		180224, 294912, 491520, 819200, 1351680, 2211840, 3670016, 6094848,
	}
	return big.NewInt(swingTable[n])
}

// sieve returns all primes up to n using the Sieve of Eratosthenes.
func sieve(n int) []int {
	isPrime := make([]bool, n+1)
	for i := 2; i <= n; i++ {
		isPrime[i] = true
	}
	for i := 2; i*i <= n; i++ {
		if isPrime[i] {
			for j := i * i; j <= n; j += i {
				isPrime[j] = false
			}
		}
	}
	primes := []int{}
	for i := 2; i <= n; i++ {
		if isPrime[i] {
			primes = append(primes, i)
		}
	}
	return primes
}