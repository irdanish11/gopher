package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"
)

type BoID int
type MeID int
type ChID int
type BookIDs []BoID
type MemberIDs []MeID
type CheckoutIDs []ChID
type BookAttributes struct {
	bookName        string
	publicationYear int
	copies          int
	author          string
	available       bool
	lender          MeID
}
type MemberAttributes struct {
	firstName string
	lastName  string
	age       int
	address   string
	booksLent []ChID
}
type CheckoutAttributes struct {
	bookID       BoID
	memberID     MeID
	bookName     string
	memberName   string
	checkOutTime string
	checkInTime  string
}
type Members map[MeID]MemberAttributes
type Books map[BoID]BookAttributes
type CheckOutTransactions map[ChID]CheckoutAttributes
type CheckOutView map[ChID]CheckoutAttributes
type LibraryInstance struct {
	bookIds              BookIDs
	memberIds            MemberIDs
	checkoutIDs          CheckoutIDs
	members              Members
	books                Books
	checkOutTransactions CheckOutTransactions
	checkOutView         CheckOutView
}

func consoleStrInput(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	strVar, _ := reader.ReadString('\n')
	strVar = strings.ReplaceAll(strVar, "\n", "")
	return strVar
}

func consoleIntInput(prompt string) int {
	var intVar int
	fmt.Print(prompt)
	_, err := fmt.Scanf("%d", &intVar)
	if err != nil {
		return 0
	}
	return intVar
}

func addMember(library *LibraryInstance, memberAttributes *MemberAttributes) {
	var memberId MeID
	numMembers := len(library.memberIds)
	if numMembers == 0 {
		memberId = 1
	} else {
		memberId = library.memberIds[numMembers-1] + 1
	}
	library.memberIds = append(library.memberIds, memberId)
	library.members[memberId] = *memberAttributes
}

func addMemberConsole(library *LibraryInstance) {
	fmt.Println("Enter Member's Details Below: ")
	var memberAttributes MemberAttributes
	memberAttributes.firstName = consoleStrInput("First Name: ")
	memberAttributes.lastName = consoleStrInput("Second Name: ")
	memberAttributes.age = consoleIntInput("Age: ")
	memberAttributes.address = consoleStrInput("Address: ")
	addMember(library, &memberAttributes)
}

func addMemberArgs(library *LibraryInstance, attributes map[string]string) {
	age, err := strconv.Atoi(attributes["age"])
	if err != nil {
		fmt.Println("Failed to convert age into integer: ", err)
	}
	var memberAttributes MemberAttributes
	memberAttributes.firstName = attributes["firstName"]
	memberAttributes.lastName = attributes["lastName"]
	memberAttributes.address = attributes["address"]
	memberAttributes.age = age
	addMember(library, &memberAttributes)
}

func addBook(library *LibraryInstance, bookAttributes *BookAttributes) {
	bookAttributes.copies = 1
	bookAttributes.available = true
	bookAttributes.lender = 0

	var bookId BoID
	numBooks := len(library.bookIds)
	if numBooks == 0 {
		bookId = 1
	} else {
		bookId = library.bookIds[numBooks-1] + 1
	}
	library.bookIds = append(library.bookIds, bookId)
	library.books[bookId] = *bookAttributes
}

func addBookConsole(library *LibraryInstance) {
	fmt.Println("Enter Book's Details Below: ")
	var bookAttributes BookAttributes
	bookAttributes.bookName = consoleStrInput("Book Name: ")
	bookAttributes.author = consoleStrInput("Author: ")
	bookAttributes.publicationYear = consoleIntInput("Publication Year: ")
	addBook(library, &bookAttributes)
}

func addBookArgs(library *LibraryInstance, attributes map[string]string) {
	publicationYear, err := strconv.Atoi(attributes["publicationYear"])
	if err != nil {
		fmt.Println("Failed to convert age into integer: ", err)
	}
	var bookAttributes BookAttributes
	bookAttributes.bookName = attributes["bookName"]
	bookAttributes.author = attributes["author"]
	bookAttributes.publicationYear = publicationYear
	addBook(library, &bookAttributes)
}

