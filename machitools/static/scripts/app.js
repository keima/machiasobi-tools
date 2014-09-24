'use strict';

angular.module('myApp', ['ngRoute', 'restangular'])
  .constant('ApiUrl', '/api/v1')
  .config(function(RestangularProvider, ApiUrl) {
    RestangularProvider.setBaseUrl(ApiUrl);
  })
  .config(function($routeProvider) {

  })
  .controller('HeaderCtrl', function(Restangular, ApiUrl) {
    var self = this;

    this.apiUrl = ApiUrl;
    this.loggedin = false;

    Restangular.all('auth').get('check')
      .then(function(result) {
        console.log(result);
        self.loggedin = true;
      }, function(reason) {
        console.log(reason);
        self.loggedin = false;
      });
  });