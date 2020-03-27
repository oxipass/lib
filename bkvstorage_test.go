package bslib

import (
	"errors"
	"fmt"
	"log"
	"os"
	"testing"
)

var globalDBFile string
var dbPass = ""

func TestLogin(t *testing.T) {
	// Getting storage instance
	bsInstance := GetInstance()

	errPass := bsInstance.Unlock(dbPass)
	if errPass != nil {
		t.Error(errPass)
		t.FailNow()
		return
	}
}

func TestDeleteItem(t *testing.T) {
	bsInstance := GetInstance()
	errPass := bsInstance.Unlock(dbPass)
	if errPass != nil {
		t.Error(errPass)
		t.FailNow()
		return
	}
	itemName := generateRandomString(20)
	response, err := bsInstance.AddNewItem(
		UpdateItemForm{
			BSItem: BSItem{
				Name: itemName,
				Icon: "fas fa-ambulance"},
		},
	)
	if err != nil {
		t.Error(err)
		t.FailNow()
		return
	}
	if response.Status != ConstSuccessResponse {
		t.Error(errors.New("response is not successful"))
		t.FailNow()
		return
	}
	errLock := bsInstance.Lock()
	if errLock != nil {
		t.Error(errLock)
		t.FailNow()
		return
	}

	errPass = bsInstance.Unlock(dbPass)
	if errPass != nil {
		t.Error(errPass)
		t.FailNow()
		return
	}

	delResponse, errDelete := bsInstance.DeleteItem(
		UpdateItemForm{
			BSItem: BSItem{
				ID: response.ItemID,
			},
		},
	)
	if errDelete != nil {
		t.Error(errPass)
		t.FailNow()
		return
	}
	if delResponse.Status != ConstSuccessResponse {
		t.Error(errors.New("deletion response is not successful"))
		t.FailNow()
		return
	}

	errLock = bsInstance.Lock()
	if errLock != nil {
		t.Error(errLock)
		t.FailNow()
		return
	}

	errPass = bsInstance.Unlock(dbPass)
	if errPass != nil {
		t.Error(errPass)
		t.FailNow()
		return
	}

	items, respErr := bsInstance.ReadAllItems(false)
	if respErr != nil {
		t.Error(respErr)
		t.FailNow()
		return
	}

	for _, item := range items {
		if item.ID == response.ItemID && item.Deleted == false {
			t.Error(errors.New("deleted item is found as not deleted"))
			t.FailNow()
			return
		}
	}

}

func TestAddItemWithNonExistingIcon(t *testing.T) {
	bsInstance := GetInstance()
	errPass := bsInstance.Unlock(dbPass)
	if errPass != nil {
		t.Error(errPass)
		t.FailNow()
		return
	}
	itemName := generateRandomString(20)
	response, err := bsInstance.AddNewItem(
		UpdateItemForm{
			BSItem: BSItem{
				Name: itemName,
				Icon: "ewfdwejfnerfkj",
			},
		},
	)
	if err == nil {
		t.Error(errors.New("no error returned in spite of the fact that icon is not existing"))
		t.FailNow()
		return
	}
	if response.Status == ConstSuccessResponse {
		t.Error(errors.New("item was added in spite of the fact that icon is not existing"))
		t.FailNow()
		return
	}
}

func TestAddField(t *testing.T) {
	bsInstance := GetInstance()
	errPass := bsInstance.Unlock(dbPass)
	if errPass != nil {
		t.Error(errPass)
		t.FailNow()
		return
	}
	itemName := generateRandomString(20)
	response, err := bsInstance.AddNewItem(
		UpdateItemForm{
			BSItem: BSItem{
				Name: itemName,
				Icon: "fas fa-ambulance",
			},
		},
	)
	if err != nil {
		t.Error(err)
		t.FailNow()
		return
	}
	if response.Status != ConstSuccessResponse {
		t.Error(errors.New("response is not successfull"))
		t.FailNow()
		return
	}

	fieldName := generateRandomString(20)
	fieldValue := generateRandomString(20)
	fieldResult, errField := bsInstance.AddNewField(
		UpdateFieldForm{
			ItemID: response.ItemID,
			BSField: BSField{
				Name:      fieldName,
				Icon:      "fab fa-amazon",
				ValueType: VTText,
				Value:     fieldValue,
			},
		},
	)
	if errField != nil {
		t.Error(errField)
		t.FailNow()
		return
	}
	if fieldResult.Status != ConstSuccessResponse {
		t.Error(errors.New("response is not successful: " + fieldResult.MsgTxt))
		t.FailNow()
		return
	}
	errLock := bsInstance.Lock()
	if errLock != nil {
		t.Error(errLock)
		t.FailNow()
		return
	}

	errPass = bsInstance.Unlock(dbPass)
	if errPass != nil {
		t.Error(errPass)
		t.FailNow()
		return
	}

	fields, errFields := bsInstance.ReadFieldsByItemID(response.ItemID)
	if errFields != nil {
		t.Error(errPass)
		t.FailNow()
		return
	}

	for _, field := range fields {
		if field.ID == fieldResult.FieldID {
			if field.Icon == "fab fa-amazon" &&
				field.Name == fieldName &&
				field.Value == fieldValue {
				return
			}
			t.Error("ID is found but content is wrong")
			t.FailNow()
			return
		}
	}
	t.Error("Saved ID is not found")
	t.FailNow()
}

