package src

import (
	"fmt"
	"path"
	"sf/cfg"
	"sf/multi"
	"sf/src/boluobao"
	"sf/src/hbooker"
	"sf/structural"
	"sf/structural/hbooker_structs"
	"sf/structural/sfacg_structs"
	"strconv"
)

type BookInits struct {
	BookID   string
	Index    int
	ShowBook bool
	Locks    *multi.GoLimit
	BookData any
}

func (books *BookInits) DownloadBookInit() Catalogue {
	if cfg.Vars.AppType == "sfacg" {
		response := boluobao.GetBookDetailedById(books.BookID)
		if response.Status.HTTPCode == 200 && response.Data.NovelName != "" {
			books.BookData = response
		} else {
			fmt.Println(books.BookID, "is not a valid book number！\nmessage:", response.Status.Msg)
			return Catalogue{TestBookResult: false}
		}
	} else if cfg.Vars.AppType == "cat" {
		response := hbooker.GetBookDetailById(books.BookID)
		if response.Code == "100000" {
			books.BookData = response
		} else {
			fmt.Println(books.BookID, "is not a valid book number！")
		}
	} else {
		panic("app type" + cfg.Vars.AppType + " is not valid!")
	}
	cfg.BookConfig.BookInfo = books.InitBookStruct()

	savePath := path.Join(cfg.Vars.SaveFile, cfg.BookConfig.BookInfo.NovelName+".txt")
	if !cfg.CheckFileExist(savePath) {
		cfg.EncapsulationWrite(savePath, books.ShowBookDetailed()+"\n\n", "w")
	} else {
		books.ShowBookDetailed()
	}
	return Catalogue{SaveTextPath: savePath, TestBookResult: true}

}
func (books *BookInits) InitBookStruct() structural.Books {
	switch books.BookData.(type) {
	case sfacg_structs.BookInfo:
		result := books.BookData.(sfacg_structs.BookInfo).Data
		return structural.Books{
			NovelName:  cfg.RegexpName(result.NovelName),
			NovelID:    strconv.Itoa(result.NovelID),
			IsFinish:   result.IsFinish,
			MarkCount:  strconv.Itoa(result.MarkCount),
			NovelCover: result.NovelCover,
			AuthorName: result.AuthorName,
			CharCount:  strconv.Itoa(result.CharCount),
			SignStatus: result.SignStatus,
		}
	case hbooker_structs.DetailStruct:
		result := books.BookData.(hbooker_structs.DetailStruct).Data.BookInfo
		return structural.Books{
			NovelName:  cfg.RegexpName(result.BookName),
			NovelID:    result.BookID,
			NovelCover: result.Cover,
			AuthorName: result.AuthorName,
			CharCount:  result.TotalWordCount,
			MarkCount:  result.UpdateStatus,
			//SignStatus: result.SignStatus,
		}
	}
	return structural.Books{}
}

func (books *BookInits) ShowBookDetailed() string {
	briefIntroduction := fmt.Sprintf(
		"Name: %v\nBookID: %v\nAuthor: %v\nCount: %v\nMark: %v\n",
		cfg.BookConfig.BookInfo.NovelName, cfg.BookConfig.BookInfo.NovelID,
		cfg.BookConfig.BookInfo.AuthorName, cfg.BookConfig.BookInfo.CharCount,
		cfg.BookConfig.BookInfo.MarkCount,
	)
	if books.ShowBook {
		fmt.Println(briefIntroduction)
	}
	return briefIntroduction
}
