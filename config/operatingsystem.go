package config

import "encoding/json"
import "os"

type OperatingSystem struct {
	Family       string
	Version      string
	Edition      string
	Architecture string
	License      string
	Metadata     map[string]string
}

func DefaultOperatingSystems() map[string]*OperatingSystem {
	operatingSystems := map[string]*OperatingSystem{
		"windows-2008-r2-standard-sp1-trial":  &OperatingSystem{Family: "windows", Version: "2008-r2-sp1", Edition: "standard", Architecture: "amd64", License: "trial", Metadata: map[string]string{"iso_url": "http://care.dlservice.microsoft.com/dl/download/7/5/E/75EC4E54-5B02-42D6-8879-D8D3A25FBEF7/7601.17514.101119-1850_x64fre_server_eval_en-us-GRMSXEVAL_EN_DVD.iso", "iso_checksum": "4263be2cf3c59177c45085c0a7bc6ca5", "iso_checksum_type": "md5"}},
		"windows-2008-r2-standard-sp1-retail": &OperatingSystem{Family: "windows", Version: "2008-r2-sp1", Edition: "standard", Architecture: "amd64", License: "retail", Metadata: map[string]string{"iso_url": "./iso/en_windows_server_2008_r2_with_sp1_x64_dvd_617601.iso", "iso_checksum": "D3FD7BF85EE1D5BDD72DE5B2C69A7B470733CD0A", "iso_checksum_type": "sha1"}},
		"windows-2008-r2-standard-sp1-volume": &OperatingSystem{Family: "windows", Version: "2008-r2-sp1", Edition: "standard", Architecture: "amd64", License: "volume", Metadata: map[string]string{"iso_url": "./iso/en_windows_server_2008_r2_with_sp1_vl_build_x64_dvd_617403.iso", "iso_checksum": "7E7E9425041B3328CCF723A0855C2BC4F462EC57", "iso_checksum_type": "sha1"}},
	}
	return operatingSystems
}

func WriteOperatingSystems(destination *os.File) {
	operatingSystems := DefaultOperatingSystems()
	b, _ := json.MarshalIndent(operatingSystems, "", "    ")
	file, err := os.Create(destination.Name())
	if err != nil {
		// handle the error here
		return
	}
	defer file.Close()

	file.Write(b)
}
