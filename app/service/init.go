package service

import (

)

// init service, for share service bettween services
// 初始化, 实例service
// 为了共享service

var notebookService *NotebookService
var noteService *NoteService
var noteContentHistoryService *NoteContentHistoryService
var trashService *TrashService
var shareService *ShareService
var userService *UserService
var tagService *TagService
var blogService *BlogService
var tokenService *TokenService
var noteImageService *NoteImageService
var fileService *FileService

func init() {
	notebookService = &NotebookService{}
	noteService = &NoteService{}
	noteContentHistoryService = &NoteContentHistoryService{}
	trashService = &TrashService{}
	shareService = &ShareService{}
	userService = &UserService{}
	tagService = &TagService{}
	blogService = &BlogService{}
	tokenService = &TokenService{}
	fileService = &FileService{}
	noteImageService = &NoteImageService{}
}