func TestAddItem(t *testing.T) {
	bsInstance := GetInstance()
	errPass := bsInstance.Unlock(dbPass)
	if errPass != nil {
		t.Error(errPass)
		t.FailNow()
		return
	}
	itemName := generateRandomString(20)
	response, err := bsInstance.AddNewItem(
		UpdateItemForm{
			BSItem: BSItem{
				Name: itemName,
				Icon: "fas fa-ambulance",
			},
		},
	)
	if err != nil {
		t.Error(err)
		t.FailNow()
		return
	}
	if response.Status != ConstSuccessResponse {
		t.Error(errors.New("response is not successfull"))
		t.FailNow()
		return
	}
	errLock := bsInstance.Lock()
	if errLock != nil {
		t.Error(errLock)
		t.FailNow()
		return
	}

	errPass = bsInstance.Unlock(dbPass)
	if errPass != nil {
		t.Error(errPass)
		t.FailNow()
		return
	}
	items, respErr := bsInstance.ReadAllItems(false)
	if respErr != nil {
		t.Error(respErr)
		t.FailNow()
		return
	}
	for _, item := range items {
		if item.ID == response.ItemID && item.Name == itemName {
			return
		}
	}
	t.Error(errors.New("created item is not found"))
	t.FailNow()

}

func TestMain(m *testing.M) {
	errSetup := setup()

	if errSetup != nil {
		fmt.Printf("Setup failure: %s\n", errSetup.Error())
		os.Exit(1)
		return
	}

	retCode := m.Run()

	errFinish := teardown()
	if errFinish != nil {
		os.Exit(2)
		return
	}

	os.Exit(retCode)
}

func setup() error {
	// Getting package instance
	log.Println("Tests initial setup")

	bsInstance := GetInstance()
	if bsInstance == nil {
		return errors.New("cannot retrieve BSLib instance")
	}
	log.Println("BSLib instance successfully retrieved")

	dbFile := getTestDbFileName()
	log.Println("Full path to database file: " + dbFile)

	if _, err := os.Stat(dbFile); err == nil {
		err = os.Remove(dbFile)
		if err != nil {
			return err
		}
	}

	errOpen := bsInstance.Open(dbFile)
	if errOpen != nil {
		return errOpen
	}
	log.Println("BSLib is initiated successfully")

	// Generating random password for the database
	dbPass = generateRandomString(12)
	log.Println("Generated DB password: " + dbPass)

	errSetPassword := bsInstance.SetNewPassword(dbPass, "AES256V1")
	if errSetPassword != nil {
		return errSetPassword
	}
	log.Println("Password is successfully set")

	errLock := bsInstance.Lock()
	if errLock != nil {
		return errLock
	}
	log.Println("Database is locked, the tests can be started")

	return nil
}

func teardown() error {
	bsInstance := GetInstance()
	err := bsInstance.Close()
	if err != nil {
		return nil
	}
	if _, err := os.Stat(getTestDbFileName()); err == nil {
		// Do not delete if local file is used (use build tag 'local' to configure)
		// Check bkv_config_global.go for default configuration
		if !isLocalFileUsed() {
			err = os.Remove(getTestDbFileName())
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func isLocalFileUsed() bool {
	return getTestDbFileName() == localTestFile
}

func getTestDbFileName() string {
	if globalDBFile != "" {
		return globalDBFile
	}
	if useLocalTestFile {
		globalDBFile = localTestFile
		return globalDBFile
	}

	globalDBFile = generateTempFilename()

	return globalDBFile
}
