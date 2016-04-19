"use strict";

angular.module('myApp.services.user', [])
  .service('User', function ($rootScope) {
    var user = {};
    var BROADCAST_NAME_CHANGED = 'UserDataIsChanged';

    function setUser (_user) {
      user = _user;
      $rootScope.$broadcast(BROADCAST_NAME_CHANGED);
    }

    function getUser () {
      return user;
    }

    function isLogin () {
      return !_.isEmpty(user)
    }

    function isAdmin () {
      return isLogin() && user.Admin;
    }

    return {
      BROADCAST_NAME_CHANGED: BROADCAST_NAME_CHANGED,

      setUser: setUser,
      getUser: getUser,
      isLogin: isLogin,
      isAdmin: isAdmin
    }
  })
;