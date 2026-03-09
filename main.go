package main

import (
	"crypto/md5"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pixb/fdf/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// 功能：查找指定目录下重复文件
// 1.命令行参数支持 --path 路径参数，指定要查找的目录,--config 配置文件参数
// 2.配置文件中主要配置文件夹中保留的权重，数字越低，权重越高，默认99
// 3.日志文件夹logs
var Version = "development" // 由构建时注入

type PriorityConfig struct {
	Info              string
	DefaultPriority   int            `mapstructure:"default_priority"`
	DirectoryPriority map[string]int `mapstructure:"directory_priority"`
}

var rootCmd = &cobra.Command{
	Use:   "fdf",
	Short: "Find duplicate files in specified directory",
	Long:  `Find duplicate files in specified directory`,
	Run: func(cmd *cobra.Command, args []string) {
		// version
		isPrintVersion := viper.GetBool("version")
		if isPrintVersion {
			fmt.Println(Version)
			return
		}

		// dry-run
		dryRun := viper.GetBool("dry-run")
		fmt.Println("dry-run?", dryRun)

		// set config
		configFilePath := viper.GetString("config")
		var priorityConfig PriorityConfig

		// check config path exists.
		if _, err := os.Stat(configFilePath); !os.IsNotExist(err) {
			fmt.Println("Config file path:", configFilePath)
			configPath := filepath.Dir(configFilePath)
			viper.SetConfigName(strings.TrimSuffix(filepath.Base(configFilePath), filepath.Ext(configFilePath)))
			viper.SetConfigType(filepath.Ext(configFilePath)[1:])
			viper.AddConfigPath(configPath)

			// read config file
			if err := viper.ReadInConfig(); err != nil {
				fmt.Printf("Error reading config file: %v", err)
				return
			}
			// 读取配置文件
			// 检查关键配置项是否存在
			requiredKeys := []string{"default_priority", "directory_priority"}
			for _, key := range requiredKeys {
				if !viper.IsSet(key) {
					fmt.Println("Missing required config key = ", key)
					return
				}
			}

			// 使用简化的Unmarshal配置
			viper.Set("Verbose", true)
			if err := viper.Unmarshal(&priorityConfig); err != nil {
				fmt.Println("Error unmarshalling config = ", err)
				return
			}
		}
		// get path
		path := viper.GetString("path")
		fmt.Println("==================================")
		fmt.Printf("path = %s\n", path)
		fmt.Println("==================================")
		// find duplicate files
		duplicateFiles := findDuplicateFiles(path)
		// print duplicate files
		if len(duplicateFiles) > 0 {
			fmt.Println("Duplicate files found:")
			for hash, paths := range duplicateFiles {
				fmt.Println("Duplicate file hash = ", hash)
				for _, path := range paths {
					fmt.Println("\tPath = ", path)
				}
			}
		} else {
			fmt.Println("No duplicate files found")
		}
		fmt.Println()
		// 处理重复文件
		duplicateFileHandler(duplicateFiles, priorityConfig, dryRun)
	},
}

func init() {
	rootCmd.PersistentFlags().StringP("path", "p", "", "Path of search for duplicate files")
	rootCmd.PersistentFlags().StringP("config", "c", "", "Path of config file,e.g. ./config.json")
	rootCmd.PersistentFlags().BoolP("version", "v", false, "Print fdf version.")
	// --dry-run -n: 参考自rsync
	rootCmd.PersistentFlags().BoolP("dry-run", "n", false, "Dry run,dont delete.")
	// --exclude -e: 排除目录
	rootCmd.PersistentFlags().StringSliceP("exclude", "e", []string{}, "Exclude directories, e.g. --exclude dir1 --exclude dir2")
	// --min-size: 最小文件大小（字节）
	rootCmd.PersistentFlags().Int64P("min-size", "m", 0, "Minimum file size in bytes")
	// --max-size: 最大文件大小（字节）
	rootCmd.PersistentFlags().Int64P("max-size", "M", 0, "Maximum file size in bytes")

	if err := viper.BindPFlag("path", rootCmd.PersistentFlags().Lookup("path")); err != nil {
		panic(err)
	}
	if err := viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config")); err != nil {
		panic(err)
	}
	if err := viper.BindPFlag("version", rootCmd.PersistentFlags().Lookup("version")); err != nil {
		panic(err)
	}
	if err := viper.BindPFlag("dry-run", rootCmd.PersistentFlags().Lookup("dry-run")); err != nil {
		panic(err)
	}
	if err := viper.BindPFlag("exclude", rootCmd.PersistentFlags().Lookup("exclude")); err != nil {
		panic(err)
	}
	if err := viper.BindPFlag("min-size", rootCmd.PersistentFlags().Lookup("min-size")); err != nil {
		panic(err)
	}
	if err := viper.BindPFlag("max-size", rootCmd.PersistentFlags().Lookup("max-size")); err != nil {
		panic(err)
	}
}

