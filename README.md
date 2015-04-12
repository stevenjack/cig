# cig
Can I go? Checks all your git repos specified in a local config and gives a 
report back of what state their in.

## Setup

First create the following yaml file in your home directory:

> ~/cig.yaml

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
* Repos that have changes that haven't been commited

To only check for a certain sub set of repos, say for example if you're at work
just put the key from the yaml file as the second param:

`$: cig work`
