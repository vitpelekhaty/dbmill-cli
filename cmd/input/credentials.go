package input

import (
	"bufio"
	"fmt"
	"os"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

// Username возвращает введенное пользователем имя пользователя БД
func Username() string {
	return readLine("Username: ")
}

// Password возвращает введенный пользователем пароль
func Password() string {
	return readPassword("Password: ")
}

func readLine(caption string) string {
	fmt.Print(caption)

	reader := bufio.NewReader(os.Stdin)

	if text, err := reader.ReadString('\n'); err == nil {
		return text
	}

	return ""
}

func readPassword(caption string) string {
	defer fmt.Println()

	fmt.Print(caption)

	if pwd, err := terminal.ReadPassword(int(syscall.Stdin)); err == nil {
		return string(pwd)
	}

	return ""
}
