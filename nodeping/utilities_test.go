package nodeping_test

import (
	"io/ioutil"
	"log"
	"os"
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
	if err != nil {
		log.Fatal(err)
	}
	for _, fileName := range []string{".terraform.lock.hcl", "main.tf",
		"terraform.tfstate", "terraform.tfstate.backup"} {
		err = os.Remove(terraformDir + "/" + fileName)
		if err != nil {
			log.Fatal(err)
		}
	}
}
