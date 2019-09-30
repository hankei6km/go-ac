# go-ac

[![Build Status](https://travis-ci.org/hankei6km/go-ac.svg?branch=master)](https://travis-ci.org/hankei6km/go-ac)

ac is helper library that creates CREDITS files(s) using `go version -m` and `gocredits`.

## Requirement

- Go 1.13
- [gocredits](https://github.com/Songmu/gocredits): if you want to use `gocredits` as extarnal program.

##  Ueage

Example of using ac in [magefile](https://github.com/magefile/mage).

input
```
  ▾ my-proj/
    ▾ dist/
      ▾ my_cmd_linux_386/
          my_cmd
      ▾ my_cmd_linux_amd64/
          my_cmd
      ▾ my_cmd_windows_amd64/
          my_cmd.exe
      go.sum
```

output(If `CREDITS_*` files are all the same, they are merged into the `CREDITS` file).
```
  ▾ my-proj/
      CREDITS_Linux_i386
      CREDITS_Linux_amd64
      CREDITS_Windows_amd64
```

magefile.go
```go
import "github.com/hankei6km/go-ac"

func Credits() error {

	// ...

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	tmpDir := filepath.Join(cwd, "tmp")
	if err := ac.ResetDir(tmpDir, os.ModePerm); err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	b := ac.NewOutputBuilder().
		GoSumFile(filepath.Join(cwd, "go.sum"))
	d := ac.NewDistBuilder().
		DistDir(filepath.Join(cwd, "dist")).
		WorkDir(tmpDir).
		ReplaceOs([][]string{
			[]string{"linux", "Linux"},
			[]string{"windows", "Windows"},
		}).
		ReplaceArch([][]string{
			[]string{"386", "i386"},
		}).
		OutDir(cwd).
		OutputBuilder(b).
		Build()
	if err := d.Run(); err != nil {
		return err
	}

	// ...

	return nil
}
```