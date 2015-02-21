'use strict';

angular.module('myApp.controller.maps', [])
  .service('MapsManager', function ($q, Restangular) {
    var idVal;
    var forceReload = false;
    var item = {};

    function reload () {
      var deferred = $q.defer();

      if (_.isUndefined(idVal)) {
        deferred.reject("MapsManager.id is not set!")
      } else {
        Restangular.all('maps').get(idVal)
          .then(function (result) {
            item = result;
            deferred.resolve(result);
          }, function (reason) {
            deferred.reject(reason);
          });
      }

      return deferred.promise;
    }

    function getMap () {
      var deferred = $q.defer();

      if (_.isEmpty(item) || forceReload) {
        reload().then(function () {
          deferred.resolve(item);
        }, function (reason) {
          deferred.reject(reason);
        });
        forceReload = false;
      } else {
        deferred.resolve(item);
      }

      return deferred.promise;
    }

    return {
      getId: function () {
        return idVal
      },
      setId: function (id) {
        idVal = id
      },
      getMap: getMap,
      reload: reload,
      forceReload: function () {
        forceReload = true;
      }
    }

  })
  .controller('MapsListCtrl', function (Restangular) {
    var self = this;

    Restangular.all('maps').getList({
      first: 0,
      size: 100,
      private: true
    }).then(function (results) {
      self.items = results;
    });

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
  .controller('MapsDetailCtrl', function ($scope, $state, $stateParams, MapsManager) {
    var self = this;

    MapsManager.setId($stateParams.id);
    MapsManager.getMap().then(function (result) {
      self.item = result;
    });

    // このViewには何もないので遷移させる
    $state.go('.markers');
  })
  .controller('MapMarkerListCtrl', function ($scope, $window, $timeout, MapsManager, Restangular) {
    var self = this;

    function init () {
      MapsManager.getMap().then(function (result) {
        self.item = result;
      });
    }

    init();

    this.deleteMarker = function (index) {
      var marker = self.item.markers[index];
      var msg = "「" + marker.name + "」を削除しても宜しいですか？\n（削除後の復元は出来ません！）";
      if ($window.confirm(msg)) {
        Restangular.one('maps', MapsManager.getId()).one('markers', marker.id).remove()
          .then(function () {
            $timeout(function () {
              MapsManager.forceReload();
              init();
            }, 1000);
          });
      }
    };
  })
  .controller('MapsEditCtrl', function ($scope, $stateParams, MapsManager) {
    var self = this;
    this.itemId = $stateParams.id;
    this.lockItemId = true;

    MapsManager.getMap().then(function(result){
      self.item = result;
    });

    // updateMap
    this.click = function () {
      self.lock = true;

      self.item.put()
        .then(function (result) {
          self.lock = false;
          self.alert = {type: 'success', msg: '編集に成功しました'};
          MapsManager.forceReload();
        }, function (reason) {
          self.lock = false;
          self.alert = {type: 'danger', msg: '登録に失敗しました:' + reason.Error}
        });
    };

  })
  .controller('MapMarkerInputCtrl', function ($scope, $timeout, $stateParams, Restangular, MapsManager) {
    var self = this;
    this.lock = false;    // form lock
    this.showMaps = false;
    this.alert = null;
    this.itemIdParam = $stateParams.id || null;

    $timeout(function () {
      self.showMaps = true;
    }, 1000);

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
          self.alert = {type: 'success', msg: '登録に成功しました'};
          MapsManager.forceReload();
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