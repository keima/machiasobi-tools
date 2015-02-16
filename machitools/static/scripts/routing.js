'use strict';

/**
 * Routing Rules
 */

angular.module('myApp')
  .config(function ($stateProvider, $urlRouterProvider) {
    $urlRouterProvider.otherwise("/");

    $urlRouterProvider.when('/event', '/event/list/day1');
    $urlRouterProvider.when('/event/list', '/event/list/day1');


    $stateProvider
      .state('root', {
        url: '/',
        templateUrl: "partials/root.html"
      })

      .state('traffic', {
        url: '/traffic',
        abstract: true,
        templateUrl: "partials/traffic/_.html"
      })
      .state('traffic.list', {
        url: '',
        templateUrl: "partials/traffic/list.html"
      })
      .state('traffic.input', {
        url: '/input',
        templateUrl: "partials/traffic/input.html"
      })

      .state('delay', {
        url: '/delay',
        abstract: true,
        templateUrl: "partials/delay/_.html"
      })
      .state('delay.list', {
        url: '',
        templateUrl: "partials/delay/list.html"
      })
      .state('delay.input', {
        url: '',
        templateUrl: "partials/delay/input.html"
      })

      .state('event', {
        url: '/event',
        templateUrl: "partials/event/_.html"
      })
      .state('event.list', {
        url: '/list',
        templateUrl: "partials/event/list.html"
      })
      .state('event.list.day', {
        url: '/:id',
        templateUrl: "partials/event/list-day.html"
      })
      .state('event.input', {
        url: '/input',
        templateUrl: "partials/event/input.html"
      })
      .state('event.edit', {
        url: '/input/:id',
        templateUrl: "partials/event/input.html"
      })

      .state('news', {
        url: '/news',
        abstract: true,
        templateUrl: "partials/news/_.html"
      })
      .state('news.list', {
        url: '/list',
        templateUrl: "partials/news/list.html"
      })
      .state('news.input', {
        url: '/input',
        templateUrl: 'partials/news/input.html'
      })
      .state('news.edit', {
        url: '/input/:id',
        templateUrl: 'partials/news/input.html'
      })

      .state('maps', {
        url: '/maps',
        abstract: true,
        templateUrl: 'partials/maps/_.html'
      })
      .state('maps.list', {
        url: '',
        templateUrl: 'partials/maps/list.html',
        controller: 'MapsListCtrl as ctrl',
        ncyBreadcrumb: {
          label: 'マップ'
        }
      })
      .state('maps.input', {
        url: '/input',
        templateUrl: 'partials/maps/input-map.html',
        ncyBreadcrumb: {
          parent: "maps.list",
          label: '作成'
        }
      })
      .state('maps.detail', {
        url: '/:id',
        templateUrl: 'partials/maps/detail.html',
        controller: 'MapsDetailCtrl as ctrl',
        ncyBreadcrumb: {
          parent: "maps.list",
          label: '詳細'
        }
      })

      //.state('maps.edit', {
      //  url: '/input/:id',
      //  templateUrl: 'partials/maps/input.html'
      //})
    ;
  })
;