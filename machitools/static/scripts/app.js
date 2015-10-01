'use strict';

angular.module('myApp', [
  // 3rd party plugins
  'ngCookies',
  'myrestangular',
  'uiGmapgoogle-maps',
  'ui.router',
  'ui.bootstrap',
  'ncy-angular-breadcrumb',
  'wu.staticGmap',
  'dndLists',

  // app component
  'myApp.calendar',
  'myApp.controller',
  'myApp.services'
])
  .config(function (uiGmapGoogleMapApiProvider) {
    uiGmapGoogleMapApiProvider.configure({
      language: 'ja',
      sensor: 'true',
      key: 'AIzaSyAxOlm0zuaBtM7D4dPcOTrUdPrzu4va1cs'
    });
  })
;
