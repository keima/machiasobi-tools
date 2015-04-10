"use strict";

angular.module("myApp.controller.steps", [])
  .controller("StepsListCtrl", function (Restangular, $timeout) {
    var self = this;

    this.lock = false;
    this.alert = null;
    this.modified = false;

    function update() {
      Restangular.all('steps').getList({
        first: 0,
        size: 100,
        private: true
      }).then(function (results) {
        self.items = results;
      });
    }
    update();

    this.move = function (index) {
      self.items.splice(index, 1);
      self.modified = true;
    };

    this.updateOrder = function(){
      this.lock = true;
      var ids = _.pluck(self.items, 'id');
      Restangular.all("steps").all("order").post(ids)
        .then(function(result){
          self.lock = false;
          self.alert = {type: 'success', msg: '登録に成功しました'};
          $timeout(update, 1000);
        }, function(reason){
          self.lock = false;
          self.alert = {type: 'danger', msg: '登録に失敗しました:' + reason.data.Error};
        });
    };

    this.closeAlert = function() {
      self.alert = null;
    };

  })
  .controller("StepsInputCtrl", function ($stateParams, Restangular) {
    var self = this;

    this.lock = false;
    this.alert = null;

    if (!_.isUndefined($stateParams.id)) {
      Restangular.all('steps').get($stateParams.id)
        .then(function (result) {
          self.item = {
            title: result.title,
            showTitle: result.showTitle,
            type: result.type,
            description: result.description,
            partialId: result.partialId,
            order: result.order,
            isPublic: result.isPublic
          }
        });
    }

    this.click = function () {
      self.lock = true;

      var req = Restangular.all('steps');
      if (!_.isUndefined($stateParams.id)) {
        req = req.all($stateParams.id);
      }
      req.post(self.item)
        .then(function () {
          self.lock = false;
          self.alert = {type: 'success', msg: '登録に成功しました'}
        }, function (reason) {
          self.lock = false;
          self.alert = {type: 'danger', msg: '登録に失敗しました:' + reason.data.Error}
        });
    };

    this.closeAlert = function() {
      self.alert = null;
    };
  });
