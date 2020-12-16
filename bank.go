package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mgutz/ansi"
)

type accountS struct {
	Username   string `json:"Username"`
	Password   string `json:"Password"`
	AccountNum string `json:"AccountNum"`
	Balance    int64  `json:"Balance"`
}

var (
	userLogged   accountS
	bothFilesEmt bool   = false
	usedFile     string = "accounts"
)

var (
	enterUsername   string = color("Enter username: ", "blue+b")
	enterPassword   string = color("Enter password: ", "blue+b")
	enterBalance    string = color("Enter balance: ", "blue+b")
	commandNotFound string = color("Command not found.", "cyan+b")
	accountNotFound string = color("Account not found.", "cyan+b")
)

func main() {
	var (
		command1 string
		running  bool = true
	)

	for running {
		fmt.Print(ansi.Color(">>", "cyan+b"))
		fmt.Scanln(&command1)

		switch command1 {
		case "sign":
			signUp()
			break
		case "login":
			if logIn() {
				loggedIn()
			}
			break
		case "quit":
			running = false
			break
		default:
			fmt.Println(commandNotFound)
			break
		}
		command1 = ""
	}
}

func signUp() {
	var (
		username       string
		password       string
		balance        int64
		accountNum     string
		isUsernameFree bool = true
		numLines       int  = numLinesInFile(usedFile)
		lineNum        int  = 0
	)

	fmt.Print(enterUsername)
	fmt.Scanln(&username)
	if username == "0" {
		main()
	}

	file, err := os.OpenFile(usedFile, os.O_APPEND|os.O_RDWR, 0644)
	isErr(err)
	reader := bufio.NewReader(file)

	for lineNum < numLines && isUsernameFree {
		accountLine, _, err := reader.ReadLine()
		isErr(err)
		if username == decodeAccount(accountLine).Username {
			isUsernameFree = false
			fmt.Println(ansi.Color("This username is taken.", "cyan+b"))
			main()
		}
		lineNum++
	}
	fmt.Print(enterPassword)
	fmt.Scanln(&password)
	if len(password) < 8 {
		fmt.Println(ansi.Color("Password need to be more than 8 charcters.", "cyan+b"))
		main()
	}
	fmt.Print(enterBalance)
	fmt.Scanln(&balance)

	accountNum = accountNumGen()
	user := accountS{Username: username, Password: password, AccountNum: accountNum, Balance: balance}

	accountJSON := encodeAccount(user)

	file.Write(accountJSON)
	file.WriteString("\n")
	file.Close()
}

func logIn() bool {

	if bothFilesEmt {
		fmt.Println(color("Error: file is empty.", "cyan+b"))
		return false
	}
	var (
		username     string
		password     string
		loggedInBool bool = false
		lineNum      int  = 0
		numLines     int  = numLinesInFile(usedFile)
	)

	fmt.Print(enterUsername)
	fmt.Scanln(&username)
	fmt.Print(enterPassword)
	fmt.Scanln(&password)

	file, err := os.OpenFile(usedFile, os.O_APPEND|os.O_RDONLY, 0664)
	isErr(err)
	reader := bufio.NewReader(file)

	for lineNum < numLines && !loggedInBool {
		accountJSON, _, err := reader.ReadLine()
		isErr(err)
		var account = decodeAccount(accountJSON)
		if username == account.Username {
			if password == account.Password {
				fmt.Println(color("Logged in succesfuly.", "cyan+h"))
				userLogged = accountS{account.Username, account.Password, account.AccountNum, account.Balance}
				loggedInBool = true
				return true

			}
		}
		lineNum++
	}

	fmt.Println(ansi.Color("Wrong username or password.", "blue+bh"))
	return false
}

func loggedIn() {
	var (
		command string
		running bool   = true
		v1      string = ansi.Color("(", "green+b")
		v2      string = ansi.Color(")", "green+b")
	)

	for running {
		fmt.Printf("%v%v%v %v", v1, ansi.Color(userLogged.Username, "red+b"), v2, ansi.Color(">>", "magenta+b"))
		fmt.Scanln(&command)

		switch command {
		case "info":
			data()
			break
		case "withdraw":
			withdraw()
			break
		case "deposit":
			deposit()
			break
		case "transfer":
			bankTransfer()
			break
		case "logout":
			fmt.Println(color("Logged out.", "cyan+h"))
			userLogged = accountS{}
			running = false
			break
		default:
			fmt.Println(commandNotFound)
			break
		}
		command = ""
	}
}

func data() {
	fmt.Println(color("Username:", "magenta+b"), userLogged.Username)
	fmt.Println(color("Account number:", "magenta+b"), userLogged.AccountNum)
	fmt.Println(color("Balance:", "magenta+b"), userLogged.Balance)
}

func withdraw() {
	var amount int64 = 0

	fmt.Print(color("Withdraw: ", "blue+b"))
	fmt.Scanln(&amount)
	if amount == 0 {
		fmt.Println(color("Amount to withdraw wasn't entered.", "cyan+b"))
		loggedIn()
	} else if amount < 0 {
		fmt.Println(color("Not valid amount.", "cyan+b"))
	}
	userLogged.Balance -= amount
	fmt.Printf("%v %v\n", color("New balance:", "blue+b"), userLogged.Balance)
	updateFile(userLogged)
}

