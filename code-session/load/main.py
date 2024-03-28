import datetime


class Loan:
    def __init__(self, amount: float = 0, deadline: datetime.date | None = None):
        # load identification which will be filled by the database.
        self.id = 0
        self.amount: float = amount
        self.deadline: datetime.date | None = deadline
        self.paid: bool = False


class Customer:
    def __init__(self):
        self.loans: list[Loan] = []
        self.credit: float = 0

    def loan(self, loan: Loan):
        self.loans.append(loan)

    def pay(self, payment: float):
        nearest_loan: Loan | None = None
        nearest_deathline = datetime.date(year=1, month=1, day=1)
        for loan in self.loans:
            if (
                loan.deadline is not None
                and (loan.deadline > nearest_deathline)
                and not loan.paid
            ):
                nearest_deathline = loan.deadline
                nearest_loan = loan

        if nearest_loan is None:
            self.credit += payment
            return

        if payment < nearest_loan.amount:
            nearest_loan.amount -= payment
        elif payment >= nearest_loan.amount:
            payment -= nearest_loan.amount
            nearest_loan.amount = 0
            nearest_loan.paid = True

            self.pay(payment)


if __name__ == "__main__":
    customer = Customer()
    customer.loan(
        Loan(
            amount=100,
            deadline=datetime.date.today() + datetime.timedelta(minutes=100),
        )
    )
    customer.loan(
        Loan(
            amount=100,
            deadline=datetime.date.today() + datetime.timedelta(minutes=1000),
        )
    )
    customer.loan(
        Loan(
            amount=100,
            deadline=datetime.date.today() + datetime.timedelta(minutes=200),
        )
    )

    customer.pay(payment=400)

    print("consumer loans after paying $400:")
    for i, loan in enumerate(customer.loans):
        print(f"{i}. ${loan.amount} {loan.deadline}")

    print(f"remaining cusmter credit: ${customer.credit}")
