'use strict';

angular.module('myApp.services.periods', [])
  .service('Periods', function (Restangular) {
    var periods = [];

    function getList() {
      if (periods.length == 0) {
        periods = Restangular.all("periods").getList().$object;
      }
      return periods;
    }

    return {
      getList: getList
    };
  });
