package samples

import (
	"strings"
	"testing"
)

func TestParseCsv(t *testing.T) {
	csv := "first_name,last_name,externalExtension,internalExtension,employee,internalPhone\n" +
		"John,Doe,4711,10,true,PHONE_12\n" +
		"Jane,Fox,0815,12,false,PHONE_13\n" +
		"Judith,Queston,2252,13,no_boolean,another_phone"
	reader := strings.NewReader(csv)
	contacts, err := ParseCsv(reader, true)
	if err != nil {
		t.Errorf("Unexpected error thrown: %v", err)
	}
	expectedNames := []string{"John", "Jane", "Judith"}
	expectedLastNames := []string{"Doe", "Fox", "Queston"}
	expectedExternal := []string{"4711", "0815", "2252"}
	expectedInternal := []string{"10", "12", "13"}
	expectedEmployee := []bool{true, false, false}
	expectedPhones := []string{"PHONE_12", "PHONE_13", "another_phone"}
	for index, contact := range contacts {
		if contact.firstName != expectedNames[index] {
			t.Errorf("First Name differs for record %d: Expected %s, got %s", index, expectedNames[index], contact.firstName)
		}
		if contact.lastName != expectedLastNames[index] {
			t.Errorf("Last Name differs for record %d: Expected %s, got %s", index, expectedLastNames[index], contact.lastName)
		}
		if contact.externalExtension != expectedExternal[index] {
			t.Errorf("External extension differs for record %d: Expected %s, got %s", index, expectedExternal[index], contact.externalExtension)
		}
		if contact.internalExtension != expectedInternal[index] {
			t.Errorf("Internal extension differs for record %d: Expected %s, got %s", index, expectedInternal[index], contact.internalExtension)
		}
		if contact.isEmployee != expectedEmployee[index] {
			t.Errorf("isEmployee differs for record %d: Expected %t, got %t", index, expectedEmployee[index], contact.isEmployee)
		}
		if contact.internalPhone != expectedPhones[index] {
			t.Errorf("Phone name differs for record %d: Expected %s, got %s", index, expectedPhones[index], contact.internalPhone)
		}
	}
	csv = csv + "\none,colomn,too,many,for,the,parser"
	_, err = ParseCsv(strings.NewReader(csv), true)
	if err == nil {
		t.Errorf("Expected error when extension of columns is wrong, but no error was returned")
	}
}
