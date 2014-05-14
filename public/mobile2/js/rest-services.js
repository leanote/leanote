'use strict';

// service 用来共享数据
angular.module('myApp.restServices', ['ngResource'])
    .factory('Employee', ['$resource',
        function ($resource) {
           var r =  $resource('http://localhost:3000/employees/:employeeId', {});
           log("r >>>");
           log(r);
           log("r <<<");
           return r;
        }])

    .factory('Report', ['$resource',
        function ($resource) {
            return $resource('http://localhost:3000/employees/:employeeId/reports', {});
        }]);


