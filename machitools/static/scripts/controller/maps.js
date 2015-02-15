'use strict';

angular.module('myApp.controller.maps', [])
  .controller('MapsListCtrl', function (Restangular, User) {
    var self = this;

    this.isAdmin = User.isAdmin();

    Restangular.all('maps').getList({
      first: 0,
      size: 100,
      private: true
    }).then(function (results) {
      self.items = results;
    }, function (reason) {
      console.log(reason);
    });

  })
  .controller('MapsDetailCtrl', function ($stateParams, Restangular, $timeout, $window) {
    var self = this;

    function init () {
      Restangular.all('maps').get($stateParams.id)
        .then(function (result) {
          self.item = result;
        });
    }

    init();

    this.deleteMarker = function (index) {
      var marker = self.item.markers[index];
      var msg = "「" + marker.name + "」を削除しても宜しいですか？\n（削除後の復元は出来ません！）";
      if ($window.confirm(msg)) {
        Restangular.one('maps', $stateParams.id).one('markers', marker.id).remove()
          .then(function () {
            $timeout(function () {
              init();
            }, 1000);
          });
      }
    };

  })
  .controller('MapsInputCtrl', function ($stateParams, Restangular) {
    var self = this;

    this.itemIdParam = $stateParams.id || null;

    this.lock = false;    // form lock
    this.alert = null;

    this.itemId = null;
    this.item = {
      'name': null,
      'isPublic': false
    };

    if (!_.isNull(this.itemIdParam)) {
      Restangular.all('maps').get(self.itemIdParam)
        .then(function (result) {
          self.itemId = result.id;
          self.item = {
            'name': result.name,
            'isPublic': result.isPublic
          }
        });
    }

    this.click = function () {
      self.lock = true;
      Restangular.all('maps').all(self.itemId).post(self.item)
        .then(function (result) {
          self.lock = false;
          self.alert = {type: 'success', msg: '登録に成功しました'}
        }, function (reason) {
          self.lock = false;
          self.alert = {type: 'danger', msg: '登録に失敗しました:' + reason.Error}
        });
    };
  })
  .controller('MapMarkerInputCtrl', function ($scope, $stateParams, Restangular) {
    var self = this;
    this.lock = false;    // form lock
    this.alert = null;
    this.itemIdParam = $stateParams.id || null;

    this.place = {
      name: null,
      description: null,
      order: null
    };

    this.map = {
      center: {latitude: 34.071144, longitude: 134.548529},
      zoom: 18
    };

    this.marker = {
      id: 0,
      coords: {
        latitude: 34.071144,
        longitude: 134.548529
      },
      options: {
        draggable: true
      },
      events: {
        'dragend': function (marker, eventName, args) {
          var lat = marker.getPosition().lat();
          var lon = marker.getPosition().lng();
        },
        'tilesloaded': function (map, eventName, args) {
          console.log(map, eventName, args);

          $scope.$apply(function () {
            google.maps.events.trigger($scope.map.control.getGMap(), 'resize');
          });
        }
      }
    };

    this.postMapMarker = function () {
      self.lock = true;
      Restangular.all('maps').all(self.itemIdParam).all('markers').post({
        name: self.place.name,
        description: self.place.description,
        order: self.place.order,
        coords: {
          latitude: self.marker.coords.latitude,
          longitude: self.marker.coords.longitude
        }
      })
        .then(function (result) {
          self.lock = false;
          self.alert = {type: 'success', msg: '登録に成功しました'}
        }, function (reason) {
          self.lock = false;
          self.alert = {type: 'danger', msg: '登録に失敗しました:' + reason.Error}
        });
    };

    this.moveMarkerToCenter = function () {
      self.marker.coords = {
        latitude: self.map.center.latitude,
        longitude: self.map.center.longitude
      };
    };

  })
;