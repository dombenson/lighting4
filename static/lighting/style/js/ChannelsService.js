var module = angular.module('ChannelsService', ['LightingAPIService', 'LightingWebSocketService']);

module.factory('ChannelsService', ['$q', '$rootScope', 'LightingAPIService', 'LightingWebSocketService', function($q, $rootScope, LightingAPIService, LightingWebSocketService) {
  var channelsCachedPromise = null;

  var channelListeners = {};

  var currentChannelLevels = {};
  var lastSeenSeqNos = {};

  LightingWebSocketService.attachSocketListener('uC', function(data){
    var channelId = data.c.i;
    var channelLevel = data.c.l;
    var seqNo = data.c.s;

    $rootScope.$apply(function() {
      notifyChannelLister(channelId, channelLevel, seqNo);
    });

  });

  function notifyChannelLister(channelId, channelLevel, seqNo) {
    var lastSeenSeqNo = 0;

    if (lastSeenSeqNos.hasOwnProperty(channelId)) {
      lastSeenSeqNo = lastSeenSeqNos[channelId];
    }
    if (seqNo > lastSeenSeqNo) {
      lastSeenSeqNos[channelId] = seqNo;
      currentChannelLevels[channelId] = channelLevel;

      if(channelListeners.hasOwnProperty(channelId)) {
        var arr = channelListeners[channelId];
        for (var i = 0; i < arr.length; i++) {
          arr[i](channelLevel);
        }
      }
    }
  }

  var Service = {};

  Service.getChannels = function(forceReload) {
    if (!channelsCachedPromise || forceReload) {
      channelsCachedPromise = LightingWebSocketService.channelState();

      channelsCachedPromise.then(function(channels) {
        currentChannelLevels = {};
        channels.forEach(function(channel) {
          var lastSeenSeqNo = 0;

          if (lastSeenSeqNos.hasOwnProperty(channel.id)) {
            lastSeenSeqNo = lastSeenSeqNos[channel.id];
          }

          if (lastSeenSeqNo > channel.seqNo) {
            channel.currentLevel = currentChannelLevels[channel.id];
          } else {
            currentChannelLevels[channel.id] = channel.currentLevel;
            lastSeenSeqNos[channel.id] = channel.seqNo;
          }
        });
      });
    }

    return channelsCachedPromise;
  };

  Service.updateChannel = function(id, level) {
    if (lastSeenSeqNos.hasOwnProperty(id)) {
      LightingWebSocketService.updateChannel(id, level, lastSeenSeqNos[id]);
    }
  };

  Service.attachChannelListener = function(channelId, callback) {
    if (!channelListeners.hasOwnProperty(channelId)) {
      channelListeners[channelId] = [];
    }
    channelListeners[channelId].push(callback);
  };
  Service.detachChannelListener = function(channelId, callback) {
    if (channelListeners.hasOwnProperty(channelId)) {
      var array = channelListeners[channelId];
      var index = array.indexOf(callback);

      if (index > -1) {
        array.splice(index, 1);
      }
    }
  };

  Service.currentChannelLevel = function(channelId) {
    return currentChannelLevels[channelId];
  }

  return Service;
}]);
