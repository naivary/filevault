package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
)

type FilevaultService interface {
	CreateFile(path string, file io.Reader) error
	RemoveFile(path string) error
	ReadFile(path string) ([]byte, error)
}

var _ FilevaultService = (*filevaultService)(nil)

type filevaultService struct {
	cfg Config
}

func NewFilevaultService(cfg Config) FilevaultService {
	return filevaultService{cfg}
}

func (f filevaultService) filePathJoin(paths ...string) string {
	p := slices.Concat([]string{f.cfg.Dir}, paths)
	return filepath.Join(p...)
}

func (f filevaultService) CreateFile(path string, r io.Reader) error {
	path = f.filePathJoin(path)
	dir, filename := filepath.Split(path)
	mode := os.FileMode(0755)
	if err := os.MkdirAll(dir, mode); err != nil {
		return err
	}
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("file with the name of %s is already existing in %s", filename, dir)
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.ReadFrom(r)
	return err
}

func (f filevaultService) RemoveFile(path string) error {
    path = f.filePathJoin(path)
	if err := os.Remove(path); err != nil {
		return err
	}
	dir := filepath.Dir(path)
	elems, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	if len(elems) == 0 {
		return os.Remove(dir)
	}
	return nil
}

func (f filevaultService) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(f.filePathJoin(path))
}
