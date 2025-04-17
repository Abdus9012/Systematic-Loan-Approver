package main

import (
	"fmt"

	"github.com/jung-kurt/gofpdf"
)

type Loan struct {
	FullName    string
	Age         int
	Income      float64
	LoanAmount  float64
	LoanType    string
	RepayOption float64
}

const (
	RetirementAge     = 62
	HouseLoanRate     = 8.75
	PersonalLoanRate  = 10.85
	EducationLoanRate = 8.25
)

func main() {
	var loan Loan
	var option int
	fmt.Println("Enter full name w/o spaces: ")
	fmt.Scanln(&loan.FullName)
	fmt.Println("Enter legal age: ")
	fmt.Scanln(&loan.Age)
	fmt.Println("Enter Inhand Salary per month as in PaySlips: ")
	fmt.Scanln(&loan.Income) //30000
	fmt.Println("Enter loan amount")
	fmt.Scanln(&loan.LoanAmount) //300000
	fmt.Println("Select loan type (house, personal, education)")
	fmt.Scanln(&loan.LoanType)
	fmt.Println("Choose Repayment Option: 1 for 10%, 2 for 20%, 3 for 25% of monthly income")
	fmt.Scanln(&option)

	switch option {
	case 1:
		loan.RepayOption = loan.Income * 0.10 //3000
	case 2:
		loan.RepayOption = loan.Income * 0.20
	case 3:
		loan.RepayOption = loan.Income * 0.25

	default:
		fmt.Println("Invalid repayment method")
		return
	}

	approved, totalPayable, numMonths, monthlyInstallemnt := loan.CalculateLoan() //here no parameters are required inside () because loan is an object which already captures all the details required for approving the loan and CalculateLoan() is designed in such a way that it can access all the details inside loan object
	GeneratePDF(loan, approved, totalPayable, numMonths, monthlyInstallemnt)      //this function is different and requires user input everytime it is called.

	if approved {
		fmt.Println("Loan Approved! Check loan_result.pdf")
	} else {
		fmt.Println("Loan rejected! Check loan_result.pdf")
	}

}

func (l Loan) CalculateLoan() (bool, float64, int, float64) { //this is a method, not just a normal function. This function "completly" works and returns only for loan object hence that "l" is mentioned
	remainingYears := RetirementAge - l.Age //32
	if remainingYears <= 0 {
		return false, 0, 0, 0
	}

	var interestRate float64
	switch l.LoanType {
	case "house":
		interestRate = HouseLoanRate
	case "personal":
		interestRate = PersonalLoanRate
	case "education":
		interestRate = EducationLoanRate

	default:
		return false, 0, 0, 0
	}

	totalPayable := l.LoanAmount + (l.LoanAmount * (interestRate / 100) * float64(remainingYears))
	numMonths := int(totalPayable / l.RepayOption)
	monthlyInstallment := l.RepayOption

	if numMonths <= remainingYears*12 {
		return true, totalPayable, numMonths, monthlyInstallment
	}
	return false, totalPayable, numMonths, monthlyInstallment
}

func GeneratePDF(loan Loan, approved bool, totalPayable float64, numMonths int, monthlyInstallment float64) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Loan Application Result") //Cell-> for printing
	pdf.Ln(12)                                  //moves the cursor down for printing

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, fmt.Sprintf("Full Name: %s", loan.FullName))
	pdf.Ln(10)

	if approved {
		pdf.Cell(40, 10, "Loan Approved")
		pdf.Ln(10)
		pdf.Cell(40, 10, fmt.Sprintf("Total Payable Amount: %.2f", totalPayable))
		pdf.Ln(10)
		pdf.Cell(40, 10, fmt.Sprintf("Monthly Installment: %.2f", monthlyInstallment))
		pdf.Ln(10)
	} else {
		pdf.Cell(40, 10, "Loan Rejected")
		pdf.Ln(10)
		pdf.Cell(40, 10, fmt.Sprintf("Total Payable Amount Needed: %.2f", totalPayable))
		pdf.Ln(10)
		pdf.Cell(40, 10, fmt.Sprintf("Months required: %d", numMonths))
		pdf.Ln(10)
	}

	err := pdf.OutputFileAndClose("loan_result.pdf")
	if err != nil {
		fmt.Println("Error generating PDF: ", err)
	}

}
