package receiptstructs

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

/*
  - @brief this function takes in a string and returns two things.
    First is the total number of alphanumeric characters/runes
    Second is not alphanumeric characters/runes.
    @Param itemName this is the string that will be scanned and told how
    many points an invalid characters are in this string
    @Return int32 is the total number of letters and numbers in the given string
    []rune is an array that holds all the invalid characters
*/
func pointForName(itemName string) (int32, []rune) {
	var count int32 = 0
	var invalidChar []rune
	for _, r := range itemName {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			count += 1
		} else if r != ' ' {
			invalidChar = append(invalidChar, r)
		}
	}
	return count, invalidChar
}

/*
@brief Takes in a string and returns a float to represent a dollar amount as a
float 64. Throws an error of -1.0 if the string could not be converted or if
the number converted is negative

@Param dAmount: this is a string that represents a number int a dollar amount

@return returns a positive float64 if the parse was successful. if Either the
parse was not successful due to formatting errors or is negative will return -1
*/
func stringToDollar(dAmount string) float64 {
	tempConv, er := strconv.ParseFloat(dAmount, 64)
	if er != nil || tempConv < 0.0 {
		return -1.0
	}
	return tempConv
}

/*
@brief This function takes a receipt and calculates the points as well as a
string that is the bulk of the breakdown

This function will calculate the total points that a receipt should earn given
the documentation.

This function assumes that the total calculated when the json was created is
correct. otherwise one would have to assume the amount of taxes for each
state (I know that this is a simplified example of how a product like this could
be used but I just wanted to make that clear)
*/
func ReceiptRewards(r Receipt) (string, int32) {
	var receiptPoints int32 = 0
	var breakDown strings.Builder
	var totalDollarAmount float64 = stringToDollar(r.Total)

	//due to the nature of 1.00 being a multiple of 0.25 it only make since to
	// call a flag to increase computational efficiency
	var quarterFlag bool = false
	//switch tab to weight

	/*First it will convert the dollar amount from a string to a float64. Then
	 check if the value is a multiple of a dollar (meaning that there is no cents).
	If there are no cents then the receipt will gain 50 points as well as making a
	flag true due to boolean logic. Due to the next statement being an ||(or)
	if statement the compiler should check if that statement is true then check
	the second statement if and only if the first bool statement is true and
	since a dollar is a multiple of 0.25 cents the flag is ticked to true
	meaning that cent checker does not have to be called again
	*/
	if centChecker(totalDollarAmount, 100) {
		breakDown.WriteString("\t50 points - total is a round dollar amount\n")
		//TODO if the rule ever changes from a dollar to any value that is not
		//a multiple of 25 cents remove this bool
		quarterFlag = true
		receiptPoints += 50
	}

	//is cent amount is a multiple of 0.25 cents
	/*Second checks if a flag was called in the previous step.
	If not called then the centChecker method is called with the cents value at 25.
	If either is correct then 25 points are added to the point value
	*/
	if quarterFlag || centChecker(totalDollarAmount, 25) {
		breakDown.WriteString("\t25 points - total is a multiple of 0.25\n")
		receiptPoints += 25
	}

	/**
	Third the program will call the pointForName function given the retailer
	value from the receipt and returns the points and the valid characters/runes
	then adds the line for the breakdown. Then checks if there were any runes
	that are invalid. If there are any invalid characters/runes then loop
	through the list of invalid runes and add them to the break down
	*/
	var retailerNamePoints, invalidR = pointForName(r.Retailer)
	var dummyString string = fmt.Sprintf("\t%d points - retailer name (%s) has %d alphanumeric characters\n", retailerNamePoints, r.Retailer, retailerNamePoints)
	if len(invalidR) > 0 {
		dummyString += "\t\t\tNote: "
		for _, r := range invalidR {
			dummyString += fmt.Sprintf("'%c' is not alphanumeric\n", r)
		}
	}
	breakDown.WriteString(dummyString)
	receiptPoints += retailerNamePoints
	dummyString = ""

	/**
	Fourth the program will call the pointsBetweenTimes method and check if the
	receipt was printed between 2pm and 4pm. if it was then add 10 points to the
	total of points earned and add to the breakdown
	*/
	timeFlag, pointsBTstring := pointsBetweenTimes(r.PurchaseTime)
	if timeFlag {
		breakDown.WriteString(pointsBTstring)
		receiptPoints += 10
	}

	/**
	Fifth checks to see if the the amount of items is divisible by two. To do
	this we first take the length of the array and do integer division to
	return the total amount of pair of items (points for ever two items). If
	there is at least two items, then add the integer division of total items
	multiplied by 5 to the points. Then add this information to the breakdown
	*/
	var itemAmountPoints = int32(len(r.Items) / 2)
	if itemAmountPoints > 0 {
		receiptPoints += itemAmountPoints * 5
		dummyString = fmt.Sprintf("\t%d points - %d items (%d pair(s) @ 5 points each)\n", (itemAmountPoints * 5), (len(r.Items)), itemAmountPoints)
		breakDown.WriteString(dummyString)
		dummyString = ""
	}

	/**
	Sixth call the descriptionMultThree method. then check if the bool returned
	is true. If it is means that there was at least one item that had a
	description thats length was a multiple of 3. intern the program will add
	the items returned to the total points earned (the ceiling(rounding up) of
	the length of description * 0.2) and add to the breakdown.
	*/
	var disBool, disText, disPoints = descriptionMultThree(r)
	if disBool {
		breakDown.WriteString(disText)
		receiptPoints += disPoints
	}

	/**
	runs the purchasedOddDay method. if the receipt was printed on an odd day
	then add 6 points to the total points earned and add the correct information
	to the breakdown
	*/
	var isOdd, oddText = purchasedOddDay(r.PurchaseDate)
	if isOdd {
		breakDown.WriteString(oddText)
		receiptPoints += 6
	}

	/**
	return statement
	*/
	return breakDown.String(), receiptPoints
}

