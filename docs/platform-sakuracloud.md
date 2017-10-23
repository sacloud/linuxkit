# Using LinuxKit on SakuraCloud 

This is a quick guide to run LinuxKit on [SakuraCloud](https://cloud.sakura.ad.jp/).

## Setup

You need to authenticate LinuxKit with your SakuraCloud APIKey. You need to set up the following environment variables:

- `SAKURACLOUD_ACCESS_TOKEN`: Your API Key(token) of SakuraCloud. 
- `SAKURACLOUD_ACCESS_TOKEN_SECRET`: Your API Key(secret) of SakuraCloud.
- `SAKURACLOUD_ZONE`: Target zone name. Must be in [is1a/is1b/tk1a/tk1v].

## Build an image

Create a new `sakuracloud.yml` file [based on the SakuraCloud example](../examples/sakuracloud.yml), generate a new SSH key and add it in the `yml`, then `moby build -format raw sakuracloud.yml`.

Note: SakuraCloud requires a `RAW` image.

```
$ moby build -format raw examples/sakuracloud.yml
```

## Push image

Do `linuxkit push sakuracloud sakuracloud.raw` to upload it to the
specified(by `SAKURACLOUD_ZONE` environment variable) zone, and create a bootable archive from the stored image.

## Create an server and connect to it

With the archive created, we can now create an server.

```
linuxkit run sakuracloud sakuracloud
```

You can edit the SakuraCloud example to allow you to SSH to your server in order to use it.
