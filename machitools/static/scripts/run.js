'use strict';

angular.module('myApp')
  .run(function ($rootScope, $state, Restangular, User) {
    // convenience state
    $rootScope.$state = $state;
    $rootScope.isAdmin = false;

    // get user data
    Restangular.all('auth').get('check')
      .then(function (result) {
        User.setUser(result);
        $rootScope.isAdmin = User.isAdmin();
      }, function (reason) {
        console.log(reason);
        User.setUser({});
      });

  })
;