/*
@brief prints the breakdown so that it is a similar format as the fetch github examples

@param breakDown: this is the string that holds all the information of all the
rules in a written format. totalPoints: is the total points earned by each rule
to be printed out
*/
func PrintBreakDown(breakDown string, totalPoints int32) {
	fmt.Printf("Total Points: %d\nBreakdown:\n", totalPoints)
	fmt.Print(breakDown)
	fmt.Printf("  + ---------\n\t= %d points\n", totalPoints)
}

/*
@brief this method take a string assumed to be in the formate of "15:04" (24
hour digital clock [Hour:Minutes]). And returns true if the time is in between
the time of 2pm and 4pm

@param sTime: passes a string that will be parsed for a time variable

@return bool: wether the time passed in (a string) is in between the time of
2pm and 4pm. string: the breakdown if the time passed in was between 2pm and
4pm. if not then returns the empty string
*/
func pointsBetweenTimes(sTime string) (bool, string) {
	var timeValue time.Time
	var err error
	var dummyString = ""
	timeValue, err = time.Parse("15:04", sTime)
	if err != nil {
		log.Fatal("Time bought was unexpected ", err)
	} else if timeValue.Hour() > 13 && timeValue.Hour() < 16 { // time bought between 2pm and 4pm
		var timeOfDay string = "AM"
		if (timeValue.Hour() / 12) > 0 {
			timeOfDay = "PM"
		}
		var hourTime = (timeValue.Hour() % 12)
		/**
		// if the hour changes to be around noon or midnight
		if (timeValue.Hour() % 12 == 0){
			hourTime = 12
		}*/

		dummyString = fmt.Sprintf("\t%d points - %d:%02d %s is between 2:00pm and 4:00pm\n", 10, hourTime, timeValue.Minute(), timeOfDay)
		return true, dummyString
	}
	return false, dummyString
}

