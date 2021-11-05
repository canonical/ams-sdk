// -*- Mode: Go; indent-tabs-mode: t -*-
/*
 * This file is part of AMS SDK
 * Copyright 2021 Canonical Ltd.
 *
 * This program is free software: you can redistribute it and/or modify it under
 * the terms of the GNU Lesser General Public License version 3, as published
 * by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful, but WITHOUT
 * ANY WARRANTY; without even the implied warranties of MERCHANTABILITY, SATISFACTORY
 * QUALITY, or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU General Public
 * License for more details.
 *
 * You should have received a copy of the Lesser GNU General Public License along
 * with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package shared

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"unicode"
	"unsafe"

	yaml "gopkg.in/yaml.v2"
)

// VarPath returns the provided path elements joined by a slash and
// appended to the end of $SNAP_COMMON, which defaults to /var/lib/ams.
func VarPath(path ...string) string {
	varDir := os.Getenv("SNAP_COMMON")
	if varDir == "" {
		varDir = "/var/lib/ams"
	}

	items := []string{varDir}
	items = append(items, path...)
	varPath := filepath.Join(items...)
	os.MkdirAll(filepath.Dir(varPath), 0755)
	return varPath
}

// GetOwnerMode retrieves the file mode and owner of a given file object
func GetOwnerMode(fInfo os.FileInfo) (os.FileMode, int, int) {
	mode := fInfo.Mode()
	uid := int(fInfo.Sys().(*syscall.Stat_t).Uid)
	gid := int(fInfo.Sys().(*syscall.Stat_t).Gid)
	return mode, uid, gid
}

// FileCopy copies a file, overwriting the target if it exists.
func FileCopy(source string, dest string) error {
	s, err := os.Open(source)
	if err != nil {
		return err
	}
	defer s.Close()

	fi, err := s.Stat()
	if err != nil {
		return err
	}

	d, err := os.Create(dest)
	if err != nil {
		if os.IsExist(err) {
			d, err = os.OpenFile(dest, os.O_WRONLY, fi.Mode())
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	defer d.Close()

	_, err = io.Copy(d, s)
	if err != nil {
		return err
	}

	// Do not follow the original file' owner until user and group management
	// is supported by snapd as the syscall chown is not allowed in the strict
	// confinement mode and snapd doesn't have seccomp argument filtering for now
	if !RunningAsSnap() {
		_, uid, gid := GetOwnerMode(fi)
		return d.Chown(uid, gid)
	}
	return err
}

// FileMove tries to move a file by using os.Rename,
// if that fails it tries to copy the file and remove the source.
func FileMove(oldPath string, newPath string) error {
	err := os.Rename(oldPath, newPath)
	if err == nil {
		return nil
	}

	err = FileCopy(oldPath, newPath)
	if err != nil {
		return err
	}

	os.Remove(oldPath)

	return nil
}

// ListFilesInDir lists all files in one specific directory recursively
func ListFilesInDir(dirPath string, recursive bool) ([]string, error) {
	filelist := []string{}
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		if f.IsDir() {
			if recursive {
				list, err := ListFilesInDir(filepath.Join(dirPath, f.Name()), true)
				if err != nil {
					return []string{}, err
				}
				filelist = append(filelist, list...)
			}
		}
		filelist = append(filelist, filepath.Join(dirPath, f.Name()))
	}

	return filelist, nil
}

// DirCopy copies a directory tree recursively
func DirCopy(source string, dest string) error {
	src := filepath.Clean(source)
	dst := filepath.Clean(dest)

	fi, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !fi.IsDir() {
		return fmt.Errorf("source is not a directory")
	}

	if _, err = os.Stat(dst); err == nil {
		return fmt.Errorf("destination already exists")
	}

	if err = os.MkdirAll(dst, fi.Mode()); err != nil {
		return err
	}

	files, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}

	for _, file := range files {
		srcPath := filepath.Join(src, file.Name())
		dstPath := filepath.Join(dst, file.Name())

		if file.IsDir() {
			if err = DirCopy(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			// Skip symlinks
			if file.Mode()&os.ModeSymlink != 0 {
				continue
			}

			if err = FileCopy(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// PathExists returns true if param path exists
func PathExists(name string) bool {
	_, err := os.Lstat(name)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

// StringInSlice returns true if the searched string key is into slice
func StringInSlice(key string, list []string) bool {
	for _, entry := range list {
		if entry == key {
			return true
		}
	}
	return false
}

// ParseByteSizeString parses a size string in bytes (e.g. 200kB or 5GB) into the number of
// bytes it represents. Supports suffixes up to EB. "" == 0.
func ParseByteSizeString(input string) (int64, error) {
	suffixLen := 2

	if input == "" {
		return 0, nil
	}

	if unicode.IsNumber(rune(input[len(input)-1])) {
		// No suffix --> bytes.
		suffixLen = 0
	} else if (len(input) >= 2) && (input[len(input)-1] == 'B') && unicode.IsNumber(rune(input[len(input)-2])) {
		// "B" suffix --> bytes.
		suffixLen = 1
	} else if strings.HasSuffix(input, " bytes") {
		// Backward compatible behaviour in case we talk to a LXD that
		// still uses GetByteSizeString() that returns "n bytes".
		suffixLen = 6
	} else if (len(input) < 3) && (suffixLen == 2) {
		return -1, fmt.Errorf("Invalid value: %s", input)
	}

	// Extract the suffix
	suffix := input[len(input)-suffixLen:]

	// Extract the value
	value := input[0 : len(input)-suffixLen]
	valueInt, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return -1, fmt.Errorf("Invalid integer: %s", input)
	}

	if valueInt < 0 {
		return -1, fmt.Errorf("Invalid value: %d", valueInt)
	}

	// The value is already in bytes.
	if suffixLen != 2 {
		return valueInt, nil
	}

	// Figure out the multiplicator
	multiplicator := int64(0)
	switch suffix {
	case "kB":
		multiplicator = 1024
	case "MB":
		multiplicator = 1024 * 1024
	case "GB":
		multiplicator = 1024 * 1024 * 1024
	case "TB":
		multiplicator = 1024 * 1024 * 1024 * 1024
	case "PB":
		multiplicator = 1024 * 1024 * 1024 * 1024 * 1024
	case "EB":
		multiplicator = 1024 * 1024 * 1024 * 1024 * 1024 * 1024
	default:
		return -1, fmt.Errorf("Unsupported suffix: %s", suffix)
	}

	return valueInt * multiplicator, nil
}

// GetByteSizeString returns a string with the given byte size formatted and the
// right unit added.
func GetByteSizeString(input int64, precision uint) string {
	if input < 1024 {
		return fmt.Sprintf("%dB", input)
	}

	value := float64(input)

	for _, unit := range []string{"kB", "MB", "GB", "TB", "PB", "EB"} {
		value = value / 1024
		if value < 1024 {
			return fmt.Sprintf("%.*f%s", precision, value, unit)
		}
	}

	return fmt.Sprintf("%.*fEB", precision, value)
}

// LoadFromFile loads a configuration from the given file path
func LoadFromFile(configPath string, cfg interface{}) error {
	bytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(bytes, cfg)
}

// ValueOrDefault returns the value if it has a size greater than zero or defaultValue otherwise
func ValueOrDefault(value string, defaultValue string) string {
	if len(value) > 0 {
		return value
	}
	return defaultValue
}

// CreateBzip2Tarball creates a bzip2 tarball file with given files path
func CreateBzip2Tarball(workingDir, outputPath string, content []string) error {
	comopressArg := []string{
		"cfj", outputPath,
		"-C", workingDir,
	}
	comopressArg = append(comopressArg, content...)
	return exec.Command("tar", comopressArg...).Run()
}

// GetFileSize returns the size of the file at the given path
func GetFileSize(path string) (int64, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return fi.Size(), nil
}

// RunningAsSnap detects if the ams service runs as a snap
func RunningAsSnap() bool {
	_, exist := os.LookupEnv("SNAP")
	return exist
}

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}
	return b, nil
}

// GenerateRandomString returns a URL-safe, base64 encoded
// securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomString(s int) (string, error) {
	b, err := GenerateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}

// RandomCryptoString returns a random base64 encoded string from crypto/rand.
func RandomCryptoString() (string, error) {
	buf := make([]byte, 32)
	n, err := rand.Read(buf)
	if err != nil {
		return "", err
	}

	if n != len(buf) {
		return "", fmt.Errorf("not enough random bytes read")
	}

	return hex.EncodeToString(buf), nil
}

// SetSize sets the size of the PTY connected with the given file descriptor
func SetSize(fd int, width int, height int) (err error) {
	var dimensions [4]uint16
	dimensions[0] = uint16(height)
	dimensions[1] = uint16(width)

	if _, _, err := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd), uintptr(syscall.TIOCSWINSZ), uintptr(unsafe.Pointer(&dimensions)), 0, 0, 0); err != 0 {
		return err
	}
	return nil
}

// GenerateFingerprintForFile generates a fingerprint (sha256) for the given file
func GenerateFingerprintForFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	return GenerateFingerprint(f)
}

// GenerateFingerprint generates a fingerprint for the given reader
func GenerateFingerprint(reader io.ReadSeeker) (string, error) {
	reader.Seek(0, 0)
	hasher := sha256.New()
	_, err := io.Copy(hasher, reader)
	if err != nil {
		return "", err
	}
	reader.Seek(0, 0)
	fingerprint := fmt.Sprintf("%x", hasher.Sum(nil))
	return fingerprint, nil
}

// StripUserPasswordFromURL removes user/password from the given URL
func StripUserPasswordFromURL(value string) string {
	u, err := url.Parse(value)
	if err != nil {
		return "invalid"
	}
	if u.User != nil {
		// We strip user/password from the URL here to ensure it's not printed
		u.User = url.UserPassword("xxx", "xxx")
	}
	return u.String()
}

type cancelableReader struct {
	ctx   context.Context
	other io.Reader
}

func (r *cancelableReader) Read(p []byte) (n int, err error) {
	select {
	case <-r.ctx.Done():
		return 0, r.ctx.Err()
	default:
		return r.other.Read(p)
	}
}

// NewCancelableReader creates a new cancelable io.Reader
func NewCancelableReader(ctx context.Context, other io.Reader) io.Reader {
	return &cancelableReader{
		ctx:   ctx,
		other: other,
	}
}

// CompareStringArrays compares two string arrays and returns true if two string
// arrays are equal to one another, otherwise returns false.
func CompareStringArrays(s1, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}

	for _, a := range s1 {
		exist := false
		for _, b := range s2 {
			if a == b {
				exist = true
				break
			}
		}

		if !exist {
			return false
		}
	}

	return true
}

// CompareIntArrays compares two int arrays and returns true if two int
// arrays are equal to one another, otherwise returns false.
func CompareIntArrays(s1, s2 []int) bool {
	if len(s1) != len(s2) {
		return false
	}

	for i := 0; i < len(s1); i++ {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}

// ImageArchToNodeArch converts the architecture of an image which is used
// in simplestream to node machine architecture
func ImageArchToNodeArch(Arch string) string {
	switch Arch {
	case "arm64":
		return "aarch64"
	case "amd64":
		return "x86_64"
	default:
		return "unknown"
	}
}

// NodeArchToImageArch converts the architecture of node machine architecture
// to the architecture of an image which is used in simplestream
func NodeArchToImageArch(Arch string) string {
	switch Arch {
	case "aarch64":
		return "arm64"
	case "x86_64":
		return "amd64"
	default:
		return "unknown"
	}
}

// BinaryEndian returns machine native binary endian
func BinaryEndian() binary.ByteOrder {
	const intSize = int(unsafe.Sizeof(0))
	i := 0x1
	bs := (*[intSize]byte)(unsafe.Pointer(&i))
	if bs[0] == 0 {
		return binary.BigEndian
	}
	return binary.LittleEndian
}
