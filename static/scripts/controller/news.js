"use strict";

angular.module("myApp.controller.news", [])
  .controller('NewsInputCtrl', function ($stateParams, Restangular) {
    var self = this;

    this.itemIdParam = $stateParams.id || null;

    this.lock = false;    // form lock
    this.alert = null;

    this.newsId = null;
    this.newsItem = {
      'Title': null,
      'Article': null,
      'IsPublic': false
    };

    if (!_.isNull(self.itemIdParam)) {
      Restangular.all('news').get(self.itemIdParam)
        .then(function (result) {
          self.newsId = result.id;
          self.newsItem = result
        });
    }

    this.click = function () {
      self.lock = true;
      var req = Restangular.all('news').all(self.newsId);
      if (_.isNull(self.itemIdParam)) {
        req = req.post(self.newsItem);
      } else {
        req = self.newsItem.put();
      }

      req.then(function (result) {
        self.lock = false;
        self.alert = {type: 'success', msg: '登録に成功しました'}
      }, function (reason) {
        self.lock = false;
        self.alert = {type: 'danger', msg: '登録に失敗しました:' + reason.Error}
      });
    };
  })
  .controller('NewsListCtrl', function (Restangular) {
    var self = this;

    Restangular.all('news').getList({
      first: 0,
      size: 100,
      private: true
    })
      .then(function (results) {
        self.items = results;
      }, function (reason) {
        console.log(reason);
      })
  })
;
