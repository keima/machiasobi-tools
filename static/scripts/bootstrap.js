(function() {
  "use strict";

  var myAppModule = angular.module("myApp");
  fetchData().then(
    bootstrapApp, function(error) {
      var $div = document.getElementById("initialize-error");
      $div.removeAttribute("style");
    }
  );

  function fetchData() {
    var initInjector = angular.injector(["ng", "myrestangular"]);
    var $q = initInjector.get("$q");
    var Restangular = initInjector.get("Restangular");

    return $q.all([
      Restangular.all("periods").getList().then(function(result) {
        myAppModule.constant("Periods", result.plain());
      }),
      Restangular.all("calendars").getList({visibility: "all"}).then(function(result) {
        myAppModule.constant("Calendars", result.plain());
      })
    ]);
  }

  function bootstrapApp() {
    angular.element(document.body).ready(function() {
      angular.bootstrap(document.body, ["myApp"]);
    });
  }

}());
