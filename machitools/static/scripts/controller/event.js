"use strict";

angular.module("myApp.controller.event", [])
  .controller('EventInputCtrl', function ($scope, $stateParams, Restangular, Periods) {
    var self = this;

    this.itemId = $stateParams.id || null;

    this.lock = false;    // form lock
    this.alert = null;

    var toDoubleDigits = function (num) {
      num += "";
      if (num.length === 1) {
        num = "0" + num;
      }
      return num;
    };

    // set startAt date block
    this.startAtElements = Periods;
    this.selectedStartAt = {
      date: null,
      time: null,
      setDate: function (dateObj) {
        var y = dateObj.getFullYear(),
          m = dateObj.getMonth() + 1,
          d = dateObj.getDate();

        this.date = new Date(y + "/" + m + "/" + d);
        this.time = toDoubleDigits(dateObj.getHours()) +
        ":" +
        toDoubleDigits(dateObj.getMinutes());
      },
      getDate: function () {
        var date = this.date;
        var time = this.time;

        if (_.isNull(date) || _.isNull(time)) {
          return null;
        }

        var y = date.getFullYear(),
          m = date.getMonth() + 1,
          d = date.getDate();

        console.log(new Date(y + "/" + m + "/" + d + " " + time));
        return new Date(y + "/" + m + "/" + d + " " + time);
      }
    };

    this.item = {
      id: null,
      title: null,
      place: null,
      message: null,

      startAt: null,

      isPublic: false,
      isRunning: false,
      isFinished: false
    };

    if (!_.isNull(this.itemId)) {
      Restangular.all('events').get(self.itemId)
        .then(function (result) {
          self.item = result;
          self.selectedStartAt.setDate(new Date(result.startAt));
        });
    }

    this.click = function () {
      self.item.startAt = self.selectedStartAt.getDate();

      self.lock = true;

      var sendEvent;
      if (_.isNull(self.itemId)) {
        sendEvent = Restangular.all('events').post(self.item);
      } else {
        sendEvent = Restangular.all('events').all(self.itemId).post(self.item);
      }

      sendEvent.then(function (result) {
        self.lock = false;
        self.alert = {type: 'success', msg: '登録に成功しました'}
      }, function (reason) {
        self.lock = false;
        self.alert = {type: 'danger', msg: '登録に失敗しました:' + reason.Error}
      });
    };
  })

  .controller('EventListCtrl', function (Restangular, User, Periods) {
    this.periods = Periods;
    this.now = new Date();
  })
  .controller('EventListDayCtrl', function ($stateParams, Restangular, User, Periods) {
    var self = this;

    this.isAdmin = User.isAdmin();

    var startAt = moment(Periods[$stateParams.id].date);
    var endAt = startAt.clone().endOf('days');

    Restangular.all('events').getList({
      first: 0,
      size: 100,
      private: true,
      startAt: startAt.toJSON(),
      endAt: endAt.toJSON()
    }).then(function (results) {
      self.items = results;
    }, function (reason) {
      console.log(reason);
    });
  })
;