'use strict';

angular.module('myApp', ['ngRoute', 'restangular', 'ui.router', 'ui.bootstrap'])
  .constant('ApiUrl', '/api/v1')
  // Restangular config
  .config(function (RestangularProvider, ApiUrl) {
    RestangularProvider.setBaseUrl(ApiUrl);
  })
  // RouteProvider config
  .config(function ($routeProvider) {

  })
  // AngularUI Router
  .config(function ($stateProvider, $urlRouterProvider) {
    $urlRouterProvider.otherwise("/");

    $urlRouterProvider.when('/traffic', '/traffic/list');
    $urlRouterProvider.when('/news', '/news/list');

    $stateProvider
      .state('root', {
        url: '/',
        templateUrl: "partials/root.html"
      })

      .state('traffic', {
        url: '/traffic',
        templateUrl: "partials/traffic.html"
      })
      .state('traffic.list', {
        url: '/list',
        templateUrl: "partials/traffic-list.html"
      })
      .state('traffic.input', {
        url: '/input',
        templateUrl: "partials/traffic-input.html"
      })

      .state('news', {
        url: '/news',
        templateUrl: "partials/news.html"
      })
      .state('news.list', {
        url: '/list',
        templateUrl: "partials/news-list.html"
      })
      .state('news.input', {
        url: '/input',
        templateUrl: 'partials/news-input.html'
      })
    ;
  })

  .service('User', function ($rootScope) {
    var user = {};
    var BROADCAST_NAME_CHANGED = 'UserDataIsChanged';

    function setUser(_user) {
      user = _user;
      $rootScope.$broadcast(BROADCAST_NAME_CHANGED);
    }

    function getUser() {
      return user;
    }

    function isLogin() {
      return !_.isEmpty(user)
    }

    function isAdmin() {
      return isLogin() && user.Admin;
    }

    return {
      BROADCAST_NAME_CHANGED: BROADCAST_NAME_CHANGED,

      setUser: setUser,
      getUser: getUser,
      isLogin: isLogin,
      isAdmin: isAdmin
    }

  })


  .run(function ($rootScope, $state, Restangular, User) {
    // convenience state
    $rootScope.$state = $state;

    // get user data
    Restangular.all('auth').get('check')
      .then(function (result) {
        console.log(result);
        User.setUser(result);
      }, function (reason) {
        console.log(reason);
        User.setUser({});
      });

  })

  .controller('RootCtrl', function ($scope, User) {
    var self = this;
    this.showError = false;

    function renderAlertIfNeeded() {
      if (User.isAdmin()) {
        self.showError = false;
      } else {
        self.showError = User.isLogin();
      }
    }

    renderAlertIfNeeded();

    $scope.$on(User.BROADCAST_NAME_CHANGED, function () {
      renderAlertIfNeeded();
    })

  })

  // Header Controller
  .controller('HeaderCtrl', function ($scope, ApiUrl, User) {
    var self = this;
    this.apiUrl = ApiUrl;
    this.loggedin = User.isLogin();

    $scope.$on(User.BROADCAST_NAME_CHANGED, function () {
      self.loggedin = User.isLogin();
    })
  })

  .controller('TabCtrl', function ($scope, User) {
    var self = this;

    self.isAdmin = User.isAdmin();

    $scope.$on(User.BROADCAST_NAME_CHANGED, function () {
      self.isAdmin = User.isAdmin();
    })
  })


  .controller('TrafficViewCtrl', function (Restangular) {
    var self = this;
    var traffics = ['ropeway', 'bus'];
    var directions = ['inbound', 'outbound'];

    this.ropeway = {
      inbound: {}, outbound: {}
    };
    this.bus = {
      inbound: {}, outbound: {}
    };

    traffics.forEach(function (traffic) {
      directions.forEach(function (direction) {
        Restangular.all('traffic').all(traffic).get(direction)
          .then(function (result) {
            self[traffic][direction] = {
              'Waiting': result.Waiting,
              'Message': result.Message
            }
          }, function () {
            self[traffic][direction] = {
              'Waiting': '---',
              'Message': 'SYSTEM: 取得に失敗しました'
            }
          });
      });
    });
  })


  .controller('TrafficInputCtrl', function ($scope, Restangular) {
    var self = this;

    // form lock
    this.lock = false;

    this.alert = null;

    this.trafficItem = {
      'Waiting': null,
      'Message': null
    };
    this.traffic;
    this.direction;

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
      getTrafficMessage();
    });

    $scope.$watch(function () {
      return self.direction;
    }, function (newVal, oldVal) {
      getTrafficMessage();
    });

    function getTrafficMessage() {
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


  .controller('NewsInputCtrl', function (Restangular) {
    var self = this;

    this.lock = false;    // form lock
    this.alert = null;

    this.newsId = null;
    this.newsItem = {
      'Title': null,
      'Article': null,
      'IsPublic': false
    };

    this.click = function () {
      self.lock = true;
      Restangular.all('news').all(self.newsId).post(self.newsItem)
        .then(function (result) {
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
        console.log(results)
        self.items = results;
      }, function (reason) {
        console.log(reason);
      })
  })
;