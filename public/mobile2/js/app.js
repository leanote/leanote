'use strict';

function log(o) {
	console.log(o);
}

angular.module('myApp', [
    'ngTouch',
    'ngRoute',
    'ngAnimate',
    'myApp.controllers',
    'myApp.services'
]).
config(['$routeProvider', '$interpolateProvider', function ($routeProvider, $interpolateProvider) {
	
	$interpolateProvider.startSymbol('[[');
	$interpolateProvider.endSymbol(']]');
	
    $routeProvider.when('/notebooks', {templateUrl: '/public/mobile2/tpl/notebooks.html', controller: 'NotebookListCtrl'});
    $routeProvider.when('/notes/:notebookId', {templateUrl: '/public/mobile2/tpl/notes.html', controller: 'NoteListCtrl'});
    $routeProvider.when('/note/:noteId', {templateUrl: '/public/mobile2/tpl/note.html', controller: 'NoteCtrl'});
    $routeProvider.otherwise({redirectTo: '/notebooks'});
}]);

