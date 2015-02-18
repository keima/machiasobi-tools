'use strict';

angular.module('myApp', [
  // 3rd party plugins
  'ngCookies',
  'restangular',
  'uiGmapgoogle-maps',
  'ui.router',
  'ui.bootstrap',
  'ncy-angular-breadcrumb',
  'wu.staticGmap',

  // app component
  //'myApp.routing',
  'myApp.calendar',
  'myApp.controller',
  'myApp.services'
])
  .constant('ApiUrl', '/api/v1')
  .config(function (uiGmapGoogleMapApiProvider) {
    uiGmapGoogleMapApiProvider.configure({
      language: 'ja',
      sensor: 'true',
      key: 'AIzaSyAxOlm0zuaBtM7D4dPcOTrUdPrzu4va1cs'
    });
  })
  .config(function (RestangularProvider, ApiUrl) {
    RestangularProvider.setBaseUrl(ApiUrl);
    RestangularProvider.setRequestInterceptor(function (elem, operation) {
      if (operation === "remove") {
        return undefined;
      }
      return elem;
    })
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


;