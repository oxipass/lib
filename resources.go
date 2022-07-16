package oxilib

import (
	"embed"
	"encoding/json"
	"errors"
	"io/fs"
)

var (
	//go:embed langs/*.json
	langsResources embed.FS

	//go:embed templates/*.json
	templatesResources embed.FS
)

var tr map[string]string

func getLangsFiles() ([]fs.DirEntry, error) {
	return langsResources.ReadDir(cLangsFolder)
}

func initLang(langCode string) error {
	files, err := getLangsFiles()
	if err != nil {
		return err
	}
	for _, file := range files {
		var lang Lang
		fileBytes, errFile := langsResources.ReadFile(cLangsFolder + "/" + file.Name())
		if errFile != nil {
			return errFile
		}
		errUnmarshal := json.Unmarshal(fileBytes, &lang)
		if errUnmarshal != nil {
			return errUnmarshal
		}
		if lang.Code == langCode {
			var translations Translations
			transBytes, errTrans := langsResources.ReadFile(cLangsFolder + "/" + file.Name())
			if errTrans != nil {
				return errTrans
			}
			errUnmarshalTrans := json.Unmarshal(transBytes, &translations)
			if errUnmarshalTrans != nil {
				return errUnmarshalTrans
			}
			tr = translations.Translations
			return nil
		}
	}
	if langCode != "en" {
		return initLang("en")
	}
	return errors.New("language not loaded")
}

func TR(key string) string {
	if tr == nil {
		err := initLang("en")
		if err != nil {
			return key
		}
	}
	return tr[key]
}

func setLang(langCode string) error {
	err := initLang(langCode)
	if err != nil {
		return err
	}
	return nil
}

func getLangs() []Lang {
	files, err := getLangsFiles()
	if err != nil {
		return nil
	}
	langs := make([]Lang, len(files))
	for i, file := range files {
		var lang Lang
		fileBytes, errFile := langsResources.ReadFile(cLangsFolder + "/" + file.Name())
		if errFile != nil {
			return nil
		}
		errUnmarshal := json.Unmarshal(fileBytes, &lang)
		if errUnmarshal != nil {
			return nil
		}
		langs[i] = lang
	}
	return langs
}
