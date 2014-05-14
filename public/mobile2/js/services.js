
// services

angular.module('myApp.services', [])
// Note
.factory('Note', ["$http", function($http) {
	var self = this;
	
	self.notesMap = {}; // noteId => note
	
	var e = {};
	
	// 得到Note
	e.getNote = function(noteId, callback) {
		if(self.notesMap[noteId]) {
			if(!self.notesMap[noteId].Content) {
				$http.get("/note/getNoteContent", {params:{noteId: noteId}}).success(function(ret) {
					$.extend(self.notesMap[noteId], ret);
					callback(self.notesMap[noteId]);
				});
			}
			callback(self.notesMap[noteId]);
		} else {
			$http.get("/note/getNoteAndContent", {params:{noteId: noteId}}).success(function(ret) {
				self.notesMap[noteId] = ret;
				callback(self.notesMap[noteId]);
			});
		}
	}
	
	// Notebook调用
	e.setCache = function(notes) {
		for(var i in notes) {
			var note = notes[i];
			self.notesMap[note.NoteId] = note;
		}
	}
	
	return e;
}])

// Notebook
.factory('Notebook', ["$http", "Note",function ($http, Note) {
	var self = this;
	self.notebooks = [];
	self.notebooksMap = {};
	
	self.notes = {}; // notebookId=>[note]
	self.notesMap = {}; // noteId => note
	
	var e = {};
	e.getNotebooks = function(callback) {
		if(!callback) {
			callback = function() {};
		}
		if(self.notebooks.length > 0) callback(self.notebooks);
		else {
			$http.get("/notebook/getNotebooks").success(function(ret) {
				self.notebooks = ret;
				for(var i in ret) {
					var notebook = ret[i];
					self.notebooksMap[notebook.NotebookId] = notebook;
				}
				callback(ret);
			});
		}
	}
	
	e.getNotebook = function(notebookId, callback) {
		if(self.notebooksMap[notebookId]) callback(self.notebooksMap[notebookId]);
		else {
			e.getNotebooks(function() {
				callback(self.notebooksMap[notebookId]);
			});
		}
	}
	
	e.getNotes = function(notebookId, callback) {
		if(self.notes[notebookId]) callback(self.notes[notebookId]);
		else {
			$http.get("/note/listNotes", {params:{notebookId: notebookId}}).success(function(ret) {
				self.notes[notebookId] = ret;
				callback(ret);
				// 存到Note中
				Note.setCache(ret);
			});
		}
	}
	
    return e;
}]);