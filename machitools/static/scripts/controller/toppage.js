"use strict";

angular.module('myApp.controller.topPage', [])
  .controller('TopPageCtrl', function ($scope, User) {
    var self = this;
    this.showError = false;

    function renderAlertIfNeeded () {
      if (User.isAdmin()) {
        self.showError = false;
      } else {
        self.showError = User.isLogin();
      }
    }

    renderAlertIfNeeded();

    $scope.$on(User.BROADCAST_NAME_CHANGED, function () {
      renderAlertIfNeeded();
    });

  })
;