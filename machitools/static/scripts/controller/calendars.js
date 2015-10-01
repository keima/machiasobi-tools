'use strict';

angular.module("myApp.controller.calendars", [])
  .controller("CalendarsListCtrl", function(Restangular, $timeout) {
    var self = this;

    self.lock = false;
    self.alert = null;
    self.modified = false;

    function fetchList() {
      Restangular.all("calendars").getList()
        .then(function(result) {
          self.items = result;
        });
    }

    fetchList();

    this.move = function(index) {
      self.items.splice(index, 1);
      self.modified = true;
    };

    this.updateOrder = function() {
      this.lock = true;
      var ids = _.pluck(self.items, 'id');
      Restangular.all("calendars").all("order").post(ids)
        .then(function(result) {
          self.lock = false;
          self.modified = false;
          self.alert = {type: 'success', msg: '登録に成功しました'};
          $timeout(fetchList, 1000);
        }, function(reason) {
          self.lock = false;
          self.alert = {type: 'danger', msg: '登録に失敗しました:' + reason.data.Error};
        });
    };

    this.closeAlert = function() {
      self.alert = null;
    };

  })
  .controller("CalendarsInputCtrl", function($stateParams, Restangular, $timeout) {
    var self = this;

    this.lock = false;
    this.alert = null;

    if (!_.isUndefined($stateParams.id)) {
      Restangular.all('calendars').get($stateParams.id)
        .then(function(result) {
          self.item = result;
        });
    } else {
      self.item = {};
    }

    this.click = function() {
      self.lock = true;

      var req;
      if (!_.isUndefined($stateParams.id)) {
        req = self.item.put();
      } else {
        req = Restangular.all('calendars').post(self.item)
      }
      req.then(function() {
        self.lock = false;
        self.alert = {type: 'success', msg: '登録に成功しました'}
      }, function(reason) {
        self.lock = false;
        self.alert = {type: 'danger', msg: '登録に失敗しました:' + reason.data.Error}
      });
    };

    this.closeAlert = function() {
      self.alert = null;
    };

  });
