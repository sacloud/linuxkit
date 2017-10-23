package main

import (
	"fmt"
	"github.com/sacloud/libsacloud/api"
	"github.com/sacloud/usacloud/helper/ftp"

	"github.com/sacloud/libsacloud/builder"
	"github.com/sacloud/libsacloud/sacloud"
	log "github.com/sirupsen/logrus"
)

type sakuraCloudClientParam struct {
	token  string
	secret string
	zone   string
}

type sakuraCloudCreateArchiveParam struct {
	name string
	path string
}

type sakuraCloudFindArchiveParam struct {
	name string
}

type sakuraCloudCreateServerParam struct {
	name            string
	core            int
	memory          int
	diskSize        int
	sourceArchiveID int64
}

type sakuraCloudResource interface {
	GetName() string
	GetID() int64
}

func newSakuraCloudClient(p sakuraCloudClientParam) *api.Client {
	c := api.NewClient(p.token, p.secret, p.zone)
	c.UserAgent = fmt.Sprintf("sacloud/linuxkit-%s", Version)
	return c
}

func uploadSakuraCloudArchive(apiParam sakuraCloudClientParam, createParam sakuraCloudCreateArchiveParam) (sakuraCloudResource, error) {

	client := newSakuraCloudClient(apiParam)
	archiveAPI := client.GetArchiveAPI()

	p := archiveAPI.New()
	p.SetName(createParam.name)
	p.SetSizeGB(20)

	// create archive
	archive, err := archiveAPI.Create(p)
	if err != nil {
		return nil, fmt.Errorf("Create Archive is failed: %s", err)
	}
	log.Debugf("Archive[%d] is created", archive.ID)

	// upload
	ftpServer, err := archiveAPI.OpenFTP(archive.ID)
	if err != nil {
		return nil, fmt.Errorf("Open to FTPS connection is failed: %s", err)
	}
	log.Debugf("Archive[%d] FTP connection is opened", archive.ID)

	ftpsClient := ftp.NewClient(ftpServer.User, ftpServer.Password, ftpServer.HostName)
	err = ftpsClient.Upload(createParam.path)
	if err != nil {
		return nil, fmt.Errorf("Upload Archive is failed: %s", err)
	}
	log.Debugf("Archive[%d] is uploaded", archive.ID)

	// close FTP
	_, err = archiveAPI.CloseFTP(archive.ID)
	if err != nil {
		return nil, fmt.Errorf("Close FTPS connection is failed: %s", err)
	}
	log.Debugf("Archive[%d] FTP connection is closed", archive.ID)

	return archive, nil
}

func findSakuraCloudArchive(apiParam sakuraCloudClientParam, findParam sakuraCloudFindArchiveParam) ([]sacloud.Archive, error) {
	client := newSakuraCloudClient(apiParam)
	archiveAPI := client.GetArchiveAPI()

	res, err := archiveAPI.Reset().WithNameLike(findParam.name).Find()
	if err != nil {
		return []sacloud.Archive{}, fmt.Errorf("Find Archive is failed: %s", err)
	}

	return res.Archives, nil
}

func createSakuraCloudServer(apiParam sakuraCloudClientParam, createParam sakuraCloudCreateServerParam) (*sacloud.Server, error) {
	client := newSakuraCloudClient(apiParam)
	builder := builder.ServerFromArchive(client, createParam.name, createParam.sourceArchiveID)
	builder.SetCore(createParam.core)
	builder.SetMemory(createParam.memory)
	builder.SetDiskSize(createParam.diskSize)
	builder.AddPublicNWConnectedNIC()

	res, err := builder.Build()
	if err != nil {
		return nil, fmt.Errorf("Create Server is failed: %s", err)
	}
	return res.Server, nil
}
