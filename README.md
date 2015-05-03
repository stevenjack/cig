<h1 align="center">cig</h1>

<p align="center">
  <a href="https://github.com/stevenjack/cig/releases" target="_blank"><img src="https://img.shields.io/github/release/stevenjack/cig.svg"></a>
  <a href="https://travis-ci.org/stevenjack/cig" target="_blank"><img src="https://travis-ci.org/stevenjack/cig.svg?branch=master"></a>
  </p>

<p align="center">
	Can I go? CLI app for checking the state of your git repositories.
</p>

![cig](https://cloud.githubusercontent.com/assets/527874/7220202/faaedf0c-e6b6-11e4-9cb8-bf62295f4128.png)

### Installation

To install cig, follow the instructions for your OS below:

#### GO

If you have [go](http://golang.org/) installed and the bin folder added to your path, just run:

```bash
$: go get github.com/stevenjack/cig
```

#### OSX/Linux

##### One line install

```bash
curl -L https://bit.ly/cig-install | sudo bash
```

> Note: this command downloads the binary and changes the execute permission on it

##### Manual install

```bash
curl -L https://github.com/stevenjack/cig/releases/download/v0.1.2/cig_`uname -s`_x86_64 > /usr/local/bin/cig
chmod +x /usr/local/bin/cig
```

### Windows

Download the following binary:

* [64bit](https://github.com/stevenjack/cig/releases/download/v0.1.2/cig_windows_amd64.exe)

Once you have the binary, run it via your cmd prompt:

```
C: cig_windows_amd64.exe
```

### Setup

Create a `.cig.yaml` configuration file within your home directory:

#### Linux

> ~/.cig.yaml

```yaml
work: /path/to/work/repos
personal: /path/to/personal/repos
```

#### Windows

> ~/Users/Steven Jack/.cig.yaml

```yaml
work: C:\path\to\work\repos
personal: C:\path\to\personal\repos
```

The configuration file defines the different folder locations that contain the `.git` repos you want cig to check for you.

### Usage

Simply run:

`$: cig`

Once executed, cig will check all your repos and the following information will be displayed:

* Repos that need pushing up to the origin: `P`
* Repos that have new, unstaged and staged changes: `M(10)`

> Note: the values will soon be replaced with something similar to:  
> `(S)taged, (M)odified, (N)ew`  
> Please see the [issues page](https://github.com/stevenjack/cig/issues) for full details

#### Filters

If you just want to check your 'work' repos for changes:

`$: cig -t work`

To filter them based on a certain string such as 'steve':

`$: cig -f steve`

You can also combine them, so to only show 'work' repos with 'steve' 
in the path:

`$: cig -t work -f steve`

### TODO

Please see [issues](https://www.github.com/stevenjack/cig/issues?utf8=✓&q=is%3Aissue+is%3Aopen+label%3ATODO)

### Contributing

1. Fork it ( http://github.com/stevenjack/cig/fork )
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request
