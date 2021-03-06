package card

import (
	"testing"

	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/customer"
	. "github.com/stripe/stripe-go/utils"
)

func init() {
	stripe.Key = GetTestKey()
}

func TestCardNew(t *testing.T) {
	customerParams := &stripe.CustomerParams{}
	customerParams.SetSource("tok_amex")

	cust, _ := customer.New(customerParams)

	cardParams := &stripe.CardParams{
		Customer: cust.ID,
		Token:    "tok_visa",
	}

	target, err := New(cardParams)

	if err != nil {
		t.Error(err)
	}

	if target.LastFour != "4242" {
		t.Errorf("Unexpected last four %q for card number %v\n", target.LastFour, cardParams.Number)
	}

	if target.Meta == nil || len(target.Meta) > 0 {
		t.Errorf("Unexpected nil or non-empty metadata in card\n")
	}

	targetCust, err := customer.Get(cust.ID, nil)

	if err != nil {
		t.Error(err)
	}

	if targetCust.Sources.Count != 2 {
		t.Errorf("Unexpected number of sources %v\n", targetCust.Sources.Count)
	}

	targetCard, err := New(&stripe.CardParams{
		Customer: targetCust.ID,
		Token:    "tok_visa",
	})

	if targetCard.LastFour != "4242" {
		t.Errorf("Unexpected last four %q for card number %v\n", targetCard.LastFour, cardParams.Number)
	}

	customer.Del(cust.ID)
}

func TestCardGet(t *testing.T) {
	customerParams := &stripe.CustomerParams{}
	customerParams.SetSource("tok_visa")

	cust, _ := customer.New(customerParams)

	target, err := Get(cust.DefaultSource.ID, &stripe.CardParams{Customer: cust.ID})

	if err != nil {
		t.Error(err)
	}

	if target.LastFour != "4242" {
		t.Errorf("Unexpected last four %q for card number %v\n", target.LastFour, customerParams.Source.Card.Number)
	}

	if target.Brand != Visa {
		t.Errorf("Card brand %q does not match expected value\n", target.Brand)
	}

	if target.Funding != Credit {
		t.Errorf("Card funding %q does not match expected value\n", target.Funding)
	}

	customer.Del(cust.ID)
}

func TestCardDel(t *testing.T) {
	customerParams := &stripe.CustomerParams{}
	customerParams.SetSource("tok_visa")

	cust, _ := customer.New(customerParams)

	cardDel, err := Del(cust.DefaultSource.ID, &stripe.CardParams{Customer: cust.ID})
	if err != nil {
		t.Error(err)
	}

	if !cardDel.Deleted {
		t.Errorf("Card id %q expected to be marked as deleted on the returned resource\n", cardDel.ID)
	}

	customer.Del(cust.ID)
}

func TestCardUpdate(t *testing.T) {
	customerParams := &stripe.CustomerParams{}
	customerParams.SetSource("tok_amex")

	cust, err := customer.New(customerParams)

	if err != nil {
		t.Error(err)
	}

	cardParams := &stripe.CardParams{
		Customer: cust.ID,
		Name:     "Updated Name",
		Month:    "10",
		Year:     "21",
	}

	target, err := Update(cust.DefaultSource.ID, cardParams)

	if err != nil {
		t.Error(err)
	}

	if target.Name != cardParams.Name {
		t.Errorf("Card name %q does not match expected name %q\n", target.Name, cardParams.Name)
	}

	if target.Month != 10 {
		t.Errorf("Unexpected expiration month %d for card where we set %q\n", target.Month, cardParams.Month)
	}

	if target.Year != 2021 {
		t.Errorf("Unexpected expiration year %d for card where we set %q\n", target.Year, cardParams.Year)
	}

	customer.Del(cust.ID)
}

func TestCardList(t *testing.T) {
	customerParams := &stripe.CustomerParams{}
	customerParams.SetSource("tok_amex")

	cust, _ := customer.New(customerParams)

	card := &stripe.CardParams{
		Customer: cust.ID,
		Token:    "tok_visa",
	}

	New(card)

	i := List(&stripe.CardListParams{Customer: cust.ID})
	for i.Next() {
		if i.Card() == nil {
			t.Error("No nil values expected")
		}

		if i.Meta() == nil {
			t.Error("No metadata returned")
		}
	}
	if err := i.Err(); err != nil {
		t.Error(err)
	}

	customer.Del(cust.ID)
}
