package service

import (

)

// init service, for share service bettween services
// 初始化, 实例service
// 为了共享service

var notebookService, NotebookS *NotebookService
var noteService, NoteS *NoteService
var noteContentHistoryService, NoteContentHistoryS *NoteContentHistoryService
var trashService, TrashS *TrashService
var shareService, ShareS *ShareService
var userService, UserS *UserService
var tagService, TagS *TagService
var blogService, BlogS *BlogService
var tokenService, TokenS *TokenService
var noteImageService, NoteImageS *NoteImageService
var fileService, FileS *FileService
var albumService, AlbumS *AlbumService
var attachService, AttachS *AttachService
var PwdS *PwdService
var SuggestionS *SuggestionService
var AuthS *AuthService

// onAppStart调用
func InitService() {
	NotebookS = &NotebookService{}
	NoteS = &NoteService{}
	NoteContentHistoryS = &NoteContentHistoryService{}
	TrashS = &TrashService{}
	ShareS = &ShareService{}
	UserS = &UserService{}
	TagS = &TagService{}
	BlogS = &BlogService{}
	TokenS = &TokenService{}
	NoteImageS = &NoteImageService{}
	FileS = &FileService{}
	AlbumS = &AlbumService{}
	AttachS = &AttachService{}
	PwdS = &PwdService{}
	SuggestionS = &SuggestionService{}
	AuthS = &AuthService{}
	
	notebookService = NotebookS
	noteService = NoteS
	noteContentHistoryService = NoteContentHistoryS
	trashService = TrashS
	shareService = ShareS
	userService = UserS
	tagService = TagS
	blogService = BlogS
	tokenService = TokenS
	noteImageService = NoteImageS
	fileService = FileS
	albumService = AlbumS
	attachService = AttachS
}