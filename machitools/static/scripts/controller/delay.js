"use strict";

angular.module("myApp.controller.delay", [])
  .controller('DelayViewCtrl', function (Restangular, Calendar) {
    var self = this;
    this.now = new Date();

    this.abs = function (value) {
      return Math.abs(value);
    };

    this.places = {
      bizan: {
        //id: "bizan",
        name: "眉山林間ステージ",
        calendarId: "p-side.net_m9s9a5ut02n6ap1s6prdj92ss4@group.calendar.google.com"
      },
      shinmachi: {
        //id: "shinmachi",
        name: "新町橋東公園",
        calendarId: "p-side.net_ctrq60t4vsvfavejbkdmbhv3k4@group.calendar.google.com"
      },
      corne: {
        //id: "corne",
        name: "コルネの泉",
        calendarId: "p-side.net_jo112m9l36p6nlkrv939sb9kr0@group.calendar.google.com"
      },
      cinema_entry: {
        //id: "cinema_entry",
        name: "CINEMA前(入り口)",
        calendarId: 'p-side.net_j3mtcq3ejulrovek8kru6vgoe8@group.calendar.google.com'
      },
      awagin: {
        //id: "awagin",
        name: "あわぎんホール小ホール",
        calendarId: 'p-side.net_oa45stb6g4h9lqiq5vd1ov844s@group.calendar.google.com'
      },
      bunka: {
        //id: "bunka",
        name: "徳島市立文化センター",
        calendarId: 'p-side.net_gocec2ij5sqho46oial3jusn1o@group.calendar.google.com'
      }
    };

    // calendar data storage
    this.calendarData = {};

    angular.forEach(this.places, function (value, key) {
      Restangular.all('delay').get(key)
        .then(function (result) {
          value.item = result;

          // 遅れている＝現在時刻から遅れ分引いたものが今やってるイベント
          // ということで -1 を掛け算している(delayは遅れが＋で進みがー)
          var time = moment().add(-1 * result.delay, "minutes");

          Calendar.getTodayData(value.calendarId, time)
            .then(function (_result) {
              self.calendarData[key] = _result;
            });
        }, function () {
          value.item = {
            error: true,
            delay: 0,
            message: 'SYSTEM: 取得に失敗しました',
            updatedAt: '---'
          }
        });
    });

  })
  .controller('DelayInputCtrl', function (Restangular) {
    var self = this;

    // form lock
    this.lock = false;
    this.alert = null;

    this.place = null;
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