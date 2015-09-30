"use strict";

angular.module("myApp.controller.periods", [])
  .controller('PeriodsCtrl', function($stateParams, Restangular, $timeout) {
    var self = this;

    this.form = {};

    function updateList() {
      Restangular.all('periods').getList( {time: new Date().getTime()})
        .then(function(result) {
          self.items = result;
        });
    }

    updateList();

    this.postPeriod = function() {
      if (_.isUndefined(self.form.date)) {
        return;
      }

      console.log(self.form.date);

      Restangular.all('periods').post({
        date: moment(self.form.date).format()
      }).then(function(result) {
        $timeout(updateList, 500);
      });
    };

    this.deactivatePeriod = function(item) {
      if (_.isUndefined(item)) {
        return;
      }

      item.post("deactivate")
        .then(function() {
          $timeout(updateList, 500);
        }, function(error) {
          console.log(error);
        })
    };

  }
);
