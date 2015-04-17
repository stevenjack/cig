<h1 align="center">cig</h1>

<p align="center">
  <a href="https://registry.hub.docker.com/u/smaj/cig" target="_blank"><img src="https://img.shields.io/badge/Docker-Hub-3a9bdc.svg?style=flat-square"></a>
  </p>

<p align="center">
	Can I go? CLI app for checking the state of your git repositories.
</p>

## Setup

First create the following yaml file in your home directory:

> ~/.cig.yaml

```yaml
work: /path/to/work/repos
personal: /path/to/personal/repos
```

this is a list of the different folders that contain your repos you want to check.

## Usage

Simply run:

`$: cig`

and all repos will be checked and the following will show up:

* Repos that need pushing up to the origin
* Repos that have new, unstaged and staged changes.

To only check for a certain sub set of repos, say for example if you're at work
just put the key from the yaml file as the second param:

`$: cig work`
