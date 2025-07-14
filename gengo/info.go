package gengo

// import (
// 	"fmt"
// 	"os"
// 	"path/filepath"
// 	"strings"

// 	"golang.org/x/mod/modfile"
// )

// type Info struct {
// 	moduleDir     string
// 	moduleName    string
// 	modelsDir     string
// 	modelsFile    string
// 	modelsImport  string
// 	modelsName    string
// 	modelsPath    string
// 	msgpackDir    string
// 	msgpackName   string
// 	msgpackImport string
// 	msgpackFile   string
// 	msgpackPath   string
// 	reformat      bool
// }

// func getInfo(opts Opts) (Info, error) {
// 	if opts.ModelOutput == "" {
// 		return Info{}, fmt.Errorf("no ModelOutput option set")
// 	}

// 	modelsDir, modelsFile, err := getFile(opts.ModelOutput, "models.go")
// 	if err != nil {
// 		return Info{}, err
// 	}
// 	moduleDir, moduleName, err := getModule(modelsDir)
// 	if err != nil {
// 		return Info{}, err
// 	}

// 	modelsName, _ := strings.CutSuffix(modelsFile, ".go")
// 	modelsPath, _ := strings.CutPrefix(modelsDir, moduleDir)
// 	modelsPath = strings.TrimPrefix(modelsPath, string(filepath.Separator))
// 	modelsImport := moduleName
// 	if modelsPath != "" {
// 		modelsImport += "/" + modelsPath
// 	}

// 	msgpackDir := modelsDir
// 	msgpackFile := "mpcodecs.go"
// 	if opts.MpOutput != "" {
// 		msgpackDir, msgpackFile, err = getFile(opts.MpOutput, msgpackFile)
// 		if err != nil {
// 			return Info{}, err
// 		}
// 	}
// 	msgpackName, _ := strings.CutSuffix(msgpackFile, ".go")
// 	msgpackPath, _ := strings.CutPrefix(msgpackDir, moduleDir)
// 	msgpackPath = strings.TrimPrefix(msgpackPath, string(filepath.Separator))
// 	msgpackImport := moduleName
// 	if msgpackPath != "" {
// 		msgpackImport += "/" + msgpackPath
// 	}

// 	info := Info{
// 		moduleDir:     moduleDir,
// 		moduleName:    moduleName,
// 		modelsDir:     modelsDir,
// 		modelsFile:    modelsFile,
// 		modelsImport:  modelsImport,
// 		modelsName:    modelsName,
// 		modelsPath:    modelsPath,
// 		msgpackDir:    msgpackDir,
// 		msgpackFile:   msgpackFile,
// 		msgpackImport: msgpackImport,
// 		msgpackName:   msgpackName,
// 		msgpackPath:   msgpackPath,
// 		reformat:      opts.Reformat,
// 	}

// 	return info, nil
// }

// func getFile(path string, name string) (string, string, error) {
// 	var dir string
// 	var file string
// 	if strings.HasSuffix(path, ".go") {
// 		dir = filepath.Dir(path)
// 		file = filepath.Base(path)
// 	} else {
// 		dir = path
// 		file = name
// 	}
// 	abs, err := filepath.Abs(dir)
// 	return abs, file, err
// }

// func getModule(path string) (string, string, error) {
// 	abs, err := filepath.Abs(path)
// 	if err != nil {
// 		return "", "", err
// 	}
// 	root := filepath.VolumeName(abs) + string(filepath.Separator)
// 	top := abs
// 	var file string
// 	for {
// 		file = filepath.Join(top, "go.mod")
// 		if _, err := os.Stat(file); err == nil {
// 			break
// 		}
// 		file = ""
// 		if top == root {
// 			break
// 		}
// 		top = filepath.Dir(top)
// 	}
// 	if file == "" {
// 		return "", "", fmt.Errorf("can't find go.mod file above %s", abs)
// 	}
// 	bytes, err := os.ReadFile(file)
// 	if err != nil {
// 		return "", "", fmt.Errorf("can't read %s file: %w", file, err)
// 	}
// 	info, err := modfile.Parse(file, bytes, nil)
// 	if err != nil {
// 		return "", "", fmt.Errorf("can't parse %s file: %w", file, err)
// 	}
// 	name := info.Module.Mod.Path
// 	// more := strings.TrimPrefix(abs, top)
// 	// more = strings.TrimPrefix(more, string(filepath.Separator))
// 	return top, name, nil
// }
