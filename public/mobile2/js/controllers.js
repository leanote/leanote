'use strict';

// module里定义controller, 注入service
// service相当于公共服务, golang的service总是与数据打交道
angular.module('myApp.controllers', [])
    .controller('MainCtrl', ['$scope', '$rootScope', '$window', '$location', function ($scope, $rootScope, $window, $location) {
        $scope.slide = '';
        
        var snapper;
		setTimeout(function() {
		    snapper = new Snap({
			    element: document.getElementById('content')
		    });
		}, 10);
		
		function openMenu() {
			// 这个时候, 可以调出slider
			if(snapper.state().state=="left" ){
				snapper.close();
			} else {
				snapper.open('left');
			}
		}
		
		$rootScope.back = function() {
          $scope.slide = 'slide-right';
          snapper.close();
          $window.history.back();
        }
        $rootScope.go = function(path){
          $scope.slide = 'slide-left';
          snapper.close();
          $location.url(path);
        }
        
        // 菜单
        $rootScope.menu = function() {
        	openMenu();
        }
        
        // 国际化
        $rootScope.msg = MSG;
    }])
    .controller('NotebookListCtrl', ['$scope', 'Notebook', function ($scope, Notebook) {
        Notebook.getNotebooks(function(notebooks) {
        	$scope.notebooks = notebooks;
        });
    }])
    .controller('NoteListCtrl', ['$scope', '$routeParams', 'Notebook', function ($scope, $routeParams, Notebook) {
    	if($routeParams.notebookId == "all") {
    		$routeParams.notebookId = "";
	        	$scope.notebook = {Title: getMsg("all")}
    	} else {
	    	Notebook.getNotebook($routeParams.notebookId, function(notebook) {
	        	$scope.notebook = notebook;
	    	});
    	}
        Notebook.getNotes($routeParams.notebookId, function(notes) {
        	$scope.notes = notes;
        });
    }])
    .controller('NoteCtrl', ['$scope','$routeParams', '$sce', 'Note', function ($scope, $routeParams, $sce, Note) {
        Note.getNote($routeParams.noteId, function(note) {
        	if(note.Content && !note.RawContent) {
	        	note.RawContent = $sce.trustAsHtml(note.Content);
        	}
        	$scope.note = note;
        });
    }])