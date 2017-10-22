package main

import (
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

// This program requires that the following environment vars are set:

// SAKURACLOUD_ACCESS_TOKEN: contains your SakuraCloud APIKey(access token)
// SAKURACLOUD_ACCESS_TOKEN_SECRET: contains your SakuraCloud APIKey(access token secret)
// SAKURACLOUD_ZONE: target zone name[is1a/is1b/tk1a/tk1v]

func runSakuraCloud(args []string) {
	flags := flag.NewFlagSet("sakuracloud", flag.ExitOnError)
	invoked := filepath.Base(os.Args[0])
	flags.Usage = func() {
		fmt.Printf("USAGE: %s run sakuracloud [options] [name]\n\n", invoked)
		fmt.Printf("'name' is the name of an SakuraCloud archive that has already been\n")
		fmt.Printf(" uploaded using 'linuxkit push'\n\n")
		fmt.Printf("Options:\n\n")
		flags.PrintDefaults()
	}

	nameFlag := flags.String("name", "", "Server Name")
	coreFlag := flags.Int("core", 1, "Number of CPU core")
	memoryFlag := flags.Int("memory", 1, "Size of memory in GB")
	diskSizeFlag := flags.Int("disk-size", 20, "Size of system disk in GB")

	if err := flags.Parse(args); err != nil {
		log.Fatalf("Unable to parse args: %s", err.Error())
	}

	remArgs := flags.Args()
	if len(remArgs) == 0 {
		fmt.Printf("Please specify the name of the image to boot\n")
		flags.Usage()
		os.Exit(1)
	}
	archiveName := remArgs[0]

	name := getStringValue("SAKURACLOUD_SERVER_NAME", *nameFlag, "")
	core := getIntValue("SAKURACLOUD_SERVER_CORE", *coreFlag, 1)
	memory := getIntValue("SAKURACLOUD_SERVER_MEMORY", *memoryFlag, 1)
	diskSize := getIntValue("SAKURACLOUD_SERVER_DISK_SIZE", *diskSizeFlag, 20)

	// create server on sakuracloud
	token := getEnvVarOrExit("SAKURACLOUD_ACCESS_TOKEN")
	secret := getEnvVarOrExit("SAKURACLOUD_ACCESS_TOKEN_SECRET")
	zone := getEnvVarOrExit("SAKURACLOUD_ZONE")
	clientParam := sakuraCloudClientParam{
		token:  token,
		secret: secret,
		zone:   zone,
	}

	// 1. Find archive by archiveName
	archives, err := findSakuraCloudArchive(clientParam, sakuraCloudFindArchiveParam{name: archiveName})
	if err != nil {
		log.Fatalf("Unable to find archive: %s", err)
	}
	if len(archives) < 1 {
		log.Fatalf("Unable to find archive with name %s", archiveName)
	}
	if len(archives) > 1 {
		log.Warnf("Found multiple archives with the same name, using the first one")
	}
	archiveID := archives[0].ID

	if name == "" {
		name = archives[0].Name
	}

	// 2. Create server
	server, err := createSakuraCloudServer(clientParam, sakuraCloudCreateServerParam{
		name:            name,
		core:            core,
		memory:          memory,
		diskSize:        diskSize,
		sourceArchiveID: archiveID,
	})
	if err != nil {
		log.Fatalf("Unable to create server: %s", err)
	}

	fmt.Printf("\nSakuraCloud server( %s[%d] ) is createdn", server.Name, server.ID)
	fmt.Printf("\nssh -i path-to-key root@%s\n\n", server.Interfaces[0].IPAddress)
}