func bookCheckOut(library *LibraryInstance) {
	bookId := BoID(consoleIntInput("Book ID: "))
	memberId := MeID(consoleIntInput("Member ID: "))
	book := library.books[bookId]
	member := library.members[memberId]
	numberCheckouts := len(library.checkoutIDs)
	var checkoutId ChID
	if numberCheckouts == 0 {
		checkoutId = 1
	} else {
		checkoutId = library.checkoutIDs[numberCheckouts-1] + 1
	}
	// updating checkout ids
	library.checkoutIDs = append(library.checkoutIDs, checkoutId)
	// creating checkout
	var checkoutAttributes CheckoutAttributes
	checkoutAttributes.bookID = bookId
	checkoutAttributes.memberID = memberId
	checkoutAttributes.bookName = book.bookName
	checkoutAttributes.memberName = member.firstName + " " + member.lastName
	checkoutAttributes.checkOutTime = time.Now().Truncate(24 * time.Hour).String()
	// updating checkout tables
	library.checkOutTransactions[checkoutId] = checkoutAttributes
	library.checkOutView[checkoutId] = checkoutAttributes
	// updating member and book table
	member.booksLent = append(member.booksLent, checkoutId)
	library.members[memberId] = member
	book.available = false
	book.lender = memberId
	library.books[bookId] = book
}

func removeCheckoutItem(itemSlice []ChID, item ChID) []ChID {
	for i, v := range itemSlice {
		if v == item {
			itemSlice = append(itemSlice[:i], itemSlice[i+1:]...)
			break
		}
	}
	return itemSlice
}

func bookCheckIn(library *LibraryInstance) {
	checkoutId := ChID(consoleIntInput("CheckOut ID: "))
	if checkoutTransaction, ok := library.checkOutView[checkoutId]; ok {
		book := library.books[checkoutTransaction.bookID]
		member := library.members[checkoutTransaction.memberID]
		// updated book availability
		book.available = true
		book.lender = 0
		library.books[checkoutTransaction.bookID] = book
		// updated member
		//member.booksLent = member.booksLent[:0]
		member.booksLent = removeCheckoutItem(member.booksLent, checkoutId)
		library.members[checkoutTransaction.memberID] = member
		// add checkin time to transaction
		checkoutTransaction.checkInTime = time.Now().Truncate(24 * time.Hour).String()
		library.checkOutTransactions[checkoutId] = checkoutTransaction
		// remove the checkout transaction from view table
		delete(library.checkOutView, checkoutId)
	} else {
		fmt.Println("Invalid CheckoutId: ", checkoutId)
		return
	}
}

func addBooksMembers(library *LibraryInstance) {
	book1 := map[string]string{
		"bookName":        "A Song of Ice & Fire",
		"author":          "George RR Martin",
		"publicationYear": "1998",
	}
	book2 := map[string]string{
		"bookName":        "A Brief History of Time",
		"author":          "Stephen Hawking",
		"publicationYear": "1989",
	}
	book3 := map[string]string{
		"bookName":        "Machine Learning System Design",
		"author":          "Chip Huyen",
		"publicationYear": "2022",
	}
	book4 := map[string]string{
		"bookName":        "Economic Hit Man",
		"author":          "XYZ",
		"publicationYear": "2005",
	}
	addBookArgs(library, book1)
	addBookArgs(library, book2)
	addBookArgs(library, book3)
	addBookArgs(library, book4)
	member1 := map[string]string{
		"firstName": "Irfan",
		"lastName":  "Danish",
		"age":       "26",
		"address":   "Piplan, Mianwali",
	}
	member2 := map[string]string{
		"firstName": "Talha",
		"lastName":  "Zaheer",
		"age":       "26",
		"address":   "DG Khan",
	}
	member3 := map[string]string{
		"firstName": "Hammad",
		"lastName":  "Munir",
		"age":       "27",
		"address":   "Lahore",
	}
	addMemberArgs(library, member1)
	addMemberArgs(library, member2)
	addMemberArgs(library, member3)
}

func logMembersTable(library *LibraryInstance) {
	// using tabwriter package to print formatted table of Books
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.TabIndent)
	_, err := fmt.Fprintln(tw, "\nMemberID\tFirst Name\tLast Name\tAge\tAddress\tNo. Books Lent")
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	for _, memberId := range library.memberIds {
		member := library.members[memberId]
		//converting age and memberId to strings
		age := strconv.Itoa(member.age)
		memberIdStr := strconv.Itoa(int(memberId))
		booksLent := strconv.Itoa(len(member.booksLent))
		prompt := memberIdStr + "\t" + member.firstName + "\t" + member.lastName + "\t" + age + "\t" + member.address + "\t" + booksLent
		_, err := fmt.Fprintln(tw, prompt)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
	}
	errFlush := tw.Flush()
	if errFlush != nil {
		return
	}
}

