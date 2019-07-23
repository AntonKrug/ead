// +Build ignore
//go:generate go run -tags=dev assets_generate.go
package main

import (
	"encoding/hex"
	"github.com/gabriel-vasile/mimetype"
	"github.com/hoisie/mustache"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type ead_file struct {
	filename string
	metadata string
}

var ead_files []ead_file
var dictionary map[string]string

func checkErr(err error) {
	if err != nil {
		log.Fatal(au.Red(err))
	}
}

func convertToSafeName(original string) (ret string) {
	reg, err := regexp.Compile("[^A-Za-z0-9]+")
	if err != nil {
		log.Println(err)
	}
	return reg.ReplaceAllString(original, "_")
}

func allSafeNames(filename string) (header string, data string, metadata string, finalFilename string) {
	safeName := convertToSafeName(filename)

	header = strings.ToUpper(safeName)
	data = "ead_content_" + strings.ToLower(safeName) + "_data"
	metadata = "ead_content_" + strings.ToLower(safeName) + "_metadata"
	finalFilename = filename + ".h"
	return
}

func producePrettyHex(input string) (ret string) {
	raw := hex.EncodeToString([]byte(input))

	bytesReg := regexp.MustCompile("(.{2})")
	hexString := bytesReg.ReplaceAllString(raw, "0x$1, ")

	if *hexDumpFoarmatingFlag {
		linesReg := regexp.MustCompile("(0x[0-9a-f][0-9a-f],.){4}")
		return linesReg.ReplaceAllString(hexString, "$0"+NEW_LINE+"    ")
	} else {
		return hexString
	}
}

func fixContentType(filename string, mine string) (ret string) {
	ret = strings.TrimSpace(strings.Split(mine, ";")[0]) // cut everything after ;

	supportedExtension := map[string]string{
		"txt":  "text/plain",
		"js":   "text/javascript",
		"css":  "text/css",
		"wasm": "application/wasm",
		"xml":  "text/xml",
		"json": "application/json",
	}

	ext := filepath.Ext(filename)[1:] // Get extension excluding the dot ".html"

	if val, ok := supportedExtension[ext]; ok {
		// if we have replacement for the extension, use it
		ret = val
	}

	return
}

func applyTemplateStandalone(contentName string, templateName string) (ret string, nameMetadata string, nameFinal string) {
	template, _ := readFileEmbedded(templateName)
	wholecontent, size := readFileReal(contentName)
	mine, _, _ := mimetype.DetectFile(contentName)
	mineFixed := fixContentType(contentName, mine)
	nameHeader, nameData, nameMetadata, nameFinal := allSafeNames(contentName)

	dictionary["FILENAME_H"] = nameHeader
	dictionary["EAD_FILENAME_VARIABLE"] = nameData
	dictionary["EAD_FILENAME_VARIABLE_ALL_CAPS"] = strings.ToUpper(nameData)
	dictionary["EAD_FILENAME_VARIABLE_METADATA"] = nameMetadata
	dictionary["ORIGINAL_PATH"] = contentName
	dictionary["DATA_SIZE"] = strconv.FormatInt(size, 10)
	dictionary["WEB_CONTENT_ENCODING"] = "EAD_CONTENT_ENCODING_NONE"
	dictionary["CONTENT-TYPE"] = mineFixed
	dictionary["DATA_CONTENT_HEX_DUMP"] = producePrettyHex(wholecontent)

	if isCompressableExtension(contentName) {
		dictionary["WEB_CONTENT_ENCODING"] = "EAD_CONTENT_ENCODING_GZIP"
	}

	ret = mustache.Render(template, dictionary)

	return
}

func getFinalOutputPath(filename string) string {
	return filepath.ToSlash(filepath.Join(outputDir, *outputContainerFlag, filename))
}

func generateFile(contentName string) {
	var templateSource string

	if *standaloneFiles {
		templateSource = "data_no_meta.h"
	} else {
		templateSource = "data.h"
	}

	content, metadata, filename := applyTemplateStandalone(contentName, templateSource)

	finalPath := getFinalOutputPath(filename)

	writeContentToFile(finalPath, content)
	log.Println("Generated include file for", au.Green(contentName))

	ead_files = append(ead_files, ead_file{metadata: metadata, filename: filename})
}

func setupPaths() {
	sourceDir = *sourceDirFlag

	cwd, err := os.Getwd()
	checkErr(err)
	if sourceDir == "" {
		if *sourceCurrentFlag {
			sourceDir = filepath.Clean(cwd)
		} else {
			log.Fatal(au.Red("You can't use current working directory as source directory without specifying -source_current_folder argument, run -h for help"))
		}
	}

	// Make sure the relative path is first converted into absolute paths
	sourceDir, err = filepath.Abs(sourceDir)
	checkErr(err)

	// Convert even Windows paths into backslash format
	sourceDir = filepath.ToSlash(sourceDir)

	if outputDir == "" {
		outputDir = filepath.Join(sourceDir, "/../")
	}
	outputDir, err = filepath.Abs(outputDir)
	checkErr(err)

	outputDir = filepath.ToSlash(outputDir)

	log.Println("Files will be fetched from", au.Bold(sourceDir), "and the output directory is", au.Bold(outputDir))
	log.Println("The includes will be saved into a container folder", au.Bold(*outputContainerFlag), "inside the output folder, while auxiliary files will get stored directly to root of the output folder")

	err = os.Chdir(sourceDir) // The chdir needs to success or we would work from wrong directories
	checkErr(err)

	log.Println("Cleaning destination include folder", au.Bold(getFinalOutputPath("")))
	directories, err := ioutil.ReadDir(getFinalOutputPath(""))
	for _, directory := range directories {
		os.RemoveAll(getFinalOutputPath(directory.Name()))
	}
}

func generateWholeDirectory() {
	err := filepath.Walk(".",
		func(path string, _ os.FileInfo, err error) error {
			checkErr(err)
			path = filepath.ToSlash(path)
			if fileExists(path) && fileIsRealFile(path) {
				generateFile(path)
			}
			return nil
		})
	checkErr(err)
}

func createAuxiliaryFile(file string, outputfile string) {
	if outputfile == "" {
		outputfile = file
	}
	final_file := filepath.ToSlash(filepath.Join(outputDir, outputfile))

	template, _ := readFileEmbedded(file)
	content := mustache.Render(template, dictionary)

	err := ioutil.WriteFile(final_file, []byte(content), 0664)
	checkErr(err)

	log.Println("Generated auxiliary file ", au.Green(file), "")
}

func createAuxiliaryIfNotExists(file string, outputFile string) {
	if outputFile == "" {
		outputFile = file
	}
	final_file := filepath.Join(outputDir, outputFile)

	if !fileExists(final_file) {
		createAuxiliaryFile(file, outputFile)
	}
}

func generateAuxialaryFiles() {
	createAuxiliaryIfNotExists("ead_helpers.c", "")
	createAuxiliaryIfNotExists("ead_helpers.h", "")
	createAuxiliaryIfNotExists("ead_structures.h", "")

	var includes string
	var metadata string

	for _, item := range ead_files {
		includes = includes + "#include \"" + filepath.ToSlash(filepath.Join(*includePrefixFlag, *outputContainerFlag, item.filename)) + "\"" + NEW_LINE
		metadata = metadata + "    " + item.metadata + "," + NEW_LINE
	}

	dictionary["INCLUDE_FILES"] = includes
	dictionary["METADATA_ENTRIES"] = metadata
	createAuxiliaryFile("ead_collection.h", "")
	createAuxiliaryIfNotExists("gitignore", ".gitignore")
}

func main() {
	dictionary = map[string]string{
		"COPYRIGHT":          "Copyright 2019 Microchip Corporation." + NEW_LINE + " *" + NEW_LINE + " * SPDX-License-Identifier: MIT",
		"EAD_COMMENT_NOTICE": "Auto-generated by EAD tool (v" + EAD_VERSION + "), do not edit these files directly." + NEW_LINE + " * Send bug reports and feature requests to: anton.krug@microchip.com",
		"GENERATED_DATE":     time.Now().Format("02/January/2006 15:04:05"),
		"INCLUDE_PREFIX":     *includePrefixFlag,
		"CONTAINER_FOLDER":   *outputContainerFlag,
	}

	if *customCopyrightFlag != "" {
		dictionary["COPYRIGHT"] = readFileRealRawToString(*customCopyrightFlag)
	}

	setupPaths()

	if *compressHtmlFlag {
		log.Println("Web compression enabled, will gzip the following file extension: js,html,htm,css")
	}

	generateWholeDirectory()

	if *outputAuxiliaryFlag && !*standaloneFiles {
		log.Println("Generating axialary files, all existing files will not get overriden except the ead_collection.h")
		generateAuxialaryFiles()
	}
}
