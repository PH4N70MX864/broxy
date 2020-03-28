package core

import (
	"encoding/pem"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)
// Config represents the global configuration
type BroxySettings struct {
	CACertificate   []byte    `xml:"CACert"`
	CAPrivateKey    []byte    `xml:"CAPvt"`
}

type GlobalSettings struct {
	GZipDecode	bool
}

var broxySettingsFileName = "broxy_settings.xml"

// LoadGlobalSettings loads the global settings
func LoadGlobalSettings(path string) *BroxySettings {

	// if path doesn't exists, create it
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, 0700)
	}

	var cfg *BroxySettings
	xmlSettingsPath := filepath.Join(path, broxySettingsFileName)
	xmlSettingsFile, err := os.Open(xmlSettingsPath)
	defer xmlSettingsFile.Close()
	if err != nil {
		// create the file and put in a new fresh default settings
		cfg = initGlobalSettings()
		saveGlobalSettings(cfg, xmlSettingsPath)
		return cfg
	}

	byteValue, err := ioutil.ReadAll(xmlSettingsFile)
	if err != nil {
		fmt.Println(err)
	}

	err = xml.Unmarshal(byteValue, &cfg)
	if err != nil {
		cfg = initGlobalSettings()
		saveGlobalSettings(cfg, xmlSettingsPath)
	}

	// TODO: add a method that checks the configuration just loaded

	return cfg

}

func initGlobalSettings() *BroxySettings {
	// generate a new CA
	// TODO: handle error generated by CreateCA
	rawPvt, rawCA, _ := CreateCA()
	pemPvt := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: rawPvt})
	pemCA := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: rawCA})
	cfg := &BroxySettings{
		CACertificate: pemCA,
		CAPrivateKey:  pemPvt,
	}
	return cfg
}

func saveGlobalSettings(cfg *BroxySettings, path string) error {
	xmlSettings, _ := xml.MarshalIndent(cfg, "", " ")
	return ioutil.WriteFile(path, xmlSettings, 0700)
}
