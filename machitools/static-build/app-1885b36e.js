"use strict";angular.module("myApp",["ngCookies","restangular","uiGmapgoogle-maps","ui.router","ui.bootstrap","ncy-angular-breadcrumb","wu.staticGmap","myApp.calendar","myApp.controller","myApp.services"]).constant("ApiUrl","/api/v1").config(["uiGmapGoogleMapApiProvider",function(t){t.configure({language:"ja",sensor:"true",key:"AIzaSyAxOlm0zuaBtM7D4dPcOTrUdPrzu4va1cs"})}]).config(["RestangularProvider","ApiUrl",function(t,e){t.setBaseUrl(e),t.setRequestInterceptor(function(t,e){return"remove"===e?void 0:t})}]).value("Periods",{day1:{name:"1日目",date:new Date("2014/10/11")},day2:{name:"2日目",date:new Date("2014/10/12")},day3:{name:"3日目",date:new Date("2014/10/13")}}),angular.module("myApp.calendar",[]).factory("CalendarRest",["Restangular",function(t){return t.withConfig(function(t){t.setBaseUrl("https://www.googleapis.com/calendar/v3/calendars"),t.setDefaultRequestParams({key:"AIzaSyCgK3kr9bdc_Qv_SnSJTxAcS1npBGqyRgw"})})}]).service("Calendar",["CalendarRest",function(t){function e(e,a){return t.all(e).get("events",{orderBy:"startTime",singleEvents:!0,timeZone:"Asia/Tokyo",timeMin:a.format(),timeMax:a.endOf("day").format(),maxResults:10}).then(function(t){for(var e=[],a=0;a<t.items.length;a++){var n=t.items[a];if(_.isUndefined(n.start.date)&&(e.push(n),3==e.length))break}return e})}return{getTodayData:e}}]),angular.module("myApp.controller",["myApp.controller.root","myApp.controller.topPage","myApp.controller.traffic","myApp.controller.delay","myApp.controller.event","myApp.controller.news","myApp.controller.maps"]),angular.module("myApp").config(["$stateProvider","$urlRouterProvider",function(t,e){e.otherwise("/"),e.when("/event","/event/list/day1"),e.when("/event/list","/event/list/day1"),t.state("root",{url:"/",templateUrl:"partials/root.html"}).state("traffic",{url:"/traffic","abstract":!0,templateUrl:"partials/traffic/_.html"}).state("traffic.list",{url:"",templateUrl:"partials/traffic/list.html"}).state("traffic.input",{url:"/input",templateUrl:"partials/traffic/input.html"}).state("delay",{url:"/delay","abstract":!0,templateUrl:"partials/delay/_.html"}).state("delay.list",{url:"",templateUrl:"partials/delay/list.html"}).state("delay.input",{url:"",templateUrl:"partials/delay/input.html"}).state("event",{url:"/event",templateUrl:"partials/event/_.html"}).state("event.list",{url:"/list",templateUrl:"partials/event/list.html"}).state("event.list.day",{url:"/:id",templateUrl:"partials/event/list-day.html"}).state("event.input",{url:"/input",templateUrl:"partials/event/input.html"}).state("event.edit",{url:"/input/:id",templateUrl:"partials/event/input.html"}).state("news",{url:"/news","abstract":!0,templateUrl:"partials/news/_.html"}).state("news.list",{url:"/list",templateUrl:"partials/news/list.html"}).state("news.input",{url:"/input",templateUrl:"partials/news/input.html"}).state("news.edit",{url:"/input/:id",templateUrl:"partials/news/input.html"}).state("maps",{url:"/maps","abstract":!0,templateUrl:"partials/maps/_.html"}).state("maps.list",{url:"",templateUrl:"partials/maps/list.html",controller:"MapsListCtrl as ctrl",ncyBreadcrumb:{label:"マップ"}}).state("maps.input",{url:"/input",templateUrl:"partials/maps/input-map.html",ncyBreadcrumb:{parent:"maps.list",label:"作成"}}).state("maps.detail",{url:"/:id",templateUrl:"partials/maps/detail.html",controller:"MapsDetailCtrl as ctrl",ncyBreadcrumb:{parent:"maps.list",label:"詳細"}}).state("maps.detail.markers",{templateUrl:"partials/maps/list-marker.html",controller:"MapMarkerListCtrl as ctrl",ncyBreadcrumb:{skip:!0}}).state("maps.detail.input-marker",{templateUrl:"partials/maps/input-marker.html",controller:"MapMarkerInputCtrl as ctrl",ncyBreadcrumb:{skip:!0}}).state("maps.detail.edit",{templateUrl:"partials/maps/input-map-form.html",controller:"MapsEditCtrl as ctrl",ncyBreadcrumb:{skip:!0}})}]),angular.module("myApp").run(["$rootScope","$state","Restangular","User",function(t,e,a,n){t.$state=e,t.isAdmin=!1,a.all("auth").get("check").then(function(e){n.setUser(e),t.isAdmin=n.isAdmin()},function(t){console.log(t),n.setUser({})})}]),angular.module("myApp.services",["myApp.services.user"]),angular.module("myApp.controller.delay",[]).controller("DelayViewCtrl",["Restangular","Calendar",function(t,e){var a=this;this.now=new Date,this.abs=function(t){return Math.abs(t)},this.places={bizan:{name:"眉山林間ステージ",calendarId:"p-side.net_m9s9a5ut02n6ap1s6prdj92ss4@group.calendar.google.com"},shinmachi:{name:"新町橋東公園",calendarId:"p-side.net_ctrq60t4vsvfavejbkdmbhv3k4@group.calendar.google.com"},corne:{name:"コルネの泉",calendarId:"p-side.net_jo112m9l36p6nlkrv939sb9kr0@group.calendar.google.com"},cinema_entry:{name:"CINEMA前(入り口)",calendarId:"p-side.net_j3mtcq3ejulrovek8kru6vgoe8@group.calendar.google.com"},awagin:{name:"あわぎんホール小ホール",calendarId:"p-side.net_oa45stb6g4h9lqiq5vd1ov844s@group.calendar.google.com"},bunka:{name:"徳島市立文化センター",calendarId:"p-side.net_gocec2ij5sqho46oial3jusn1o@group.calendar.google.com"}},this.calendarData={},angular.forEach(this.places,function(n,r){t.all("delay").get(r).then(function(t){n.item=t;var l=moment().add(-1*t.delay,"minutes");e.getTodayData(n.calendarId,l).then(function(t){a.calendarData[r]=t})},function(){n.item={error:!0,delay:0,message:"SYSTEM: 取得に失敗しました",updatedAt:"---"}})})}]).controller("DelayInputCtrl",["Restangular",function(t){var e=this;this.lock=!1,this.alert=null,this.place=null,this.item={delay:0,message:"",isPostponed:!1},this.click=function(){e.lock=!0,t.all("delay").all(e.place).post(e.item).then(function(){e.lock=!1,e.alert={type:"success",msg:"登録に成功しました"}},function(t){e.lock=!1,e.alert={type:"danger",msg:"登録に失敗しました:"+t.Error}})},this.closeAlert=function(){e.alert=null}}]),angular.module("myApp.controller.event",[]).controller("EventInputCtrl",["$scope","$stateParams","Restangular","Periods",function(t,e,a,n){var r=this;this.itemId=e.id||null,this.lock=!1,this.alert=null;var l=function(t){return t+="",1===t.length&&(t="0"+t),t};this.startAtElements=n,this.selectedStartAt={date:null,time:null,setDate:function(t){var e=t.getFullYear(),a=t.getMonth()+1,n=t.getDate();this.date=new Date(e+"/"+a+"/"+n),this.time=l(t.getHours())+":"+l(t.getMinutes())},getDate:function(){var t=this.date,e=this.time;if(_.isNull(t)||_.isNull(e))return null;var a=t.getFullYear(),n=t.getMonth()+1,r=t.getDate();return console.log(new Date(a+"/"+n+"/"+r+" "+e)),new Date(a+"/"+n+"/"+r+" "+e)}},this.item={id:null,title:null,place:null,message:null,startAt:null,isPublic:!1,isRunning:!1,isFinished:!1},_.isNull(this.itemId)||a.all("events").get(r.itemId).then(function(t){r.item=t,r.selectedStartAt.setDate(new Date(t.startAt))}),this.click=function(){r.item.startAt=r.selectedStartAt.getDate(),r.lock=!0;var t;t=_.isNull(r.itemId)?a.all("events").post(r.item):a.all("events").all(r.itemId).post(r.item),t.then(function(){r.lock=!1,r.alert={type:"success",msg:"登録に成功しました"}},function(t){r.lock=!1,r.alert={type:"danger",msg:"登録に失敗しました:"+t.Error}})}}]).controller("EventListCtrl",["Restangular","User","Periods",function(t,e,a){this.periods=a,this.now=new Date}]).controller("EventListDayCtrl",["$stateParams","Restangular","User","Periods",function(t,e,a,n){var r=this;this.isAdmin=a.isAdmin();var l=moment(n[t.id].date),i=l.clone().endOf("days");e.all("events").getList({first:0,size:100,"private":!0,startAt:l.toJSON(),endAt:i.toJSON()}).then(function(t){r.items=t},function(t){console.log(t)})}]),angular.module("myApp.controller.maps",[]).service("MapsManager",["$q","Restangular",function(t,e){function a(){var a=t.defer();return _.isUndefined(r)?a.reject("MapsManager.id is not set!"):e.all("maps").get(r).then(function(t){i=t,a.resolve(t)},function(t){a.reject(t)}),a.promise}function n(){var e=t.defer();return _.isEmpty(i)||l?(a().then(function(){e.resolve(i)},function(t){e.reject(t)}),l=!1):e.resolve(i),e.promise}var r,l=!1,i={};return{getId:function(){return r},setId:function(t){r=t},getMap:n,reload:a,forceReload:function(){l=!0}}}]).controller("MapsListCtrl",["Restangular",function(t){var e=this;t.all("maps").getList({first:0,size:100,"private":!0}).then(function(t){e.items=t})}]).controller("MapsInputCtrl",["$stateParams","Restangular",function(t,e){var a=this;this.itemIdParam=t.id||null,this.lock=!1,this.alert=null,this.itemId=null,this.item={name:null,isPublic:!1},_.isNull(this.itemIdParam)||e.all("maps").get(a.itemIdParam).then(function(t){a.itemId=t.id,a.item={name:t.name,isPublic:t.isPublic}}),this.click=function(){a.lock=!0,e.all("maps").all(a.itemId).post(a.item).then(function(){a.lock=!1,a.alert={type:"success",msg:"登録に成功しました"}},function(t){a.lock=!1,a.alert={type:"danger",msg:"登録に失敗しました:"+t.Error}})}}]).controller("MapsDetailCtrl",["$scope","$state","$stateParams","MapsManager",function(t,e,a,n){var r=this;n.setId(a.id),n.getMap().then(function(t){r.item=t}),e.go(".markers")}]).controller("MapMarkerListCtrl",["$scope","$window","$timeout","MapsManager","Restangular",function(t,e,a,n,r){function l(){n.getMap().then(function(t){i.item=t})}var i=this;l(),this.deleteMarker=function(t){var s=i.item.markers[t],o="「"+s.name+"」を削除しても宜しいですか？\n（削除後の復元は出来ません！）";e.confirm(o)&&r.one("maps",n.getId()).one("markers",s.id).remove().then(function(){a(function(){n.forceReload(),l()},1e3)})}}]).controller("MapsEditCtrl",["$scope","$stateParams","MapsManager",function(t,e,a){var n=this;this.itemId=e.id,this.lockItemId=!0,a.getMap().then(function(t){n.item=t}),this.click=function(){n.lock=!0,n.item.put().then(function(){n.lock=!1,n.alert={type:"success",msg:"編集に成功しました"},a.forceReload()},function(t){n.lock=!1,n.alert={type:"danger",msg:"登録に失敗しました:"+t.Error}})}}]).controller("MapMarkerInputCtrl",["$scope","$timeout","$stateParams","Restangular","MapsManager",function(t,e,a,n,r){var l=this;this.lock=!1,this.showMaps=!1,this.alert=null,this.itemIdParam=a.id||null,e(function(){l.showMaps=!0},1e3),this.place={name:null,description:null,order:null},this.map={center:{latitude:34.071144,longitude:134.548529},zoom:18},this.marker={id:0,coords:{latitude:34.071144,longitude:134.548529},options:{draggable:!0}},this.postMapMarker=function(){l.lock=!0,n.all("maps").all(l.itemIdParam).all("markers").post({name:l.place.name,description:l.place.description,order:l.place.order,coords:{latitude:l.marker.coords.latitude,longitude:l.marker.coords.longitude}}).then(function(){l.lock=!1,l.alert={type:"success",msg:"登録に成功しました"},r.forceReload()},function(t){l.lock=!1,l.alert={type:"danger",msg:"登録に失敗しました:"+t.Error}})},this.moveMarkerToCenter=function(){l.marker.coords={latitude:l.map.center.latitude,longitude:l.map.center.longitude}}}]),angular.module("myApp.controller.news",[]).controller("NewsInputCtrl",["$stateParams","Restangular",function(t,e){var a=this;this.itemIdParam=t.id||null,this.lock=!1,this.alert=null,this.newsId=null,this.newsItem={Title:null,Article:null,IsPublic:!1},_.isNull(this.itemIdParam)||e.all("news").get(a.itemIdParam).then(function(t){a.newsId=t.Id,a.newsItem={Title:t.Title,Article:t.Article,IsPublic:t.IsPublic}}),this.click=function(){a.lock=!0,e.all("news").all(a.newsId).post(a.newsItem).then(function(){a.lock=!1,a.alert={type:"success",msg:"登録に成功しました"}},function(t){a.lock=!1,a.alert={type:"danger",msg:"登録に失敗しました:"+t.Error}})}}]).controller("NewsListCtrl",["Restangular","User",function(t,e){var a=this;this.isAdmin=e.isAdmin(),t.all("news").getList({first:0,size:100,"private":!0}).then(function(t){a.items=t},function(t){console.log(t)})}]),angular.module("myApp.controller.root",[]).controller("RootCtrl",["$scope","$cookies",function(t,e){var a=this;this.hideHeader="true"===e.hideHeader,t.$watch(function(){return a.hideHeader},function(t){e.hideHeader=t})}]).controller("HeaderCtrl",["$scope","ApiUrl","User",function(t,e,a){var n=this;this.apiUrl=e,this.loggedin=a.isLogin(),t.$on(a.BROADCAST_NAME_CHANGED,function(){n.loggedin=a.isLogin()})}]).controller("TabCtrl",["$scope","User",function(t,e){var a=this;a.isAdmin=e.isAdmin(),t.$on(e.BROADCAST_NAME_CHANGED,function(){a.isAdmin=e.isAdmin()})}]),angular.module("myApp.controller.topPage",[]).controller("TopPageCtrl",["$scope","User",function(t,e){function a(){n.showError=e.isAdmin()?!1:e.isLogin()}var n=this;this.showError=!1,a(),t.$on(e.BROADCAST_NAME_CHANGED,function(){a()})}]),angular.module("myApp.controller.traffic",[]).controller("TrafficViewCtrl",["Restangular",function(t){this.now=new Date,this.transits=[{name:"ロープウェイ乗り場",id:"ropeway",places:[{name:"山麓駅(阿波おどり会)",direction:"inbound"},{name:"山頂駅",direction:"outbound"}]},{name:"シャトルバス乗り場",id:"bus",places:[{name:"山麓駅(阿波踊り会館 前)",direction:"inbound"},{name:"山頂駅(かんぽの宿 前)",direction:"outbound"}]}],this.transits.forEach(function(e){var a=e.id;e.places.forEach(function(e){var n=e.direction;t.all("traffic").all(a).get(n).then(function(t){e.item=t},function(){e.item={Waiting:"---",Message:"SYSTEM: 取得に失敗しました",updatedAt:"---"}})})})}]).controller("TrafficInputCtrl",["$scope","$cookies","Restangular",function(t,e,a){function n(){_.isUndefined(r.traffic)||_.isUndefined(r.direction)||a.all("traffic").all(r.traffic).get(r.direction).then(function(t){r.trafficItem.Message=t.Message},function(){r.trafficItem.Message=""})}var r=this;this.lock=!1,this.alert=null,this.trafficItem={Waiting:null,Message:null},this.traffic=e.traffic,this.direction=e.direction,this.click=function(){r.lock=!0,a.all("traffic").all(r.traffic).all(r.direction).post(r.trafficItem).then(function(){r.lock=!1,r.alert={type:"success",msg:"登録に成功しました"}},function(t){r.lock=!1,r.alert={type:"danger",msg:"登録に失敗しました:"+t.Error}})},this.closeAlert=function(){r.alert=null},t.$watch(function(){return r.traffic},function(t){e.traffic=t,n()}),t.$watch(function(){return r.direction},function(t){e.direction=t,n()})}]),angular.module("myApp.services.user",[]).service("User",["$rootScope",function(t){function e(e){l=e,t.$broadcast(i)}function a(){return l}function n(){return!_.isEmpty(l)}function r(){return n()&&l.Admin}var l={},i="UserDataIsChanged";return{BROADCAST_NAME_CHANGED:i,setUser:e,getUser:a,isLogin:n,isAdmin:r}}]);