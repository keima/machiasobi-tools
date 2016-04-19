(function() {
  "use strict";

  angular.module("myrestangular", ["restangular"])
    .constant('ApiUrl', '/api/v1')
    .config(function(RestangularProvider, ApiUrl) {
      RestangularProvider.setBaseUrl(ApiUrl);
      RestangularProvider.setRequestInterceptor(function(elem, operation) {
        if (operation === "remove") {
          return undefined;
        }
        return elem;
      })
    });
}());
