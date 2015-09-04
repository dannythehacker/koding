# fuseklient

Prototype that integrates [Fuse](https://github.com/bazil/fuse) and [Klient](https://github.com/koding/klient). Currently for OSX only.

## WARNING

  Use a new VM.

## Steps to get started:

    # Download fuseklient from https://github.com/koding/fuseklient/releases.

    # Get the IP of your VM and folder in that VM you want to mount locally.
    #
    # Or you can use mine:
    #   user         : sent-hil
    #   klientip     : 52.7.78.76
    #   externalpath : /home/sent-hil/fusemount

    # Create a folder in local to mount external folder:
    mkdir -p <fullpath>/local

    # Start daemon:
    ./fuseklient --klientip=... --externalpath=/path/to/external --internalpath=/path/to/local --user=<koding username>

    # In another terminal:
    ls -alh /path/to/local

    # If you get `Device not configured` when trying to access mount when daemon is not running:
    diskutil unmount force <folder>

    # If you get `mount point <folder> is itself on a OSXFUSE volume`:
    diskutil unmount force <folder>

## Following commands work on a mounted folder.

  * ls -alh folder/ file
  * cd folder/
  * mkdir -p folder/nested
  * mv folder/nested/ nested/
  * rm -rf folder/ file
  * cat file
  * touch file
  * echo 1 >> file
  * find ... # recursive search will be slow
  * grep ... # recursive search will be slow

## Milestones:

    * ALPHA
        read operations
        klient authentication
        write operations
    * BETA
        klient running on OSX
        integrate fuseklient into klient
          `klient mount --vm`
        invalidate local cache on file changes in user VM
        kd ... - run entire command on VM, return results
        shell hooks: fish, bash
    * 1.0
        klient ps - return list of user VMs to mount
        lock resources in VM on open or write operations
        remaining FUSE operations

## Notes:

  * Use fullpath in arguments.
  * Mounting on an existing folder won't overwrite contents, but they won't be visible while Fuse is running.
  * Do `go generate` if you make changes to `install-alpha.sh`
  * Use `tar -cvf fuseklient_OSX.tar fuseklient Readme.md` for distribution.
