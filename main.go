package main

import (
	"fmt"
)

func main() {
	result := twoSum([]int{3, 10, 5, 4, -2, 15}, 1)
	fmt.Printf("[%d, %d]\n", result[0], result[1])
}

func twoSum(nums []int, target int) []int {

	count_numbers := len(nums)

	for i, first_value := range nums[:count_numbers-1] {
		for j, second_value := range nums[i+1 : count_numbers] {

			if first_value+second_value == target {
				return []int{i, j + i + 1}
			}

		}
	}

	panic("twoSum: incorect input")
}
