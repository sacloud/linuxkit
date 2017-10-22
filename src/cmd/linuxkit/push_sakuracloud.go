package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

func pushSakuraCloud(args []string) {
	flags := flag.NewFlagSet("sakuracloud", flag.ExitOnError)
	invoked := filepath.Base(os.Args[0])
	flags.Usage = func() {
		fmt.Printf("USAGE: %s push sakuracloud [options] path\n\n", invoked)
		fmt.Printf("'path' is the full path to a SakuraCloud image. It will be uploaded to SakuraCloud and Archive will be created from it.\n")
		fmt.Printf("Options:\n\n")
		flags.PrintDefaults()
	}
	nameFlag := flags.String("img-name", "", "Overrides the name used to identify the file in SakuraCloud Archive. Defaults to the base of 'path' with the file extension removed.")

	if err := flags.Parse(args); err != nil {
		log.Fatal("Unable to parse args")
	}

	remArgs := flags.Args()
	if len(remArgs) == 0 {
		fmt.Printf("Please specify the path to the image to push\n")
		flags.Usage()
		os.Exit(1)
	}
	path := remArgs[0]

	name := getStringValue(nameVar, *nameFlag, "")

	if name == "" {
		name = strings.TrimSuffix(path, filepath.Ext(path))
		name = filepath.Base(name)
	}

	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer f.Close()

	if _, err := f.Stat(); err != nil {
		log.Fatalf("Error reading file information: %v", err)
	}

	// upload to sakuracloud
	token := getEnvVarOrExit("SAKURACLOUD_ACCESS_TOKEN")
	secret := getEnvVarOrExit("SAKURACLOUD_ACCESS_TOKEN_SECRET")
	zone := getEnvVarOrExit("SAKURACLOUD_ZONE")
	clientParam := sakuraCloudClientParam{
		token:  token,
		secret: secret,
		zone:   zone,
	}

	res, err := uploadSakuraCloudArchive(clientParam, sakuraCloudCreateArchiveParam{
		name: name,
		path: path,
	})

	if err != nil {
		log.Fatalf("Error uploading file: %v", err)
	}

	log.Infof("Created Ardchive: %d", res.GetID())
}
