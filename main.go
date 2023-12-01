package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Input: ")
	if scanner.Scan() {
		input := scanner.Text()
		result, err := calculate(input)
		if err != nil {
			fmt.Println("Output:", err)
		} else {
			fmt.Println("Output:", result)
		}
	}
}

func calculate(input string) (string, error) {
	input = strings.ToUpper(strings.TrimSpace(input))
	match, _ := regexp.MatchString(`^(I{1,3}|IV|V|VI{0,3}|IX|X|[1-9]|10) [\+\-\*\/] (I{1,3}|IV|V|VI{0,3}|IX|X|[1-9]|10)$`, input)
	if !match {
		return "", fmt.Errorf("Ошибка: неверный формат ввода")
	}
	parts := strings.Split(input, " ")
	op1, op2 := parts[0], parts[2]
	operator := parts[1]
	isRoman := strings.ContainsAny(op1, "IVX")
	if (isRoman && strings.ContainsAny(op2, "0123456789")) || (!isRoman && strings.ContainsAny(op2, "IVX")) {
		return "", fmt.Errorf("Ошибка: разные системы счисления")
	}
	var num1, num2 int
	var err error
	if isRoman {
		num1, err = romanToArabic(op1)
		if err != nil {
			return "", err
		}
		num2, err = romanToArabic(op2)
		if err != nil {
			return "", err
		}
	} else {
		num1, err = strconv.Atoi(op1)
		if err != nil {
			return "", fmt.Errorf("Ошибка: неверное число %s", op1)
		}
		num2, err = strconv.Atoi(op2)
		if err != nil {
			return "", fmt.Errorf("Ошибка: неверное число %s", op2)
		}
	}
	result := 0
	switch operator {
	case "+":
		result = num1 + num2
	case "-":
		result = num1 - num2
		if isRoman && result < 1 {
			return "", fmt.Errorf("Ошибка: в римской системе нет отрицательных чисел")
		}
	case "*":
		result = num1 * num2
	case "/":
		if num2 == 0 {
			return "", fmt.Errorf("Ошибка: деление на ноль")
		}
		result = num1 / num2
	default:
		return "", fmt.Errorf("Ошибка: неверный оператор")
	}
	if isRoman {
		return arabicToRoman(result), nil
	}
	return strconv.Itoa(result), nil
}

func romanToArabic(roman string) (int, error) {
	romanNumerals := map[rune]int{'I': 1,
		'V': 5, 'X': 10,
	}
	var arabic int
	var lastValue int
	for i := len(roman) - 1; i >= 0; i-- {
		value, exists := romanNumerals[rune(roman[i])]
		if !exists {
			return 0, fmt.Errorf("Ошибка: неверный символ %c", roman[i])
		}
		if value < lastValue {
			arabic -= value
		} else {
			arabic += value
		}
		lastValue = value
	}
	if arabic > 10 {
		return 0, fmt.Errorf("Ошибка: число больше 10")
	}
	return arabic, nil
}

func arabicToRoman(number int) string {
	if number < 1 {
		return "Ошибка: римские числа не могут быть меньше 1"
	}
	romanNumerals := []struct {
		Value  int
		Symbol string
	}{
		{1000, "M"}, {900, "CM"}, {500, "D"}, {400, "CD"},
		{100, "C"}, {90, "XC"}, {50, "L"}, {40, "XL"},
		{10, "X"}, {9, "IX"}, {5, "V"}, {4, "IV"}, {1, "I"},
	}
	var roman string
	for _, numeral := range romanNumerals {
		for number >= numeral.Value {
			roman += numeral.Symbol
			number -= numeral.Value
		}
	}
	return roman
}
