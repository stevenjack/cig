<h1 align="center">cig</h1>

<p align="center">
  <a href="https://registry.hub.docker.com/u/smaj/cig" target="_blank"><img src="https://img.shields.io/badge/Docker-Hub-3a9bdc.svg?style=flat-square"></a>
  </p>

<p align="center">
	Can I go? CLI app for checking the state of your git repositories.
</p>

### Installation

To install cig follow the instructions for your os below:

#### OSX/Linux

##### One line install

You can run the following that downloads the binary and changes the execute
permission on it:

```bash
curl -L http://bit.ly/1Hm2Llh | sudo bash
```

##### Manual install

```bash
curl -L https://s3-eu-west-1.amazonaws.com/cig.steven-jack.co.uk/0.1.0/cig_`uname -s`_`uname -m` > /usr/local/bin/cig
chmod +x /usr/local/bin/cig
```

### Windows

Download the following binary:

* [32bit](https://s3-eu-west-1.amazonaws.com/cig.steven-jack.co.uk/0.1.0/cig_windows_386.exe)
* [64bit](https://s3-eu-west-1.amazonaws.com/cig.steven-jack.co.uk/0.1.0/cig_windows_amd64.exe)

then run it in a cmd prompt:

```
C: cig_windows_amd64.exe
```

### Setup

First create the following yaml file in your home directory:

> ~/.cig.yaml

```yaml
work: /path/to/work/repos
personal: /path/to/personal/repos
```

this is a list of the different folders that contain your repos you want to check.

### Usage

Simply run:

`$: cig`

and all repos will be checked and the following will show up:

* Repos that need pushing up to the origin
* Repos that have new, unstaged and staged changes.

#### Filters

If you just want to check your 'work' repos for changes:

`$: cig -t work`

To filter them based on a certain string such as 'steve':

`$: cig -f steve`

You can also combine them, so to only show 'work' repos with 'steve' 
in the path:

`$: cig -t work -f steve`
