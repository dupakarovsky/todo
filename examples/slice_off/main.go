package main

import "fmt"

func main() {

	sliceOff := func(slice []int, pos int) ([]int, error) {
		if pos <= 0 || pos > len(slice) {
			return nil, fmt.Errorf("invalid position %d", pos)
		}

		left := slice[:pos-1] // 0 1
		right := slice[pos:]  // 3 4 5

		fmt.Printf("left %+v\n", left)
		fmt.Printf("right %+v\n", right)

		sliced := append(left, right...)

		return sliced, nil
	}

	s, err := sliceOff([]int{1, 2, 3, 4, 5, 6}, 3)
	if err != nil {
		fmt.Print(err)
	}

	fmt.Printf("sliced: %+v", s)

}