/*
@brief runs through all the items in the receipt and checks if the trimmed
length is a multiple of 3

@param is a receipt struct object  that will have its items checked on
their description length

@return bool: is if any of the items have a description length of mod 3 or not,
string: is the breakdown from all items with descriptions of lengths of mod 3,
int32: is the total amount of points earned from all item description
*/
func descriptionMultThree(r Receipt) (bool, string, int32) {
	var disLengthPos []int32
	var counter int32
	var dummyString string = ""
	var pointTotal int32 = 0
	for _, i := range r.Items {
		i.ShortDescription = strings.Trim(i.ShortDescription, " ")
		if (utf8.RuneCountInString(i.ShortDescription) % 3) == 0 {
			disLengthPos = append(disLengthPos, counter)

		}
		counter++
	}
	if counter > 0 {
		for _, i := range disLengthPos {
			var bRoundedPoints float64 = math.Ceil(stringToDollar(r.Items[i].Price) * 0.2)
			var disPoints int32 = int32(bRoundedPoints)
			pointTotal += disPoints
			r.Items[i].ShortDescription = strings.Trim(r.Items[i].ShortDescription, " ")
			dummyString += fmt.Sprintf("\t%d points - \"%s\" is %d characters (a multiple of 3)\n", disPoints, r.Items[i].ShortDescription, utf8.RuneCountInString(r.Items[i].ShortDescription))
			dummyString += fmt.Sprintf("\t\titem price of %s * 0.2 = %g rounded up is %d\n", r.Items[i].Price, bRoundedPoints, disPoints)

		}
		return true, dummyString, pointTotal
	}
	return false, dummyString, pointTotal
}

/*
@brief takes in a string to represent the time and if the day variable is odd
then return true as well as the breakdown string

Using the time function parse we use the assumed format of "2006-01-02" (Year,
Month, Day) then use the date value to see if it an odd day by using the
modulus operator.

@param dTime is a string that should be in the format of "2006-01-02"

@return bool is whether or not the the purchase day was odd, string is the
breakdown if the day was odd. is the empty string if not
*/
func purchasedOddDay(dTime string) (bool, string) {
	var dateValue, err = time.Parse("2006-01-02", dTime)
	var oddBool = false
	var dummyString string = ""
	if err != nil {
		log.Fatal("Date Was not entered as expected ", err)
	} else if (dateValue.Day() % 2) == 1 {
		oddBool = true
		dummyString = fmt.Sprintf("\t%d points - purchase day is odd\n", 6)
	}
	return oddBool, dummyString
}

/*
@brief This function takes in a float and an int. the float represents a dollar
amount and the int is meant to represent a multiple in cents. Intern the
function check to see if the float when multiplied by a 100 is a multiple of
the int

@param dollarAmount: represents an amount of money int the format of (1.25 a
dollar and 25 cents)

@param multipleInCents: represents the number as a mod for the cents to check
if they are a multiple of

@return a value if the float value is a multiple of int using the modulus operator
*/
func centChecker(dollarAmount float64, multipleInCents int) bool {
	return (int(dollarAmount*100))%(multipleInCents) == 0
}

/*
@brief this function take a string in and prints out
the receipt from the json from the input file location

If the file location is invalid or is not a json file the program will through
a log.Fatal saying that the JSON file could not be opened properly

@pram jsonFileLocation: a string that should be a json's file location.
*/
func PrintReceiptFJson(jsonFileLocation string) {
	fmt.Println("___________________________________")
	jsonFile, err := os.Open(jsonFileLocation)
	if err != nil {
		log.Fatal("JSON file could not be opened properly", err)
	}

	byteValue, _ := io.ReadAll(jsonFile)
	var itemsJ1 Receipt
	json.Unmarshal(byteValue, &itemsJ1)

	PrintBreakDown(ReceiptRewards(itemsJ1))
	fmt.Println("___________________________________")
	jsonFile.Close()
}
