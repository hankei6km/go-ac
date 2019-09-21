# go-ac

ac is helper library that creates CREDITS files(s) using `go version -m` and `gocredits`.

## Requirement

- Go 1.13
- gocredits

##  Ueage

as is used mainly in magefile.go.

input
```
  ▾ my-proj/
    ▾ dist/
      ▾ my_cmd_linux_386/
          my_cmd
      ▾ my_cmd_linux_amd64
          my_cmd
      ▾ my_cmd_windows_amd64
          my_cmd.exe
      go.sum
```

output(If `CREDITS_*` files are all the same, they are merged into the `CREDITS` file)
```
  ▾ my-proj/
      CREDITS_linux_386
      CREDITS_linux_amd64
      CREDITS_windows_amd64
```

code(magefile.go)
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
		GoSumFile(filepath.Join(cwd, "go.sum")).
		Prog("./gocredits")
	d := ac.NewDistBuilder().
		DistDir(filepath.Join(cwd, "dist")).
		WorkDir(tmpDir).
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