func deposit() {
	var amount int64 = 0

	fmt.Print(color("Deposit: ", "blue+b"))
	fmt.Scanln(&amount)
	if amount == 0 {
		fmt.Println(color("Amount to deposit wasn't entered.", "cyan+b"))
		loggedIn()
	} else if amount < 0 {
		fmt.Println(color("Not valid amount.", "cyan+b"))
	}
	userLogged.Balance += amount
	fmt.Printf("%v %v\n", color("New balance:", "blue+b"), userLogged.Balance)
	updateFile(userLogged)
}

func bankTransfer() {
	var (
		amount         int64 = 0
		toAccountNum   string
		toAccount      accountS = accountS{}
		toAccountFound bool     = false
		lineNum        int      = 0
		numLines       int      = numLinesInFile(usedFile)
	)

	fmt.Print(color("Transfer to: ", "blue+b"))
	fmt.Scanln(&toAccountNum)

	file, err := os.OpenFile(usedFile, os.O_RDWR, 0644)
	isErr(err)
	reader := bufio.NewReader(file)

	for lineNum < numLines && !toAccountFound {
		accountJSON, _, err := reader.ReadLine()
		isErr(err)
		toAccount = decodeAccount(accountJSON)
		if toAccount.AccountNum == toAccountNum {
			toAccountFound = true
		}
		lineNum++
	}
	if !toAccountFound {
		fmt.Println(accountNotFound)
		loggedIn()
	}
	fmt.Print(color("Enter amount to transfer: ", "blue+b"))
	fmt.Scanln(&amount)
	if amount == 0 {
		fmt.Println(color("You didn't enter amount to transfer.", "cyan+b"))
		loggedIn()
	} else if amount < 0 {
		fmt.Println(color("Not valid amount.", "cyan+b"))
		loggedIn()
	} else if amount < userLogged.Balance {
		fmt.Println(color("You don't have enough money.", "cyan+b"))
		loggedIn()
	}

	userLogged.Balance -= amount
	toAccount.Balance += amount

	updateFile(userLogged, toAccount)
	fmt.Println(color("Transfer went succesfuly.", "blue+b"))
	file.Close()
}

func updateFile(userToUpdate ...accountS) {
	var (
		user          accountS = accountS{}
		lineNum       int      = 0
		numLines      int      = numLinesInFile(usedFile)
		usersUsername []string
	)

	for i := 0; i < len(userToUpdate); i++ {
		usersUsername = append(usersUsername, userToUpdate[i].Username)
	}

	file, err := os.OpenFile(usedFile, os.O_RDWR, 0644)
	isErr(err)
	reader := bufio.NewReader(file)

	fileToUpdate, err := os.OpenFile("accounts1", os.O_RDWR|os.O_CREATE, 0644)
	isErr(err)

	for lineNum < numLines {
		accountLine, _, err := reader.ReadLine()
		isErr(err)
		user = decodeAccount(accountLine)

		indexInslice := findInSlice(usersUsername, user.Username)
		if indexInslice >= 0 {
			fileToUpdate.Write(encodeAccount(userToUpdate[indexInslice]))
			fileToUpdate.WriteString("\n")
		} else {
			fileToUpdate.Write(accountLine)
			fileToUpdate.WriteString("\n")
		}
		lineNum++
	}
	err = os.Remove(usedFile)
	isErr(err)
	file.Close()
	os.Rename("accounts1", "accounts")
	fileToUpdate.Close()
}

func accountNumGen() string {
	var (
		accountNum               = strings.Builder{}
		allAccountaNums          = []string{}
		user            accountS = accountS{}
	)

	file, err := os.OpenFile(usedFile, os.O_RDONLY, 0644)
	isErr(err)
	scanner := bufio.NewScanner(file)

	for i := 0; i < numLinesInFile(usedFile); i++ {
		scanner.Scan()
		accountJSON := scanner.Bytes()
		user = decodeAccount(accountJSON)
		allAccountaNums = append(allAccountaNums, user.AccountNum)
	}

	accountNum = createAccountNum()
	for findInSlice(allAccountaNums, accountNum.String()) >= 0 {
		accountNum = createAccountNum()

	}
	return accountNum.String()
}

func findInSlice(slice []string, val string) int {
	for i := 0; i < len(slice); i++ {
		if slice[i] == val {
			return i
		}
	}
	return -1
}

func createAccountNum() strings.Builder {
	var accountNum = strings.Builder{}

	rand.Seed(time.Now().UTC().UnixNano())

	for i := 0; i < 8; i++ {
		dig := rand.Intn(10)
		accountNum.WriteString(strconv.Itoa(dig))
	}
	return accountNum
}

func numLinesInFile(filename string) int {
	file, err := os.Open(filename)
	isErr(err)
	fileScanner := bufio.NewScanner(file)
	lineCount := 0
	for fileScanner.Scan() {
		lineCount++
	}
	file.Close()
	return lineCount
}

func isErr(err error) {
	if err != nil {
		panic(err)
	}
}

func encodeAccount(account accountS) []byte {
	accountJSON, err := json.Marshal(account)
	if err != nil {
		panic(err)
	}
	return accountJSON
}

func decodeAccount(accountJSON []byte) accountS {
	var account accountS
	err := json.Unmarshal([]byte(accountJSON), &account)
	if err != nil {
		panic(err)
	}
	return account
}

func color(str string, flag string) string {
	return ansi.Color(str, flag)
}
