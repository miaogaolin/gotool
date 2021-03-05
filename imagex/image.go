package imagex

import (
	"bytes"
	"errors"
	"github.com/miaogaolin/gotool/rest"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"github.com/miaogaolin/gotool/filex"

	"github.com/disintegration/imaging"
)

var defaultTransferImageSize = 150
var ext = []string{
	"ase",
	"art",
	"bmp",
	"blp",
	"cd5",
	"cit",
	"cpt",
	"cr2",
	"cut",
	"dds",
	"dib",
	"djvu",
	"egt",
	"exif",
	"gif",
	"gpl",
	"grf",
	"icns",
	"ico",
	"iff",
	"jng",
	"jpeg",
	"jpg",
	"jfif",
	"jp2",
	"jps",
	"lbm",
	"max",
	"miff",
	"mng",
	"msp",
	"nitf",
	"ota",
	"pbm",
	"pc1",
	"pc2",
	"pc3",
	"pcf",
	"pcx",
	"pdn",
	"pgm",
	"PI1",
	"PI2",
	"PI3",
	"pict",
	"pct",
	"pnm",
	"pns",
	"ppm",
	"psb",
	"psd",
	"pdd",
	"psp",
	"px",
	"pxm",
	"pxr",
	"qfx",
	"raw",
	"rle",
	"sct",
	"sgi",
	"rgb",
	"int",
	"bw",
	"tga",
	"tiff",
	"tif",
	"vtf",
	"xbm",
	"xcf",
	"xpm",
	"3dv",
	"amf",
	"ai",
	"awg",
	"cgm",
	"cdr",
	"cmx",
	"dxf",
	"e2d",
	"egt",
	"eps",
	"fs",
	"gbr",
	"odg",
	"svg",
	"stl",
	"vrml",
	"x3d",
	"sxd",
	"v2d",
	"vnd",
	"wmf",
	"emf",
	"art",
	"xar",
	"png",
	"webp",
	"jxr",
	"hdp",
	"wdp",
	"cur",
	"ecw",
	"iff",
	"lbm",
	"liff",
	"nrrd",
	"pam",
	"pcx",
	"pgf",
	"sgi",
	"rgb",
	"rgba",
	"bw",
	"int",
	"inta",
	"sid",
	"ras",
	"sun",
	"tga",
}

type Stat struct {
	Mime string
	Ext  string
}

func IsExt(e string) bool {
	e = strings.Trim(e, ".")
	for i := range ext {
		if ext[i] == e {
			return true
		}
	}
	return false
}

func Mime(file io.ReadSeeker) (*Stat, error) {
	file.Seek(0, 0)
	defer file.Seek(0, 0)
	all, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	fileCopy := bytes.NewReader(all)

	buff := make([]byte, 512)
	_, err = fileCopy.Read(buff)
	if err != nil {
		return nil, err
	}

	fileType := http.DetectContentType(buff)

	for i := 0; i < len(ext); i++ {
		if strings.Contains(ext[i], fileType[6:]) {
			return &Stat{
				Mime: fileType,
				Ext:  ext[i],
			}, nil
		}
	}

	return nil, errors.New("invalid image type")

}

// min(width, height) = size
func TransferJPEG(img io.ReadSeeker, size int) (*os.File, error) {
	imgContainer, err := imaging.Decode(img)
	if err != nil {
		return nil, err
	}
	bounds := imgContainer.Bounds()
	var width, height int
	if bounds.Dx() > bounds.Dy() {
		height = size
	} else {
		width = size
	}

	dstImage := imaging.Resize(imgContainer, width, height, imaging.Lanczos)
	dstFilename := TempDir() + filex.RandomName() + ".jpg"
	err = imaging.Save(dstImage, dstFilename)
	if err != nil {
		return nil, err
	}
	return os.Open(dstFilename)
}

func TransferDefault(img io.ReadSeeker) (*os.File, error) {
	return TransferJPEG(img, defaultTransferImageSize)
}

func Download(imgUrl string) (*os.File, error) {
	imgUrl = strings.Replace(imgUrl, "_.webp", "", 1)
	imgUrl = strings.Replace(imgUrl, "_webp", "", 1)
	ext := filepath.Ext(imgUrl)
	resp, err := rest.Client.Get(rest.Request{
		Url: imgUrl,
	})

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	imageName := TempDir() + filex.RandomName() + ext
	err = ioutil.WriteFile(imageName, data, 0644)
	if err != nil {
		return nil, err
	}
	return os.Open(imageName)
}

func TempDir() string {
	dir := ".temp/image/"
	filex.CreateDir(dir)
	return dir
}
