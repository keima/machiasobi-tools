"use strict";

angular.module("myApp.controller.delay", [])
  .value("DelayPlaces", [ // arrayで書かないと順番保証されない
    {
      id: "shinmachi",
      name: "新町橋公園",
      calendarId: "p-side.net_ctrq60t4vsvfavejbkdmbhv3k4@group.calendar.google.com"
    },
    {
      id: "ryogoku",
      name: "両国橋公園",
      calendarId: "p-side.net_timelrcritenrfmn86lco3qt9o@group.calendar.google.com"
    },
    {
      id: "bizan",
      name: "眉山林間ステージ",
      calendarId: "p-side.net_m9s9a5ut02n6ap1s6prdj92ss4@group.calendar.google.com"
    },
    {
      id: "corne",
      name: "コルネの泉",
      calendarId: "p-side.net_jo112m9l36p6nlkrv939sb9kr0@group.calendar.google.com"
    },
    {
      id: "cinema_entry",
      name: "CINEMA前(入り口)",
      calendarId: 'p-side.net_j3mtcq3ejulrovek8kru6vgoe8@group.calendar.google.com'
    },
    {
      id: "awagin",
      name: "あわぎんホール小ホール",
      calendarId: 'p-side.net_oa45stb6g4h9lqiq5vd1ov844s@group.calendar.google.com'
    },
    {
      id: "bunka",
      name: "徳島市立文化センター",
      calendarId: 'p-side.net_gocec2ij5sqho46oial3jusn1o@group.calendar.google.com'
    }
  ])
  .controller('DelayViewCtrl', function (Restangular, Calendar, DelayPlaces) {
    var self = this;
    this.now = new Date();

    this.abs = function (value) {
      return Math.abs(value);
    };

    this.places = DelayPlaces;

    // calendar data storage
    this.calendarData = {};

    angular.forEach(this.places, function (place, i) {
      Restangular.all('delay').get(place.id)
        .then(function (result) {
          place.item = result;

          // 遅れている＝現在時刻から遅れ分引いたものが今やってるイベント
          // subtract = 減算、引く
          var time = moment().subtract(result.delay, "minutes");

          Calendar.getTodayData(place.calendarId, time)
            .then(function (_result) {
              self.calendarData[place.id] = _result;
            });
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
  .controller('DelayInputCtrl', function (Restangular, DelayPlaces) {
    var self = this;

    this.places = DelayPlaces;

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
