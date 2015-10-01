'use strict';

angular.module('myApp')
  .run(function ($rootScope, $state, Restangular, User, Calendars, Periods) {
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

    console.log(Calendars, Periods)

  })
;
