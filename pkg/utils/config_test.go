package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSetLocation(t *testing.T) {
	home, err := os.UserHomeDir()
	require.NoError(t, err)
	SetLocation()
	require.Equal(t, home+"/.harbor", HarborFolder)
	require.Equal(t, home+"/.harbor/config.yaml", DefaultConfigPath)
}

func TestGetCurrentCredentialName(t *testing.T) {
	hc := HarborConfig{
		CurrentCredentialName: "TestCred",
	}

	require.Equal(t, "TestCred", hc.GetCurrentCredentialName())
}

func TestCreateConfigFile(t *testing.T) {
	SetLocation()
	err := CreateConfigFile()
	require.NoError(t, err)

	// Check if the file exists
	_, err = os.Stat(DefaultConfigPath)
	require.NoError(t, err)

	// Clean up
	err = os.Remove(DefaultConfigPath)
	require.NoError(t, err)

	err = os.Remove(HarborFolder)
	require.NoError(t, err)
}

func TestAddCredentialsToConfigFile(t *testing.T) {
	SetLocation()
	err := CreateConfigFile()
	require.NoError(t, err)

	credential := Credential{
		Name:          "TestCred",
		Username:      "TestUser",
		Password:      "TestPass",
		ServerAddress: "https://test.com",
	}

	err = AddCredentialsToConfigFile(credential, DefaultConfigPath)
	require.NoError(t, err)

	// Check if the credential is added
	credentialFile, err := os.ReadFile(DefaultConfigPath)
	require.NoError(t, err)

	// Ensure the content of the file is correct
	require.Equal(t, "credentials:\n    - name: TestCred\n      username: TestUser\n      password: TestPass\n      serveraddress: https://test.com\ncurrent-credential-name: TestCred\n", string(credentialFile))

	credential2 := Credential{
		Name:          "TestCred2",
		Username:      "TestUser2",
		Password:      "TestPass2",
		ServerAddress: "https://test2.com",
	}

	err = AddCredentialsToConfigFile(credential2, DefaultConfigPath)
	require.NoError(t, err)

	// Check if the credential is added & new credential is set as current
	credentialFile, err = os.ReadFile(DefaultConfigPath)
	require.NoError(t, err)

	require.Equal(t, "credentials:\n    - name: TestCred\n      username: TestUser\n      password: TestPass\n      serveraddress: https://test.com\n    - name: TestCred2\n      username: TestUser2\n      password: TestPass2\n      serveraddress: https://test2.com\ncurrent-credential-name: TestCred2\n", string(credentialFile))

	// Clean up
	err = os.Remove(DefaultConfigPath)
	require.NoError(t, err)

	err = os.Remove(HarborFolder)
	require.NoError(t, err)
}

func TestGetCredentials(t *testing.T) {
	SetLocation()
	err := CreateConfigFile()
	require.NoError(t, err)

	credential := Credential{
		Name:          "TestCred",
		Username:      "TestUser",
		Password:      "TestPass",
		ServerAddress: "https://test.com",
	}

	credential2 := Credential{
		Name:          "TestCred2",
		Username:      "TestUser2",
		Password:      "TestPass2",
		ServerAddress: "https://test2.com",
	}

	err = AddCredentialsToConfigFile(credential, DefaultConfigPath)
	require.NoError(t, err)

	err = AddCredentialsToConfigFile(credential2, DefaultConfigPath)
	require.NoError(t, err)

	// Get credential with name TestCred
	credentialFound, err := GetCredentials(credential.Name)
	require.NoError(t, err)
	require.Equal(t, credential, credentialFound)

	// Get empty credential if not found
	credentialFound, err = GetCredentials("TestCred3")
	require.NoError(t, err)
	require.Equal(t, Credential{}, credentialFound)

	// Clean up
	err = os.Remove(DefaultConfigPath)
	require.NoError(t, err)

	err = os.Remove(HarborFolder)
	require.NoError(t, err)
}
