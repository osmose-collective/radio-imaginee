package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
	"time"

	cli "gopkg.in/urfave/cli.v2"
)

type PlaylistsInfo struct {
	Global struct {
		NbDirs      int
		NbFiles     int
		TotalSizeMb int64
	}
	ByCategory        map[string]int
	ByExtension       map[string]int
	LastModifiedFiles []string
	BiggestFiles      []string
}

type playlistFile struct {
	path string
	info os.FileInfo
}

func getPlaylistsInfo(c *cli.Context) (*PlaylistsInfo, error) {
	baseDir := c.String("playlists-dir")
	pi := PlaylistsInfo{
		ByCategory:  map[string]int{},
		ByExtension: map[string]int{},
	}
	allFiles := []playlistFile{}

	// walk fs
	if err := filepath.Walk(
		baseDir,
		func(fullPath string, info os.FileInfo, err error) error {
			relPath := strings.TrimPrefix(fullPath, baseDir+"/")
			if info.IsDir() {
				pi.Global.NbDirs++
			} else {
				allFiles = append(allFiles, playlistFile{
					path: relPath,
					info: info,
				})
			}
			return nil
		}); err != nil {
		return nil, err
	}

	for _, file := range allFiles {
		//fmt.Println(relPath, info, err)
		// compute files by categories
		category := strings.Split(file.path, "/")[0]
		if _, exists := pi.ByCategory[category]; !exists {
			pi.ByCategory[category] = 0
		}
		pi.ByCategory[category]++
		// compute files by extensions
		extension := path.Ext(file.path)
		if _, exists := pi.ByExtension[extension]; !exists {
			pi.ByExtension[extension] = 0
		}
		pi.ByExtension[extension]++
		// global files
		pi.Global.NbFiles++
		pi.Global.TotalSizeMb += file.info.Size()
	}
	pi.Global.TotalSizeMb /= (1024 * 1024)

	// biggest files
	sort.Slice(allFiles, func(i, j int) bool {
		return allFiles[i].info.Size() > allFiles[j].info.Size()
	})
	pi.BiggestFiles = []string{}
	for _, file := range allFiles {
		if file.IsTrack() {
			pi.BiggestFiles = append(
				pi.BiggestFiles,
				fmt.Sprintf("%s (%dMb)", file.path, file.info.Size()/1024/1024),
			)
			if len(pi.BiggestFiles) >= 20 {
				break
			}
		}
	}
	// last modified files
	sort.Slice(allFiles, func(i, j int) bool {
		aStat := allFiles[i].info.Sys().(*syscall.Stat_t)
		bStat := allFiles[j].info.Sys().(*syscall.Stat_t)
		aTime := time.Unix(int64(aStat.Atim.Sec), int64(aStat.Atim.Nsec))
		bTime := time.Unix(int64(bStat.Atim.Sec), int64(bStat.Atim.Nsec))
		return aTime.Before(bTime)
	})
	pi.LastModifiedFiles = []string{}
	for _, file := range allFiles {
		if file.IsTrack() {
			pi.LastModifiedFiles = append(
				pi.LastModifiedFiles,
				fmt.Sprintf("%s (%v)", file.path, file.info.ModTime()),
			)
			if len(pi.LastModifiedFiles) >= 20 {
				break
			}
		}
	}

	return &pi, nil
}

func (f playlistFile) IsTrack() bool {
	switch path.Ext(f.path) {
	case ".mp3", ".MP3", ".WAV":
		return true
	}
	return false
}
