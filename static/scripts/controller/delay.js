"use strict";

angular.module("myApp.controller.delay", [])
  .controller('DelayViewCtrl', function (Restangular, Calendar, Calendars) {
    var self = this;
    this.now = new Date(); // TODO: 不要

    this.abs = function (value) {
      return Math.abs(value);
    };

    this.places = Calendars;

    // calendar data storage
    this.calendarData = {};

    angular.forEach(this.places, function (place, i) {
      Restangular.all('delay').get(place.id)
        .then(function (result) {
          place.item = result;

          // FIXME: カレンダー表示しなくなったので削除予定
          // 遅れている＝現在時刻から遅れ分引いたものが今やってるイベント
          // subtract = 減算、引く
          // var time = moment().subtract(result.delay, "minutes");
          //
          // Calendar.getTodayData(place.calendarId, time)
          //   .then(function (_result) {
          //     self.calendarData[place.id] = _result;
          //   });
        }, function () {
          place.item = {
            error: true,
            delay: 0,
            message: '取得に失敗しました',
            updatedAt: '---'
          }
        });
    });

  })
  .controller('DelayInputCtrl', function($stateParams, Restangular, Calendars) {
    var self = this;

    this.places = Calendars;

    // form lock
    this.lock = false;
    this.alert = null;

    this.place = $stateParams.id || null;
    this.item = {
      delay: 0,
      message: "",
      isPostponed: false
    };

    this.click = function () {
      self.lock = true;
      Restangular.all('delay').all(self.place).post(self.item)
        .then(function () {
          self.lock = false;
          self.alert = {type: 'success', msg: '登録に成功しました'}
        }, function (reason) {
          self.lock = false;
          self.alert = {type: 'danger', msg: '登録に失敗しました:' + reason.Error}
        })
    };

    this.closeAlert = function () {
      self.alert = null;
    };

  })
;
