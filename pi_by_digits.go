// Computing pi with arbitrary precision using Machin's formula.
// Algorithm is taken from:
// http://en.literateprograms.org/Pi_with_Machin%27s_formula_%28Python%29

package main

import (
    "fmt"
    "math/big"
    "os"
    "path/filepath"
    "strconv"
)

func main() {
    places := handleCommandLine(1000)  // 1000 digits is the default
    scaledPi := fmt.Sprint(π(places))
    fmt.Printf("3.%s\n", scaledPi[1:])
}

func handleCommandLine(defaultValue int) int {
    if len(os.Args) > 1 {
        if os.Args[1] == "-h" || os.Args[1] == "--help" {
            // handle call for help
            usage := "usage: %s [digits]\n e.g.: %s 10000"
            app := filepath.Base(os.Args[0])
            fmt.Fprintln(os.Stderr, fmt.Sprintf(usage, app, app))
            os.Exit(1)
        }
        
        if x, err := strconv.Atoi(os.Args[1]); err != nil {
            fmt.Fprintf(os.Stderr, "ignoring invalid number of " +
                "digits: will display %d\n", defaultValue)
        } else {
            return x
        }
    }
    
    return defaultValue
}

func π(places int) *big.Int {
    digits := big.NewInt(int64(places))
    unity := big.NewInt(0)
    ten := big.NewInt(10)
    exponent := big.NewInt(0)
    
    // Compute the unity scaling factor, add extra 10 digits 
    // to avoid rounding errors
    // unity = 10**(digits + 10)
    unity.Exp(ten, exponent.Add(digits, ten), nil)
    
    // Start approximation of pi with 4
    pi := big.NewInt(4)
    
    // Machin's formula
    // pi = 4 * (4 * arccot(5) - arccot(239))
    
    // Left part of Machin's formula
    left := arccot(big.NewInt(5), unity)
    left.Mul(left, big.NewInt(4))
    
    // Right part of Machin's formula
    right := arccot(big.NewInt(239), unity)
    
    // Subtract right from left and save result in left
    left.Sub(left, right)
    
    // Bring it all together to compute pi: pi = 4 * left
    pi.Mul(pi, left)
    
    // Remove the extra 10 digits
    // pi = pi / 10**10
    pi.Div(pi, big.NewInt(0).Exp(ten, ten, nil))
    
    return pi
}

// Compute arccot with a given precision
//
//             1     1     1     1
// arccot(x) = -  - --- + --- - --- + ...
//              1     3     5     7
//             x    3x    5x    7x
//
// To calculate arccot of an argument x, we start by dividing the number 1 
// (represented by 10n, which we provide as the argument unity) by x to obtain
// the first term. We then repeatedly divide by x**2 and a counter value that 
// runs over 3, 5, 7, ..., to obtain each next term. The summation is stopped 
// at the first zero term, which in this fixed-point representation corresponds 
// to a real value less than 10-n.

func arccot(x, unity *big.Int) *big.Int {
    // Init sum with 1/x
    sum := big.NewInt(0)
    sum.Div(unity, x)
    
    // Init xpower with 1/x
    xpower := big.NewInt(0)
    xpower.Div(unity, x)
    
    // Init n with 3, sign with -1, zero with 0 and square with x*x
    n := big.NewInt(3)
    sign := big.NewInt(-1)
    zero := big.NewInt(0)
    square := big.NewInt(0)
    square.Mul(x, x)
    
    // Compute successive terms until first term is 0
    for {
        // xpower = xpower / x*x
        xpower.Div(xpower, square)
        
        //         1
        // term = ---
        //          n
        //        nx
        term := big.NewInt(0)
        term.Div(xpower, n)
        
        if term.Cmp(zero) == 0 { // term == 0
            break
        }
        
        // sum = sum + sign*term
        addend := big.NewInt(0)
        sum.Add(sum, addend.Mul(sign, term))
        
        // Prepare for next iteration
        // sign = -sign
        // n = n + 2
        sign.Neg(sign)
        n.Add(n, big.NewInt(2))
    }
    
    return sum
}