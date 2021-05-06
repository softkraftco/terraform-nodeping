package nodeping_test

import (
	"errors"
	"io/ioutil"
	"log"
	"os"

	apiClient "terraform-nodeping/nodeping_api_client"
)

func copyFile(src, dst string) {
	data, err := ioutil.ReadFile(src)
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(dst, data, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func cleanupTerraformDir(terraformDir string) {
	err := os.RemoveAll(terraformDir + "/.terraform")
	// don't mind if this is already missing
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Fatal(err)
	}
	for _, fileName := range []string{".terraform.lock.hcl", "main.tf",
		"terraform.tfstate", "terraform.tfstate.backup"} {
		err = os.Remove(terraformDir + "/" + fileName)
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			log.Fatal(err)
		}
	}
}

func getClient() *apiClient.Client {
	token := os.Getenv("NODEPING_API_TOKEN")
	return apiClient.NewClient(token)
}
