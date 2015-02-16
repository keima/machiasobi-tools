'use strict';

angular.module('myApp.controller.root', [])
  // ヘッダの折りたたみに特化したCtrl
  .controller('RootCtrl', function ($scope, $cookies) {
    var self = this;
    this.hideHeader = ($cookies.hideHeader === 'true');

    $scope.$watch(function () {
      return self.hideHeader;
    }, function (newVal) {
      $cookies.hideHeader = newVal;
    });
  })

  // Header Controller
  .controller('HeaderCtrl', function ($scope, ApiUrl, User) {
    var self = this;
    this.apiUrl = ApiUrl;
    this.loggedin = User.isLogin();

    $scope.$on(User.BROADCAST_NAME_CHANGED, function () {
      self.loggedin = User.isLogin();
    })
  })

  .controller('TabCtrl', function ($scope, User) {
    var self = this;

    self.isAdmin = User.isAdmin();

    $scope.$on(User.BROADCAST_NAME_CHANGED, function () {
      self.isAdmin = User.isAdmin();
    })
  })
;