// 查找重复文件函数
// 遍历数据目录, 根据文件路径调用获取文件hash函数，取得hash值。
// 存入map集合中，key为hash值，value为文件路径列表，因为可能有重复的文件。
// 遍历map集合，如果value的长度大于1，则输出日志。
// 返回重复文件大于1的文件列表组成的map集合。
func findDuplicateFiles(path string) map[string][]string {
	// 获取排除目录列表
	excludeDirs := viper.GetStringSlice("exclude")
	// 获取文件大小限制
	minSize := viper.GetInt64("min-size")
	maxSize := viper.GetInt64("max-size")

	fmt.Println("Scanning directory for files...")
	// 遍历数据目录，收集所有文件路径
	var filePaths []string
	err := filepath.WalkDir(path, func(path string, info os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// 检查是否为目录
		if info.IsDir() {
			// 检查是否需要排除该目录
			dirName := filepath.Base(path)
			for _, exclude := range excludeDirs {
				if dirName == exclude {
					return filepath.SkipDir
				}
			}
			return nil
		}

		// 检查是否为@eaDir目录下的文件
		if strings.Contains(path, "@eaDir") {
			return nil
		}

		// 检查文件大小
		fileInfo, err := info.Info()
		if err != nil {
			fmt.Printf("Error getting file info for %s: %v\n", path, err)
			return nil
		}

		fileSize := fileInfo.Size()
		// 检查最小文件大小
		if minSize > 0 && fileSize < minSize {
			return nil
		}
		// 检查最大文件大小
		if maxSize > 0 && fileSize > maxSize {
			return nil
		}

		filePaths = append(filePaths, path)
		return nil
	})
	if err != nil {
		fmt.Println("Error while traversing directory:", err)
		return make(map[string][]string)
	}

	fmt.Printf("Found %d files to process\n", len(filePaths))
	if len(filePaths) == 0 {
		fmt.Println("No files found to process")
		return make(map[string][]string)
	}

	// 并发处理文件
	const maxWorkers = 4 // 并发数
	fileChan := make(chan string, len(filePaths))
	resultChan := make(chan struct {
		hash string
		path string
		err  error
	}, len(filePaths))

	// 启动工作协程
	for i := 0; i < maxWorkers; i++ {
		go func() {
			hasher := md5.New() // 每个协程使用自己的哈希器
			for path := range fileChan {
				hash, err := util.CalculateHash(path, hasher)
				resultChan <- struct {
					hash string
					path string
					err  error
				}{hash, path, err}
			}
		}()
	}

	// 发送文件路径到通道
	for _, path := range filePaths {
		fileChan <- path
	}
	close(fileChan)

	// 收集结果
	fileMap := make(map[string][]string)
	processed := 0
	total := len(filePaths)
	for i := 0; i < total; i++ {
		result := <-resultChan
		processed++
		// 每处理10%的文件显示一次进度
		if processed%((total/10)+1) == 0 {
			fmt.Printf("Processing files: %.1f%%\r", float64(processed)/float64(total)*100)
		}
		if result.err != nil {
			fmt.Printf("Error calculating hash for %s: %v\n", result.path, result.err)
			continue
		}
		fileMap[result.hash] = append(fileMap[result.hash], result.path)
	}
	fmt.Println("Processing files: 100.0%")

	// 遍历map集合，输出日志
	duplicateFiles := make(map[string][]string)
	for hash, paths := range fileMap {
		if len(paths) > 1 {
			duplicateFiles[hash] = paths
		}
	}

	fmt.Printf("Found %d duplicate file groups\n", len(duplicateFiles))
	return duplicateFiles
}

// 处理重复文件函数
// 遍历重复文件列表，根据文件路径调用处理重复文件函数，处理重复文件。
// 根据配置文件列表中权重，删除重复文件。
// 根据文件路径权限对文件列表进行排序。
// 保留排序后列表的第一个文件，获取第一个文件的权重。
// 处理第一个文件之后的文件列表，根据权重，删除文件。
func duplicateFileHandler(duplicateFiles map[string][]string, priorityConfig PriorityConfig, dryRun bool) {
	if len(duplicateFiles) == 0 {
		return
	}

	for hash, paths := range duplicateFiles {
		fmt.Printf("Processing duplicate files with hash: %s\n", hash)
		// 对paths进行排序，排序规则是文件最后路径名的权重
		sort.Slice(paths, func(i, j int) bool {
			// 获取文件最后路径名的权重, 优先级越高，权重越低
			iPath := filepath.Base(filepath.Dir(paths[i]))
			jPath := filepath.Base(filepath.Dir(paths[j]))
			iWeight, ok := priorityConfig.DirectoryPriority[iPath]
			if !ok {
				iWeight = priorityConfig.DefaultPriority
			}
			jWeight, ok := priorityConfig.DirectoryPriority[jPath]
			if !ok {
				jWeight = priorityConfig.DefaultPriority
			}
			return iWeight < jWeight
		})
		// keep the first file
		firstFile := paths[0]
		// 获取第一个文件的权重
		firstPath := filepath.Dir(firstFile)
		firstWeight, ok := priorityConfig.DirectoryPriority[filepath.Base(firstPath)]
		if !ok {
			firstWeight = priorityConfig.DefaultPriority
		}
		fmt.Printf("\tKeeping file: %s (priority: %d)\n", firstFile, firstWeight)
		// 处理第一个文件之后的文件列表
		for i := 1; i < len(paths); i++ {
			path := paths[i]
			// 获取文件路径权限
			_, err := os.Stat(path)
			if err != nil {
				fmt.Printf("Error getting file info for %s: %v\n", path, err)
				continue
			}
			// 获取文件路径权重
			pathWeight, ok := priorityConfig.DirectoryPriority[filepath.Base(filepath.Dir(path))]
			if !ok {
				pathWeight = priorityConfig.DefaultPriority
			}
			// 如果文件路径权重大于等于第一个文件路径权重，则删除文件
			if pathWeight >= firstWeight {
				fmt.Printf("\tDeleting file: %s (priority: %d)\n", path, pathWeight)
				if !dryRun {
					err := os.Remove(path)
					if err != nil {
						fmt.Printf("\tError deleting file %s: %v\n", path, err)
						continue
					}
					fmt.Printf("\tSuccessfully deleted file: %s\n", path)
				}
			} else {
				fmt.Printf("\tKeeping file: %s (priority: %d)\n", path, pathWeight)
			}
		}
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Error executing command: %v\n", err)
		os.Exit(1)
	}
}
