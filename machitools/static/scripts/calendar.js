'use strict';

angular.module('myApp.calendar',[])
  .factory('CalendarRest', function (Restangular) {
    return Restangular.withConfig(function (config) {
      config.setBaseUrl('https://www.googleapis.com/calendar/v3/calendars');
      config.setDefaultRequestParams({
        key: "AIzaSyCgK3kr9bdc_Qv_SnSJTxAcS1npBGqyRgw"
      });
    });
  }).service('Calendar', function(CalendarRest, Restangular) {

    /**
     * calendarIdの指定時刻から今日の終わりまでのイベント(終日イベントは除く)を取得する
     * @param calId calendar id
     * @param time moment-ed time
     */
    function getTodayData(calId, time) {
      return CalendarRest.all(calId).get('events', {
        orderBy: 'startTime',
        singleEvents: true,
        timeZone: 'Asia/Tokyo',
        timeMin: time.format(),
        timeMax: time.endOf('day').format(),
        maxResults: 10
      }).then(function (result) {
        // 終日イベントは除外する
        var list = [];
        for (var i = 0; i < result.items.length; i++) {
          var item = result.items[i];

          if (!_.isUndefined(item.start.date)) {
            continue;
          }

          list.push(item);

          if (list.length == 3) {
            break;
          }
        }
        return list;
      })
    }

    return {
      getTodayData: getTodayData
    }
  });
