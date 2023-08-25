package console

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"simplepractice.com/wp-post-migrator/migrators"
)

func printOptions() {
	fmt.Println()
	fmt.Println("1: Migrate Authors")
	fmt.Println("2: Migrate Blog Posts")
	fmt.Println("3: Delete Authors")
	fmt.Println("4: Delete Blog Posts")
	fmt.Println("5: Print Options")
	fmt.Println("6: Exit program")
	fmt.Println()
}

func optionsMessage() {
	fmt.Println("\n5 to see options.")
	fmt.Println()
}

func runPostsMigrator(s *bufio.Scanner) {
	fmt.Print("Posts per page: ")
	s.Scan()
	pppString := s.Text()
	postsPerPage, err := strconv.Atoi((pppString))
	if err != nil {
		fmt.Println("Invalid entry")
		return 
	}
	fmt.Println()
	fmt.Print("Page number: ")
	s.Scan()
	pnString := s.Text()
	pageNumber, err := strconv.Atoi(pnString)
	if err != nil {
		fmt.Println("Invalid entry")
		return 
	}
	migrators.MigratePosts(postsPerPage, pageNumber)
}

func StartProgram() {
	s := bufio.NewScanner(os.Stdin)
	var o string
	fmt.Println()
	fmt.Println()
	fmt.Println("*************** WP -> Builder Migrator ***************")
	fmt.Println()
	fmt.Println("****** Please select one of the below options ******")
	printOptions()

	for o != "6" {
		s.Scan()
		o = s.Text()
		switch o {
		case "1":
			migrators.MigrateAuthors()
		case "2":
			runPostsMigrator(s)
			printOptions()
		case "5":
			printOptions()
		case "6":
			fmt.Println("\nExiting program...")
		default:
			fmt.Println()
			fmt.Println("Unknown Command")
			optionsMessage()
		}
	}
}
