# git-get-pr-id

Go version a part of the `open-pull-request` command of [this site](https://techlife.cookpad.com/entry/2015/11/17/151426)

## Usage

```sh
$ git-get-pr-id <commit hash>
```

## Tips

```sh
$ gh pr view (git-get-pr-id <commit hash>) -w
```

### tig

You can open the PR from tig with `Shit - P` command use by below config.

```
bind generic P @sh -c 'gh pr view $(git-get-pr-id %(commit)) -w'
```


