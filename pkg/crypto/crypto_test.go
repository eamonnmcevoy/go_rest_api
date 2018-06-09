package crypto_test

import (
	"go_web_server/pkg/crypto"
	"testing"
)

func Test_Hash(t *testing.T) {
	t.Run("Can hash and compare", should_be_able_to_hash_and_compare_strings)
	t.Run("Can detect unequal hashes", should_return_error_when_comparing_unequal_hashes)
	t.Run("Generates a different salt every time", should_generate_a_different_salt_each_time)
}

func should_be_able_to_hash_and_compare_strings(t *testing.T) {
	//Arrange
	c := crypto.Hash{}
	testInput := "testInput"

	//Act
	generatedHash, generateError := c.Generate(testInput)
	compareError := c.Compare(generatedHash, testInput)

	//Assert
	if generateError != nil {
		t.Error("Error generating hash")
	}
	if testInput == generatedHash {
		t.Error("Generated hash is the same as input")
	}
	if compareError != nil {
		t.Error("Error comparing hash to input")
	}
}

func should_return_error_when_comparing_unequal_hashes(t *testing.T) {
	//Arrange
	c := crypto.Hash{}
	testInput := "testInput"
	testCompare := "testCompare"

	//Act
	generatedHash, generateError := c.Generate(testInput)
	compareError := c.Compare(generatedHash, testCompare)

	//Assert
	if generateError != nil {
		t.Error("Error generating hash")
	}
	if testInput == generatedHash {
		t.Error("Generated hash is the same as input")
	}
	if compareError == nil {
		t.Error("Compare should not have been successful")
	}
}

func should_generate_a_different_salt_each_time(t *testing.T) {
	//Arrange
	c := crypto.Hash{}
	testInput := "testInput"

	hash1, _ := c.Generate(testInput)
	hash2, _ := c.Generate(testInput)

	if hash1 == hash2 {
		t.Error("Subsequent hashes should not be equal")
	}
}
