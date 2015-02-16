"use strict";

angular.module("myApp.controller.traffic", [])
  .controller('TrafficViewCtrl', function (Restangular) {
    var self = this;
    this.now = new Date();

    this.transits = [
      {
        name: 'ロープウェイ乗り場',
        id: 'ropeway',
        places: [
          {
            name: '山麓駅(阿波おどり会)',
            direction: 'inbound'
          },
          {
            name: '山頂駅',
            direction: 'outbound'
          }
        ]
      },
      {
        name: 'シャトルバス乗り場',
        id: 'bus',
        places: [
          {
            name: '山麓駅(阿波踊り会館 前)',
            direction: 'inbound'
          },
          {
            name: '山頂駅(かんぽの宿 前)',
            direction: 'outbound'
          }
        ]
      }
    ];

    this.transits.forEach(function (transit) {
      var traffic = transit.id;

      transit.places.forEach(function (place) {
        var direction = place.direction;

        Restangular.all('traffic').all(traffic).get(direction)
          .then(function (result) {
            place.item = result;
          }, function () {
            place.item = {
              'Waiting': '---',
              'Message': 'SYSTEM: 取得に失敗しました',
              updatedAt: '---'
            }
          });
      });
    });
  })


  .controller('TrafficInputCtrl', function ($scope, $cookies, Restangular) {
    var self = this;

    // form lock
    this.lock = false;

    this.alert = null;

    this.trafficItem = {
      'Waiting': null,
      'Message': null
    };

    this.traffic = $cookies.traffic;
    this.direction = $cookies.direction;

    this.click = function () {
      self.lock = true;
      Restangular.all('traffic').all(self.traffic).all(self.direction).post(self.trafficItem)
        .then(function (result) {
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

    // メッセージを引き継ぐ
    $scope.$watch(function () {
      return self.traffic;
    }, function (newVal, oldVal) {
      $cookies.traffic = newVal;
      getTrafficMessage();
    });

    $scope.$watch(function () {
      return self.direction;
    }, function (newVal, oldVal) {
      $cookies.direction = newVal;
      getTrafficMessage();
    });

    function getTrafficMessage () {
      if (!_.isUndefined(self.traffic) && !_.isUndefined(self.direction)) {
        Restangular.all('traffic').all(self.traffic).get(self.direction)
          .then(function (result) {
            self.trafficItem.Message = result.Message;
          }, function () {
            self.trafficItem.Message = "";
          });
      }
    }
  })
;