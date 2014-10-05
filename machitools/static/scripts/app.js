'use strict';

angular.module('myApp', ['ngCookies', 'restangular', 'ui.router', 'ui.bootstrap'])
  .constant('ApiUrl', '/api/v1')
  // Restangular config
  .config(function (RestangularProvider, ApiUrl) {
    RestangularProvider.setBaseUrl(ApiUrl);
  })
  // AngularUI Router
  .config(function ($stateProvider, $urlRouterProvider) {
    $urlRouterProvider.otherwise("/");

    $urlRouterProvider.when('/traffic', '/traffic/list');
    $urlRouterProvider.when('/event', '/event/list/day1');
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

      .state('event', {
        url: '/event',
        templateUrl: "partials/event.html"
      })
      .state('event.list', {
        url: '/list',
        templateUrl: "partials/event-list.html"
      })
      .state('event.list.day', {
        url: '/:id',
        templateUrl: "partials/event-list-day.html"
      })
      .state('event.input', {
        url: '/input',
        templateUrl: "partials/event-input.html"
      })
      .state('event.edit', {
        url: '/input/:id',
        templateUrl: "partials/event-input.html"
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
      .state('news.edit', {
        url: '/input/:id',
        templateUrl: 'partials/news-input.html'
      })
    ;
  })

  .service('User', function ($rootScope) {
    var user = {};
    var BROADCAST_NAME_CHANGED = 'UserDataIsChanged';

    function setUser (_user) {
      user = _user;
      $rootScope.$broadcast(BROADCAST_NAME_CHANGED);
    }

    function getUser () {
      return user;
    }

    function isLogin () {
      return !_.isEmpty(user)
    }

    function isAdmin () {
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

    function renderAlertIfNeeded () {
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

  .value('Periods', {
    day1: {
      name: "1日目", date: new Date("2014/10/11")
    },
    day2: {
      name: "2日目", date: new Date("2014/10/12")
    },
    day3: {
      name: "3日目", date: new Date("2014/10/13")
    }
  })
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
  })
  .controller('EventListDayCtrl', function ($stateParams, Restangular, User, Periods) {
    var self = this;

    this.isAdmin = User.isAdmin();

    var startAt = Periods[$stateParams.id].date;
    var endAt = new Date(startAt.getTime());
    endAt.setDate(startAt.getDate() + 1);

    Restangular.all('events').getList({
      first: 0,
      size: 100,
      private: true,
      startAt: startAt,
      endAt: endAt
    }).then(function (results) {
      self.items = results;
    }, function (reason) {
      console.log(reason);
    });
  })

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

    if (!_.isNull(this.itemIdParam)) {
      Restangular.all('news').get(self.itemIdParam)
        .then(function (result) {
          self.newsId = result.Id;
          self.newsItem = {
            'Title': result.Title,
            'Article': result.Article,
            'IsPublic': result.IsPublic
          }
        });
    }

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


  .controller('NewsListCtrl', function (Restangular, User) {
    var self = this;

    this.isAdmin = User.isAdmin();

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