func logBooksTable(library *LibraryInstance) {
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.TabIndent)
	_, err := fmt.Fprintln(tw, "\nBookID\tBook\tAuthor\tPublication Year\tCount\tAvailable")
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	for _, bookId := range library.bookIds {
		book := library.books[bookId]
		publicationYear := strconv.Itoa(book.publicationYear)
		count := strconv.Itoa(book.copies)
		bookIdStr := strconv.Itoa(int(bookId))
		available := strconv.FormatBool(book.available)
		prompt := bookIdStr + "\t" + book.bookName + "\t" + book.author + "\t" + publicationYear + "\t" + count + "\t" + available
		_, err := fmt.Fprintln(tw, prompt)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
	}
	errFlush := tw.Flush()
	if err != nil {
		fmt.Println("Error", errFlush)
		return
	}
}

func logCheckoutTableContents(checkoutTablePtr *CheckOutTransactions, library *LibraryInstance) {
	checkoutTable := *checkoutTablePtr
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.TabIndent)
	_, err := fmt.Fprintln(tw, "\nCheckOut ID\tBook ID\tBook Name\tMember ID\tMember Name\tCheckOut Time\tCheckIn Time")
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	for _, checkoutId := range library.checkoutIDs {
		if checkoutTransaction, ok := checkoutTable[checkoutId]; ok {
			checkoutIdStr := strconv.Itoa(int(checkoutId))
			bookID := strconv.Itoa(int(checkoutTransaction.bookID))
			memberID := strconv.Itoa(int(checkoutTransaction.memberID))
			part1 := checkoutIdStr + "\t" + bookID + "\t" + checkoutTransaction.bookName + "\t" + memberID + "\t"
			part2 := checkoutTransaction.memberName + "\t" + checkoutTransaction.checkOutTime
			prompt := part1 + part2 + "\t" + checkoutTransaction.checkInTime
			_, err := fmt.Fprintln(tw, prompt)
			if err != nil {
				fmt.Println("Error: ", err)
				return
			}
		}
	}
	errFlush := tw.Flush()
	if errFlush != nil {
		fmt.Println("Error: ", errFlush)
		return
	}
}

func logCheckoutTables(library *LibraryInstance, transactions bool) {
	if transactions {
		fmt.Printf("\nCheckout Transactions Table: \n")
		checkoutTable := library.checkOutTransactions
		logCheckoutTableContents(&checkoutTable, library)
	} else {
		fmt.Printf("\nCurrent Checkout View Table: \n")
		checkoutTable := CheckOutTransactions(library.checkOutView)
		logCheckoutTableContents(&checkoutTable, library)
	}
}

func menu() int {
	fmt.Printf("\n====================================================\n")
	fmt.Printf("\t\t\tMenu\n")
	fmt.Println("====================================================")
	fmt.Println("1. Add Member")
	fmt.Println("2. Add Book")
	fmt.Println("3. Check Out A Book")
	fmt.Println("4. Return A Book")
	fmt.Println("5. View Members")
	fmt.Println("6. View Books")
	fmt.Println("7. View Check Out Transactions History")
	fmt.Println("8. View Books Currently in Check Out")
	fmt.Println("9. Exit")
	fmt.Printf("====================================================\n")
	choice := consoleIntInput("Enter Your Choice: ")
	fmt.Printf("\n====================================================\n")
	return choice
}

func main() {
	//	instantiate library
	var library LibraryInstance
	// initializing maps in the struct, to avoid following RunTime Error:
	// panic: assignment to entry in nil map
	library.members = make(Members)
	library.books = make(Books)
	library.checkOutTransactions = make(CheckOutTransactions)
	library.checkOutView = make(CheckOutView)

	for {
		choice := menu()
		switch choice {
		case 0:

		case 1:
			addMemberConsole(&library)
		case 2:
			addBookConsole(&library)
		case 3:
			bookCheckOut(&library)
		case 4:
			bookCheckIn(&library)
		case 5:
			logMembersTable(&library)
		case 6:
			logBooksTable(&library)
		case 7:
			logCheckoutTables(&library, true)
		case 8:
			logCheckoutTables(&library, false)
		case 9:
			os.Exit(0)
		default:
			fmt.Printf("\nWrong Choice!!!\n")
		}
	}
}
